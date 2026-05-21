package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/wahyunurdian26/product-service/config"
	"github.com/wahyunurdian26/product-service/endpoint"
	postgresrepo "github.com/wahyunurdian26/product-service/repository/postgres"
	redisrepo "github.com/wahyunurdian26/product-service/repository/redis"
	"github.com/wahyunurdian26/product-service/service"
	httptransport "github.com/wahyunurdian26/product-service/transport/http"
)

func main() {
	cfg := config.LoadConfigs()

	// Setup Database
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Setup Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})
	defer redisClient.Close()

	// Run Goose Migrations
	log.Println("Running database migrations...")
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully.")

	// Repositories
	productRepo := postgresrepo.NewProductRepository(db)
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	// Service
	productService := service.NewProductService(productRepo, cacheRepo)

	// Endpoints
	endpoints := endpoint.MakeEndpoints(productService)

	// Transport / Gateway Handler
	httpHandler := httptransport.NewHTTPHandler(endpoints)

	// Start Server
	port := cfg.HttpPort
	log.Printf("Server listening on port %s", port)
	
	errs := make(chan error)
	go func() {
		errs <- http.ListenAndServe(":"+port, httpHandler)
	}()

	log.Fatalf("Server exit: %s", <-errs)
}
