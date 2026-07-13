package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"minikatalog/internal/config"
	"minikatalog/internal/handler"
	"minikatalog/internal/middleware"
	"minikatalog/internal/repository"
	"minikatalog/internal/usecase"
)

func main() {
	// =====================================================================
	// INISIALISASI KONEKSI
	// =====================================================================
	config.InitMongo()
	config.InitRedis()
	defer config.CloseConnections()

	// =====================================================================
	// DEPENDENCY INJECTION
	// =====================================================================
	// Repositories
	userRepo := repository.NewUserRepository()
	productRepo := repository.NewProductRepository()
	redisRepo := repository.NewRedisRepository()

	// Usecases
	authUC := usecase.NewAuthUsecase(userRepo, redisRepo)
	productUC := usecase.NewProductUsecase(productRepo, redisRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC)
	productHandler := handler.NewProductHandler(productUC, userRepo)

	// =====================================================================
	// ROUTER & ROUTES
	// =====================================================================
	r := mux.NewRouter()

	// Middleware CORS dan JSON untuk semua route
	r.Use(corsMiddleware)

	// Serve frontend (folder static)
	r.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir("./frontend/"))))
	r.Handle("/", http.RedirectHandler("/app/", http.StatusMovedPermanently))

	// ===== AUTH ROUTES (Public) =====
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST", "OPTIONS")

	// ===== PRODUCT ROUTES =====
	// GET /api/products - Public (semua bisa lihat produk, termasuk efek caching)
	api.HandleFunc("/products", productHandler.GetAll).Methods("GET", "OPTIONS")
	api.HandleFunc("/products/{id}", productHandler.GetByID).Methods("GET", "OPTIONS")

	// Protected routes - harus login
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(authUC))
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET", "OPTIONS")
	protected.HandleFunc("/products", productHandler.Create).Methods("POST", "OPTIONS")
	protected.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE", "OPTIONS")

	// =====================================================================
	// START SERVER
	// =====================================================================
	port := ":8080"
	fmt.Println("\n" + repeatStr("=", 60))
	fmt.Println("  🚀 MiniKatalog Server")
	fmt.Println(repeatStr("=", 60))
	fmt.Printf("  Backend API : http://localhost%s/api\n", port)
	fmt.Printf("  Frontend    : http://localhost%s/app/\n", port)
	fmt.Println(repeatStr("=", 60))
	fmt.Println("\n  Konsep Redis yang aktif:")
	fmt.Println("  ✅ Session Management (Login/Logout)")
	fmt.Println("  ✅ Rate Limiting (POST /api/auth/login)")
	fmt.Println("  ✅ Cache-Aside Pattern (GET /api/products)")
	fmt.Println(repeatStr("=", 60) + "\n")

	log.Fatal(http.ListenAndServe(port, r))
}

// corsMiddleware mengizinkan request dari browser (frontend development)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "X-Cache, X-Response-Time")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}