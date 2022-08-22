package go_redis_v8_impl

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClientImplGoRedis struct {
	client *redis.Client
}

func (r *RedisClientImplGoRedis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *RedisClientImplGoRedis) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}
