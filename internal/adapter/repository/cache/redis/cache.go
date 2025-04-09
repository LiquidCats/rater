package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	baseKey string
	client  *redis.Client
}

func NewCacheRepository(cfg configs.RedisConfig, baseKey string) (*CacheRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: string(cfg.Password),
		DB:       cfg.DB,
		Protocol: 3, // nolint:mnd
	})

	cmd := client.Ping(context.Background())

	return &CacheRepository{
		baseKey: baseKey,
		client:  client,
	}, cmd.Err()
}

func (c *CacheRepository) GetRate(ctx context.Context, pair entity.Pair) (*entity.Rate, error) {
	b, err := c.client.Get(ctx, c.makeRateKey(pair)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "get rate from cache")
	}

	var rate entity.Rate
	decoder := json.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&rate); err != nil {
		return nil, errors.Wrap(err, "decode rate from cache")
	}

	return &rate, nil
}

func (c *CacheRepository) PutRate(ctx context.Context, rate entity.Rate, expire time.Duration) error {
	buff := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buff).Encode(rate); err != nil {
		return errors.Wrap(err, "encode rate to cache")
	}

	cmd := c.client.Set(ctx, c.makeRateKey(rate.Pair), buff.Bytes(), expire)
	if err := cmd.Err(); err != nil {
		return errors.Wrap(err, "save rate to cache")
	}

	return nil
}

func (c *CacheRepository) makeRateKey(pair entity.Pair) string {
	return strings.ToLower(
		fmt.Sprintf(
			"%s:rate:from:%s:to:%s",
			c.baseKey,
			pair.From,
			pair.To,
		),
	)
}
