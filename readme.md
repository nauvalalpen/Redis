# MiniKatalog 🗂️

> **Aplikasi web katalog produk** sebagai proyek UAS Topik Khusus, mendemonstrasikan implementasi **Redis** dalam arsitektur web modern (Go + MongoDB + Redis).

## Tech Stack

| Layer | Teknologi |
|---|---|
| **Backend** | Go 1.21, gorilla/mux |
| **Database Utama** | MongoDB (localhost:27017) |
| **In-Memory Cache** | Redis (localhost:6379) |
| **Frontend** | HTML5 + Vanilla CSS + Vanilla JS |
| **Arsitektur** | Clean Architecture (domain → repository → usecase → handler) |

---

## Konsep Redis yang Diimplementasikan

| Konsep | Endpoint | Cara Kerja |
|---|---|---|
| **Session Management** | `POST /api/auth/login` | Token random disimpan di Redis dengan TTL 2 jam (`SET session:{token} {userID} EX 7200`) |
| **Rate Limiting** | `POST /api/auth/login` | INCR atomic counter per username, window 1 menit, max 5 percobaan (`INCR ratelimit:login:{username}` + `EXPIRE`) |
| **Cache-Aside Pattern** | `GET /api/products` | Cek Redis dulu - hit: return cache; miss: query MongoDB, simpan ke Redis TTL 5 menit |

---

## Prasyarat

- [Go 1.21+](https://go.dev/dl/)
- [MongoDB](https://www.mongodb.com/try/download/community) berjalan di `localhost:27017`
- [Redis](https://redis.io/download/) berjalan di `localhost:6379`

### Instalasi MongoDB & Redis (Windows)
```powershell
# MongoDB via winget:
winget install MongoDB.Server

# Redis via WSL atau installer resmi:
# https://github.com/microsoftarchive/redis/releases
```

---

## Cara Menjalankan

### 1. Clone & Download Dependencies
```powershell
# Masuk ke folder project
cd "d:\Kuliah SMT6 - TRPL 3D\Topik Khusus\Praktek\redis_nauval_pertemuan_2"

# Download semua dependencies Go
go mod tidy
```

### 2. Pastikan MongoDB & Redis Berjalan
```powershell
# Cek Redis
redis-cli ping
# Output: PONG

# Cek MongoDB  
mongosh --eval "db.runCommand({ ping: 1 })"
```

### 3. Jalankan Backend + Frontend
```powershell
go run main.go
```

Output yang muncul:
```
Koneksi MongoDB berhasil!
Koneksi Redis berhasil!

Backend API : http://localhost:8080/api
Frontend    : http://localhost:8080/app/
```

### 4. Buka Browser
- **Frontend**: http://localhost:8080/app/
- **API Base**: http://localhost:8080/api

---

## API Reference (untuk Postman)

### Auth Endpoints
```
POST /api/auth/register     - Daftar user baru
POST /api/auth/login        - Login (rate limited: 5x/menit)
POST /api/auth/logout       - Logout (hapus session Redis)
GET  /api/auth/me           - Validasi session (butuh token)
```

### Product Endpoints
```
GET    /api/products         - Ambil semua produk (cached di Redis)
GET    /api/products/{id}    - Ambil satu produk
POST   /api/products         - Buat produk baru (butuh token)
PUT    /api/products/{id}    - Update produk (butuh token)
DELETE /api/products/{id}    - Hapus produk (butuh token)
```

### Contoh Request Postman

**Register:**
```json
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "username": "nauval",
  "password": "password123"
}
```

**Login:**
```json
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "username": "nauval",
  "password": "password123"
}
```

**Buat Produk (gunakan token dari Login):**
```json
POST http://localhost:8080/api/products
Authorization: Bearer {token_dari_login}
Content-Type: application/json

{
  "name": "Laptop Gaming ASUS ROG",
  "price": 18000000,
  "category": "Elektronik",
  "description": "Laptop gaming dengan RTX 4060, RAM 16GB",
  "stock": 10
}
```

---

## Cara Membuktikan Efek Caching

### Via Postman (Paling Jelas)

1. **Request pertama**: `GET http://localhost:8080/api/products`
   - Response body: `"cache_status": "CACHE_MISS"`
   - Header: `X-Cache: CACHE_MISS`
   - Catat `response_time` (biasanya lebih lambat, ~10-50ms)

2. **Request kedua** (langsung tanpa jeda):
   - Response body: `"cache_status": "CACHE_HIT"`
   - Header: `X-Cache: CACHE_HIT`
   - `response_time` jauh lebih cepat (~1-5ms)

3. **Tambah produk baru** maka cache di-invalidate
4. **Request lagi** kembali `CACHE_MISS`

### Via Redis CLI (Real-time Monitoring)

```bash
# Monitor semua perintah Redis secara real-time
redis-cli MONITOR

# Cek apakah cache ada
redis-cli GET "cache:products:all"

# Cek TTL cache
redis-cli TTL "cache:products:all"

# Lihat semua key session yang aktif
redis-cli KEYS "session:*"

# Lihat rate limit counter
redis-cli KEYS "ratelimit:*"
```

---

## Cara Membuktikan Rate Limiting

Kirim login dengan password salah berulang kali (max 5x per menit).
Pada percobaan ke-6, respons akan berupa HTTP 429:
```json
{
  "error": "terlalu banyak percobaan login. Coba lagi dalam 58 detik"
}
```

---

## Struktur Project

```
redis_nauval_pertemuan_2/
├── main.go                          # Entry point (routing + DI)
├── go.mod / go.sum                  # Dependencies
├── frontend/
│   ├── index.html                   # UI utama
│   ├── style.css                    # Dark theme styling
│   └── app.js                       # Logic fetch API + demo cache
├── internal/
│   ├── config/
│   │   └── config.go                # Init MongoDB & Redis
│   ├── domain/
│   │   ├── user.go                  # Entitas User + DTO auth
│   │   └── product.go               # Entitas Product + DTO
│   ├── repository/
│   │   ├── user_repository.go       # CRUD user ke MongoDB
│   │   ├── product_repository.go    # CRUD produk ke MongoDB
│   │   └── redis_repository.go      # Session, Rate Limit, Cache ke Redis
│   ├── usecase/
│   │   ├── auth_usecase.go          # Logic auth + session + rate limit
│   │   └── product_usecase.go       # Logic CRUD produk + cache-aside
│   ├── handler/
│   │   ├── auth_handler.go          # HTTP handler auth
│   │   └── product_handler.go       # HTTP handler produk
│   └── middleware/
│       └── auth.go                  # Bearer token middleware
├── rencana-uas.md                   # Rencana dan arsitektur
├── progress.md                      # Log progress pengerjaan
└── README.md                        # File ini
```

---

## Database Schema

### MongoDB Collection: `users`
```json
{
  "_id": "ObjectId",
  "username": "string (unique)",
  "password_hash": "string (bcrypt)"
}
```

### MongoDB Collection: `products`
```json
{
  "_id": "ObjectId",
  "user_id": "ObjectId",
  "created_by": "string (username)",
  "name": "string",
  "price": "number",
  "category": "string",
  "description": "string",
  "stock": "number"
}
```

### Redis Key Design
```
session:{token}            -> {userID}          TTL: 2 jam
ratelimit:login:{username} -> {jumlah_coba}     TTL: 1 menit (window)
cache:products:all         -> {JSON string}     TTL: 5 menit
```
