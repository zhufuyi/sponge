package hystrix

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryClientInterceptor(t *testing.T) {
	interceptor := UnaryClientInterceptor("hystrix",
		WithStatsDCollector("localhost:5555", "hystrix", 0.5, 2048))
	assert.NotNil(t, interceptor)

	err := interceptor(context.Background(), "test.ping", nil, nil, nil, clientInvoker)
	assert.NoError(t, err)
}

func TestStreamClientInterceptor(t *testing.T) {
	interceptor := StreamClientInterceptor("hystrix",
		WithStatsDCollector("localhost:5555", "hystrix", 0.5, 2048))
	assert.NotNil(t, interceptor)

	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return &clientStream{}, nil
	}
	_, err := interceptor(context.Background(), nil, nil, "test.ping", streamer)
	assert.NoError(t, err)
}

func Test_durationToInt(t *testing.T) {
	durationToInt(10*time.Second, time.Second)
}

// ------------------------------------------------------------------------------------------

var clientInvoker = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	return nil
}

type clientStream struct {
}

func (s clientStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s clientStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s clientStream) CloseSend() error {
	return nil
}

func (s clientStream) Context() context.Context {
	return context.Background()
}

func (s clientStream) SendMsg(m interface{}) error {
	return nil
}

func (s clientStream) RecvMsg(m interface{}) error {
	return nil
}
