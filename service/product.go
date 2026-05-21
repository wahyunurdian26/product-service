package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/wahyunurdian26/product-service/model"
	"github.com/wahyunurdian26/product-service/repository"
	"github.com/wahyunurdian26/product-service/util"
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

	_ = s.cacheRepo.InvalidateCache(ctx, util.MakeProductsPattern())

	return product, nil
}

func (s *productService) ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error) {

	cacheKey := util.MakeListProductsKey(req)

	cachedProducts, err := s.cacheRepo.GetProductsCache(ctx, cacheKey)
	if err == nil && cachedProducts != nil {
		return cachedProducts, nil
	} else if err != nil && err != redis.Nil {

	}

	products, err := s.productRepo.ListProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	_ = s.cacheRepo.SetProductsCache(ctx, cacheKey, products, 5*time.Minute)

	return products, nil
}
