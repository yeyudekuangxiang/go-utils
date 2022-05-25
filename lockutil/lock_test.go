package lockutil

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if rdb.Ping(context.Background()).Err() != nil {
		log.Println("redis not collection skip lock test")
		return
	}
	lock := RedisDistributedLock{
		Redis:  rdb,
		Prefix: "test",
	}

	lock.UnLock("unit_lock_test")
	var lockNum int64
	timer := time.NewTimer(time.Millisecond * 5000)
	for j := 0; j < 10000; j++ {
		go func() {
			if lock.Lock("unit_lock_test", time.Second*8) {
				atomic.AddInt64(&lockNum, 1)
			}
		}()
	}
	<-timer.C
	assert.Equal(t, int64(1), lockNum)
}
