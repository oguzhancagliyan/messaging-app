package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type redisCache struct {
	client *redis.Client
}

type Cache interface {
	SaveMessageID(messageID string, sentAt time.Time) error
}

func NewRedisCache(addr string) *redisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &redisCache{client: rdb}
}

func (r *redisCache) SaveMessageID(messageID string, sentAt time.Time) error {
	key := fmt.Sprintf("message:%s", messageID)
	value := sentAt.Format(time.RFC3339)
	return r.client.Set(ctx, key, value, 0).Err()
}
