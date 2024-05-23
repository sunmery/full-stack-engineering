package data

import (
	"backend/internal/conf"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	slog "log"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB, NewCache)

// Data .
type Data struct {
	db *gorm.DB
	rdb *redis.Client
}

// NewData .
func NewData(
	c *conf.Data,
	logger log.Logger,
	db *gorm.DB,
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

func NewDB(c *conf.Data) *gorm.DB {
	// 终端打印输入sql执行记录
	sqlLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢sql查询阈值
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		})

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			c.Database.Host,
			c.Database.Username,
			c.Database.Password,
			c.Database.Dbname,
			c.Database.Port,
			c.Database.SslMode,
			c.Database.TimeZone,
		),
		PreferSimpleProtocol: false, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger:                                   sqlLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用使用单数数据库表名, true:是, false: 不是
		},
	})
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		panic("failed to connect database")
	}

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
