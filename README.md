# Super Indo Product API Service

This is a microservice developed for the Super Indo coding test. It is built using Golang and follows a modern `grpc-gateway` based architecture, providing both gRPC and REST HTTP interfaces seamlessly on the same port.

## Features
1. Add new products.
2. List all products.
3. Search products by Name and ID.
4. Filter products by Type (Sayuran, Protein, Buah, Snack).
5. Sort products by Date, Price, or Name.
6. Robust caching using Redis (Cache-aside pattern).
7. Postgres SQL Database with Goose Migrations.

## Tech Stack
- **Language**: Golang 1.20
- **Database**: PostgreSQL 14
- **Cache**: Redis 7
- **Routing/Transport**: gRPC + `grpc-gateway` (REST reverse-proxy multiplexing)
- **Containerization**: Docker & Docker Compose

## Getting Started

### Prerequisites
- Docker and Docker Compose installed
- Go 1.20+ (if running locally without Docker)
- `make` and `protoc` (if you want to re-generate the proto files)

### Running with Docker Compose
This is the easiest way to get the service running with its database and cache.

1. Start the services:
   ```bash
   docker-compose up -d --build
   ```
2. The service will be available at `http://localhost:6668` (serving both gRPC and REST).
3. The database will automatically be created and seeded with migrations upon startup.

### Running Locally (Without Docker)
4. Run the Go app:
   ```bash
   go run main.go
   ```

### REST Gateway (Port 7070)

The REST API is exposed via gRPC-Gateway.

#### 1. Health Check
```bash
curl -X GET http://localhost:7070/health
```

#### 2. Create Product
```bash
curl -X POST http://localhost:7070/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Susu UHT", "price":18000, "type":"Protein"}'
```

#### 3. List Products (With Filter & Sort)
```bash
curl -X GET "http://localhost:7070/api/v1/products?search=Susu&type=Protein&sort_by=price&order=ASC"
```

**Query Parameters yang Tersedia:**
- `search` (string): Mencari produk berdasarkan nama (menggunakan `ILIKE`).
- `type` (string): Filter berdasarkan kategori (contoh: `Protein`, `Sayuran`, `Buah`, `Snack`).
- `sort_by` (string): Kolom yang digunakan untuk mengurutkan (contoh: `price`, `created_at`, `name`).
- `order` (string): Arah urutan, `ASC` (naik) atau `DESC` (turun).

## API Documentation

### 1. Add Product
- **Endpoint**: `POST /api/v1/products`
- **Body**:
  ```json
  {
      "name": "Susu UHT",
      "price": 18000,
      "type": "Protein"
  }
  ```

### 2. List, Search, Filter, Sort Products
- **Endpoint**: `GET /api/v1/products`
- **Query Parameters**:
  - `search` (Optional): Search by product name or id.
  - `type` (Optional): Filter by type (Sayuran, Protein, Buah, Snack).
  - `sort_by` (Optional): Sort by field (`created_at`, `price`, `name`). Default is `created_at`.
  - `order` (Optional): Sort order (`ASC`, `DESC`). Default is `DESC`.
- **Example**: `GET /api/v1/products?search=Apel&type=Buah&sort_by=price&order=ASC`

### Swagger Documentation
The Swagger specification is automatically generated and can be found in `contract/swagger/product.swagger.json`.

## Architecture Overview
- `config`: Handles environment variables.
- `contract`: Contains the protobuf definition (`product.proto`), the gRPC-Gateway REST routing configuration (`product.yaml`), the generated clients, and Swagger specs.
- `model`: Entities and structs.
- `repository`: DB (Postgres) and Cache (Redis) abstractions.
- `service`: Core business logic (validation, caching logic).
- `transport`: gRPC Server logic and `grpc-gateway` initialization multiplexer.
