package gotest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func newService() *Service {
	var testData = map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}

	// 初始化mock cache
	c := NewCache(map[string]interface{}{"no cache": testData})
	c.ICache = struct{}{} // instantiated cache interface

	// 初始化mock dao
	d := NewDao(c, testData)
	d.IDao = struct{}{} // instantiated dao interface

	// 初始化mock handler
	h := NewService(d, testData)
	h.IServiceClient = struct{}{} // instantiated handler interface

	return h
}

func TestNewService(t *testing.T) {
	s := newService()
	assert.NotNil(t, s)
	defer s.Close()
}

func TestService_GetClientConn(t *testing.T) {
	s := newService()
	assert.NotNil(t, s)
	defer s.Close()

	conn := s.GetClientConn()
	assert.NotNil(t, conn)
}

func TestService_GoGrpcServer(t *testing.T) {
	defer func() { recover() }()

	s := newService()
	assert.NotNil(t, s)
	defer s.Close()

	s.GoGrpcServer()
	time.Sleep(time.Millisecond * 100)
}

func TestServiceError(t *testing.T) {
	var testData = map[string]interface{}{
		"1": "foo",
		"2": "bar",
	}

	s := NewService(nil, testData)
	s.clientAddr = ":0"
	s.GetClientConn()
	time.Sleep(time.Millisecond * 10)

	defer func() { recover() }()
	s.clientConn = &grpc.ClientConn{}
	s.Close()
}
