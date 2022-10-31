package service

import (
	"context"

	serverNameExampleV1 "github.com/zhufuyi/sponge/api/serverNameExample/v1"
	"github.com/zhufuyi/sponge/internal/rpcclient"
)

var _ serverNameExampleV1.UserExampleServiceLogicer = (*userExampleServiceClient)(nil)

type userExampleServiceClient struct {
	userExampleServiceCli serverNameExampleV1.UserExampleServiceClient
	// If required, fill in the definition of the other service client code here.
}

// NewUserExampleServiceClient creating rpc clients
func NewUserExampleServiceClient() serverNameExampleV1.UserExampleServiceLogicer {
	return &userExampleServiceClient{
		userExampleServiceCli: serverNameExampleV1.NewUserExampleServiceClient(rpcclient.GetServerNameExampleRPCConn()),
		// If required, fill in the code to implement other service clients here.
	}
}

func (c *userExampleServiceClient) Create(ctx context.Context, req *serverNameExampleV1.CreateUserExampleRequest) (*serverNameExampleV1.CreateUserExampleReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other microservices here.
	return c.userExampleServiceCli.Create(ctx, req)
}

func (c *userExampleServiceClient) DeleteByID(ctx context.Context, req *serverNameExampleV1.DeleteUserExampleByIDRequest) (*serverNameExampleV1.DeleteUserExampleByIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other microservices here.
	return c.userExampleServiceCli.DeleteByID(ctx, req)
}

func (c *userExampleServiceClient) UpdateByID(ctx context.Context, req *serverNameExampleV1.UpdateUserExampleByIDRequest) (*serverNameExampleV1.UpdateUserExampleByIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other microservices here.
	return c.userExampleServiceCli.UpdateByID(ctx, req)
}

func (c *userExampleServiceClient) GetByID(ctx context.Context, req *serverNameExampleV1.GetUserExampleByIDRequest) (*serverNameExampleV1.GetUserExampleByIDReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other microservices here.
	return c.userExampleServiceCli.GetByID(ctx, req)
}

func (c *userExampleServiceClient) ListByIDs(ctx context.Context, req *serverNameExampleV1.ListUserExampleByIDsRequest) (*serverNameExampleV1.ListUserExampleByIDsReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other microservices here.
	return c.userExampleServiceCli.ListByIDs(ctx, req)
}

func (c *userExampleServiceClient) List(ctx context.Context, req *serverNameExampleV1.ListUserExampleRequest) (*serverNameExampleV1.ListUserExampleReply, error) {
	// implement me
	// If required, fill in the code to fetch data from other microservices here.
	return c.userExampleServiceCli.List(ctx, req)
}
