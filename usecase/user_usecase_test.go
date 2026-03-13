package usecase

import (
	"errors"
	"testing"
	"redis/domain"

	"github.com/stretchr/testify/assert"
)

// Mock
// mockUserRepo adalah implementasi manual dari domain.UserRepository
// yang memungkinkan tiap test mengontrol perilaku setiap method.
type mockUserRepo struct {
	saveFn     func(user domain.User) error
	findByIDFn func(id string) (*domain.User, error)
	deleteFn   func(id string) error
}

func (m *mockUserRepo) Save(user domain.User) error {
	return m.saveFn(user)
}

func (m *mockUserRepo) FindByID(id string) (*domain.User, error) {
	return m.findByIDFn(id)
}

func (m *mockUserRepo) Delete(id string) error {
	return m.deleteFn(id)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := &mockUserRepo{
		saveFn: func(user domain.User) error {
			return nil // Simulasi sukses
		},
	}
	uc := NewUserUsecase(mockRepo)

	err := uc.CreateUser(domain.User{ID: "1", Name: "Test"})
	assert.NoError(t, err)
}

func TestCreateUser_Error(t *testing.T) {
	mockRepo := &mockUserRepo{
		saveFn: func(user domain.User) error {
			return errors.New("redis error") // Simulasi error
		},
	}
	uc := NewUserUsecase(mockRepo)

	err := uc.CreateUser(domain.User{ID: "1", Name: "Test"})
	assert.Error(t, err)
}