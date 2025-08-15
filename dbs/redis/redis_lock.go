package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"live-user/consts"
	"live-user/dbs"
)

type RedisLock struct {
	client *redis.Client
	key    string
	value  string
	expire time.Duration
}

func NewRedisLock(client *redis.Client, key string, expire time.Duration) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		value:  fmt.Sprintf("%d", time.Now().UnixNano()),
		expire: expire,
	}
}

// Acquire 获取锁
func (l *RedisLock) Acquire(ctx context.Context) (bool, error) {
	return l.client.SetNX(ctx, l.key, l.value, l.expire).Result()
}

// Release 释放锁
func (l *RedisLock) Release(ctx context.Context) error {
	// 使用Lua脚本保证原子性
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	_, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	return err
}

func NewEmailSendLock(email string) *RedisLock {
	ac := NewRedisLock(dbs.RedisEngine, consts.BuildEmailKey(email), consts.LIMIT_SEND_EMAIL)
	return ac
}
