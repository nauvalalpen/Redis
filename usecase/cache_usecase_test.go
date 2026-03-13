package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock
// mockCacheRepo adalah implementasi manual dari domain.CacheRepository
type mockCacheRepo struct {
	setFn       func(key string, value string, ttl time.Duration) error
	getFn       func(key string) (string, error)
	deleteFn    func(key string) error
	incrementFn func(key string) (int64, error)
	pushFn      func(q string, v ...string) error
	popFn       func(q string) (string, error)
}

func (m *mockCacheRepo) Set(k string, v string, ttl time.Duration) error { return m.setFn(k, v, ttl) }
func (m *mockCacheRepo) Get(k string) (string, error) { return m.getFn(k) }
func (m *mockCacheRepo) Delete(k string) error { return m.deleteFn(k) }
func (m *mockCacheRepo) Increment(k string) (int64, error) { return m.incrementFn(k) }
func (m *mockCacheRepo) PushQueue(q string, v ...string) error { return m.pushFn(q, v...) }
func (m *mockCacheRepo) PopQueue(q string) (string, error) { return m.popFn(q) }

// 1. Skenario Sukses
func TestCacheSet_Success(t *testing.T) {
	mockRepo := &mockCacheRepo{
		setFn: func(key string, value string, ttl time.Duration) error {
			return nil
		},
	}
	
	err := mockRepo.Set("test_key", "test_val", 0)
	assert.NoError(t, err) // Berharap tidak ada error
}

// 2. Skenario Error (Di sini package "errors" digunakan!)
func TestCacheSet_Error(t *testing.T) {
	mockRepo := &mockCacheRepo{
		setFn: func(key string, value string, ttl time.Duration) error {
			return errors.New("redis connection timeout") // Menggunakan package errors
		},
	}
	
	err := mockRepo.Set("test_key", "test_val", 0)
	assert.Error(t, err) // Berharap harusnya menghasilkan error
}