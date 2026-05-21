package service

import (
	"context"
	"strings"
	"testing"

	"github.com/wahyunurdian26/product-service/mock"
	"github.com/wahyunurdian26/product-service/model"
	pb "github.com/wahyunurdian26/product-service/contract/client"
	"go.uber.org/mock/gomock"
)

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		req            *pb.CreateProductRequest
		mockSetup      func(repo *mock.MockProductRepository, cache *mock.MockCacheRepository)
		expectedName   string
		expectedErrMsg string
	}{
		{
			name: "Success - Create Product",
			req: &pb.CreateProductRequest{
				Name:  "Apel",
				Price: 15000,
				Type:  string(model.TypeBuah),
			},
			mockSetup: func(repo *mock.MockProductRepository, cache *mock.MockCacheRepository) {
				repo.EXPECT().CreateProduct(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, product *model.Product) error {
					if product.Name != "Apel" {
						t.Errorf("expected product name Apel, got %s", product.Name)
					}
					return nil
				})
				cache.EXPECT().InvalidateCache(ctx, gomock.Any()).Return(nil)
			},
			expectedName: "Apel",
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
			resp, err := svc.CreateProduct(ctx, tt.req)

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
				if resp == nil || resp.Product.Name != tt.expectedName {
					t.Errorf("expected product name %s, got %+v", tt.expectedName, resp)
				}
			}
		})
	}
}
