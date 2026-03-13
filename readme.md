# 🚀 Golang + Redis (Clean Architecture)

Tugas Topik Khusus: Implementasi Caching Redis menggunakan Golang dengan pendekatan Clean Architecture dan Unit Testing.

---

## 💬 Prompt Engineering History

Project ini di-generate menggunakan AI Assistant (GitHub Copilot)

### Tahap 1: Setup Struktur & Domain

_Di tahap awal, saya minta AI buatin kerangka dasar dan kontrak interface-nya dulu biar struktur Clean Architecture-nya kebentuk._

> **Prompt:**
> Buatin kerangka project Golang pake Clean Architecture buat konek ke Redis. Bikin folder `domain`, `repository/redis`, sama `usecase`.
> Di folder `domain`, tolong buatin:
>
> 1. `user.go` (isi struct User dan interface UserRepository)
> 2. `cache.go` (isi interface CacheRepository buat operasi Redis dasar kayak Set, Get, Delete, Increment, PushQueue, sama PopQueue).

### Tahap 2: Implementasi Repository (Redis)

_Setelah kerangka beres, saya suruh AI ngisi kodingan aslinya yang langsung nembak ke database Redis._

> **Prompt:**
> Oke bagus. Sekarang isi folder `repository/redis`-nya. Tolong pake library `github.com/redis/go-redis/v9` ya.
> Buatin `client.go` buat koneksi ke localhost:6379 karena saya menggunakan redis dari docker. Terus buatin juga implementasinya di `user_repository.go` dan `cache_repository.go` sesuai interface yang tadi udah dibikin.

### Tahap 3: Usecase & Unit Testing

_Nah, ini bagian krusial karena dosen minta mocking-nya dibikin manual, jadi saya harus wanti-wanti AI-nya biar gak pake tools otomatis._

> **Prompt:**
> Lanjut bikin file `user_usecase.go` yang manggil UserRepository.
> Kalau udah, buatin file unit test-nya (`user_usecase_test.go` & `cache_usecase_test.go`) pake library `testify/assert`. Buatin test case untuk skenario sukses dan skenario error (pake package `errors`).

### Tahap 4: (Main.go)

_Terakhir, minta AI buatin script buat ngetest semua fungsinya jalan atau nggak di terminal._

> **Prompt:**
> "Terakhir, buatin `main.go` di folder root untuk ngetest semuanya.
> Buat output simulasi yang rapi nantinya

---
