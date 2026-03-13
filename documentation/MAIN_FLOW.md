# 📊 Alur Bisnis & Flow main.go

## Overview

File `main.go` adalah demonstrasi lengkap dari **5 operasi Redis utama** yang ditampilkan secara berurutan dengan penjelasan alur bisnis yang jelas.

---

## Struktur Output

Output program diorganisir dalam **5 DEMO** yang masing-masing menunjukkan use case nyata:

```
┌─────────────────────────────────────────────────────┐
│           DEMO 1: SETUP & KONEKSI                   │
├─────────────────────────────────────────────────────┤
│           DEMO 2: QUEUE (FIFO)                      │
├─────────────────────────────────────────────────────┤
│           DEMO 3: STRING (SET/GET)                  │
├─────────────────────────────────────────────────────┤
│           DEMO 4: COUNTER (INCR)                    │
├─────────────────────────────────────────────────────┤
│           DEMO 5: TTL/EXPIRY                        │
├─────────────────────────────────────────────────────┤
│           SUMMARY & CONCLUSION                      │
└─────────────────────────────────────────────────────┘
```

---

## DEMO 1: SETUP & KONEKSI KE REDIS

### Alur Bisnis

Tahap pertama program adalah **menginisialisasi koneksi** ke Redis server. Ini adalah langkah kritis karena semua operasi bergantung pada koneksi yang valid.

### Flow Code

```go
rdb := repo.NewRedisClient()                    // Step 1: Buat koneksi
cacheRepo := repo.NewCacheRepository(rdb)       // Step 2: Buat repository
```

### Output yang Diharapkan

```
================================================================================
  DEMO SISTEM REDIS - GOLANG CLEAN ARCHITECTURE
================================================================================

--------------------------------------------------------------------------------
1. SETUP & KONEKSI KE REDIS
--------------------------------------------------------------------------------
Alur Bisnis: Menginisialisasi koneksi ke Redis server

✅ Status: Terhubung ke Redis

✅ Status: CacheRepository siap digunakan
```

### Penjelasan

- **NewRedisClient():** Membuat koneksi TCP ke `localhost:6379`
- **NewCacheRepository(rdb):** Wrap koneksi Redis ke repository layer
- **Error Handling:** Jika koneksi gagal, program akan panic (untuk demo)

---

## DEMO 2: QUEUE/ANTRIAN - LPUSH & RPOP (FIFO)

### Alur Bisnis

**Use Case:** Sistem antrian tugas untuk memproses pekerjaan secara berurutan (First In First Out).

**Skenario Nyata:**

- Tugas masuk dari berbagai sistem
- Diproses satu per satu sesuai urutan kedatangan
- Tugas yang penting (DARURAT) diproses terlebih dahulu (masuk terakhir, dikeluarkan pertama)

### Flow Code

```go
// Step 1: Clear queue lama
cacheRepo.Delete("antrian_tugas")

// Step 2: Push 4 tugas ke queue (LPUSH = masukkan ke head)
cacheRepo.PushQueue("antrian_tugas", "tugas-C", "tugas-B", "tugas-A", "tugas-DARURAT")
//                                     ↑                                    ↑
//                                  last                                 first

// Step 3: Pop tugas dari queue (RPOP = ambil dari tail)
diproses, _ := cacheRepo.PopQueue("antrian_tugas")  // Returns "tugas-DARURAT"
```

### Visualisasi Queue

```
LPUSH dengan urutan: C, B, A, DARURAT

Setelah LPUSH:
┌──────────┬──────────┬──────────┬───────────────┐
│ tugas-C  │ tugas-B  │ tugas-A  │ tugas-DARURAT │
└──────────┴──────────┴──────────┴───────────────┘
   tail                                  head
   (RPOP)                              (LPUSH)

Setelah RPOP 1x:
┌──────────┬──────────┬──────────┐
│ tugas-C  │ tugas-B  │ tugas-A  │
└──────────┴──────────┴──────────┘
```

### Output yang Diharapkan

