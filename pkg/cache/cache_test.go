package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

type cacheUser struct {
	ID   uint64
	Name string
}

func newCache() *gotest.Cache {
	record1 := &cacheUser{
		ID:   1,
		Name: "foo",
	}
	record2 := &cacheUser{
		ID:   2,
		Name: "bar",
	}

	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	cachePrefix := ""
	DefaultClient = NewRedisCache(c.RedisClient, cachePrefix, encoding.JSONEncoding{}, func() interface{} {
		return &cacheUser{}
	})
	c.ICache = DefaultClient

	return c
}

func TestCache(t *testing.T) {
	c := newCache()
	defer c.Close()
	testData := c.TestDataSlice[0].(*cacheUser)

	key := utils.Uint64ToStr(testData.ID)
	err := Set(c.Ctx, key, c.TestDataMap[key], time.Minute)
	assert.NoError(t, err)

	val := &cacheUser{}
	err = Get(c.Ctx, key, val)
	assert.NoError(t, err)
	assert.Equal(t, testData.Name, val.Name)

	err = Del(c.Ctx, key)
	assert.NoError(t, err)

	err = MultiSet(c.Ctx, c.TestDataMap, time.Minute)
	assert.NoError(t, err)

	var keys []string
	for k := range c.TestDataMap {
		keys = append(keys, k)
	}
	vals := make(map[string]*cacheUser)
	err = MultiGet(c.Ctx, keys, vals)
	assert.NoError(t, err)
	assert.Equal(t, len(c.TestDataSlice), len(vals))

	err = SetCacheWithNotFound(c.Ctx, "not_found")
	assert.NoError(t, err)
}
