.PHONY: run build test docker-up docker-down clean mock

# Variabel
APP_NAME=product-service

# Menjalankan service 
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

# Membuat file mock untuk keperluan unit testing
mock:
	mockgen -source=repository/interface.go -destination=mock/repository_mock.go -package=mock

# Generate Protobuf, gRPC, grpc-gateway, dan swagger
protogen:
	@protoc -I contract \
		--go_out=contract/client --go_opt=paths=source_relative \
		--go-grpc_out=contract/client --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=contract/client \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=grpc_api_configuration=contract/product.yaml \
		--openapiv2_out=contract/swagger \
		--openapiv2_opt=logtostderr=true \
		--openapiv2_opt=grpc_api_configuration=contract/product.yaml \
		contract/product.proto

