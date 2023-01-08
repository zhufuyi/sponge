package cache

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

type redisUser struct {
	ID   uint64
	Name string
}

func newRedisCache() *gotest.Cache {
	record1 := &redisUser{
		ID:   1,
		Name: "foo",
	}
	record2 := &redisUser{
		ID:   2,
		Name: "bar",
	}

	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	cachePrefix := ""
	c.ICache = NewRedisCache(c.RedisClient, cachePrefix, encoding.JSONEncoding{}, func() interface{} {
		return &redisUser{}
	})

	return c
}

func TestRedisCache(t *testing.T) {
	c := newRedisCache()
	defer c.Close()
	testData := c.TestDataSlice[0].(*redisUser)
	iCache := c.ICache.(Cache)

	key := utils.Uint64ToStr(testData.ID)
	err := iCache.Set(c.Ctx, key, c.TestDataMap[key], time.Minute)
	assert.NoError(t, err)
	err = iCache.Set(c.Ctx, key, c.TestDataMap[key], 0)
	assert.NoError(t, err)

	val := &redisUser{}
	err = iCache.Get(c.Ctx, key, val)
	assert.NoError(t, err)
	assert.Equal(t, testData.Name, val.Name)

	err = iCache.Del(c.Ctx, key)
	assert.NoError(t, err)

	err = iCache.MultiSet(c.Ctx, c.TestDataMap, time.Minute)
	assert.NoError(t, err)
	err = iCache.MultiSet(c.Ctx, c.TestDataMap, 0)
	assert.NoError(t, err)

	var keys []string
	for k := range c.TestDataMap {
		keys = append(keys, k)
	}
	vals := make(map[string]*redisUser)
	err = iCache.MultiGet(c.Ctx, keys, vals)
	assert.NoError(t, err)
	assert.Equal(t, len(c.TestDataSlice), len(vals))

	err = iCache.SetCacheWithNotFound(c.Ctx, "not_found")
	assert.NoError(t, err)
}

func TestRedisCacheError(t *testing.T) {
	c := newRedisCache()
	defer c.Close()
	testData := c.TestDataSlice[0].(*redisUser)
	iCache := c.ICache.(Cache)

	// Set empty key error test
	key := utils.Uint64ToStr(testData.ID)
	err := iCache.Set(c.Ctx, "", c.TestDataMap[key], time.Minute)
	assert.Error(t, err)

	// Set empty value error test
	key = utils.Uint64ToStr(testData.ID)
	err = iCache.Set(c.Ctx, key, nil, time.Minute)
	assert.Error(t, err)

	// Get empty key error test
	val := &redisUser{}
	err = iCache.Get(c.Ctx, "", val)
	assert.Error(t, err)

	// Get empty result  test
	key = utils.Uint64ToStr(testData.ID)
	err = iCache.Get(c.Ctx, key, val)
	assert.Error(t, err)

	// Get result error test
	key = utils.Uint64ToStr(testData.ID)
	_ = iCache.Set(c.Ctx, key, c.TestDataMap[key], time.Minute)
	time.Sleep(time.Millisecond)
	err = iCache.Get(c.Ctx, key, nil)
	assert.Error(t, err)

	_ = iCache.MultiSet(c.Ctx, nil, time.Minute)
	_ = iCache.MultiGet(c.Ctx, nil, time.Minute)

	// Del empty key error test
	err = iCache.Del(c.Ctx)
	assert.NoError(t, err)
	err = iCache.Del(c.Ctx, "")
	assert.NoError(t, err)
}

func TestBuildCacheKey(t *testing.T) {
	_, err := BuildCacheKey("", "")
	assert.Error(t, err)
	_, err = BuildCacheKey("foo", "bar")
	assert.NoError(t, err)
}
