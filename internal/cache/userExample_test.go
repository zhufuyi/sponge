package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/sponge/internal/model"
)

type _userExampleCache struct {
	ctx         context.Context
	testDatas   []*model.UserExample
	redisServer *miniredis.Miniredis
	iCache      UserExampleCache
}

func newUserExampleCache() *_userExampleCache {
	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	cache := NewUserExampleCache(redis.NewClient(&redis.Options{Addr: redisServer.Addr()}))
	record1 := &model.UserExample{}
	record1.ID = 1
	record2 := &model.UserExample{}
	record2.ID = 2

	return &_userExampleCache{
		ctx:         context.Background(),
		testDatas:   []*model.UserExample{record1, record2},
		redisServer: redisServer,
		iCache:      cache,
	}
}

func (c *_userExampleCache) close() {
	c.redisServer.Close()
}

func (c *_userExampleCache) getIDs() []uint64 {
	var ids []uint64
	for _, v := range c.testDatas {
		ids = append(ids, v.ID)
	}
	return ids
}

func (c *_userExampleCache) getExpected() map[string]*model.UserExample {
	expected := make(map[string]*model.UserExample)
	for _, v := range c.testDatas {
		record := &model.UserExample{}
		record.ID = v.ID
		expected[fmt.Sprintf("%d", v.ID)] = record
	}
	return expected
}

func Test_userExampleCache_Set(t *testing.T) {
	c := newUserExampleCache()
	defer c.close()

	record := c.testDatas[0]
	err := c.iCache.Set(c.ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_Get(t *testing.T) {
	c := newUserExampleCache()
	defer c.close()

	record := c.testDatas[0]
	err := c.iCache.Set(c.ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.iCache.Get(c.ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)
}

func Test_userExampleCache_MultiGet(t *testing.T) {
	c := newUserExampleCache()
	defer c.close()

	err := c.iCache.MultiSet(c.ctx, c.testDatas, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.iCache.MultiGet(c.ctx, c.getIDs())
	if err != nil {
		t.Fatal(err)
	}
	expected := c.getExpected()
	assert.Equal(t, expected, got)
}

func Test_userExampleCache_MultiSet(t *testing.T) {
	c := newUserExampleCache()
	defer c.close()

	err := c.iCache.MultiSet(c.ctx, c.testDatas, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_Del(t *testing.T) {
	c := newUserExampleCache()
	defer c.close()

	record := c.testDatas[0]
	err := c.iCache.Del(c.ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}
