package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDB     *mongo.Client
	RedisClient *redis.Client
)

// InitMongo menginisialisasi koneksi ke MongoDB lokal
func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Gagal koneksi ke MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB tidak dapat diakses: %v", err)
	}

	MongoDB = client
	fmt.Println("✅ Koneksi MongoDB berhasil!")
}

// InitRedis menginisialisasi koneksi ke Redis lokal
func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Gagal koneksi ke Redis: %v", err)
	}

	RedisClient = rdb
	fmt.Println("✅ Koneksi Redis berhasil!")
}

// GetCollection adalah helper untuk mengambil collection MongoDB dengan mudah
func GetCollection(name string) *mongo.Collection {
	return MongoDB.Database("minikatalog").Collection(name)
}

// CloseConnections menutup semua koneksi saat aplikasi berhenti
func CloseConnections() {
	if MongoDB != nil {
		MongoDB.Disconnect(context.Background())
	}
	if RedisClient != nil {
		RedisClient.Close()
	}
}
