# 🏆Sistem Pelaporan Prestasi Mahasiswa

Aplikasi backend untuk pencatatan prestasi mahasiswa berbasis Go Fiber, PostgreSQL, dan MongoDB, dengan otentikasi JWT dan role-based access sesuai SRS.

## 📌 Fitur Utama

- 🔐 Autentikasi & Otorisasi JWT
- 👨‍🎓 Manajemen Mahasiswa
- 👨‍🏫 Manajemen Dosen & Advisees (Mahasiswa Bimbingan)
- 🏅 Manajemen Prestasi Mahasiswa
- 📎 Upload & Manajemen Lampiran Prestasi
- ✔️ Verifikasi & Penolakan Prestasi
- 👑 Role-Based Access Control (Admin, Mahasiswa, Dosen Wali)

## 👤Data Diri
- Nama    : Shendy Tria Amelyana
- NIM     : 434231003
- Kelas   : C1

## 📚 1. Teknologi yang Digunakan
| Layer | Teknologi |
|--------|---------|
| Backend Framework  | Go Fiber   |
| Database Relasional  | PostgreSQL   |
| Database Non-Relasional  | MongoDB   |
| Auth  | JWT   |	
| ORM / DB Access  | pgx + Mongo Driver   |
| Dependency  | Go Modules   |	

## 👥 2. Role & Hak Akses

🔸 Admin
- CRUD User
- Mengubah role user
- Melihat semua mahasiswa
- Melihat semua prestasi
- Melihat semua prestasi mahasiswa tertentu
- Verifikasi / reject prestasi
- Menambah/mengubah advisor mahasiswa

🔸 Dosen Wali
- Melihat data diri sebagai dosen wali
- Melihat daftar mahasiswa bimbingan
- Melihat prestasi mahasiswa bimbingan
- Verifikasi / reject prestasi mahasiswa bimbingan

🔸 Mahasiswa
- Melihat profil diri
- Membuat prestasi (draft)
- Update draft prestasi
- Menghapus draft prestasi
- Submit prestasi ke dosen wali
- Melihat riwayat prestasi
- Menambah lampiran prestasi

## 🔐 3. Autentikasi & JWT

API membutuhkan token JWT pada semua endpoint **kecuali login**.

### 🔑 Login
````md
POST /api/auth/login
````

### ♻️ Refresh Token

```bash
POST /api/auth/refresh
```

### 👤 Get Profile

```bash
GET /api/auth/profile
```

Semua endpoint lain **memerlukan header**:

```
Authorization: Bearer <token>
```

## 📁 4. Struktur Endpoint
### A. Auth
| Method | Endpoint | Deskripsi |
|--------|---------|---------|
| POST  | /api/auth/login   | Login dan generate token   |
| POST  | /api/auth/refresh   | Refresh token   |
| GET  | /api/auth/profile   | Mendapatkan profil user   |
| POST  | /api/auth/logout   |	 Logout   |

### B. Users (Admin Only)
| Method | Endpoint | Deskripsi |
|--------|---------|---------|
| GET  | /api/users   | List semua user   |
| GET  | /api/users/:id   | Detail user   |
| POST  | /api/users   | Create user   |
| PUT  | /api/users/:id   |	 Update user   |
| DELETE  | /api/users/:id   |	 Hapus user   |
| PUT  | /api/users/:id/role   |	 Update peran user   |

### C. Achievements
🔸 Admin
- Melihat semua prestasi

🔸 Mahasiswa
- CRUD draft prestasi
- Submit prestasi
- Upload lampiran

🔸 Dosen Wali
- Melihat prestasi mahasiswa bimbingan
- Verifikasi / reject prestasi

| Method | Endpoint | Deskripsi |
|--------|---------|---------|
| GET  | /api/achievements   | Admin: semua prestasi — Mahasiswa: punya sendiri — Dosen: milik advisees   |
| POST  | /api/achievements   | Mahasiswa membuat prestasi   |
| GET  | /api/achievements/:id   | Detail prestasi   |
| PUT  | /api/achievements/:id   |	 Update draft prestasi   |
| DELETE  | /api/achievements/:id   |	 Delete draft   |
| POST  | /api/achievements/:id/submit   |	 Submit draft   |
| POST  | /api/achievements/:id/verify   |	 Verifikasi   |
| POST  | /api/achievements/:id/reject   |	 Reject prestasi   |
| GET  | /api/achievements/:id/history   |	 Riwayat perubahan status   |
| POST  | /api/achievements/:id/attachments   |	 Tambah lampiran   |

### D. Students
| Method | Endpoint | Deskripsi |
|--------|---------|---------|
| GET  | /api/students   | Admin: list semua mahasiswa   |
| GET  | /api/students/:id   | Admin/dosen   |
| GET  | /api/students/:id/achievements   | Admin/dosen   |
| PUT  | /api/students/:id/advisor   |	 Admin ganti dosen wali   |

### E. Lecturers
| Method | Endpoint | Deskripsi |
|--------|---------|---------|
| GET  | /api/lecturers   | Admin list dosen   |
| GET  | /api/lecturers/:id/advisees   | Admin / Dosen (hanya melihat mahasiswa bimbingannya sendiri)   |

### F. Reports & Analytics 
| Method | Endpoint | Deskripsi |
|--------|---------|---------|
| GET  | /api/reports/statistics   | Admin melihat list prestasi mahasiswa: total draft, submitted, verified, rejected.   |
| GET  | /api/reports/student/:id   | Admin melihat detail prestasi mahasiswa: total draft, submitted, verified, rejected.   |

## 🧠 5. Business Flow 
- Mahasiswa membuat prestasi (status: draft)
- Mahasiswa submit prestasi ke dosen wali
- Dosen wali melihat list prestasi mahasiswa bimbingan
- Dosen wali verify atau reject
- Admin bisa melakukan override untuk semua data
- Sistem menyimpan riwayat perubahan status (history timeline)

## 🗄️ 6. Struktur Database
PostgreSQL (Relasional)
- users
- students
- lecturers
- achievement_references
- permissions
- role_permissions

MongoDB (Dokumen)
- achievements

### 🏗️ 7. Cara Menjalankan Project
1. Clone Repository
````md
git clone <repo-url>
cd prestasi_mhs
````
2. Install Dependencies
````bash
go mod tidy
````
3. Jalankan Server
````bash
go run main.go
````
4. Server Running
````bash
http://localhost:8080
````

## 🧪 8. Postman Collection
Tersedia pengujian untuk:
- Admin
- Mahasiswa
- Dosen Wali

Masing-masing mencakup:
- Auth test
- Users test
- Achievement test
- Student test
- Lecturer test
- Report testt