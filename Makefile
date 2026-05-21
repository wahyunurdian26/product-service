.PHONY: run build test docker-up docker-down clean

# Variabel
APP_NAME=product-service

# Menjalankan service secara lokal
run:
	go run main.go

# Build aplikasi ke executable
build:
	go build -o $(APP_NAME) .

# Menjalankan unit test
test:
	go test ./... -v

# Menjalankan menggunakan docker compose (database, redis, app)
docker-up:
	docker-compose up -d --build

# Mematikan service docker
docker-down:
	docker-compose down

# Membersihkan file build
clean:
	rm -f $(APP_NAME)
