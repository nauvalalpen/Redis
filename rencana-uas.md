# Rencana Proyek UAS Topik Khusus: "MiniKatalog"

## 1. Deskripsi Aplikasi
**MiniKatalog** adalah aplikasi web sederhana untuk manajemen dan penelusuran katalog produk. Aplikasi ini dipilih karena domainnya lugas, mudah dikerjakan dalam waktu singkat, dan sangat cocok untuk mendemonstrasikan keunggulan Redis dalam skenario *high-read traffic* (banyak yang melihat produk) serta pengamanan endpoint sensitif.

## 2. Fitur & Mapping ke Konsep Redis

| Fitur | Implementasi & Database | Konsep Redis yang Dipakai |
| --- | --- | --- |
| **Autentikasi (Login/Register)** | Session management. User data di MongoDB, Session di Redis. | **String Key-Value (SET/GET/TTL)**: Token session sebagai key, User ID sebagai value. Otomatis kedaluwarsa (TTL). |
| **Rate Limiting Login** | Membatasi maksimal 5x gagal login per menit dari satu IP/Username. | **Counter (INCR/EXPIRE)**: Atomic counter untuk melacak jumlah percobaan per rentang waktu. |
| **Manajemen Produk (CRUD)** | Menyimpan data produk yang dibuat user. Disimpan di MongoDB. | *Tidak pakai Redis, murni MongoDB sebagai primary persistence.* |
| **Caching List Produk** | Mempercepat *response time* saat user melihat daftar produk. | **Cache-Aside Pattern (SET/GET/TTL)**: Cek cache Redis, jika *miss*, ambil dari Mongo, simpan ke Redis dengan TTL, lalu *return*. |
| *(Opsional)* **Notifikasi Real-time** | Memberitahu pengunjung lain saat ada produk baru yang di-post. | **Pub/Sub (PUBLISH/SUBSCRIBE)**: Mengirim event "NewProduct" ke klien yang sedang membuka web. |

## 3. Tech Stack
Berdasarkan evaluasi repository `redis_nauval_pertemuan_2` dan `nosql_mongodb_nauval_pertemuan_4`:
*   **Backend**: Golang. Memakai framework routing `gorilla/mux` (diadaptasi dari repo mongodb) karena ringan dan cepat dikerjakan.
*   **Database Utama**: **MongoDB** (`go.mongodb.org/mongo-driver`). Dipilih karena sudah ada basis kodenya di pertemuan 4 (skema product & customer).
*   **In-Memory DB**: **Redis** (`go-redis/v9`). Dipakai sesuai target utama UAS.
*   **Arsitektur Backend**: Menggunakan pola *Clean Architecture* sederhana (seperti di repo pertemuan 2: `domain`, `repository`, `usecase`).
*   **Frontend**: HTML5, Vanilla CSS, Vanilla JS (Fetch API). Pendekatan Monolith (disajikan langsung dari backend Go sebagai folder statis) agar tidak perlu repot setup *build tools* seperti React/Vite. Prioritas utama adalah fungsi, bukan desain.

## 4. Struktur Project
Direkomendasikan meletakkan semuanya dalam satu repository agar mudah dijalankan:
```text
/
├── cmd/
│   └── main.go               # Entry point
├── frontend/
│   ├── index.html            # UI Utama
│   ├── style.css             # Styling sederhana
│   └── app.js                # Logika fetch API
├── internal/
│   ├── config/               # Setup koneksi Mongo & Redis
│   ├── domain/               # Struct/Entitas dan interface
│   ├── repository/           # Implementasi query Mongo & perintah Redis
│   ├── usecase/              # Logika bisnis (cache-aside, rate limit)
│   └── delivery/http/        # Route handler (gorilla/mux)
├── go.mod
└── rencana-uas.md
```

## 5. Daftar Entitas Data (Untuk SKPL/Perancangan)

**1. Collection `users` (MongoDB)**
*   `_id`: ObjectID
*   `username`: String (Unique)
*   `password_hash`: String

**2. Collection `products` (MongoDB)**
*   `_id`: ObjectID
*   `user_id`: ObjectID (Reference ke pembuat)
*   `name`: String
*   `price`: Number
*   `description`: String

**3. Struktur Data di Redis (Key Design)**
*   Session: `session:{token}` -> `user_id` (TTL: 2 jam)
*   Rate Limit: `ratelimit:login:{ip_address}` -> `jumlah_percobaan` (TTL: 1 menit)
*   Cache Produk: `cache:products:all` -> `JSON string` (TTL: 5 menit)

## 6. Timeline Realistis (3 Hari) & Strategi Eksekusi
*Sangat Ketat - Fokus pada KEPASTIAN SELESAI*

*   **Hari 1: Fondasi & Primary Database (MongoDB)**
    *   Setup kerangka folder Clean Architecture.
    *   Setup koneksi MongoDB.
    *   Buat endpoint API Register (Insert MongoDB) dan CRUD Produk sederhana (Insert & Get All ke MongoDB dulu, tanpa Redis).
*   **Hari 2: Injeksi Redis (Core Requirement)**
    *   Setup koneksi Redis.
    *   Buat endpoint API Login: implementasi simpan token ke Redis (Session) & validasi lewat Middleware.
    *   Tambahkan logika Rate Limiting pada API Login memakai INCR.
    *   Modifikasi endpoint Get All Products: bungkus dengan logika *Cache-Aside* menggunakan Redis. Siapkan endpoint timer/logger untuk melihat beda performa *cache hit* vs *cache miss*.
*   **Hari 3: Frontend Minimalis & Penyerahan**
    *   Buat `index.html` dengan form login dan tabel/list produk.
    *   Sambungkan UI dengan Fetch API ke endpoint backend.
    *   Buktikan fitur (Demo siap): Coba brute force login, coba lihat produk berulang kali untuk membuktikan cache jalan.
    *   *Jika sisa waktu:* Tambah notifikasi Pub/Sub sederhana, atau merapikan UI.

## 7. Kriteria Minimum (MVP) vs Boleh Dikorbankan
Mengingat waktu kurang dari 3 hari, inilah penentuan prioritas:

**🟢 WAJIB ADA (Kriteria Selesai Minimum):**
1. Bisa Register dan Login (Session di Redis).
2. API ambil daftar produk **di-cache** di Redis (bukti response time lebih cepat).
3. API Login memiliki pembatasan percobaan gagal (Rate Limit Redis).
4. Bisa tambah (Create) produk baru yang tersimpan di MongoDB.
5. Frontend *sangat sederhana* bisa mendemokan fungsi login, tambah produk, dan lihat produk.

**🔴 BOLEH DIKORBANKAN (Skip jika waktu habis di Hari ke-3):**
1. **Fitur Update dan Delete Produk**. Jika waktu habis, sisakan hanya Create dan Read. Itu sudah cukup untuk CRUD dasar dan demonstrasi cache.
2. **UI/UX yang Kompleks/Estetik**. Pakai alert native browser (`alert()`) untuk error form. Tampilan *bare-bones* tabel HTML biasa tidak masalah, asalkan *flow* backend terbukti jalan.
3. **Redis Pub/Sub (Notifikasi real-time)**. Ini adalah fitur yang paling memakan waktu integrasinya di frontend (harus pakai Websocket/SSE). Segera coret fitur ini jika Hari ke-2 backend belum beres 100%.
4. **Relasi antar entitas yang rumit**. Hindari operasi *join*/`$lookup` kompleks di MongoDB, cukup simpan data secara datar (flat).
