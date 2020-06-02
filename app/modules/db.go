package modules

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Dao interface {
	Set(key string, value string, expirationSecond int64) error
	Get(key string) (string, error)
	GetAllKeys() ([]string, error)
	Del(key string) error
}

type RedisClient struct {
	client  *redis.Client
	context context.Context
}

func NewRedisClient(address string, db int) (*RedisClient, error) {
	r := new(RedisClient)
	r.client = redis.NewClient(&redis.Options{
		Addr: address,
		DB: db,
	})
	r.context = context.Background()
	// TODO(kenji-kondo): check connectivity to redis before returning
	return r, nil
}

func (r *RedisClient) Set(key string, value string, expirationSecond int64) error {
	return r.client.Set(r.context, key, value, time.Duration(expirationSecond) * time.Second).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.context, key).Result()
}

func (r *RedisClient) GetAllKeys() ([]string, error) {
	return r.client.Keys(r.context, "*").Result()
}

func (r *RedisClient) Del(key string) error {
	return r.client.Del(r.context, key).Err()
}