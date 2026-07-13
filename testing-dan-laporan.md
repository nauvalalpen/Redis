# LAPORAN PRAKTIKUM TOPIK KHUSUS
## PERTEMUAN 2: IMPLEMENTASI CACHING, ANTRIAN (FIFO QUEUE), ATOMIC COUNTER, DAN TTL EXPIRE MENGGUNAKAN REDIS DAN GOLANG CLEAN ARCHITECTURE

**Dosen Pengampu:** [PERLU DIISI MANUAL: Nama Dosen]

---

### **IDENTITAS MAHASISWA**
*   **Nama:** [PERLU DIISI MANUAL: Nama Lengkap Anda]
*   **NIM:** [PERLU DIISI MANUAL: NIM Anda]
*   **Kelas:** TRPL 3D
*   **Program Studi:** D4 Teknologi Rekayasa Perangkat Lunak
*   **Jurusan:** Teknologi Informasi

---

### **I. TUJUAN PRAKTIKUM**

1.  Menganalisis dan menerapkan **Redis** sebagai *in-memory data structure store* yang berfungsi sebagai sistem caching berlapis, antrian pesan asinkron, dan pencatat status atomik pada arsitektur web modern.
2.  Mengimplementasikan **Clean Architecture** secara ketat di Go, dengan menegaskan batasan (*boundaries*) antara lapisan *Domain* (entitas dan kontrak), *Repository* (infrastruktur data), *Usecase* (logika bisnis), dan *Presentation/Main* (orkestrasi jalannya aplikasi).
3.  Memvalidasi integrasi operasi inti Redis langsung dari dalam kode Go melalui driver `go-redis/v9`:
    *   **FIFO Queue:** Menggunakan struktur data List (via `LPUSH` dan `RPOP`) untuk antrian tugas atau pesan.
    *   **String Key-Value:** Menggunakan `SET` dan `GET` sebagai mekanisme penyimpanan status sederhana atau caching data berulang.
    *   **Atomic Counter:** Memanfaatkan `INCR` guna memastikan keamanan konkuren (*thread-safe*) saat mengubah nilai penghitung tanpa perlu mekanisme *locking* tambahan.
    *   **TTL / Auto-Expiry:** Menetapkan siklus hidup data (*Time-To-Live*) secara imperatif untuk mengelola sesi atau cache kadaluwarsa otomatis.
4.  Menjalankan pengujian komponen terisolasi menggunakan teknik *manual mocking* pada Unit Test layer Usecase. Praktik ini membuktikan sistem dapat dites meskipun layanan eksternal (Redis server) tidak beroperasi atau tidak tersedia.

---

### **II. SKENARIO PENGUJIAN MENDALAM**

Bagian ini merinci instruksi teknis untuk menguji seluruh komponen yang telah dibuat di proyek `redis_nauval_pertemuan_2`. 

#### **1. Pengujian Koneksi Inisial ke Redis (`client.go`)**
**Tujuan:** Memastikan aplikasi Go dapat membuat TCP connection yang valid ke service Redis Docker.
*   **Precondition (State Awal):** 
    *   Docker Desktop harus berjalan.
    *   Container Redis telah dijalankan dengan command `docker run -d --name redis-praktikum -p 6379:6379 redis:alpine`
    *   Tidak ada proxy/firewall lokal yang memblokir akses ke `localhost:6379`.
*   **Command Persis:**
    ```powershell
    cd "d:\Kuliah SMT6 - TRPL 3D\Topik Khusus\Praktek\redis_nauval_pertemuan_2"
    go run main.go
    ```