```
--------------------------------------------------------------------------------
2. QUEUE/ANTRIAN - LPUSH & RPOP (FIFO)
--------------------------------------------------------------------------------
Alur Bisnis: Sistem antrian untuk memproses tugas secara berurutan (FIFO)

[Step 1] Clear queue lama (memastikan fresh start)
Command: DEL antrian_tugas

[Step 2] Masukkan 4 tugas dengan LPUSH (push ke head/kiri)
Command: LPUSH antrian_tugas tugas-DARURAT tugas-A tugas-B tugas-C
Struktur Queue: [tugas-C] ← [tugas-B] ← [tugas-A] ← [tugas-DARURAT]

[Step 3] Ambil tugas pertama dengan RPOP (pop dari tail/kanan)
Command: RPOP antrian_tugas
Tugas Diproses: tugas-DARURAT
Sisa Queue: 3 tugas
Struktur Queue Sekarang: [tugas-C] ← [tugas-B] ← [tugas-A]
```

### Penjelasan

- **LPUSH:** Push dari LEFT (head) - tugas baru masuk di depan
- **RPOP:** Pop dari RIGHT (tail) - ambil dari belakang (FIFO)
- **DEL:** Clear semua data lama untuk fresh start
- **Output:** Menunjukkan struktur queue sebelum dan sesudah operasi

---

## DEMO 3: STRING KEYING - SET & GET

### Alur Bisnis

**Use Case:** Menyimpan dan mengambil data string sederhana untuk metadata/tagging sistem.

**Skenario Nyata:**

- Menyimpan tag untuk kategorisasi
- Menyimpan config value
- Menyimpan metadata aplikasi
- TTL = 0 berarti data **tidak pernah expire** (permanent)

### Flow Code

```go
// Step 1: Set value untuk key
cacheRepo.Set("golang_tag", "redis backend", 0)
//          ↑                ↑                ↑
//         key              value            TTL (0 = permanent)

// Step 2: Get value dari key
val, err := cacheRepo.Get("golang_tag")
//          ↑                ↑
//       result         retrieve
```

### Output yang Diharapkan

```
--------------------------------------------------------------------------------
3. STRING KEYING - SET & GET
--------------------------------------------------------------------------------
Alur Bisnis: Menyimpan dan mengambil data string sederhana (Tagging)

[Step 1] Set value untuk key 'golang_tag'
Command: SET golang_tag 'redis backend'
├─ Waktu Expired: NEVER (TTL = 0 = permanent)
└─ Status: Value tersimpan

[Step 2] Get value dari key 'golang_tag'
Command: GET golang_tag
Result: 'redis backend'
Key exists? true
└─ Status: Data berhasil diambil
```

### Penjelasan

- **SET:** Menyimpan key-value pair dengan TTL
- **GET:** Mengambil value berdasarkan key
- **TTL = 0:** Data bersifat permanent (tidak auto-expire)
- **Error Check:** Verifikasi bahwa key dan value berhasil disimpan

---

## DEMO 4: COUNTER/INCREMENT - INCR (Atomic Operation)

### Alur Bisnis

**Use Case:** Menghitung pengunjung website secara real-time dengan operasi atomic (tidak ada race condition).

**Skenario Nyata:**

- Counter pengunjung/views
- Counter transaksi
- Counter download
- Hits per day/hour
- Operasi INCR adalah **atomic** = thread-safe tanpa lock

### Flow Code

```go
// Step 1: Reset counter
cacheRepo.Delete("pengunjung_counter")

// Step 2: Simulate 5 visitors (INCR)
for i := 1; i <= 5; i++ {
    total, _ := cacheRepo.Increment("pengunjung_counter")
    //                      ↑
    //        Atomic increment, no race condition
}
```

### Visualisasi Counter

```
Initial State:
pengunjung_counter = (not exists) = 0

Pengunjung #1 → INCR → 1
Pengunjung #2 → INCR → 2
Pengunjung #3 → INCR → 3
Pengunjung #4 → INCR → 4
Pengunjung #5 → INCR → 5

Final State:
pengunjung_counter = 5
```

### Output yang Diharapkan

