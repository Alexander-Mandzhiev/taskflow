package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func newTestClient(t *testing.T) (RedisClient, *miniredis.Miniredis, func()) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	client := NewClient(rdb, &logger.NoopLogger{}, 2*time.Second, "test.cache")
	return client, mr, func() { mr.Close(); _ = rdb.Close() }
}

func TestClient_Set_Get(t *testing.T) {
	client, _, cleanup := newTestClient(t)
	defer cleanup()
	ctx := context.Background()

	err := client.Set(ctx, "k1", "v1", time.Minute)
	require.NoError(t, err)
	got, err := client.Get(ctx, "k1")
	require.NoError(t, err)
	assert.Equal(t, []byte("v1"), got)
	// промах
	got2, err := client.Get(ctx, "missing")
	require.NoError(t, err)
	assert.Nil(t, got2)
}

func TestClient_Del(t *testing.T) {
	client, _, cleanup := newTestClient(t)
	defer cleanup()
	ctx := context.Background()
	_ = client.Set(ctx, "k", "v", time.Minute)
	err := client.Del(ctx, "k")
	require.NoError(t, err)
	got, _ := client.Get(ctx, "k")
	assert.Nil(t, got)
}

func TestClient_Ping(t *testing.T) {
	client, _, cleanup := newTestClient(t)
	defer cleanup()
	err := client.Ping(context.Background())
	require.NoError(t, err)
}

func TestClient_MGet(t *testing.T) {
	client, _, cleanup := newTestClient(t)
	defer cleanup()
	ctx := context.Background()
	_ = client.Set(ctx, "a", "1", time.Minute)
	_ = client.Set(ctx, "b", "2", time.Minute)

	// пустой список
	out, err := client.MGet(ctx)
	require.NoError(t, err)
	assert.Nil(t, out)

	out, err = client.MGet(ctx, "a", "b", "c")
	require.NoError(t, err)
	require.Len(t, out, 3)
	assert.Equal(t, []byte("1"), out[0])
	assert.Equal(t, []byte("2"), out[1])
	assert.Nil(t, out[2])
}

func TestClient_HSet_HGet_HGetAll_HDel(t *testing.T) {
	client, _, cleanup := newTestClient(t)
	defer cleanup()
	ctx := context.Background()

	err := client.HSet(ctx, "h1", map[string]interface{}{"f1": "v1", "f2": "v2"})
	require.NoError(t, err)
	val, err := client.HGet(ctx, "h1", "f1")
	require.NoError(t, err)
	assert.Equal(t, "v1", val)
	val, err = client.HGet(ctx, "h1", "missing")
	require.NoError(t, err)
	assert.Equal(t, "", val)

	all, err := client.HGetAll(ctx, "h1")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"f1": "v1", "f2": "v2"}, all)

	err = client.HDel(ctx, "h1", "f1")
	require.NoError(t, err)
	val, _ = client.HGet(ctx, "h1", "f1")
	assert.Equal(t, "", val)
}

func TestClient_Expire(t *testing.T) {
	client, mr, cleanup := newTestClient(t)
	defer cleanup()
	ctx := context.Background()
	_ = client.Set(ctx, "k", "v", time.Minute)
	err := client.Expire(ctx, "k", time.Second)
	require.NoError(t, err)
	mr.FastForward(2 * time.Second)
	got, _ := client.Get(ctx, "k")
	assert.Nil(t, got)
}
