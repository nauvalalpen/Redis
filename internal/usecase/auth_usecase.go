package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"minikatalog/internal/domain"
	"minikatalog/internal/repository"
)

const (
	MaxLoginAttempts = 5               // Maksimal percobaan login
	SessionTTL       = 2 * time.Hour  // Durasi session
	TokenLength      = 32             // Panjang token (dalam byte hex)
)

// AuthUsecase menangani logika bisnis autentikasi
type AuthUsecase struct {
	userRepo  *repository.UserRepository
	redisRepo *repository.RedisRepository
}

func NewAuthUsecase(userRepo *repository.UserRepository, redisRepo *repository.RedisRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		redisRepo: redisRepo,
	}
}

// Register membuat akun user baru
func (uc *AuthUsecase) Register(req domain.RegisterRequest) error {
	if req.Username == "" || req.Password == "" {
		return errors.New("username dan password tidak boleh kosong")
	}
	if len(req.Password) < 6 {
		return errors.New("password minimal 6 karakter")
	}

	// Hash password dengan bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal memproses password: %v", err)
	}

	user := domain.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		PasswordHash: string(hash),
	}

	return uc.userRepo.Create(user)
}

// Login memverifikasi kredensial dan membuat session di Redis
// Menerapkan Rate Limiting: max 5x gagal per menit per username
func (uc *AuthUsecase) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username dan password tidak boleh kosong")
	}

	// ===== RATE LIMITING (KONSEP REDIS #1) =====
	attempts, err := uc.redisRepo.GetLoginAttempts(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error cek rate limit: %v", err)
	}
	if attempts >= MaxLoginAttempts {
		ttl, _ := uc.redisRepo.GetLoginAttemptTTL(req.Username)
		return nil, fmt.Errorf("terlalu banyak percobaan login. Coba lagi dalam %d detik", ttl)
	}

	// Cari user di MongoDB
	user, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		// Tambah counter meskipun username tidak ditemukan (mencegah enumeration)
		uc.redisRepo.IncrLoginAttempt(req.Username)
		return nil, errors.New("username atau password salah")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Tambah counter percobaan gagal
		newAttempts, _ := uc.redisRepo.IncrLoginAttempt(req.Username)
		remaining := MaxLoginAttempts - newAttempts
		if remaining <= 0 {
			return nil, fmt.Errorf("password salah. Akun diblokir sementara (terlalu banyak percobaan)")
		}
		return nil, fmt.Errorf("password salah. Sisa percobaan: %d", remaining)
	}

	// ===== SESSION MANAGEMENT (KONSEP REDIS #2) =====
	// Generate token unik menggunakan ObjectID (cukup untuk demo)
	token := primitive.NewObjectID().Hex() + primitive.NewObjectID().Hex()

	if err := uc.redisRepo.SetSession(token, user.ID.Hex(), SessionTTL); err != nil {
		return nil, fmt.Errorf("gagal membuat session: %v", err)
	}

	return &domain.LoginResponse{
		Token:    token,
		Username: user.Username,
		Message:  "Login berhasil",
	}, nil
}

// Logout menghapus session dari Redis
func (uc *AuthUsecase) Logout(token string) error {
	return uc.redisRepo.DeleteSession(token)
}

// ValidateSession memvalidasi token dan mengembalikan userID
func (uc *AuthUsecase) ValidateSession(token string) (string, error) {
	userID, err := uc.redisRepo.GetSession(token)
	if err == redis.Nil || userID == "" {
		return "", errors.New("session tidak valid atau sudah kadaluwarsa")
	}
	if err != nil {
		return "", fmt.Errorf("error validasi session: %v", err)
	}
	return userID, nil
}
