package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"minikatalog/internal/domain"
	"minikatalog/internal/middleware"
	"minikatalog/internal/repository"
	"minikatalog/internal/usecase"
)

// ProductHandler menangani endpoint produk
type ProductHandler struct {
	productUC *usecase.ProductUsecase
	userRepo  *repository.UserRepository
}

func NewProductHandler(productUC *usecase.ProductUsecase, userRepo *repository.UserRepository) *ProductHandler {
	return &ProductHandler{
		productUC: productUC,
		userRepo:  userRepo,
	}
}

// GetAll godoc
// GET /api/products
// Mengembalikan semua produk. Menggunakan Cache-Aside Pattern.
// Response Header X-Cache: HIT / MISS untuk membuktikan cache bekerja
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	start := time.Now()
	products, cacheStatus, err := h.productUC.GetAll()
	elapsed := time.Since(start)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Header tambahan untuk bukti caching
	w.Header().Set("X-Cache", cacheStatus)
	w.Header().Set("X-Response-Time", elapsed.String())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":          products,
		"cache_status":  cacheStatus,
		"response_time": elapsed.String(),
		"count":         len(products),
	})
}

// Create godoc
// POST /api/products (Protected - memerlukan Auth)
// Body: {"name": "...", "price": 0, "category": "...", "description": "...", "stock": 0}
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := middleware.GetUserIDFromContext(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tidak terautentikasi"})
		return
	}

	// Ambil username dari MongoDB untuk disimpan sebagai created_by
	user, err := h.userRepo.FindByID_Str(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "User tidak ditemukan"})
		return
	}

	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format request tidak valid"})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Nama produk tidak boleh kosong"})
		return
	}

	product, err := h.productUC.Create(req, userID, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetByID godoc
// GET /api/products/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	product, err := h.productUC.GetByID(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Produk tidak ditemukan"})
		return
	}

	json.NewEncoder(w).Encode(product)
}

// Update godoc
// PUT /api/products/{id} (Protected)
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format request tidak valid"})
		return
	}

	if err := h.productUC.Update(params["id"], req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil diperbarui"})
}

// Delete godoc
// DELETE /api/products/{id} (Protected)
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if err := h.productUC.Delete(params["id"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berhasil dihapus"})
}