```
--------------------------------------------------------------------------------
4. COUNTER/INCREMENT - INCR (Atomic Operation)
--------------------------------------------------------------------------------
Alur Bisnis: Menghitung pengunjung website secara real-time (atomic)

[Step 1] Reset counter ke 0
Command: DEL pengunjung_counter

[Step 2] Simulasi 5 pengunjung website (INCR)
  Pengunjung #1 masuk → Total counter: 1
  Pengunjung #2 masuk → Total counter: 2
  Pengunjung #3 masuk → Total counter: 3
  Pengunjung #4 masuk → Total counter: 4
  Pengunjung #5 masuk → Total counter: 5

Command: INCR pengunjung_counter (executed 5x)
└─ Status: Counter atomic, tidak ada race condition
```

### Penjelasan

- **INCR:** Operasi increment yang atomic (thread-safe)
- **No Race Condition:** Bahkan di distributed system, INCR aman
- **Use Case:** Perfect untuk counters, analytics, metrics
- **Performance:** Sangat cepat karena single Redis command

---

## DEMO 5: TTL/EXPIRY - AUTO-EXPIRE DATA

### Alur Bisnis

**Use Case:** Menyimpan session token yang otomatis hilang setelah expire time.

**Skenario Nyata:**

- Session token (30 menit)
- OTP code (5 menit)
- Password reset link (1 jam)
- Temporary cache (5 menit)
- Auto-cleanup tanpa perlu background job

### Flow Code

```go
// Step 1: Set dengan TTL 5 detik
cacheRepo.Set("sesi_token", "abc123xyz", 5*time.Second)
//                                          ↑
//                                    TTL (auto-expire)

// Step 2: Get sebelum expire
val, _ := cacheRepo.Get("sesi_token")  // SUCCESS: "abc123xyz"

// Step 3: Tunggu 6 detik
time.Sleep(6 * time.Second)

// Step 4: Get setelah expire
val, err := cacheRepo.Get("sesi_token")  // ERROR: key not found
```

### Timeline Visualisasi

```
Timeline (detik)
0 ──────────────────────────────────────── 5 ───────────────── 11
│                                          │                    │
SET sesi_token                        Auto-Expire             GET
(TTL = 5s)                          (delete by Redis)      (tidak ditemukan)

┌─── Data Valid ───────────┐
│ sesi_token = abc123xyz   │  ← GET di sini: VALID ✅
└──────────────────────────┘
           │
           └─→ Expire setelah 5s
               (Redis otomatis delete)
                     │
                     └─→ GET di sini: NOT FOUND ❌
```

### Output yang Diharapkan

```
--------------------------------------------------------------------------------
5. TTL/EXPIRY - AUTO-EXPIRE DATA
--------------------------------------------------------------------------------
Alur Bisnis: Menyimpan session token yang otomatis expire setelah 5 detik

[Step 1] Set token dengan TTL 5 detik
Command: SET sesi_token 'abc123xyz' EX 5
├─ Token Value: abc123xyz
└─ Expiry: 5 detik

[Step 2] Get token SEBELUM expire (Token masih valid)
Command: GET sesi_token
Result: 'abc123xyz' ✅ (Token VALID)
└─ Status: Session masih aktif

[Step 3] Tunggu 6 detik hingga token expire...
Countdown: 6 5 4 3 2 1
└─ Waktu berlalu: 6 detik

[Step 4] Get token SETELAH expire (Token sudah hilang)
Command: GET sesi_token
Result: (nil) ❌ (Token EXPIRED)
└─ Status: Session berakhir, user harus login ulang
```

### Penjelasan

- **TTL (Time To Live):** Waktu dalam detik sebelum key auto-delete
- **EX 5:** Set expiry 5 detik (format Redis)
- **Auto-Cleanup:** Redis handle delete, tidak perlu app logic
- **Use Case:** Perfect untuk temporary data seperti session, OTP, cache
- **Countdown:** Program sleep 6 detik sambil menunjukkan countdown

---

## SUMMARY & CONCLUSION

### Output yang Diharapkan

