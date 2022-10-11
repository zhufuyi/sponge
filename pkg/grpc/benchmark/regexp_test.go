package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = []byte(`
syntax = "proto3";

package api.use.v1;

option go_package = "./v1;v1";

service useService {
  rpc Create(CreateUseRequest) returns (CreateUseReply) {}
  rpc DeleteByID(DeleteUseByIDRequest) returns (DeleteUseByIDReply) {}
}
`)

func Test_getName(t *testing.T) {
	actual := getName(testData, packagePattern)
	assert.Equal(t, "api.use.v1", actual)

	actual = getName(testData, servicePattern)
	assert.Equal(t, "useService", actual)
}

func Test_getMethodNames(t *testing.T) {
	actual := getMethodNames(testData, methodPattern)
	assert.EqualValues(t, []string{"Create", "DeleteByID"}, actual)
}

func Test_matchName(t *testing.T) {
	methodNames := []string{"Create", "DeleteByID"}
	actual := matchName(methodNames, "Create")
	assert.NotEmpty(t, actual)

	actual = matchName(methodNames, "a")
	assert.Empty(t, actual)
}
