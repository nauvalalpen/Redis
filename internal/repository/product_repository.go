package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"minikatalog/internal/config"
	"minikatalog/internal/domain"
)

// ProductRepository menangani operasi database untuk entitas Product
type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// Create menyimpan produk baru ke MongoDB
func (r *ProductRepository) Create(product domain.Product) (*domain.Product, error) {
	col := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product.ID = primitive.NewObjectID()
	_, err := col.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll mengambil semua produk dari MongoDB, diurutkan dari terbaru
func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	col := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Urutkan dari terbaru (_id descending)
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	// Kembalikan slice kosong, bukan nil, agar JSON encode jadi [] bukan null
	if products == nil {
		products = []domain.Product{}
	}
	return products, nil
}

// GetByID mengambil satu produk berdasarkan ID
func (r *ProductRepository) GetByID(id primitive.ObjectID) (*domain.Product, error) {
	col := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product domain.Product
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update memperbarui data produk berdasarkan ID
func (r *ProductRepository) Update(id primitive.ObjectID, updates domain.Product) error {
	col := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateDoc := bson.M{
		"$set": bson.M{
			"name":        updates.Name,
			"price":       updates.Price,
			"category":    updates.Category,
			"description": updates.Description,
			"stock":       updates.Stock,
		},
	}

	_, err := col.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	return err
}

// Delete menghapus produk berdasarkan ID
func (r *ProductRepository) Delete(id primitive.ObjectID) error {
	col := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
