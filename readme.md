# 📝 Todo-App-API

Project backend REST API sederhana menggunakan **Golang (Fiber + GORM)** dengan dukungan **PostgreSQL** sebagai database utama dan **SQLite (in-memory)** untuk unit testing.

---

## 🚀 Fitur Utama

- CRUD Todo (Create, Read, Update, Delete)
- Validasi input menggunakan `go-playground/validator`
- Error handling dengan middleware Fiber
- Unit test lengkap untuk Controller, Service, Repository, Helper, dan Exception
- Test coverage 72% menggunakan `testify` dan `mock`

---

## 🧩 Teknologi yang Digunakan

- [Go 1.24](https://go.dev/)
- [Fiber v2](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Testify](https://github.com/stretchr/testify)
- SQLite (untuk testing)
- PostgreSQL (untuk production)

---

## 🛠️ Cara Menjalankan Project

### 1. Clone repository

```bash
git clone https://github.com/username/todo-app-api.git
cd todo-app-api
```

### 2. Jalankan module

```bash
go mod tidy
```

---

### 3. Siapkan file `.env`

Buat file `.env` di root project:

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=todo_db
APP_PORT=8080
```

---

### 4. Jalankan aplikasi

```bash
go run main.go
```

---

## 🧪 Menjalankan Unit Test

```bash
go test ./... -v -coverpkg=./...
```

📊 Hasil coverage: ~72%

---

## 🧱 Struktur Folder

```
todo-app-api/
├── controller/ # Fiber controllers (request handlers)
├── service/ # Business logic layer
├── repository/ # Database access (GORM)
├── helper/ # Utility & response helper
├── exception/ # Error handling
├── test/ # Unit tests
├── main.go # Entry point
└── .env.example # Contoh environment
```

---

## 🧑‍💻 Author

**Dhahika Rahmadani**  
Backend Developer • Go Enthusiast  
📧 [dhahikardani@gmail.com](mailto:dhahikardani@gmail.com)
