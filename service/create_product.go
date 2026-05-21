package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/wahyunurdian26/product-service/model"
	pb "github.com/wahyunurdian26/product-service/contract/client"
	"github.com/wahyunurdian26/product-service/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *productService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := &model.Product{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Price:     req.Price,
		Type:      model.ProductType(req.Type),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	_ = s.cacheRepo.InvalidateCache(ctx, util.MakeProductsPattern())

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Type:      string(product.Type),
			CreatedAt: timestamppb.New(product.CreatedAt),
			UpdatedAt: timestamppb.New(product.UpdatedAt),
		},
	}, nil
}
