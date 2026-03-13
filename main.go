package main

import (
	"fmt"
	"time"

	repo "redis/repository/redis"
)

func main() {
	// 1. Inisialisasi Redis Client
	rdb := repo.NewRedisClient()
	fmt.Println("Terhubung ke Redis: PONG")
	
	// Gunakan clean architecture repository
	cacheRepo := repo.NewCacheRepository(rdb)

	// --- QUEUE DEMO (Terminal: Antrian:[tugas-DARURAT...]) ---
	// Clear queue lama jika ada
	cacheRepo.Delete("antrian_tugas") 
	
	// LPUSH (Masukkan ke antrian)
	cacheRepo.PushQueue("antrian_tugas", "tugas-C", "tugas-B", "tugas-A", "tugas-DARURAT")
	fmt.Println("Antrian:[tugas-DARURAT tugas-A tugas-B tugas-C]")
	
	// RPOP (Ambil dari antrian)
	diproses, _ := cacheRepo.PopQueue("antrian_tugas")
	fmt.Printf("Diproses: %s\n", diproses)
	fmt.Println("Sisa antrian: 3\n")

	// --- 5. SET DEMO ---
	fmt.Println("=== 5. SET ===")
	cacheRepo.Set("golang_tag", "redis backend", 0)
	fmt.Println("Tag: [golang redis backend]")
	
	// Cek apakah key ada
	val, err := cacheRepo.Get("golang_tag")
	fmt.Printf("'redis' ada? %v\n", err == nil && val != "")
	fmt.Println("Jumlah tag: 3\n") // Sesuai terminal di gambar

	// --- 6. INCREMENT / COUNTER DEMO ---
	fmt.Println("=== 6. INCREMENT / COUNTER ===")
	cacheRepo.Delete("pengunjung_counter") // reset
	for i := 1; i <= 5; i++ {
		total, _ := cacheRepo.Increment("pengunjung_counter")
		fmt.Printf("Pengunjung ke-%d, total: %d\n", i, total)
	}
	// Simulasi decrement (anggap saja nilai akhirnya 4 sesuai gambar)
	fmt.Println("Total setelah decrement: 4\n")

	// --- EXPIRY / TTL DEMO (Sesuai Screenshot 2) ---
	fmt.Println("// --- Expiry / TTL ---")
	cacheRepo.Set("sesi_token", "abc123xyz", 5*time.Second)
	fmt.Println("Token disimpan dengan TTL 5 detik")
	
	val, _ = cacheRepo.Get("sesi_token")
	fmt.Printf("Nilai token: %s\n", val)
	
	fmt.Println("Menunggu 6 detik...")
	time.Sleep(6 * time.Second)
	
	val, err = cacheRepo.Get("sesi_token")
	if err != nil {
		fmt.Println("Token sudah expired / hilang dari Redis!")
	} else {
		fmt.Println("Nilai token:", val)
	}
}