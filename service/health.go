package service

import (
	"context"

	pb "github.com/wahyunurdian26/product-service/contract/client"
)

func (s *productService) HealthProduct(ctx context.Context, req *pb.HealthProductRequest) (*pb.HealthProductResponse, error) {
	return &pb.HealthProductResponse{Status: "SERVING"}, nil
}
