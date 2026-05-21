package repository

import (
	"context"
	"time"

	"github.com/wahyunurdian26/product-service/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error)
}

type CacheRepository interface {
	SetProductsCache(ctx context.Context, key string, products []model.Product, ttl time.Duration) error
	GetProductsCache(ctx context.Context, key string) ([]model.Product, error)
	InvalidateCache(ctx context.Context, pattern string) error
}
