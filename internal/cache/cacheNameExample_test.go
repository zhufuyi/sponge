package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/zhufuyi/sponge/internal/model"
)

type cacheNameExampleData struct {
	ID    uint64
	Key   interface{}
	Value interface{}
}

func newCacheNameExampleCache() *gotest.Cache {
	// change the type of the value before testing
	var (
		key = "foo1"
		val = "bar1"
	)

	record1 := &cacheNameExampleData{ID: 1, Key: key, Value: val}
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewCacheNameExampleCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_cacheNameExampleCache_Set(t *testing.T) {
	c := newCacheNameExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*cacheNameExampleData)
	key := record.Key.(keyTypeExample)
	value := record.Value.(valueTypeExample)
	err := c.ICache.(CacheNameExampleCache).Set(c.Ctx, key, value, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_cacheNameExampleCache_Get(t *testing.T) {
	c := newCacheNameExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*cacheNameExampleData)
	key := record.Key.(keyTypeExample)
	value := record.Value.(valueTypeExample)
	err := c.ICache.(CacheNameExampleCache).Set(c.Ctx, key, value, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CacheNameExampleCache).Get(c.Ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, value, got)
}

func Test_cacheNameExampleCache_Del(t *testing.T) {
	c := newCacheNameExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*cacheNameExampleData)
	key := record.Key.(keyTypeExample)
	err := c.ICache.(CacheNameExampleCache).Del(c.Ctx, key)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCacheNameExampleCache(t *testing.T) {
	c := NewCacheNameExampleCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)

	defer func() {
		_ = recover()
	}()
	c = NewCacheNameExampleCache(&model.CacheType{
		CType: "",
	})
}
