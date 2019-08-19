package db

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
	"log"
	"os"
	"sync"
)

type RedisConnectPool struct{}

var redisInstance *RedisConnectPool
var redisOnce sync.Once

var client *redis.Client

func init() {
	redis.SetLogger(log.New(os.Stderr, "redis: ", log.LstdFlags))
}

func GetRedisConnection() *RedisConnectPool {
	redisOnce.Do(func() {
		redisInstance = &RedisConnectPool{}
	})
	return redisInstance
}

func (r *RedisConnectPool) InitialRedisClient() bool {
	host := beego.AppConfig.String("redis_host")
	db, _ := beego.AppConfig.Int("redis_db")

	client = redis.NewClient(&redis.Options{
		Addr:     host, // use default Addr
		Password: "",   // no password set
		DB:       db,   // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

func (r *RedisConnectPool) GetRedisClient() *redis.Client {
	return client
}
