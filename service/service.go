package service

import (
	"github.com/wahyunurdian26/product-service/repository"
)

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
