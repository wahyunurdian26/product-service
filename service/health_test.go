package service

import (
	"context"
	"testing"

	"github.com/wahyunurdian26/product-service/mock"
	pb "github.com/wahyunurdian26/product-service/contract/client"
	"go.uber.org/mock/gomock"
)

func TestHealthProduct(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProductRepository(ctrl)
	mockCache := mock.NewMockCacheRepository(ctrl)

	svc := NewProductService(mockRepo, mockCache)

	req := &pb.HealthProductRequest{}
	resp, err := svc.HealthProduct(ctx, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Status != "SERVING" {
		t.Errorf("expected status SERVING, got %s", resp.Status)
	}
}
