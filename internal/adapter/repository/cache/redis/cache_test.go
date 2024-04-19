package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type MockRedis struct {
	mock.Mock
}

func (r *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := r.Mock.Called(ctx, key, value, expiration)

	return args.Get(0).(*redis.StatusCmd)
}

func (r *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	args := r.Mock.Called(ctx, key)

	return args.Get(0).(*redis.StringCmd)
}

func (r *MockRedis) Exists(ctx context.Context, key ...string) *redis.IntCmd {
	args := r.Mock.Called(ctx, key)

	return args.Get(0).(*redis.IntCmd)
}

func (r *MockRedis) Del(ctx context.Context, key ...string) *redis.IntCmd {
	args := r.Mock.Called(ctx, key)

	return args.Get(0).(*redis.IntCmd)
}

var key = "test"

func TestCacheRepository_Set(t *testing.T) {
	ctx := context.Background()

	client := &MockRedis{}

	repo := CacheRepository{
		baseKey: "rater",
		client:  client,
	}

	cmd := redis.NewStatusCmd(ctx)

	client.On("Set", ctx, "rater:"+key, "test-value", 150*time.Second).Return(cmd).Once()

	err := repo.Set(ctx, key, "test-value", 150*time.Second)

	require.NoError(t, err)
}

func TestCacheRepository_Has(t *testing.T) {
	ctx := context.Background()

	client := &MockRedis{}

	repo := CacheRepository{
		baseKey: "rater",
		client:  client,
	}

	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(1)

	client.On("Exists", ctx, []string{"rater:" + key}).Return(cmd).Once()

	val := repo.Has(ctx, key)

	require.True(t, val)
}

func TestCacheRepository_Del(t *testing.T) {
	ctx := context.Background()

	client := &MockRedis{}

	repo := CacheRepository{
		baseKey: "rater",
		client:  client,
	}

	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(1)

	client.On("Del", ctx, []string{"rater:" + key}).Return(cmd).Once()

	err := repo.Del(ctx, key)

	require.NoError(t, err)
}

func TestCacheRepository_Get(t *testing.T) {
	ctx := context.Background()

	client := &MockRedis{}

	repo := CacheRepository{
		baseKey: "rater",
		client:  client,
	}

	cmd := redis.NewStringCmd(ctx)
	cmd.SetVal("test-value")

	client.On("Get", ctx, "rater:"+key).Return(cmd).Once()

	result, err := repo.Get(ctx, key)

	require.NoError(t, err)

	require.Equal(t, "test-value", result)

	//key := "cache:repository"
	//value := strconv.FormatInt(rand.Int63(), 10)
	//
	//require.False(t, repo.Has(ctx, key))
	//
	//err := repo.Set(ctx, key, value, 1*time.Second)
	//require.NoError(t, err)
	//
	//require.True(t, repo.Has(ctx, key))
	//
	//cached, err := repo.Get(ctx, key)
	//require.NoError(t, err)
	//require.Equal(t, value, cached)
	//
	//err = repo.Del(ctx, key)
	//require.NoError(t, err)
	//
	//require.False(t, repo.Has(ctx, key))
}
