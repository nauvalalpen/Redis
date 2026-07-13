package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"minikatalog/internal/config"
)

// RedisRepository menangani semua operasi Redis
type RedisRepository struct{}

func NewRedisRepository() *RedisRepository {
	return &RedisRepository{}
}

func (r *RedisRepository) client() *redis.Client {
	return config.RedisClient
}

// ===== SESSION MANAGEMENT =====

// SetSession menyimpan token session dengan TTL
func (r *RedisRepository) SetSession(token string, userID string, ttl time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", token)
	return r.client().Set(ctx, key, userID, ttl).Err()
}

// GetSession mengambil userID dari token session
func (r *RedisRepository) GetSession(token string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", token)
	return r.client().Get(ctx, key).Result()
}

// DeleteSession menghapus token session (logout)
func (r *RedisRepository) DeleteSession(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", token)
	return r.client().Del(ctx, key).Err()
}

// ===== RATE LIMITING =====

// IncrLoginAttempt menambah counter percobaan login dan return jumlahnya
// Key otomatis kadaluwarsa setelah TTL (window time)
func (r *RedisRepository) IncrLoginAttempt(identifier string) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("ratelimit:login:%s", identifier)

	pipe := r.client().Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 1*time.Minute) // Window: 1 menit
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// GetLoginAttempts mengambil jumlah percobaan login saat ini
func (r *RedisRepository) GetLoginAttempts(identifier string) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("ratelimit:login:%s", identifier)
	val, err := r.client().Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// GetLoginAttemptTTL mengambil sisa TTL rate limit dalam detik
func (r *RedisRepository) GetLoginAttemptTTL(identifier string) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("ratelimit:login:%s", identifier)
	ttl, err := r.client().TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int64(ttl.Seconds()), nil
}

// ===== CACHING =====

// SetCache menyimpan data ke Redis cache dengan TTL
func (r *RedisRepository) SetCache(key string, value string, ttl time.Duration) error {
	ctx := context.Background()
	return r.client().Set(ctx, key, value, ttl).Err()
}

// GetCache mengambil data dari Redis cache
// Return (value, nil) jika ada, atau ("", redis.Nil) jika cache miss
func (r *RedisRepository) GetCache(key string) (string, error) {
	ctx := context.Background()
	return r.client().Get(ctx, key).Result()
}

// DeleteCache menghapus cache berdasarkan key (dipakai saat data berubah)
func (r *RedisRepository) DeleteCache(key string) error {
	ctx := context.Background()
	return r.client().Del(ctx, key).Err()
}
