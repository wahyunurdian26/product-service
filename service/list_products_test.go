package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wahyunurdian26/product-service/mock"
	"github.com/wahyunurdian26/product-service/model"
	pb "github.com/wahyunurdian26/product-service/contract/client"
	"go.uber.org/mock/gomock"
)

func TestListProducts(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		req            *pb.ListProductsRequest
		mockSetup      func(repo *mock.MockProductRepository, cache *mock.MockCacheRepository)
		expectedCount  int
		expectedErrMsg string
	}{
		{
			name: "Success - Cache Hit",
			req: &pb.ListProductsRequest{
				Search: "Apel",
			},
			mockSetup: func(repo *mock.MockProductRepository, cache *mock.MockCacheRepository) {
				cachedData := []model.Product{
					{Name: "Apel Fuji", Price: 15000, Type: model.TypeBuah},
				}
				cache.EXPECT().GetProductsCache(ctx, gomock.Any()).Return(cachedData, nil)
			},
			expectedCount: 1,
		},
		{
			name: "Success - Cache Miss (Fetch DB)",
			req: &pb.ListProductsRequest{
				Search: "Apel",
			},
			mockSetup: func(repo *mock.MockProductRepository, cache *mock.MockCacheRepository) {
				cache.EXPECT().GetProductsCache(ctx, gomock.Any()).Return(nil, redis.Nil)
				
				dbData := []model.Product{
					{Name: "Apel Malang", Price: 10000, Type: model.TypeBuah},
				}
				repo.EXPECT().ListProducts(ctx, gomock.Any()).Return(dbData, nil)
				
				cache.EXPECT().SetProductsCache(ctx, gomock.Any(), dbData, 5*time.Minute).Return(nil)
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockProductRepository(ctrl)
			mockCache := mock.NewMockCacheRepository(ctrl)

			tt.mockSetup(mockRepo, mockCache)

			svc := NewProductService(mockRepo, mockCache)
			resp, err := svc.ListProducts(ctx, tt.req)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got none", tt.expectedErrMsg)
				} else if !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("expected error containing %q, got %q", tt.expectedErrMsg, err.Error())
				}
				if resp != nil {
					t.Errorf("expected response to be nil, got %+v", resp)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if resp == nil || len(resp.Products) != tt.expectedCount {
					t.Errorf("expected %d products, got %d", tt.expectedCount, len(resp.Products))
				}
			}
		})
	}
}
