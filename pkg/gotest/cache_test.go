package gotest

import (
	"testing"
)

func getTestData() map[string]interface{} {
	return map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}
}

func TestNewRedisCache(t *testing.T) {
	testData := getTestData()
	c := NewCache(testData)
	c.ICache = struct{}{}
	c.Close()
}

func TestRedisCache_GetIDs(t *testing.T) {
	testData := getTestData()
	c := NewCache(testData)
	c.ICache = struct{}{}
	c.Close()
	t.Log(c.GetIDs())
}

func TestRedisCache_GetTestData(t *testing.T) {
	testData := getTestData()
	c := NewCache(testData)
	c.ICache = struct{}{}
	c.Close()
	t.Log(c.GetTestData())
}

func TestNewRedisClusterCache(t *testing.T) {
	testData := getTestData()
	c := NewCache(testData)
	c.ICache = struct{}{}
	c.Close()
}

func TestRedisClusterCache_GetIDs(t *testing.T) {
	testData := getTestData()
	rc := NewRCCache(testData)
	rc.ICache = struct{}{}
	rc.Close()
	t.Log(rc.GetIDs())
}

func TestRedisClusterCache_GetTestData(t *testing.T) {
	testData := getTestData()
	rc := NewRCCache(testData)
	rc.ICache = struct{}{}
	rc.Close()
	t.Log(rc.GetTestData())
}
