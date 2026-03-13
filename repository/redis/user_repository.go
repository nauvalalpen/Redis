package redis

import (
	"context"
	"encoding/json"
	"redis/domain"

	"github.com/redis/go-redis/v9"
)

type userRepository struct {
	client *redis.Client
}

func NewUserRepository(client *redis.Client) domain.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Save(user domain.User) error {
	ctx := context.Background()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// Simpan data user sebagai string JSON tanpa kadaluarsa (0)
	return r.client.Set(ctx, "user:"+user.ID, data, 0).Err()
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, "user:"+id).Result()
	if err != nil {
		return nil, err // Return error jika tidak ditemukan
	}

	var user domain.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(id string) error {
	ctx := context.Background()
	return r.client.Del(ctx, "user:"+id).Err()
}