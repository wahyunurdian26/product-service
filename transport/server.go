package transport

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/wahyunurdian26/product-service/config"
	"github.com/wahyunurdian26/product-service/gateway/handler"
	"github.com/wahyunurdian26/product-service/gateway/kit"
	postgresrepo "github.com/wahyunurdian26/product-service/repository/postgres"
	redisrepo "github.com/wahyunurdian26/product-service/repository/redis"
	"github.com/wahyunurdian26/product-service/service"
)

type HttpServer struct {
	router http.Handler
	port   string
	close  func()
}

func NewHTTPServer() *HttpServer {
	cfg := config.LoadConfigs()

	// Setup Database
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Setup Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})

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

	// Router setup via Gateway Handler
	router := mux.NewRouter()
	kitRouter := kit.NewRouter(router)
	
	productHandler := handler.NewProductHandler(productService)
	productHandler.InitializeRoute(kitRouter)

	// Apply default HTTP middlewares from kit
	finalHandler := kit.DefaultHTTPHandler(router)

	return &HttpServer{
		router: finalHandler,
		port:   cfg.HttpPort,
		close: func() {
			db.Close()
			redisClient.Close()
		},
	}
}
