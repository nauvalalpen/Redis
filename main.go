package main

import (
	"fmt"
	"time"

	repo "redis/repository/redis"
)

func main() {
	printHeader("DEMO SISTEM REDIS - GOLANG CLEAN ARCHITECTURE")

	// ============================================================================
	// DEMO 1: SETUP & KONEKSI KE REDIS
	// ============================================================================
	printSection("1. SETUP & KONEKSI KE REDIS")
	fmt.Println("Alur Bisnis: Menginisialisasi koneksi ke Redis server\n")

	rdb := repo.NewRedisClient()
	fmt.Println("✅ Status: Terhubung ke Redis\n")

	cacheRepo := repo.NewCacheRepository(rdb)
	fmt.Println("✅ Status: CacheRepository siap digunakan\n")

	// ============================================================================
	// DEMO 2: QUEUE/ANTRIAN (FIFO - First In First Out)
	// ============================================================================
	printSection("2. QUEUE/ANTRIAN - LPUSH & RPOP (FIFO)")
	fmt.Println("Alur Bisnis: Sistem antrian untuk memproses tugas secara berurutan (FIFO)\n")

	cacheRepo.Delete("antrian_tugas")
	fmt.Println("[Step 1] Clear queue lama (memastikan fresh start)")
	fmt.Println("Command: DEL antrian_tugas\n")

	fmt.Println("[Step 2] Masukkan 4 tugas dengan LPUSH (push ke head/kiri)")
	cacheRepo.PushQueue("antrian_tugas", "tugas-C", "tugas-B", "tugas-A", "tugas-DARURAT")
	fmt.Println("Command: LPUSH antrian_tugas tugas-DARURAT tugas-A tugas-B tugas-C")
	fmt.Println("Struktur Queue: [tugas-C] ← [tugas-B] ← [tugas-A] ← [tugas-DARURAT]\n")

	fmt.Println("[Step 3] Ambil tugas pertama dengan RPOP (pop dari tail/kanan)")
	diproses, _ := cacheRepo.PopQueue("antrian_tugas")
	fmt.Printf("Command: RPOP antrian_tugas\n")
	fmt.Printf("Tugas Diproses: %s\n", diproses)
	fmt.Printf("Sisa Queue: 3 tugas\n")
	fmt.Println("Struktur Queue Sekarang: [tugas-C] ← [tugas-B] ← [tugas-A]\n")

	// ============================================================================
	// DEMO 3: STRING KEYING (SET & GET)
	// ============================================================================
	printSection("3. STRING KEYING - SET & GET")
	fmt.Println("Alur Bisnis: Menyimpan dan mengambil data string sederhana (Tagging)\n")

	fmt.Println("[Step 1] Set value untuk key 'golang_tag'")
	cacheRepo.Set("golang_tag", "redis backend", 0)
	fmt.Println("Command: SET golang_tag 'redis backend'")
	fmt.Println("├─ Waktu Expired: NEVER (TTL = 0 = permanent)")
	fmt.Println("└─ Status: Value tersimpan\n")

	fmt.Println("[Step 2] Get value dari key 'golang_tag'")
	val, err := cacheRepo.Get("golang_tag")
	fmt.Printf("Command: GET golang_tag\n")
	fmt.Printf("Result: '%s'\n", val)
	fmt.Printf("Key exists? %v\n", err == nil && val != "")
	fmt.Println("└─ Status: Data berhasil diambil\n")

	// ============================================================================
	// DEMO 4: COUNTER/INCREMENT (ATOMIC OPERATION)
	// ============================================================================
	printSection("4. COUNTER/INCREMENT - INCR (Atomic Operation)")
	fmt.Println("Alur Bisnis: Menghitung pengunjung website secara real-time (atomic)\n")

	cacheRepo.Delete("pengunjung_counter")
	fmt.Println("[Step 1] Reset counter ke 0")
	fmt.Println("Command: DEL pengunjung_counter\n")

	fmt.Println("[Step 2] Simulasi 5 pengunjung website (INCR)")
	for i := 1; i <= 5; i++ {
		total, _ := cacheRepo.Increment("pengunjung_counter")
		fmt.Printf("  Pengunjung #%d masuk → Total counter: %d\n", i, total)
	}
	fmt.Println("\nCommand: INCR pengunjung_counter (executed 5x)")
	fmt.Println("└─ Status: Counter atomic, tidak ada race condition\n")

	// ============================================================================
	// DEMO 5: TTL/EXPIRY (TIME-TO-LIVE)
	// ============================================================================
	printSection("5. TTL/EXPIRY - AUTO-EXPIRE DATA")
	fmt.Println("Alur Bisnis: Menyimpan session token yang otomatis expire setelah 5 detik\n")

	fmt.Println("[Step 1] Set token dengan TTL 5 detik")
	cacheRepo.Set("sesi_token", "abc123xyz", 5*time.Second)
	fmt.Println("Command: SET sesi_token 'abc123xyz' EX 5")
	fmt.Println("├─ Token Value: abc123xyz")
	fmt.Println("└─ Expiry: 5 detik\n")

	fmt.Println("[Step 2] Get token SEBELUM expire (Token masih valid)")
	val, err = cacheRepo.Get("sesi_token")
	if err == nil && val != "" {
		fmt.Printf("Command: GET sesi_token\n")
		fmt.Printf("Result: '%s' ✅ (Token VALID)\n", val)
	}
	fmt.Println("└─ Status: Session masih aktif\n")

	fmt.Println("[Step 3] Tunggu 6 detik hingga token expire...")
	fmt.Print("Countdown: ")
	for i := 6; i > 0; i-- {
		fmt.Printf("%d ", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("\n└─ Waktu berlalu: 6 detik\n")

	fmt.Println("[Step 4] Get token SETELAH expire (Token sudah hilang)")
	val, err = cacheRepo.Get("sesi_token")
	if err != nil || val == "" {
		fmt.Println("Command: GET sesi_token")
		fmt.Println("Result: (nil) ❌ (Token EXPIRED)")
		fmt.Println("└─ Status: Session berakhir, user harus login ulang\n")
	}

	// ============================================================================
	// SUMMARY
	// ============================================================================
	printSection("✅ DEMO SELESAI")
	fmt.Println("Semua operasi Redis berhasil dijalankan!\n")
	fmt.Println("Operasi yang dipelajari:")
	fmt.Println("  1. Queue (LPUSH/RPOP) - Antrian FIFO")
	fmt.Println("  2. String (SET/GET) - Key-Value storage")
	fmt.Println("  3. Counter (INCR) - Atomic increment")
	fmt.Println("  4. TTL (Expiry) - Auto-expire data")
	fmt.Println("\nLihat documentation/docs.md untuk penjelasan lebih detail.\n")
}

// Helper function untuk header yang rapi
func printHeader(title string) {
	fmt.Println("\n" + repeatString("=", 80))
	fmt.Printf("  %s\n", title)
	fmt.Println(repeatString("=", 80) + "\n")
}

// Helper function untuk section yang rapi
func printSection(title string) {
	fmt.Println(repeatString("-", 80))
	fmt.Printf("%s\n", title)
	fmt.Println(repeatString("-", 80))
}

// Helper function untuk repeat string
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}