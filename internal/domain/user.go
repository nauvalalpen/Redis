package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// User adalah entitas untuk autentikasi
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string             `bson:"username" json:"username"`
	PasswordHash string             `bson:"password_hash" json:"-"` // Tidak dikirim ke frontend
}

// RegisterRequest adalah payload untuk registrasi
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest adalah payload untuk login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse adalah respons sukses login
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
