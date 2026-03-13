# 📚 Dokumentasi Sistem Redis - Golang Clean Architecture

## Daftar Isi

1. [Overview Sistem](#overview-sistem)
2. [Arsitektur Clean Architecture](#arsitektur-clean-architecture)
3. [Struktur Folder Proyek](#struktur-folder-proyek)
4. [Komponen Utama](#komponen-utama)
5. [Operasi Redis](#operasi-redis)
6. [Setup & Instalasi](#setup--instalasi)
7. [Cara Menggunakan](#cara-menggunakan)
8. [Testing](#testing)
9. [Troubleshooting](#troubleshooting)

---

## Overview Sistem

Sistem ini adalah implementasi **Golang + Redis** menggunakan prinsip **Clean Architecture** untuk mengelola data user dan caching secara efisien. Project ini mendemonstrasikan:

- ✅ Koneksi ke Redis Server
- ✅ Operasi cache dasar (Set, Get, Delete)
- ✅ Counter/Increment dalam Redis
- ✅ Queue/Antrian menggunakan Lists
- ✅ Manajemen user dengan caching
- ✅ Unit testing dengan testify/assert

**Stack Teknologi:**

- **Bahasa:** Go (Golang) 1.18+
- **Database:** Redis (via Docker)
- **Testing:** testify/assert
- **Library:** github.com/redis/go-redis/v9

---

## Arsitektur Clean Architecture

Sistem ini mengikuti prinsip **Clean Architecture** yang membagi kode menjadi beberapa layer:

```
┌─────────────────────────────────────┐
│         MAIN (Presentation)         │  ← Orchestration & Testing
└────────────────┬────────────────────┘
                 │
┌─────────────────▼────────────────────┐
│         USECASE (Business Logic)     │  ← Definisikan logika aplikasi
└────────────────┬────────────────────┘
                 │
┌─────────────────▼────────────────────┐
│    DOMAIN (Business Rules/Entities)  │  ← Definisikan interface & entity
└────────────────┬────────────────────┘
                 │
┌─────────────────▼────────────────────┐
│ REPOSITORY (Data Access/Persistence) │  ← Implementasi konkret ke Redis
└─────────────────────────────────────┘
```

**Manfaat Clean Architecture:**

- 🔄 **Loosely Coupled:** Setiap layer independen
- 🧪 **Testable:** Mudah membuat unit test
- 🔧 **Maintainable:** Mudah di-maintain dan extend
- 📦 **Reusable:** Bisa ditukar implementasi storage (Redis → DB lain)

---

## Struktur Folder Proyek

```
Redis/
├── main.go                           # Entrypoint aplikasi (demo)
├── go.mod                            # Module definition
├── readme.md                         # Penjelasan prompt engineering
├── documentation/
│   └── docs.md                       # File ini (dokumentasi lengkap)
├── domain/                           # Layer Domain (Business Rules)
│   ├── user.go                       # Entity User & Interface UserRepository
│   └── cache.go                      # Interface CacheRepository
├── repository/
│   └── redis/                        # Implementasi dengan Redis
│       ├── client.go                 # Koneksi ke Redis
│       ├── user_repository.go        # Implementasi UserRepository
│       └── cache_repository.go       # Implementasi CacheRepository
└── usecase/                          # Layer Usecase (Business Logic)
    ├── user_usecase.go               # Logika bisnis User
    ├── user_usecase_test.go          # Unit test UserUsecase
    └── cache_usecase_test.go         # Unit test Cache
```

---

## Komponen Utama

### 1. Domain Layer

#### `domain/user.go`

- **User Struct:** Entitas represensasi user (ID, Name, Email, Age)
- **UserRepository Interface:** Kontrak untuk operasi user
  - `Save(user User) error` - Simpan user ke Redis
  - `FindByID(id string) (*User, error)` - Cari user berdasarkan ID
  - `Delete(id string) error` - Hapus user dari Redis

#### `domain/cache.go`

- **CacheRepository Interface:** Kontrak untuk operasi cache umum
  - `Set(key string, value string, ttl time.Duration) error`
  - `Get(key string) (string, error)`
  - `Delete(key string) error`
  - `Increment(key string) (int64, error)` - Counter operation
  - `PushQueue(queueName string, values ...string) error` - Masukkan ke Queue
  - `PopQueue(queueName string) (string, error)` - Ambil dari Queue

### 2. Repository Layer (Redis)

#### `repository/redis/client.go`

**Fungsi:** Menghubungkan aplikasi ke Redis server

```go
NewRedisClient() *redis.Client
```

- **Konfigurasi:**
  - Host: `localhost`
  - Port: `6379`
  - Database: `0`
  - Password: Kosong (default)
- **Error Handling:** Panic jika koneksi gagal

#### `repository/redis/user_repository.go`

**Implementasi UserRepository** menggunakan Redis:

- Menyimpan user sebagai JSON string di Redis
- Menggunakan ID user sebagai key: `user:{id}`
- Operasi atomic untuk konsistensi data

#### `repository/redis/cache_repository.go`

**Implementasi CacheRepository** dengan fitur:

- **String Operations:** SET, GET, DEL
- **Counter:** INCR untuk increment value
- **Queue (List):** LPUSH (masukkan kiri), RPOP (ambil kanan)
- **TTL:** Time-to-live untuk auto-expire key

### 3. Usecase Layer

#### `usecase/user_usecase.go`

Logika bisnis untuk mengelola user:

```go
type UserUsecase struct {
    userRepo domain.UserRepository
}
```

Methods:

- `SaveUser(user domain.User) error` - Simpan user baru
- `GetUser(id string) (*domain.User, error)` - Ambil user by ID
- `DeleteUser(id string) error` - Hapus user

### 4. Main Layer

#### `main.go`

Demo dan testing semua fitur:

- Inisialisasi Redis Client
- Test Queue operation (antrian tugas)
- Test String operation (tag)
- Test Counter operation (pengunjung)

---

## Operasi Redis

### 1. **String Operations** (Key-Value)

```
SET key value [EX seconds]    # Set value dengan optional TTL
GET key                        # Ambil value dari key
DEL key                        # Hapus key
```

**Contoh di code:**

```go
cacheRepo.Set("golang_tag", "redis backend", 0)
val, _ := cacheRepo.Get("golang_tag")
```

### 2. **Counter/Increment**

```
INCR key                       # Increment counter (atomic)
DECR key                       # Decrement counter
```

**Contoh di code:**

```go
total, _ := cacheRepo.Increment("pengunjung_counter")  // 1, 2, 3, 4, 5
```

### 3. **Queue/List Operations**

```
LPUSH key value [value ...]   # Push ke kiri (head)
RPOP key                       # Pop dari kanan (tail)
LLEN key                       # Jumlah element di list
```

**Contoh di code:**

```go
cacheRepo.PushQueue("antrian_tugas", "tugas-C", "tugas-B", "tugas-A")
diproses, _ := cacheRepo.PopQueue("antrian_tugas")
```

**Karakteristik Queue FIFO:**

```
Push: A, B, C        →  [C, B, A]  (LPUSH ke kiri)
                        ↓
Pop:                     A  (RPOP dari kanan)
                        [C, B]
```

---

## Setup & Instalasi

### Prerequisites

- Go 1.18 atau lebih baru
- Redis Server (disarankan via Docker)
- VSCode atau text editor apapun

### Step 1: Setup Redis dengan Docker

```bash
# Pull Redis image
docker pull redis:latest

# Jalankan Redis container
docker run --name redis-container -p 6379:6379 -d redis:latest

# Verify Redis berjalan
redis-cli ping
# Output: PONG
```

### Step 2: Clone/Setup Project

```bash
# Clone atau setup project di folder lokal
cd "Kuliah SMT6 - TRPL 3D/Topik Khusus/Redis"

# Download dependencies
go mod download
```

### Step 3: Jalankan Project

```bash
# Run main.go
go run main.go

# Output:
# Terhubung ke Redis: PONG
# Antrian:[tugas-DARURAT tugas-A tugas-B tugas-C]
# Diproses: tugas-DARURAT
# ...
```

---

## Cara Menggunakan

### Contoh 1: Menyimpan & Mengambil User

```go
package main

import (
    "fmt"
    "redis/domain"
    "redis/repository/redis"
    "redis/usecase"
)

func main() {
    // 1. Setup Redis
    rdb := redis.NewRedisClient()

    // 2. Buat repository
    userRepo := redis.NewUserRepository(rdb)

    // 3. Buat usecase
    userUC := usecase.NewUserUsecase(userRepo)

    // 4. Simpan user
    user := domain.User{
        ID:    "001",
        Name:  "Budi Santoso",
        Email: "budi@example.com",
        Age:   25,
    }
    userUC.SaveUser(user)

    // 5. Ambil user
    retrieved, _ := userUC.GetUser("001")
    fmt.Printf("User: %+v\n", retrieved)

    // 6. Hapus user
    userUC.DeleteUser("001")
}
```

### Contoh 2: Menggunakan Cache dengan TTL

```go
cacheRepo := redis.NewCacheRepository(rdb)

// Set dengan TTL 5 menit
cacheRepo.Set("user:profile", jsonData, 5*time.Minute)

// Get
profile, err := cacheRepo.Get("user:profile")
if err != nil {
    fmt.Println("Cache miss atau expired")
}

// Delete manual
cacheRepo.Delete("user:profile")
```

### Contoh 3: Menggunakan Queue

```go
// Push tasks ke queue
cacheRepo.PushQueue("tasks",
    "task-3",
    "task-2",
    "task-1",
)

// Process (FIFO)
task1, _ := cacheRepo.PopQueue("tasks")  // task-1
task2, _ := cacheRepo.PopQueue("tasks")  // task-2
task3, _ := cacheRepo.PopQueue("tasks")  // task-3
```

### Contoh 4: Counter untuk Metrics

```go
// Track pengunjung
cacheRepo.Increment("page_visits")
cacheRepo.Increment("page_visits")

// Get value
visits, _ := cacheRepo.Get("page_visits")
// visits = "2"
```

---

## Testing

### Unit Test Structure

Project ini menggunakan **testify/assert** untuk testing yang clean dan readable.

#### `usecase/user_usecase_test.go`

Tests untuk UserUsecase:

```go
func TestSaveUser_Success(t *testing.T) {
    // Setup mock repository
    // Execute SaveUser
    // Assert hasil
}

func TestSaveUser_Error(t *testing.T) {
    // Setup mock dengan error
    // Execute SaveUser
    // Assert error handling
}
```

#### `usecase/cache_usecase_test.go`

Tests untuk Cache operations:

- String operations (Set, Get, Delete)
- Counter operations
- Queue operations

### Menjalankan Unit Tests

```bash
# Run semua tests
go test ./usecase -v

# Run test spesifik
go test -run TestSaveUser_Success ./usecase -v

# Run dengan coverage
go test ./usecase -cover
```

### Test Output Format

```
=== RUN   TestSaveUser_Success
    user_usecase_test.go:15: Save user success
--- PASS: TestSaveUser_Success (0.05s)

=== RUN   TestSaveUser_Error
    user_usecase_test.go:35: Save user failed
--- PASS: TestSaveUser_Error (0.03s)

PASS
ok      redis/usecase   0.091s
```

---

## Troubleshooting

### Issue 1: Connection Refused (localhost:6379)

**Masalah:** `Error: Failed to connect to localhost:6379`

**Solusi:**

```bash
# Cek Redis running
docker ps | grep redis

# Jika tidak running, start Redis
docker run --name redis-container -p 6379:6379 -d redis:latest

# Verify koneksi
redis-cli ping
# Seharusnya output: PONG
```

### Issue 2: Permission Denied di Windows

**Masalah:** Port 6379 sudah dipakai

**Solusi:**

```bash
# Gunakan port berbeda
docker run --name redis-new -p 6380:6379 -d redis:latest

# Update client.go:
// Addr: "localhost:6380"
```

### Issue 3: Module Not Found (go-redis)

**Masalah:** `cannot find package "github.com/redis/go-redis/v9"`

**Solusi:**

```bash
# Download module
go get github.com/redis/go-redis/v9

# Verify di go.mod
cat go.mod
```

### Issue 4: Test Failures

**Masalah:** Unit tests fail

**Solusi:**

```bash
# Run debug mode dengan verbose
go test ./usecase -v

# Check mock setup correct
# Verify Redis mock behavior

# Run single test
go test -run TestSpecific ./usecase -v
```

---

## Best Practices

### 1. **Connection Management**

✅ DO:

```go
rdb := redis.NewRedisClient()  // Single instance (reuse)
defer rdb.Close()               // Close saat selesai
```

❌ DON'T:

```go
// Jangan buat client baru setiap kali
rdb := redis.NewClient(...)
```

### 2. **Key Naming Convention**

✅ DO:

```go
"user:001"          // prefix:id
"cache:profile"     // prefix:data_type
"queue:tasks"       // prefix:queue_name
```

❌ DON'T:

```go
"user001"           // Tidak jelas tipe data
"u_1_p_n"          // Abbreviation konfusing
```

### 3. **TTL Management**

✅ DO:

```go
cacheRepo.Set("session", data, 30*time.Minute)  // Explicit TTL
cacheRepo.Set("permanent", data, 0)             // 0 = no expiry
```

❌ DON'T:

```go
cacheRepo.Set("key", data, -1)  // Invalid TTL
```

### 4. **Error Handling**

✅ DO:

```go
if err != nil {
    log.Printf("Error: %v", err)
    return err  // Return error untuk upstream
}
```

❌ DON'T:

```go
cacheRepo.Get("key")  // Ignore error
```

### 5. **Testing**

✅ DO:

```go
// Mock repository
// Test isolation
// Assert behavior
```

❌ DON'T:

```go
// Bergantung ke real Redis
// Test terikat pada order
```

---

## Referensi & Resources

### Dokumentasi Official

- [Go Redis Library](https://github.com/redis/go-redis)
- [Redis Commands](https://redis.io/commands)
- [Go Testing Package](https://golang.org/pkg/testing)
- [Testify Assert](https://github.com/stretchr/testify/assert)

### Artikel Berguna

- Clean Architecture di Go
- Redis Data Structures
- Unit Testing Best Practices

### Tools & Scripts

```bash
# Inspect Redis dari CLI
redis-cli

# Commands
redis-cli PING                    # Test connection
redis-cli KEYS "*"               # List semua keys
redis-cli GET user:001           # Get specific key
redis-cli LLEN antrian_tugas      # Check queue length
redis-cli FLUSHDB                # Clear database
```

---

## Summary

| Aspek            | Deskripsi                                                 |
| ---------------- | --------------------------------------------------------- |
| **Arsitektur**   | Clean Architecture (Domain → Usecase → Repository → Main) |
| **Database**     | Redis (Key-Value, Counter, Queue)                         |
| **Framework**    | Go 1.18+, go-redis/v9                                     |
| **Testing**      | Unit tests dengan testify/assert                          |
| **Deployment**   | Docker Redis                                              |
| **Key Features** | User management, Caching, Counter, Queue                  |

---

**Last Updated:** March 2026  
**Version:** 1.0  
**Status:** Production Ready ✅

---

_Dokumentasi ini dibuat sebagai bagian dari Tugas Topik Khusus Semester 6 TRPL 3D_
