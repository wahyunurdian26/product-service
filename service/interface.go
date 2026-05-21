package service

import (
	"context"

	pb "github.com/wahyunurdian26/product-service/contract/client"
)

type ProductService interface {
	HealthProduct(ctx context.Context, req *pb.HealthProductRequest) (*pb.HealthProductResponse, error)
	CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error)
	ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
}
