package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/wahyunurdian26/product-service/model"
	"github.com/wahyunurdian26/product-service/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error)
	ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
	cacheRepo   repository.CacheRepository
}

func NewProductService(productRepo repository.ProductRepository, cacheRepo repository.CacheRepository) ProductService {
	return &productService{
		productRepo: productRepo,
		cacheRepo:   cacheRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	product := &model.Product{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Price:     req.Price,
		Type:      req.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	// Invalidate cache after creating a new product
	_ = s.cacheRepo.InvalidateCache(ctx, "products_list_*")

	return product, nil
}

func (s *productService) ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error) {
	// Generate cache key based on request parameters
	reqBytes, _ := json.Marshal(req)
	hash := sha256.Sum256(reqBytes)
	cacheKey := "products_list_" + hex.EncodeToString(hash[:])

	// Try to get from cache
	cachedProducts, err := s.cacheRepo.GetProductsCache(ctx, cacheKey)
	if err == nil && cachedProducts != nil {
		return cachedProducts, nil
	} else if err != nil && err != redis.Nil {
		// Log error if it's not simply a miss
		// Here we ignore for simplicity
	}

	// Fetch from DB
	products, err := s.productRepo.ListProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Save to cache
	_ = s.cacheRepo.SetProductsCache(ctx, cacheKey, products, 5*time.Minute)

	return products, nil
}
