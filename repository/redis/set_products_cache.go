package redisrepo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wahyunurdian26/product-service/model"
)

func (r *cacheRepository) SetProductsCache(ctx context.Context, key string, products []model.Product, ttl time.Duration) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}
