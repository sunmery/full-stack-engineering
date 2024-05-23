package data

import (
	"backend/internal/conf"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun/extra/bunotel"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB, NewCache)

// Data .
type Data struct {
	db *bun.DB
	rdb *redis.Client
}

// NewData .
func NewData(
	c *conf.Data,
	logger log.Logger,
	db *bun.DB,
	rdb *redis.Client,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:db,
		rdb:rdb,
	}, cleanup, nil
}

func NewDB(c *conf.Data) *bun.DB {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&TimeZone=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Dbname,
		c.Database.SslMode,
		c.Database.TimeZone,
	))))
	// bun.WithDiscardUnknownColumns() 若要使应用在迁移过程中对错误具有更强的弹性，可以调整 Bun 以丢弃生产中的未知列：
	db := bun.NewDB(sqlDb, pgdialect.New(), bun.WithDiscardUnknownColumns())

	// Go 程序的并发执行使用的 CPU 核心数设置, 0为最大可用值，并且保持当前的设置
	db.SetMaxOpenConns(int(c.Database.MaxOpenConn) * runtime.GOMAXPROCS(0))
	db.SetMaxIdleConns(int(c.Database.MaxIdleConn) * runtime.GOMAXPROCS(0))

	// otel 可观测性套件
	db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(c.Database.Dbname)))

	return db
}

func NewCache(c *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Protocol:     3,
		Addr:         c.Redis.Addr,
		Username:     c.Redis.User,
		Password:     c.Redis.Pass,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})

	return rdb
}
