package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFieldRequestIDFromContext(t *testing.T) {
	field := FieldRequestIDFromContext(&gin.Context{})
	assert.NotNil(t, field)
}

func TestFieldRequestIDFromHeader(t *testing.T) {
	field := FieldRequestIDFromHeader(&gin.Context{
		Request: &http.Request{
			Header: map[string][]string{},
		},
	})

	assert.NotNil(t, field)
}

func TestGetRequestIDFromContext(t *testing.T) {
	str := GetRequestIDFromContext(&gin.Context{})
	t.Log(str)
}

func TestGetRequestIDFromHeaders(t *testing.T) {
	str := GetRequestIDFromHeaders(&gin.Context{
		Request: &http.Request{
			Header: map[string][]string{},
		},
	})
	t.Log(str)
}
