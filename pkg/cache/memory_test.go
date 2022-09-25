package cache

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
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
