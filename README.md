# 🏛️ Cultour Backend: Menjelajahi Warisan Budaya Indonesia 🇮🇩

## 🌟 Deskripsi Proyek

Cultour adalah platform inovatif yang bertujuan untuk melestarikan dan mempromosikan kekayaan budaya Indonesia melalui teknologi modern. Backend ini dibangun dengan Go (Golang), dirancang untuk memberikan pengalaman yang kaya dan mendalam tentang warisan budaya nusantara.

### 🚀 Fitur Utama

- **Manajemen Lokasi Budaya**: Dokumentasi detail lokasi bersejarah
- **Cerita Lokal**: Perpustakaan digital kisah-kisah tradisional
- **Acara Budaya**: Informasi mendalam tentang event dan festival
- **Autentikasi Aman**: Sistem keamanan berbasis Supabase
- **Dokumentasi API Komprehensif**: Swagger untuk kemudahan integrasi

## 🛠️ Teknologi Utama

- **Bahasa**: Go (Golang) 1.20+
- **Database**: Supabase (PostgreSQL)
- **Autentikasi**: Supabase Auth
- **Framework Web**: Gin
- **Logging**: Structured logging dengan `slog`
- **Dokumentasi API**: Swagger

## 📦 Struktur Proyek

```
cultour-backend/
├── cmd/                # Titik masuk aplikasi
├── configs/            # Konfigurasi aplikasi
├── internal/           # Logika bisnis internal
│   ├── cultural/       # Modul budaya
│   ├── location/       # Manajemen lokasi
│   ├── place/          # Informasi tempat
│   └── supabase/       # Integrasi Supabase
├── pkg/                # Paket utilitas yang dapat digunakan ulang
│   ├── logger/         # Sistem logging
│   └── response/       # Utilitas respons API
└── docs/               # Dokumentasi Swagger
```

## 🔧 Prasyarat

- Go 1.20+
- Supabase Account
- PostgreSQL

## 🚀 Instalasi & Pengaturan

1. Clone repositori
```bash
git clone https://github.com/holycann/cultour-backend.git
cd cultour-backend
```

2. Instal dependensi
```bash
go mod tidy
```

3. Salin dan edit konfigurasi
```bash
cp .env.example .env
# Edit .env dengan kredensial Anda
```

4. Jalankan migrasi database
```bash
go run cmd/migrate/main.go
```

5. Jalankan server
```bash
go run cmd/main.go
# Atau gunakan air untuk development
air
```

## 📘 Dokumentasi API

Akses dokumentasi Swagger di:
`http://localhost:8181/docs/index.html`

## 🔐 Autentikasi

Cultour menggunakan Supabase untuk autentikasi. Setiap endpoint yang memerlukan otentikasi membutuhkan Bearer Token.

## 🧪 Testing

Jalankan test:
```bash
go test ./...
```

## 🤝 Kontribusi

1. Fork repositori
2. Buat branch fitur (`git checkout -b fitur/AturanBaru`)
3. Commit perubahan (`git commit -m 'Tambah fitur baru'`)
4. Push ke branch (`git push origin fitur/AturanBaru`)
5. Buka Pull Request

## 📊 Statistik Proyek

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Supabase](https://img.shields.io/badge/Supabase-3ECF8E?style=for-the-badge&logo=supabase&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)

## 🏆 Tim Pengembang

- **Holycann Team** - Pencipta platform pelestarian budaya

## 📜 Lisensi

Proyek ini dilisensikan di bawah MIT License.

---

🌍 **Cultour: Melestarikan Warisan, Menginspirasi Generasi** 