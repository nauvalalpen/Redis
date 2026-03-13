package domain

import "time"

// CacheRepository mendefinisikan kontrak untuk operasi cache umum (String, Counter, TTL)
type CacheRepository interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Increment(key string) (int64, error)
	PushQueue(queueName string, values ...string) error
	PopQueue(queueName string) (string, error)
}