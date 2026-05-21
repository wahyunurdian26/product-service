package redisrepo

import (
	"github.com/go-redis/redis/v8"
	"github.com/wahyunurdian26/product-service/repository"
)

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) repository.CacheRepository {
	return &cacheRepository{client: client}
}