```
--------------------------------------------------------------------------------
✅ DEMO SELESAI
--------------------------------------------------------------------------------
Semua operasi Redis berhasil dijalankan!

Operasi yang dipelajari:
  1. Queue (LPUSH/RPOP) - Antrian FIFO
  2. String (SET/GET) - Key-Value storage
  3. Counter (INCR) - Atomic increment
  4. TTL (Expiry) - Auto-expire data

Lihat documentation/docs.md untuk penjelasan lebih detail.
```

---

## Tabel Perbandingan Operasi Redis

| Demo | Operasi | Command    | Use Case                   | TTL          |
| ---- | ------- | ---------- | -------------------------- | ------------ |
| 2    | Queue   | LPUSH/RPOP | Task queue, Job processing | No           |
| 3    | String  | SET/GET    | Config, Metadata, Tags     | Configurable |
| 4    | Counter | INCR       | Views, Clicks, Analytics   | No           |
| 5    | Expiry  | EX (TTL)   | Session, OTP, Temp Cache   | Yes (Auto)   |

---

## Kode Helper Functions

File `main.go` juga memiliki **3 helper functions** untuk membuat output lebih rapi:

### 1. `printHeader(title string)`

```go
func printHeader(title string) {
    fmt.Println("\n" + repeatString("=", 80))
    fmt.Printf("  %s\n", title)
    fmt.Println(repeatString("=", 80) + "\n")
}
```

**Output:**

```
================================================================================
  DEMO SISTEM REDIS - GOLANG CLEAN ARCHITECTURE
================================================================================
```

### 2. `printSection(title string)`

```go
func printSection(title string) {
    fmt.Println(repeatString("-", 80))
    fmt.Printf("%s\n", title)
    fmt.Println(repeatString("-", 80))
}
```

**Output:**

```
--------------------------------------------------------------------------------
1. SETUP & KONEKSI KE REDIS
--------------------------------------------------------------------------------
```

### 3. `repeatString(s string, count int) string`

```go
func repeatString(s string, count int) string {
    result := ""
    for i := 0; i < count; i++ {
        result += s
    }
    return result
}
```

**Fungsi:** Repeat string untuk membuat separator line yang rapi

---

## Tips Menjalankan Program

### Standar Output (default)

```bash
go run main.go
```

### Dengan Timing (lihat berapa lama eksekusi)

```bash
time go run main.go
```

### Redirect Output ke File

```bash
go run main.go > output.log
```

### Run dan Lihat STDERR (debugging)

```bash
go run main.go 2>&1
```

---

## Troubleshooting Output

### Masalah 1: Output tidak rapi / ada karakter aneh

**Penyebab:** Terminal tidak mendukung 80 character width
**Solusi:** Resize terminal window lebih lebar

### Masalah 2: Program hang di DEMO 5 (Timer)

**Penyebab:** Normal - program sleep 6 detik
**Solusi:** Tunggu atau ctrl+C untuk stop (untuk development)

### Masalah 3: Redis connection error

**Penyebab:** Redis server tidak running
**Solusi:** Jalankan Redis (Docker)

```bash
docker run --name redis-container -p 6379:6379 -d redis:latest
```

---

## Summary Alur Bisnis

| #   | Demo    | Input           | Proses           | Output                   | Status    |
| --- | ------- | --------------- | ---------------- | ------------------------ | --------- |
| 1   | Setup   | -               | Connect ke Redis | ✅ Connected             | Ready     |
| 2   | Queue   | 4 tasks         | LPUSH + RPOP     | 1 processed, 3 remaining | Running   |
| 3   | String  | key-value       | SET + GET        | Value retrieved          | Permanent |
| 4   | Counter | 5 visitors      | 5x INCR          | Counter = 5              | Atomic    |
| 5   | TTL     | token + 5s wait | SET EX + GET     | Token expired            | Deleted   |

---

**Created:** March 2026  
**Version:** 1.0  
**Status:** Complete ✅

_Dokumen ini menjelaskan alur bisnis lengkap dari main.go dengan contoh output, visualisasi, dan penjelasan setiap tahap._