*   **Expected Output:**
    ```text
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
*   **Edge Case / Uji Skenario Negatif:** 
    *   *Kondisi:* Hentikan paksa container Redis: `docker stop redis-praktikum`. Lalu jalankan `go run main.go`.
    *   *Expected behavior:* Aplikasi langsung *panic* (*crash* secara disengaja) pada baris `repo.NewRedisClient()`. Ini memvalidasi desain *fail-fast*, di mana aplikasi menolak untuk startup jika dependency utama (Redis) tidak tersedia.

#### **2. Pengujian FIFO Queue / Antrian (`cache_repository.go` - PushQueue & PopQueue)**
**Tujuan:** Memvalidasi operasi `LPUSH` dan `RPOP` berjalan secara First-In-First-Out, memastikan urutan eksekusi tugas dalam antrian.
*   **Precondition:** Redis berjalan normal. Key `"antrian_tugas"` dalam status kosong (atau akan dihapus di step 1).
*   **Command Persis:**
    Output ini otomatis tereksekusi pada kelanjutan `go run main.go`. Secara eksplisit, kode di `main.go` baris 36 mengeksekusi: `cacheRepo.PushQueue("antrian_tugas", "tugas-C", "tugas-B", "tugas-A", "tugas-DARURAT")`.
*   **Expected Output:**
    ```text
    --------------------------------------------------------------------------------
    2. QUEUE/ANTRIAN - LPUSH & RPOP (FIFO)
    --------------------------------------------------------------------------------
    Alur Bisnis: Sistem antrian untuk memproses tugas secara berurutan (FIFO)

    [Step 1] Clear queue lama (memastikan fresh start)
    Command: DEL antrian_tugas

    [Step 2] Masukkan 4 tugas dengan LPUSH (push ke head/kiri)
    Command: LPUSH antrian_tugas tugas-DARURAT tugas-A tugas-B tugas-C
    Struktur Queue: [tugas-C] <- [tugas-B] <- [tugas-A] <- [tugas-DARURAT]

    [Step 3] Ambil tugas pertama dengan RPOP (pop dari tail/kanan)
    Command: RPOP antrian_tugas
    Tugas Diproses: tugas-C
    Sisa Queue: 3 tugas
    Struktur Queue Sekarang: [tugas-C] <- [tugas-B] <- [tugas-A]
    ```
*   **Edge Case / Uji Skenario Negatif:**
    *   *Kondisi:* Antrian kosong (dipanggil `PopQueue` berturut-turut lebih dari 3 kali tanpa diisi).
    *   *Expected behavior:* `go-redis` mengembalikan error tipe `redis.Nil`. Karena di `main.go` error ini di-ignore dengan wildcard `_`, log `Tugas Diproses:` akan mencetak string kosong. Ini mengidentifikasi batasan kode saat ini, di mana sistem belum menangani gracefull response jika antrian sudah habis.

#### **3. Pengujian String Key-Value (`cache_repository.go` - Set & Get dengan TTL=0)**
**Tujuan:** Menguji primitif data struktur sederhana Redis untuk menyimpan tag statis.
*   **Precondition:** Key `"golang_tag"` bisa berisi atau belum. Akan di-*overwrite* oleh `SET`.
*   **Command Persis:** Dilanjutkan dalam eksekusi `go run main.go`. `Set("golang_tag", "redis backend", 0)` dijalankan.
*   **Expected Output:**
    ```text
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
*   **Edge Case / Uji Skenario Negatif:**
    *   *Kondisi:* Mengambil key yang tidak pernah diset sebelumnya (misal: `Get("kunci_hilang")`).
    *   *Expected behavior:* Fungsi `Get()` mengembalikan error `redis.Nil`. Hasil return value dari *method* akan berupa string kosong `""` dan pointer error terisi. Aplikasi tidak boleh panic, tapi merespons bahwa data *cache miss*.

#### **4. Pengujian Atomic Counter (`cache_repository.go` - Increment)**
**Tujuan:** Mensimulasikan *race condition prevention* menggunakan operasi atomik `INCR`.
*   **Precondition:** Key `"pengunjung_counter"` dibersihkan di awal eksekusi (`DEL pengunjung_counter`).
*   **Command Persis:** Dieksekusi otomatis dalam *loop* iterasi `for i := 1; i <= 5` di `main.go`.
*   **Expected Output:**
    ```text
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
*   **Edge Case / Uji Skenario Negatif:**
    *   *Kondisi:* Key sudah berisi tipe data non-integer (contoh: `SET pengunjung_counter "satu"` lalu di `INCR`).
    *   *Expected behavior:* Redis menolak operasi (mengembalikan `ERR value is not an integer or out of range`). Fungsi `Increment()` akan me-*return* *error*, total *counter* menjadi 0 dan operasi dibatalkan di sisi Go.

#### **5. Pengujian TTL / Expiry Sesi (`cache_repository.go` - Set dengan Time-to-Live)**
**Tujuan:** Membuktikan eviksi otomatis data di sisi server (Redis) setelah interval waktu tertentu.
*   **Precondition:** Key `"sesi_token"` di-SET dengan `TTL 5*time.Second`.
*   **Command Persis:** Bagian akhir dari `go run main.go`. Tunggu selama 6 detik setelah output countdown selesai dicetak.
*   **Expected Output:**
    ```text
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
*   **Edge Case / Uji Skenario Negatif:**
    *   *Kondisi:* Server Go mengalami *delay* (misalnya proses berat lain) sehingga `Get` SEBELUM expire tidak bisa berjalan tepat waktu (misal dieksekusi detik ke-6 alih-alih detik ke-1).
    *   *Expected behavior:* Token tetap akan terbaca *EXPIRED* karena state timer dijaga oleh server Redis, bukan oleh waktu runtime Golang. Data benar-benar konsisten dan tidak bergantung pada seberapa cepat golang merespons.

