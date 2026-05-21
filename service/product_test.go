package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wahyunurdian26/product-service/model"
)

type mockProductRepo struct {
	products []model.Product
	err      error
}

func (m *mockProductRepo) CreateProduct(ctx context.Context, product *model.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products = append(m.products, *product)
	return nil
}

func (m *mockProductRepo) ListProducts(ctx context.Context, req model.ListProductRequest) ([]model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.products, nil
}

type mockCacheRepo struct {
	data map[string][]model.Product
	err  error
}

func (m *mockCacheRepo) SetProductsCache(ctx context.Context, key string, products []model.Product, ttl time.Duration) error {
	if m.data == nil {
		m.data = make(map[string][]model.Product)
	}
	m.data[key] = products
	return m.err
}

func (m *mockCacheRepo) GetProductsCache(ctx context.Context, key string) ([]model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	if p, ok := m.data[key]; ok {
		return p, nil
	}
	return nil, errors.New("redis: nil")
}

func (m *mockCacheRepo) InvalidateCache(ctx context.Context, pattern string) error {
	m.data = make(map[string][]model.Product)
	return m.err
}

func TestCreateProduct_Success(t *testing.T) {
	repo := &mockProductRepo{}
	cache := &mockCacheRepo{}
	svc := NewProductService(repo, cache)

	req := &model.CreateProductRequest{
		Name:  "Apel",
		Price: 15000,
		Type:  model.TypeBuah,
	}

	product, err := svc.CreateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.Name != "Apel" {
		t.Errorf("expected Apel, got %s", product.Name)
	}
}


