package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient membuat koneksi baru ke Redis
func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Host dan port Redis standar
		Password: "",               // Tidak ada password default
		DB:       0,                // Gunakan DB 0 default
	})

	// Test koneksi
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Gagal terhubung ke Redis: %v", err))
	}
	
	return rdb
}