#### **6. Pengujian Lapisan Logika Bisnis Terisolasi (Unit Testing Usecase & Domain)**
**Tujuan:** Memvalidasi pemisahan *concern* dari Clean Architecture. Lapisan bisnis (usecase) harus lulus uji verifikasi menggunakan manual mock *tanpa membutuhkan dependensi eksternal (Redis server)*.
*   **Precondition:** Container Redis boleh menyala atau **dimatikan**. Hasil test wajib tidak terpengaruh oleh ketersediaan service Redis.
*   **Command Persis:**
    ```powershell
    cd "d:\Kuliah SMT6 - TRPL 3D\Topik Khusus\Praktek\redis_nauval_pertemuan_2"
    go test ./usecase/... -v -cover
    ```
*   **Expected Output:**
    ```text
    === RUN   TestCacheSet_Success
    --- PASS: TestCacheSet_Success (0.00s)
    === RUN   TestCacheSet_Error
    --- PASS: TestCacheSet_Error (0.00s)
    === RUN   TestCreateUser_Success
    --- PASS: TestCreateUser_Success (0.00s)
    === RUN   TestCreateUser_Error
    --- PASS: TestCreateUser_Error (0.00s)
    PASS
    coverage: [persentase]% of statements
    ok      redis/usecase   0.XXXs
    ```
*   **Edge Case / Uji Skenario Negatif:**
    *   *Kondisi:* Mock pada `TestCreateUser_Error` (`saveFn`) sengaja di-*inject* untuk mengembalikan suatu fatal error tak terduga.
    *   *Expected behavior:* Aplikasi lulus testing (`PASS`) di case Error tersebut, membuktikan layer *usecase* tidak "menelan" error (menyembunyikan) dan tetap mempropagasi error infrastruktur (di-*simulate* mock) langsung ke layer teratas untuk dilaporkan.

---

### **III. CHECKLIST SCREENSHOT**

Berikut adalah pedoman presisi (*granular*) pengumpulan dokumentasi tangkapan layar, memastikan bukti yang terlampir komprehensif.

| No | Disarankan Nama File | Sumber Tangkapan (Source) | Area Fokus yang Harus Tampak Jelas | Penggunaan di Bab Laporan |
|:---|:---|:---|:---|:---|
| 1 | `ss01_redis_container.png` | **Docker Desktop GUI / Terminal** (`docker ps`) | Baris container `redis-praktikum`, image `redis:alpine`, Status `Up` dan Port mapping `6379->6379/tcp`. | **Langkah Implementasi** (Konfigurasi Environment awal). |
| 2 | `ss02_koneksi_dan_queue.png` | **Terminal PowerShell** (Output awal `go run main.go`) | Header "DEMO SISTEM REDIS", baris konfirmasi `Status: Terhubung ke Redis`, dan keseluruhan blok simulasi **Step 2 (Queue LPUSH/RPOP)** yang menunjukkan "Tugas Diproses: tugas-C". | **Hasil Praktikum** (Skenario Demo 1 & Demo 2). |
| 3 | `ss03_string_dan_counter.png` | **Terminal PowerShell** (Lanjutan eksekusi) | Blok log **Step 3 (String SET & GET)** menunjukkan *Key exists? true*, dan blok **Step 4 (Atomic INCR)** menunjukan *Pengunjung #1* s.d *#5* secara inkremental tanpa angka meloncat. | **Hasil Praktikum** (Skenario Demo 3 & Demo 4). |
| 4 | `ss04_ttl_countdown_awal.png` | **Terminal PowerShell** (Saat program masih jalan) | Blok log **Step 5 (TTL/Expiry)**. Tangkap detik-detik proses "Tunggu 6 detik hingga token expire..." saat countdown masih mencetak angka `6 5 4` dan status masih `✅ (Token VALID)`. | **Dasar Teori / Hasil Praktikum** (Menjelaskan fenomena TTL internal Redis vs Aplikasi Go). |
| 5 | `ss05_ttl_expired_summary.png` | **Terminal PowerShell** (Akhir eksekusi) | Output baris "Get token SETELAH expire" menunjukkan `Result: (nil) ❌ (Token EXPIRED)`, diikuti oleh header kesimpulan `✅ DEMO SELESAI`. | **Hasil Praktikum** (Kesimpulan Demo 5). |
| 6 | `ss06_unit_test_isolated.png` | **Terminal PowerShell** (`go test ./usecase/... -v -cover`) | Memperlihatkan secara rinci 4 fungsi Test (`TestCacheSet_Success`, `Error`, `TestCreateUser_Success`, `Error`) semua tertulis `PASS`, serta persentase *coverage statements*. | **Analisis Arsitektur / Hasil Praktikum** (Bukti validasi *Dependency Inversion Principle* layer Usecase). |
| 7 | `ss07_redis_cli_manual.png` (Opsional) | **Terminal PowerShell 2** (`docker exec -it redis-praktikum redis-cli`) | Eksekusi manual CLI `KEYS *`, kemudian `GET pengunjung_counter` (harus berisi "5"), lalu `TTL sesi_token` (membuktikan return -2 alias expired). | **Kendala & Solusi / Eksplorasi Ekstra** (Verifikasi *Data Consistency* server-side). |

