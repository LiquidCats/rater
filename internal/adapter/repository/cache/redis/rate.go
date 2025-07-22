package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/rotisserie/eris"
)

type RateKey struct {
	From string
	To   string
}

func (p RateKey) String() string {
	return strings.ToLower(
		fmt.Sprintf(
			"rate:from:%s:to:%s",
			p.From,
			p.To,
		),
	)
}

func (c *Repository) GetRate(ctx context.Context, key RateKey) (Rate, error) {
	var value Rate

	err := c.client.Get(ctx, c.key(key), &value)
	if err != nil {
		return value, eris.Wrap(err, "get value from cache")
	}

	return value, nil
}

func (c *Repository) PutRate(ctx context.Context, key RateKey, value Rate) error {
	err := c.client.Set(&cache.Item{
		Ctx:   ctx,
		Key:   c.key(key),
		Value: value,
		TTL:   5 * time.Second, //nolint:mnd
	})
	if err != nil {
		return eris.Wrap(err, "put value to cache")
	}

	return nil
}
