package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type Keys struct {
	Host string
	Port int32
	Username string
	Password string
	Dbname   string
	SslMode string
}

func TestNewDB(t *testing.T) {
	keys := Keys{
		Host:     "192.168.2.182",
		Port: 5432,
		Username: "postgres",
		Password: "msdnmm",
		Dbname:   "postgres",
		SslMode :"disable",
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%v TimeZone=Asia/Shanghai", keys.Host, keys.Username, keys.Password, keys.Dbname, keys.Port, keys.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Ping数据库
	sqlDB, err := db.DB()
	if err != nil {
		panic("无法获取数据库连接")
	}
	err = sqlDB.Ping()
	if err != nil {
		panic("无法ping数据库")
	}
}

// 参考 https://redis.uptrace.dev/guide/go-redis-sentinel.html#redis-server-client
func TestNewCache(t *testing.T) {
	// 默认方式, 稳定:
	// client := redis.NewFailoverClient(&redis.FailoverOptions{
	// 	MasterName:    "master",
	// 	SentinelAddrs: []string{"192.168.2.155:26379", "192.168.2.158:26379", "192.168.2.152:26379"},
	// 	Password:      "263393", // 如果有密码，请填写
	// 	DB:            0,
	// })

	// 从 v8 开始，您可以使用实验性 NewFailoverClusterClient 命令将只读命令路由到从节点
	client := redis.NewClient(&redis.Options{
		// Addr:     "192.168.2.158:6379",
		Protocol: 3,
		Addr:     "192.168.2.192:6379",
		Username: "default", // redis实例的用户名, 非哨兵节点名
		Password: "msdnmm,.",  // redis实例的用户密码, 如果有密码，请填写
		DB:       0,
	})
	// client := redis.NewFailoverClusterClient(&redis.FailoverOptions{
	// 	// MasterName:              "master1",                                                                     // 主节点master名
	// 	// SentinelAddrs:           []string{"192.168.2.155:6379", "192.168.2.158:6379", "192.168.2.152:6379"}, // 哨兵节点
	// 	// ClientName:              "",
	// 	// SentinelUsername:        "master1", // 哨兵节点的账号
	// 	// SentinelPassword:        "263393",  // 哨兵节点的密码
	// 	// RouteByLatency:          true, // 将只读命令路由到从节点
	// 	// RouteRandomly:           true, // 将只读命令路由到从节点
	//
	// 	ReplicaOnly:             false,
	// 	UseDisconnectedReplicas: false,
	// 	Dialer:                  nil,
	// 	OnConnect:               nil,
	// 	Protocol:                0,
	// 	Username:                "master1", // redis实例的用户名, 非哨兵节点名
	// 	Password:                "263393",  // redis实例的用户密码, 如果有密码，请填写
	// 	DB:                      0,
	// 	MaxRetries:              0,
	// 	MinRetryBackoff:         0,
	// 	MaxRetryBackoff:         0,
	// 	DialTimeout:             0,
	// 	ReadTimeout:             0,
	// 	WriteTimeout:            0,
	// 	ContextTimeoutEnabled:   false,
	// 	PoolFIFO:                false,
	// 	PoolSize:                0,
	// 	PoolTimeout:             0,
	// 	MinIdleConns:            0,
	// 	MaxIdleConns:            0,
	// 	MaxActiveConns:          0,
	// 	ConnMaxIdleTime:         0,
	// 	ConnMaxLifetime:         0,
	// 	TLSConfig:               nil,
	// 	DisableIndentity:        false,
	// })

	// 设置一个键值对
	err := client.Set(context.Background(), "example_key", "example_value2", 0).Err()
	err = client.Set(context.Background(), "example_key2", "example_value3", 0).Err()
	err = client.Set(context.Background(), "example_key3", "example_value31", 0).Err()
	err = client.Set(context.Background(), "example_key4", "example_value34", 0).Err()
	if err != nil {
		panic(err)
	}

	// 获取键值对
	val, err := client.Get(context.Background(), "example_key").Result()
	val2, err := client.Get(context.Background(), "example_key2").Result()
	val3, err := client.Get(context.Background(), "example_key3").Result()
	val4, err := client.Get(context.Background(), "example_key4").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("example_key", val)
	fmt.Println("example_key2", val2)
	fmt.Println("example_key3", val3)
	fmt.Println("example_key4", val4)

	// 关闭连接
	if err := client.Close(); err != nil {
		panic(err)
	}
}
