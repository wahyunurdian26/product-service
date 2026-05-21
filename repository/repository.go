package repository

import (
	"context"

	"github.com/wahyunurdian26/product-service/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error)
}
