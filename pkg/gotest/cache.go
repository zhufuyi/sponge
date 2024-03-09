package gotest

import (
	"context"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"github.com/zhufuyi/sponge/pkg/utils"
)

// Cache redis cache
type Cache struct {
	Ctx           context.Context
	TestDataSlice []interface{}
	TestDataMap   map[string]interface{}
	RedisClient   *redis.Client
	redisServer   *miniredis.Miniredis
	ICache        interface{}
}

// NewCache instantiated redis cache
func NewCache(testDataMap map[string]interface{}) *Cache {
	var tds []interface{}
	for _, data := range testDataMap {
		tds = append(tds, data)
	}

	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return &Cache{
		Ctx:           context.Background(),
		TestDataSlice: tds,
		TestDataMap:   testDataMap,
		RedisClient:   redis.NewClient(&redis.Options{Addr: redisServer.Addr()}),
		redisServer:   redisServer,
	}
}

// Close redis server
func (c *Cache) Close() {
	c.redisServer.Close()
}

// GetIDs get test data ids
func (c *Cache) GetIDs() []uint64 {
	var ids []uint64
	for idStr := range c.TestDataMap {
		ids = append(ids, utils.StrToUint64(idStr))
	}
	return ids
}

// GetTestData get test data
func (c *Cache) GetTestData() map[string]interface{} {
	return c.TestDataMap
}
