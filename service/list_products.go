package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wahyunurdian26/product-service/model"
	pb "github.com/wahyunurdian26/product-service/contract/client"
	"github.com/wahyunurdian26/product-service/util"
	"github.com/wahyunurdian26/util/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *productService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	logger.Info("ListProducts request - search: ", req.Search, " type: ", req.Type, " sort_by: ", req.SortBy, " order: ", req.Order)

	modelReq := model.ListProductRequest{
		Search: req.Search,
		Type:   req.Type,
		SortBy: req.SortBy,
		Order:  req.Order,
	}

	cacheKey := util.MakeListProductsKey(modelReq)

	cachedProducts, err := s.cacheRepo.GetProductsCache(ctx, cacheKey)
	if err == nil && cachedProducts != nil {
		logger.Info("ListProducts cache HIT - key: ", cacheKey)
		return mapProductsToPb(cachedProducts), nil
	} else if err != nil && err != redis.Nil {
		logger.Warn("ListProducts cache error (fallback to DB): ", err)
	}

	logger.Info("ListProducts cache MISS - fetching from DB")
	products, err := s.productRepo.ListProducts(ctx, modelReq)
	if err != nil {
		logger.Error("ListProducts DB error: ", err)
		return nil, err
	}

	_ = s.cacheRepo.SetProductsCache(ctx, cacheKey, products, 5*time.Minute)

	return mapProductsToPb(products), nil
}

func mapProductsToPb(products []model.Product) *pb.ListProductsResponse {
	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:        p.ID,
			Name:      p.Name,
			Price:     p.Price,
			Type:      string(p.Type),
			CreatedAt: timestamppb.New(p.CreatedAt),
			UpdatedAt: timestamppb.New(p.UpdatedAt),
		})
	}
	return &pb.ListProductsResponse{
		Products: pbProducts,
	}
}
