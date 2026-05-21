package transport

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wahyunurdian26/util/logger"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/wahyunurdian26/product-service/config"
	"github.com/wahyunurdian26/product-service/endpoint"
	postgresrepo "github.com/wahyunurdian26/product-service/repository/postgres"
	redisrepo "github.com/wahyunurdian26/product-service/repository/redis"
	"github.com/wahyunurdian26/product-service/service"

	pb "github.com/wahyunurdian26/product-service/contract/client"

	"github.com/go-kit/kit/transport/grpc"
	sgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

type GrpcServer struct {
	pb.UnimplementedProductServiceServer
	handler       *sgrpc.Server
	healthServer  *health.Server
	healthProduct grpc.Handler
	createProduct grpc.Handler
	listProducts  grpc.Handler
	close         func()
}

func NewGRPCServer() *GrpcServer {
	cfg := config.LoadConfigs()

	// Setup Database
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		logger.Error("failed to open db connection: ", err)
		panic(err)
	}
	if err := db.Ping(); err != nil {
		logger.Error("failed to ping db: ", err)
		panic(err)
	}

	// Setup Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})

	// Run Goose Migrations
	logger.Info("Running database migrations...")
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("failed to set goose dialect: ", err)
		panic(err)
	}
	if err := goose.Up(db, "db/migrations"); err != nil {
		logger.Error("failed to run migrations: ", err)
		panic(err)
	}
	logger.Info("Migrations completed successfully.")

	// Repositories
	productRepo := postgresrepo.NewProductRepository(db)
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	// Service
	productService := service.NewProductService(productRepo, cacheRepo)

	// Endpoints
	endpoints := endpoint.MakeEndpoints(productService)

	return &GrpcServer{
		healthServer: health.NewServer(),
		healthProduct: grpc.NewServer(
			endpoints.HealthProduct,
			decodeHealthCheckRequest,
			encodeHealthCheckResponse,
		),
		createProduct: grpc.NewServer(
			endpoints.CreateProductEndpoint,
			decodeCreateProductRequest,
			encodeCreateProductResponse,
		),
		listProducts: grpc.NewServer(
			endpoints.ListProductsEndpoint,
			decodeListProductsRequest,
			encodeListProductsResponse,
		),
		close: func() {
			db.Close()
			redisClient.Close()
		},
	}
}

func (s *GrpcServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	_, resp, err := s.createProduct.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateProductResponse), nil
}

func (s *GrpcServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	_, resp, err := s.listProducts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ListProductsResponse), nil
}
