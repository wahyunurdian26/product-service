package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wahyunurdian26/product-service/model"
)

type CacheRepository interface {
	SetProductsCache(ctx context.Context, key string, products []model.Product, ttl time.Duration) error
	GetProductsCache(ctx context.Context, key string) ([]model.Product, error)
	InvalidateCache(ctx context.Context, pattern string) error
}

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) CacheRepository {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) SetProductsCache(ctx context.Context, key string, products []model.Product, ttl time.Duration) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *cacheRepository) GetProductsCache(ctx context.Context, key string) ([]model.Product, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var products []model.Product
	if err := json.Unmarshal([]byte(val), &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *cacheRepository) InvalidateCache(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			r.client.Del(ctx, key)
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
