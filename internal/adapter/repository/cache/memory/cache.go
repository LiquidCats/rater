package memory

import (
	"context"
	"time"
)

type CacheRepository struct {
	memory map[string]string
}

func NewCacheRepository() *CacheRepository {
	return &CacheRepository{
		memory: make(map[string]string),
	}
}

func (c *CacheRepository) Set(_ context.Context, key string, value string, _ time.Duration) error {
	c.memory["rater:"+key] = value

	return nil
}

func (c *CacheRepository) Get(_ context.Context, key string) (string, error) {
	return c.memory["rater:"+key], nil
}

func (c *CacheRepository) Del(_ context.Context, key string) error {
	delete(c.memory, "rater:"+key)

	return nil
}

func (c *CacheRepository) Has(_ context.Context, key string) bool {
	_, ok := c.memory["rater:"+key]

	return ok
}
