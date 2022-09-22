package gotest

import (
	"testing"
)

func TestNewRedisCache(t *testing.T) {
	var testData = map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}

	c := NewCache(testData)
	c.ICache = struct{}{}

	defer c.Close()
}

func TestRedisCache_GetIDs(t *testing.T) {
	var testData = map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}

	c := NewCache(testData)
	c.ICache = struct{}{}
	defer c.Close()

	t.Log(c.GetIDs())
}

func TestRedisCache_GetTestData(t *testing.T) {
	var testData = map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}

	c := NewCache(testData)
	c.ICache = struct{}{}
	defer c.Close()

	t.Log(c.GetTestData())
}
