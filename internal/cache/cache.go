package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{
		rdb: rdb,
	}
}

func (c *Cache) Set(userId, key, value string) {
	c.rdb.Set(context.Background(), userId+key, value, 0)
}

func (c *Cache) Get(userId, key string) string {
	r, _ := c.rdb.Get(context.Background(), userId+key).Result()
	return r
}
