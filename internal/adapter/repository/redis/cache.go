package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"rater/configs"
	"strings"
	"time"
)

type CacheRepository struct {
	baseKey string
	client  *redis.Client
}

func NewCacheRepository(cfg configs.RedisConfig, baseKey string) *CacheRepository {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		DB:   cfg.DB,
	})

	return &CacheRepository{
		baseKey: baseKey,
		client:  client,
	}
}

func (c *CacheRepository) Has(ctx context.Context, key string) bool {
	cmd := c.client.Exists(ctx, c.makeKey(key))

	if err := cmd.Err(); nil != err {
		return false
	}

	return cmd.Val() > 0
}

func (c *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	cmd := c.client.Get(ctx, c.makeKey(key))

	if err := cmd.Err(); nil != err {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", err
	}

	return cmd.Val(), nil
}

func (c *CacheRepository) Set(ctx context.Context, key string, value string, expire time.Duration) error {
	cmd := c.client.Set(ctx, c.makeKey(key), value, expire)
	if err := cmd.Err(); nil != err {
		return err
	}

	return nil
}

func (c *CacheRepository) Del(ctx context.Context, key string) error {
	cmd := c.client.Del(ctx, c.makeKey(key))
	if err := cmd.Err(); nil != err && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (c *CacheRepository) makeKey(key string) string {
	return strings.ToLower(fmt.Sprintf("%s:%s", c.baseKey, key))
}