---

### **IV. OUTLINE LAPORAN (PANDUAN MENDALAM)**

Gunakan kerangka outline berikut untuk penulisan laporan PDF/Word yang memiliki bobot analitis (*analytical depth*) level senior/engineer.

#### **A. Pendahuluan**
*   **Latar Belakang Teknologi:** Jelaskan korelasi bottleneck *I/O Bound* pada Relational Database Management System (RDBMS) biasa jika digunakan sebagai cache/session layer, dan *mengapa* In-Memory database menjadi krusial di backend modern (misal: menangani ribuan *request per second* tanpa saturasi disk). 
*   **Relevansi Kombinasi Go & Redis:** Go terkenal dengan performa *concurrent* ekstrim (goroutine). Redis terkenal dengan kecepatan ekstrim yang bersifat single-threaded per task. Keduanya adalah kombinasi sempurna bagi arsitektur berlatensi sangat rendah (sub-millisecond latencies).
*   **Tujuan Laporan:** Menguraikan hasil porting teori tersebut ke dalam program Go dengan aturan Clean Architecture (pemisahan infra dan bisnis rule).

#### **B. Dasar Teori & Implementasi Domain (*Deep Dive*)**
*   **Konsep Arsitektur:** Jangan hanya sekedar mendeskripsikan "apa" itu Clean Architecture. Melainkan jelaskan "MENGAPA" layer terbagi 4. Sentuh konsep krusial **Dependency Inversion Principle (DIP)**: Layer logik bisnis (*usecase*) tidak boleh mengimpor langsung framework eksternal (`go-redis/v9`). Usecase harus hanya bergantung pada *contract* (`domain.CacheRepository`). Implementasi nyatanya (*redis.client*) ada di Layer terluar (*repository*), yang baru dikawinkan lewat Dependency Injection di file `main.go`.
*   **Operasi Dasar Redis:** 
    *   *TTL & Memory Eviction:* Bagaimana jika memory penuh? Redis memiliki *eviction policy* (contoh: *allkeys-lru*).
    *   *INCR Thread Safety:* Jelaskan mengapa operasi atomik `INCR` lebih baik daripada mengambil nilai (GET), menambah di Go (val+1), lalu menaruh kembali (SET) (mencegah *race-condition* saat >100 user hit endpoint bersamaan).
    *   *FIFO via List:* Menjelaskan struktur internal `List` Redis (linked list ganda). `LPUSH` menaruh antrian di head (kiri), dan `RPOP` mengambil dari tail (kanan), menghasilkan prinsip *First-In-First-Out*. Secara intuitif, argument awal yang di-push masuk ke head dan terdorong ke posisi tail, dan dikeluarkan lebih dulu.

#### **C. Arsitektur dan Desain Sistem**
Laporan Anda WAJIB menggambar minimal dua (2) diagram berikut (bisa pakai draw.io, Mermaid, atau visio):
1.  **Diagram Lapisan Clean Architecture (Onion Architecture):**
    *   *Komponen:* Lingkaran terdalam (Domain -> entitas User, interface UserRepository, CacheRepository), Lingkaran tengah (Usecase -> UserUsecase), Lingkaran terluar (Repository -> `cache_repository.go` + Redis Server).
    *   *Arah Ketergantungan:* Gambar panah yang secara eksplisit HANYA masuk ke arah dalam.
