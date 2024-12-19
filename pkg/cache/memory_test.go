package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

type memoryUser struct {
	ID   uint64
	Name string
}

func newMemoryCache() *gotest.Cache {
	record1 := &memoryUser{
		ID:   1,
		Name: "foo",
	}
	record2 := &memoryUser{
		ID:   2,
		Name: "bar",
	}

	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	cachePrefix := ""
	c.ICache = NewMemoryCache(cachePrefix, encoding.JSONEncoding{}, func() interface{} {
		return &memoryUser{}
	})

	return c
}

func TestMemoryCache(t *testing.T) {
	c := newMemoryCache()
	defer c.Close()
	testData := c.TestDataSlice[0].(*memoryUser)
	iCache := c.ICache.(Cache)

	key := utils.Uint64ToStr(testData.ID)
	err := iCache.Set(c.Ctx, key, c.TestDataMap[key], time.Minute)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond)
	val := &memoryUser{}
	err = iCache.Get(c.Ctx, key, val)
	assert.NoError(t, err)
	assert.Equal(t, testData.Name, val.Name)

	err = iCache.Del(c.Ctx, key)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond)
	err = iCache.MultiSet(c.Ctx, c.TestDataMap, time.Minute)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond)
	var keys []string
	for k := range c.TestDataMap {
		keys = append(keys, k)
	}
	vals := make(map[string]*memoryUser)
	err = iCache.MultiGet(c.Ctx, keys, vals)
	assert.NoError(t, err)
	assert.Equal(t, len(c.TestDataSlice), len(vals))

	err = iCache.SetCacheWithNotFound(c.Ctx, "not_found")
	assert.NoError(t, err)
}

func TestMemoryCacheError(t *testing.T) {
	c := newMemoryCache()
	defer c.Close()
	testData := c.TestDataSlice[0].(*memoryUser)
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
	val := &memoryUser{}
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

	// Del empty key error test
	err = iCache.Del(c.Ctx)
	assert.NoError(t, err)
	err = iCache.Del(c.Ctx, "")
	assert.Error(t, err)

	// empty key test
	err = iCache.SetCacheWithNotFound(c.Ctx, "")
	assert.Error(t, err)
}
