package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product adalah entitas produk yang disimpan di MongoDB
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedBy   string             `bson:"created_by" json:"created_by"` // Username pembuatnya
	Name        string             `bson:"name" json:"name"`
	Price       int                `bson:"price" json:"price"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Stock       int                `bson:"stock" json:"stock"`
}

// CreateProductRequest adalah payload untuk membuat produk baru
type CreateProductRequest struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}
