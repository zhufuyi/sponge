package cache

import (
	"testing"
	"time"

	"github.com/zhufuyi/sponge/internal/model"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func newUserExampleCache() *gotest.Cache {
	record1 := &model.UserExample{}
	record1.ID = 1
	record2 := &model.UserExample{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewUserExampleCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_userExampleCache_Set(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserExample)
	err := c.ICache.(UserExampleCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(UserExampleCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_userExampleCache_Get(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserExample)
	err := c.ICache.(UserExampleCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserExampleCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(UserExampleCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_userExampleCache_MultiGet(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	var testData []*model.UserExample
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserExample))
	}

	err := c.ICache.(UserExampleCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserExampleCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.UserExample))
	}
}

func Test_userExampleCache_MultiSet(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	var testData []*model.UserExample
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserExample))
	}

	err := c.ICache.(UserExampleCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_Del(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserExample)
	err := c.ICache.(UserExampleCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_SetCacheWithNotFound(t *testing.T) {
	c := newUserExampleCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserExample)
	err := c.ICache.(UserExampleCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewUserExampleCache(t *testing.T) {
	c := NewUserExampleCache(&model.CacheType{
		CType: "memory",
	})

	assert.NotNil(t, c)
}
