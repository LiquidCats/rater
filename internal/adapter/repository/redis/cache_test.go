package redis

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/rand"
	"rater/configs"
	"strconv"
	"testing"
	"time"
)

func TestCacheRepository(t *testing.T) {
	ctx := context.Background()
	cfg := configs.RedisConfig{
		Address: "cache:6379",
		DB:      0,
	}

	repo := NewCacheRepository(cfg, "tests")

	key := "cache:repository"
	value := strconv.FormatInt(rand.Int63(), 10)

	require.False(t, repo.Has(ctx, key))

	err := repo.Set(ctx, key, value, 1*time.Second)
	require.NoError(t, err)

	require.True(t, repo.Has(ctx, key))

	cached, err := repo.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, value, cached)

	err = repo.Del(ctx, key)
	require.NoError(t, err)

	require.False(t, repo.Has(ctx, key))
}
