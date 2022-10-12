package v1

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/utils"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGRPCServer() string {
	port, _ := utils.GetAvailablePort()
	clientAddr := fmt.Sprintf("127.0.0.1:%d", port)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerRecovery(),
	))

	RegisterUserExampleServiceServer(server, &UnimplementedUserExampleServiceServer{})

	go func() {
		err = server.Serve(list)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Millisecond * 200)

	return clientAddr
}

func newGRPCClient(addr string) UserExampleServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return NewUserExampleServiceClient(conn)
}

func TestUserExampleService(t *testing.T) {
	addr := newGRPCServer()
	cli := newGRPCClient(addr)
	ctx := context.Background()

	_, err := cli.Create(ctx, &CreateUserExampleRequest{})
	assert.Error(t, err)

	_, err = cli.DeleteByID(ctx, &DeleteUserExampleByIDRequest{})
	assert.Error(t, err)

	_, err = cli.UpdateByID(ctx, &UpdateUserExampleByIDRequest{})
	assert.Error(t, err)

	_, err = cli.GetByID(ctx, &GetUserExampleByIDRequest{})
	assert.Error(t, err)

	_, err = cli.ListByIDs(ctx, &ListUserExampleByIDsRequest{})
	assert.Error(t, err)

	_, err = cli.List(ctx, &ListUserExampleRequest{})
	assert.Error(t, err)
}
