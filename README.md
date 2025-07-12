
# Students REST API (Go + SQLite)

This is a simple and clean RESTful API written in Go for managing student records. It supports full CRUD operations, uses SQLite for storage.

## 📦 Features

- ✅ Create, Read, Update, Delete (CRUD) operations on students
- ✅ SQLite database integration
- ✅ Graceful server shutdown and structured logging
- ✅ Minimal dependencies, built using Go standard library

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/atharv3221/rest-api.git
cd rest-api
```
### 2. Install Dependencies

```bash
go mod tidy
```
### 3. Starting the server

```bash
go run cmd/rest-api/main.go -config config/local.yaml
```

