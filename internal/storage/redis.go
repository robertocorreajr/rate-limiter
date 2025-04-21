package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(addr string) *RedisStorage {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStorage{client: client, ctx: ctx}
}

func (r *RedisStorage) Increment(key string) (int, error) {
	res, err := r.client.Incr(r.ctx, key).Result()
	return int(res), err
}

func (r *RedisStorage) TTL(key string) (int, error) {
	d, err := r.client.TTL(r.ctx, key).Result()
	return int(d.Seconds()), err
}

func (r *RedisStorage) SetTTL(key string, seconds int) error {
	return r.client.Expire(r.ctx, key, time.Duration(seconds)*time.Second).Err()
}

func (r *RedisStorage) GetLimit(key string) (int, error) {
	res, err := r.client.Get(r.ctx, "limit:"+key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (r *RedisStorage) SetLimit(key string, value int) error {
	return r.client.Set(r.ctx, "limit:"+key, value, 0).Err()
}

func (r *RedisStorage) Block(key string, seconds int) error {
	return r.client.Set(r.ctx, "block:"+key, "1", time.Duration(seconds)*time.Second).Err()
}

func (r *RedisStorage) IsBlocked(key string) (bool, error) {
	exists, err := r.client.Exists(r.ctx, "block:"+key).Result()
	return exists == 1, err
}
