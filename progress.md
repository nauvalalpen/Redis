# Progress Log — MiniKatalog UAS

> Update file ini setiap fitur selesai. Format: [x] = selesai, [/] = sedang dikerjakan, [ ] = belum.

---

## Status: 🚧 Dalam Pengerjaan

**Tanggal mulai**: 2026-07-13  
**Deadline target**: 3 hari kerja

---

## Fitur Wajib (HARUS SELESAI)

- [x] **Struktur Project** — Clean Architecture (config, domain, repository, usecase, handler, middleware)
- [x] **Koneksi Redis** — `internal/config/config.go` `InitRedis()`
- [x] **Koneksi MongoDB** — `internal/config/config.go` `InitMongo()`
- [x] **Domain entities** — `User`, `Product`, request/response DTOs
- [x] **Redis Repository** — Session, Rate Limit, Cache (SET/GET/DEL, INCR, EXPIRE)
- [x] **User Repository** — MongoDB CRUD (Create, FindByUsername, FindByID)
- [x] **Product Repository** — MongoDB CRUD (Create, GetAll, GetByID, Update, Delete)
- [x] **Auth Usecase** — Register + Login (Rate Limiting) + Logout + ValidateSession
- [x] **Product Usecase** — CRUD + Cache-Aside Pattern (GetAll dengan HIT/MISS)
- [x] **Auth Handler** — `POST /api/auth/register`, `POST /api/auth/login`, `POST /api/auth/logout`
- [x] **Product Handler** — `GET|POST /api/products`, `GET|PUT|DELETE /api/products/{id}`
- [x] **Auth Middleware** — Bearer token validation dari Redis session
- [x] **Main entry point** — Router gorilla/mux, serve frontend statis, CORS middleware
- [x] **Frontend HTML** — index.html dengan struktur lengkap
- [x] **Frontend CSS** — Dark theme profesional
- [x] **Frontend JS** — Fetch API, auth flow, cache demo panel
- [x] **README.md** — Cara install, jalankan, dan demo caching

## Perlu Dilakukan Berikutnya

- [ ] **go mod tidy** — Download semua dependencies
- [ ] **Build & Test kompilasi** — `go build ./...`
- [ ] **Test manual Register + Login** — via Postman atau curl
- [ ] **Test Rate Limiting** — Login salah 6x, pastikan HTTP 429
- [ ] **Test CRUD Produk** — Create, baca daftar, hapus
- [ ] **Test Cache Effect** — 2x GET /api/products, bandingkan MISS vs HIT
- [ ] **Test Frontend di browser** — Login, tambah produk, lihat cache panel
- [ ] **Push ke GitHub**

---

## Fitur Nice-to-Have (Skip jika waktu habis)

- [ ] ~~Redis Pub/Sub (Notifikasi real-time)~~ — SKIP untuk sementara
- [ ] Update produk via UI frontend — edit form di-frontend (saat ini PUT endpoint ada di backend)
- [ ] Pagination produk
- [ ] Filter/search produk

---

## Log Masalah & Solusi

| Tanggal | Masalah | Solusi | Status |
|---|---|---|---|
| — | — | — | — |

*(Isi jika ada blocker saat testing)*

---

## Catatan Testing

### Test 1: Register & Login
```
Tanggal: ____
Endpoint: POST /api/auth/register, POST /api/auth/login
Hasil: ____
Token didapat: ____
```

### Test 2: Rate Limiting  
```
Tanggal: ____
Percobaan gagal ke-1 s/d ke-5: ____
Percobaan ke-6 (harusnya 429): ____
```

### Test 3: Cache-Aside Pattern
```
Tanggal: ____
Request ke-1 (CACHE_MISS) response_time: ____ms
Request ke-2 (CACHE_HIT) response_time:  ____ms
Selisih: ____ms (____x lebih cepat)
```
