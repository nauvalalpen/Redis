package handler

import (
	"encoding/json"
	"net/http"

	"minikatalog/internal/domain"
	"minikatalog/internal/usecase"
)

// AuthHandler menangani endpoint autentikasi
type AuthHandler struct {
	authUC *usecase.AuthUsecase
}

func NewAuthHandler(authUC *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Register godoc
// POST /api/auth/register
// Body: {"username": "...", "password": "..."}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format request tidak valid"})
		return
	}

	if err := h.authUC.Register(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registrasi berhasil! Silakan login."})
}

// Login godoc
// POST /api/auth/login
// Body: {"username": "...", "password": "..."}
// Rate limited: max 5x gagal per menit
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format request tidak valid"})
		return
	}

	resp, err := h.authUC.Login(req)
	if err != nil {
		// Status 429 untuk rate limit, 401 untuk kredensial salah
		status := http.StatusUnauthorized
		errMsg := err.Error()
		if len(errMsg) > 15 && errMsg[:15] == "terlalu banyak" {
			status = http.StatusTooManyRequests
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// Logout godoc
// POST /api/auth/logout
// Header: Authorization: Bearer <token>
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 {
		token := authHeader[7:] // Hapus "Bearer "
		h.authUC.Logout(token)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Logout berhasil"})
}

// Me godoc
// GET /api/auth/me
// Header: Authorization: Bearer <token>
// Mengembalikan info user yang sedang login (dipakai untuk validasi sesi di frontend)
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// userID sudah divalidasi oleh middleware dan ada di context
	// Handler ini hanya perlu konfirmasi bahwa session valid
	json.NewEncoder(w).Encode(map[string]string{"message": "Session valid", "status": "authenticated"})
}
