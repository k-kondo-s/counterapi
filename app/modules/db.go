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
	Exists(key string) (int64, error)
}

type RedisClient struct {
	client                             *redis.Client
	context                            context.Context
	redisConnectionRetryIntervalSecond int
	redisConnectionRetryNum            int
}

func NewRedisClient(address string, db int) (*RedisClient, error) {
	r := new(RedisClient)
	r.client = redis.NewClient(&redis.Options{
		Addr: address,
		DB: db,
	})
	r.context = context.Background()
	// TODO(kenji-kondo) These params should be set by user with, for instance, environment variables.
	r.redisConnectionRetryIntervalSecond = 5
	r.redisConnectionRetryNum = 6
	// Check connectivity to redis before returning
	_, err := r.client.Ping(r.context).Result()
	if err != nil {
		// If it failed to connect redis, further try to do for several times.
		for i := 0; i < r.redisConnectionRetryNum - 1; i++ {
			time.Sleep(time.Duration(r.redisConnectionRetryIntervalSecond) * time.Second)
			_, errFinal := r.client.Ping(r.context).Result()
			if errFinal == nil {
				return r, nil
			}
		}
		// Give up
		return nil, err

	}
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

func (r *RedisClient) Exists(key string) (int64, error) {
	// "1" means the key exists in Redis, otherwise doesn't exist.
	return r.client.Exists(r.context, key).Result()
}