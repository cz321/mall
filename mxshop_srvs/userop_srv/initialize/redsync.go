package initialize

import (
	"fmt"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"mxshop_srvs/userop_srv/global"
)

func InitRedsync() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will use while communicating with Redis.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d",global.ServerConfig.RedisConfig.Host,global.ServerConfig.RedisConfig.Port),
	})
	pool := goredis.NewPool(client)

	// Create an instance of redisync to be used to obtain a mutual exclusion lock.
	global.Redsync = redsync.New(pool)
}