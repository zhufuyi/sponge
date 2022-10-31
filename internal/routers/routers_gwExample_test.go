package routers

import (
	"context"
	"testing"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"

	"github.com/gin-gonic/gin"
)

func TestNewRouter_gwExample(t *testing.T) {
	err := config.Init(configs.Path("serverNameExample.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = true
	config.Get().App.EnableTracing = true
	config.Get().App.EnablePprof = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableCircuitBreaker = true

	gin.SetMode(gin.ReleaseMode)
	r := NewRouter_gwExample()

	defer func() { recover() }()
	userExampleRouter_gwExample(r, &mockGw{})
}

type mockGw struct{}

func (m mockGw) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	return nil, nil
}

func (m mockGw) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	return nil, nil
}

func (m mockGw) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	return nil, nil
}

func (m mockGw) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	return nil, nil
}

func (m mockGw) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	return nil, nil
}

func (m mockGw) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	return nil, nil
}
