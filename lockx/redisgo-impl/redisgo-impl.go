package redisgo_impl

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisClientImplRedisGo struct {
	conn redis.Conn
}

func NewRedisClientImplRedisGo(conn redis.Conn) *RedisClientImplRedisGo {
	return &RedisClientImplRedisGo{conn: conn}
}

func (r *RedisClientImplRedisGo) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return redis.Bool(r.conn.Do("setnx", key, value, expiration))
}

func (r *RedisClientImplRedisGo) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	var p = []interface{}{script, len(keys)}
	for _, key := range keys {
		p = append(p, key)
	}

	for _, arg := range args {
		p = append(p, arg)
	}

	return r.conn.Do("eval", p...)
}
