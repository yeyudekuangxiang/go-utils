package lockutil

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

// DistributedLock 分布式锁
type DistributedLock interface {
	// Lock 对key进行加唯一锁
	//加锁成功后expireIn为锁自动失效时间
	//加锁成功返回true 加锁失败返回false
	Lock(key string, expireIn time.Duration) bool
	// LockWait 对key进行加唯一锁
	//加锁失败时将会重试直到加锁成功
	//加锁成功后expireIn为锁自动失效时间
	LockWait(key string, expireIn time.Duration)
	// LockNum 对key进行加锁
	//锁不存在时直接加锁成功 并且可加锁数量设为 num-1
	//后续每次加锁成功后可加锁数量减1 直到数量为0且时加锁失败
	//加锁成功返回true 加锁失败返回false
	//加锁成功后expireIn为锁自动失效时间
	LockNum(key string, num int, expireIn time.Duration) bool
	// LockNumWait 对key进行加锁
	//锁不存在时直接加锁成功 并且可加锁数量设为 num-1
	//后续每次加锁成功后可加锁数量减1 直到数量为0且时加锁失败
	//此方法加锁失败时将会重试直到加锁成功
	//加锁成功后expireIn为锁自动失效时间
	LockNumWait(key string, num int, expireIn time.Duration)
	// UnLock 删除锁
	UnLock(key string)
}

type RedisDistributedLock struct {
	Redis  *redis.Client
	Prefix string
}

func (lock *RedisDistributedLock) FormatKey(key string) string {
	return fmt.Sprintf(lock.Prefix + key)
}
func (lock *RedisDistributedLock) Lock(key string, expireIn time.Duration) bool {
	key = lock.FormatKey(key)
	ok, err := lock.Redis.SetNX(context.Background(), key, 1, expireIn).Result()
	if err != nil {
		panic(err)
	}
	if !ok {
		return false
	}
	return true
}
func (lock *RedisDistributedLock) LockWait(key string, expireIn time.Duration) {
	key = lock.FormatKey(key)
	for !lock.Lock(key, expireIn) {
		time.Sleep(time.Millisecond * 100)
	}
}
func (lock *RedisDistributedLock) UnLock(key string) {
	key = lock.FormatKey(key)
	_, err := lock.Redis.Del(context.Background(), key).Result()
	if err != nil {
		panic(err)
	}
}

const redisLua = `local num = redis.pcall("GET",KEYS[1])
if num == false then
redis.pcall("SETEX",KEYS[1],ARGV[2],ARGV[1]-1)
return 1
elseif tonumber(num) <= 0 then
return 0
else
redis.pcall("DECR",KEYS[1])
return 1
end`

// LockNum expireIn四舍五入到秒
func (lock *RedisDistributedLock) LockNum(key string, num int, expireIn time.Duration) bool {
	key = lock.FormatKey(key)
	result, err := lock.Redis.Eval(context.Background(), redisLua, []string{key}, num, int(math.Round(expireIn.Seconds()))).Int()
	if err != nil {
		panic(err)
	}
	return result == 1
}

// LockNumWait expireIn四舍五入到秒
func (lock *RedisDistributedLock) LockNumWait(key string, num int, expireIn time.Duration) {
	key = lock.FormatKey(key)
	for !lock.LockNum(key, num, expireIn) {
		time.Sleep(time.Millisecond * 100)
	}
}
