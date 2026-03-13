package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"redis/domain"
)

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) domain.CacheRepository {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) Set(key string, value string, ttl time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *cacheRepository) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *cacheRepository) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *cacheRepository) Increment(key string) (int64, error) {
	ctx := context.Background()
	return r.client.Incr(ctx, key).Result()
}

func (r *cacheRepository) PushQueue(queueName string, values ...string) error {
	ctx := context.Background()
	// Gunakan interface{} untuk variadic arguments
	var interfaceVals[]interface{}
	for _, v := range values {
		interfaceVals = append(interfaceVals, v)
	}
	return r.client.LPush(ctx, queueName, interfaceVals...).Err()
}

func (r *cacheRepository) PopQueue(queueName string) (string, error) {
	ctx := context.Background()
	return r.client.RPop(ctx, queueName).Result()
}