2.  **Diagram Sekuensial State Queue Redis:**
    *   *State 1:* Kosong.
    *   *State 2 (LPUSH):* Array/List berisi: `[Head] "tugas-DARURAT" -> "tugas-A" -> "tugas-B" -> "tugas-C" [Tail]`.
    *   *State 3 (RPOP):* Tail `"tugas-C"` ditarik (return ke aplikasi Go), sisa List adalah 3 elemen teratas.

#### **D. Hasil Praktikum (Observasi)**
*   Tempelkan screenshot berdasarkan tabel *Checklist Screenshot* (Bab III). 
*   **Wajib Diisi:** Tambahkan opini hasil pengamatan di bawah gambar. (Contoh: "Pada screenshot 5, terbukti bahwa delay yang dilakukan fungsi sleep Go langsung berdampak hilangnya cache data pada Redis server secara *absolute* ketika diakses.").

#### **E. Analisis Kendala Nyata & Solusi (*Real World Limitations*)**
Identifikasi batasan riil dan spesifik (bukan generalisasi internet mati) dari *source code* proyek `redis_nauval_pertemuan_2`:

1.  **Kendala 1: Tidak Adanya Validasi Struct Entitas (Business Logic Flaw).**
    *   *Analisis:* Dalam `user_usecase.go`, ketika method `CreateUser` dipanggil, struct parameter `domain.User` dilempar secara buta ke repo `Save(user)`. Apabila `ID` kosong, sistem Redis menimpa semua data user ke dalam satu kunci `"user:"` statis.
    *   *Solusi Konkret:* Menambahkan validasi `if user.ID == ""` me-return standard error pada layer usecase, mencegah data rusak sampai ke server Redis.
    *   *[PERLU DIISI MANUAL: Tulis pengamatan pribadi: Apakah layer Usecase di kode praktikum ini masih tergolong sangat "tipis" (*anemic*)?]*

2.  **Kendala 2: Silent Error / Error Handling Diabaikan pada Operator RPOP.**
    *   *Analisis:* Pada `main.go`, kode memanggil `diproses, _ := cacheRepo.PopQueue(...)`. Ini berbahaya (*silent error masking*). Jika antrian kosong, Redis mengembalikan error `redis.Nil`, tapi program Go justru meneruskannya sebagai eksekusi wajar mencetak tugas kosong.
    *   *Solusi Konkret:* Seharusnya error tidak di *blank-identifier* (`_`). Harus diproses: `if errors.Is(err, redis.Nil) { log.Print("Antrian habis, worker standby...") }`.
    *   *[PERLU DIISI MANUAL: Apakah Anda melakukan verifikasi apa yang terjadi jika RPOP di-run ketika queue sudah kosong (0 sisa tugas)?]*

3.  **Kendala 3: Mekanisme Fail-Fast yang Agresif pada Koneksi.**
    *   *Analisis:* Fungsi `NewRedisClient()` di layer repository berani menggunakan `panic(err)` jika `Ping()` gagal (Docker down). Di development sangat membantu, tapi di production ini berakibat satu *node worker* Go padam total dan perlu intervensi sistem operasi (misal Docker Swarm / Kubernetes restart loop).
    *   *Solusi Konkret:* Mengganti panic dengan implementasi *Retry-Backoff Algorithm* atau membiarkan koneksi di-*retry* secara perlahan menggunakan `Ping()` secara berkala.

#### **F. Kesimpulan**
*   Simpulkan pencapaian. Rangkum bagaimana *interface* Go sangat sakti untuk memungkinkan *Mocking* tanpa perlu Redis nyata saat test.
*   Konfirmasi bahwa kecepatan *memory RAM* Redis dipadukan eksekusi *compiled-binary* Golang menjadikan arsitektur ini standar de-facto untuk fitur *rate-limiter*, *session store*, maupun *message broker*.
*   *[PERLU DIISI MANUAL: Masukan personal terkait kemudahan library `go-redis/v9`]*

---

### **V. REFERENSI**
1. Redis Documentation Team. (2025). *Redis Data Types & Operations (List, String, Expire, Incr)*. https://redis.io/docs
2. Uptrace. (2025). *go-redis/v9 Official Guide*. https://redis.uptrace.dev/
3. Martin, Robert C. (2017). *Clean Architecture: A Craftsman's Guide to Software Structure and Design*. Prentice Hall.
4. Go Dev. (2025). *Testing in Go: Interfaces and Mocking Strategies*. https://pkg.go.dev/testing
5. Docker Inc. (2025). *Redis Official Image Deployment*. https://hub.docker.com/_/redis
