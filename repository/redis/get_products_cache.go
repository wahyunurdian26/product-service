package redisrepo

import (
	"context"
	"encoding/json"

	"github.com/wahyunurdian26/product-service/model"
)

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
