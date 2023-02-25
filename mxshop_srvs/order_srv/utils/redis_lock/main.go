package main

import (
	"fmt"
	"sync"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.10.105:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "421"

	gNum := 2
	var wg sync.WaitGroup
	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func() {
			defer wg.Done()
			mutex := rs.NewMutex(mutexname)

			fmt.Println("开始获取锁")
			// Obtain a lock for our given mutex. After this is successful, no one else
			// can obtain the same lock (the same mutex name) until we unlock it.
			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			// Do your work that requires the lock.
			fmt.Println("获取锁success")
			time.Sleep(8 * time.Second)

			fmt.Println("开始释放锁")
			// Release the lock so other processes or threads can obtain a lock.
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			fmt.Println("释放锁成功")
		}()
	}
	wg.Wait()
}