package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/LiquidCats/rater/configs"
	"github.com/bytedance/sonic"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	baseKey string
	client  *cache.Cache
}

func New(cfg configs.RedisConfig) *Repository {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: string(cfg.Password),
		DB:       cfg.DB,
		Protocol: cfg.Protocol, // nolint:mnd
	})

	return &Repository{
		baseKey: configs.AppName,
		client: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(10, time.Second*5), //nolint:mnd

			Marshal:   sonic.Marshal,
			Unmarshal: sonic.Unmarshal,
		}),
	}
}

func (c *Repository) key(k fmt.Stringer) string {
	return strings.ToLower(fmt.Sprintf("%s:%s", c.baseKey, k))
}
