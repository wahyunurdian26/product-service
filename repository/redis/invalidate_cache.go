package redisrepo

import (
	"context"
)

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
