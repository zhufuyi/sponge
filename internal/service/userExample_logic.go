package service

import (
	"context"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/rpcclient"
)

var _ serverNameExampleV1.UserExampleLogicer = (*userExampleClient)(nil)

type userExampleClient struct {
	userExampleCli serverNameExampleV1.UserExampleClient
	// If required, fill in the definition of the other service client code here.
}

// NewUserExampleClient creating rpc clients
func NewUserExampleClient() serverNameExampleV1.UserExampleLogicer {
	return &userExampleClient{
		userExampleCli: serverNameExampleV1.NewUserExampleClient(rpcclient.GetServerNameExampleRPCConn()),
		// If required, fill in the code to implement other service clients here.
	}
}

func (c *userExampleClient) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.Create(ctx, req)
}

func (c *userExampleClient) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.DeleteByID(ctx, req)
}

func (c *userExampleClient) DeleteByIDs(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDsRequest) (*serverNameExampleV1.DeleteUserExampleByIDsReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.DeleteByIDs(ctx, req)
}

func (c *userExampleClient) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.UpdateByID(ctx, req)
}

func (c *userExampleClient) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.GetByID(ctx, req)
}

func (c *userExampleClient) GetByCondition(ctx context.Context, req *serverNameExampleV1.GetUserExampleByConditionRequest) (*serverNameExampleV1.GetUserExampleByConditionReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.GetByCondition(ctx, req)
}

func (c *userExampleClient) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.ListByIDs(ctx, req)
}

func (c *userExampleClient) ListByLastID(ctx context.Context, req *serverNameExampleV1.ListUserExampleByLastIDRequest) (*serverNameExampleV1.ListUserExampleByLastIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.ListByLastID(ctx, req)
}

func (c *userExampleClient) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other rpc servers here.
	return c.userExampleCli.List(ctx, req)
}
