package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"minikatalog/internal/config"
	"minikatalog/internal/domain"
)

// UserRepository menangani operasi database untuk entitas User
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create menyimpan user baru ke MongoDB
func (r *UserRepository) Create(user domain.User) error {
	col := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cek apakah username sudah ada
	var existing domain.User
	err := col.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return errors.New("username sudah digunakan")
	}

	_, err = col.InsertOne(ctx, user)
	return err
}

// FindByUsername mencari user berdasarkan username
func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	col := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return &user, nil
}

// FindByID mencari user berdasarkan ObjectID
func (r *UserRepository) FindByID(id primitive.ObjectID) (*domain.User, error) {
	col := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return &user, nil
}

// FindByID_Str mencari user berdasarkan ID dalam bentuk string (hex ObjectID)
// Digunakan untuk lookup dari session Redis yang menyimpan userID sebagai string
func (r *UserRepository) FindByID_Str(idStr string) (*domain.User, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, errors.New("format ID tidak valid")
	}
	return r.FindByID(id)
}
