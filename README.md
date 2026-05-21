# Super Indo Product API Service

This is a microservice developed for the Super Indo coding test. It is built using Golang and follows the `go-kit` based architecture common in corporate payment microservices. The HTTP Gateway is integrated directly into the service using `gorilla/mux`.

## Features
1. Add new products.
2. List all products.
3. Search products by Name and ID.
4. Filter products by Type (Sayuran, Protein, Buah, Snack).
5. Sort products by Date, Price, or Name.
6. Caching using Redis.
7. Postgres SQL Database with Migrations and Seeders.

## Tech Stack
- **Language**: Golang 1.20
- **Database**: PostgreSQL 14
- **Cache**: Redis 7
- **Routing/Transport**: `gorilla/mux` + `go-kit`
- **Containerization**: Docker & Docker Compose

## Getting Started

### Prerequisites
- Docker and Docker Compose installed
- Go 1.20+ (if running locally without Docker)

### Running with Docker Compose
This is the easiest way to get the service running with its database and cache.

1. Start the services:
   ```bash
   docker-compose up -d --build
   ```
2. The service will be available at `http://localhost:8080`.
3. The database will automatically be created and seeded with dummy data upon startup.

### Running Locally (Without Docker)
1. Start a local Postgres and Redis instance.
2. Copy `.env.example` to `.env` and adjust the credentials.
3. Import the SQL files located in `db/migrations/` and `db/seeds/` to your Postgres DB manually.
4. Run the Go app:
   ```bash
   go run main.go
   ```

## API Documentation

### 1. Add Product
- **Endpoint**: `POST /product`
- **Body**:
  ```json
  {
      "name": "Susu UHT",
      "price": 18000,
      "type": "Protein"
  }
  ```

### 2. List, Search, Filter, Sort Products
- **Endpoint**: `GET /product`
- **Query Parameters**:
  - `search` (Optional): Search by product name or id.
  - `type` (Optional): Filter by type (Sayuran, Protein, Buah, Snack).
  - `sort_by` (Optional): Sort by field (`created_at`, `price`, `name`). Default `created_at`.
  - `order` (Optional): Sort order (`ASC`, `DESC`). Default `DESC`.
- **Example**: `GET /product?search=Apel&type=Buah&sort_by=price&order=ASC`

## Architecture Overview
- `config`: Handles environment variables.
- `model`: Entities and structs.
- `repository`: DB (Postgres) and Cache (Redis) abstractions.
- `service`: Core business logic (validation, caching logic).
- `endpoint`: `go-kit` endpoints.
- `transport/http`: `gorilla/mux` based HTTP Gateway handling JSON payload decode/encode.
