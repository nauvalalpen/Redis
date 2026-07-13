package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"minikatalog/internal/domain"
	"minikatalog/internal/repository"
)

const (
	CacheKeyAllProducts = "cache:products:all"
	CacheTTL            = 5 * time.Minute // TTL cache produk
)

// ProductUsecase menangani logika bisnis produk + caching Redis
type ProductUsecase struct {
	productRepo *repository.ProductRepository
	redisRepo   *repository.RedisRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository, redisRepo *repository.RedisRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
		redisRepo:   redisRepo,
	}
}

// Create membuat produk baru dan membersihkan cache
func (uc *ProductUsecase) Create(req domain.CreateProductRequest, userID string, username string) (*domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("userID tidak valid: %v", err)
	}

	product := domain.Product{
		UserID:      objID,
		CreatedBy:   username,
		Name:        req.Name,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
		Stock:       req.Stock,
	}

	result, err := uc.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	// Invalidate cache saat ada produk baru (cache-aside pattern)
	if cacheErr := uc.redisRepo.DeleteCache(CacheKeyAllProducts); cacheErr != nil {
		log.Printf("⚠️  Gagal hapus cache setelah create: %v", cacheErr)
	}

	return result, nil
}

// GetAll mengambil semua produk dengan CACHE-ASIDE PATTERN
// 1. Cek Redis cache dulu
// 2. Jika cache HIT → return langsung dari Redis (cepat)
// 3. Jika cache MISS → ambil dari MongoDB, simpan ke Redis, return
func (uc *ProductUsecase) GetAll() ([]domain.Product, string, error) {
	// ===== CACHE-ASIDE PATTERN (KONSEP REDIS #3) =====

	// Step 1: Cek cache
	cached, err := uc.redisRepo.GetCache(CacheKeyAllProducts)
	if err == nil && cached != "" {
		// CACHE HIT: data ada di Redis
		var products []domain.Product
		if jsonErr := json.Unmarshal([]byte(cached), &products); jsonErr == nil {
			log.Println("✅ [CACHE HIT] Data produk diambil dari Redis")
			return products, "CACHE_HIT", nil
		}
	}
	if err != nil && err != redis.Nil {
		log.Printf("⚠️  Error saat cek cache: %v", err)
	}

	// CACHE MISS: ambil dari MongoDB
	log.Println("🔄 [CACHE MISS] Mengambil data dari MongoDB...")
	products, err := uc.productRepo.GetAll()
	if err != nil {
		return nil, "ERROR", err
	}

	// Simpan ke Redis cache untuk request berikutnya
	jsonBytes, err := json.Marshal(products)
	if err == nil {
		if cacheErr := uc.redisRepo.SetCache(CacheKeyAllProducts, string(jsonBytes), CacheTTL); cacheErr != nil {
			log.Printf("⚠️  Gagal simpan ke cache: %v", cacheErr)
		} else {
			log.Printf("💾 [CACHED] Data produk disimpan ke Redis (TTL: %v)", CacheTTL)
		}
	}

	return products, "CACHE_MISS", nil
}

// GetByID mengambil satu produk berdasarkan ID (tanpa cache, data spesifik)
func (uc *ProductUsecase) GetByID(idStr string) (*domain.Product, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, fmt.Errorf("format ID tidak valid")
	}
	return uc.productRepo.GetByID(id)
}

// Update memperbarui produk dan membersihkan cache
func (uc *ProductUsecase) Update(idStr string, req domain.CreateProductRequest) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return fmt.Errorf("format ID tidak valid")
	}

	updates := domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
		Stock:       req.Stock,
	}

	if err := uc.productRepo.Update(id, updates); err != nil {
		return err
	}

	// Invalidate cache setelah update
	uc.redisRepo.DeleteCache(CacheKeyAllProducts)
	return nil
}

// Delete menghapus produk dan membersihkan cache
func (uc *ProductUsecase) Delete(idStr string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return fmt.Errorf("format ID tidak valid")
	}

	if err := uc.productRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache setelah delete
	uc.redisRepo.DeleteCache(CacheKeyAllProducts)
	return nil
}
