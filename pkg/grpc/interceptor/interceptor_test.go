package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var unaryClientInvoker = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	return nil
}

type streamClient struct {
}

func (s streamClient) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s streamClient) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s streamClient) CloseSend() error {
	return nil
}

func (s streamClient) Context() context.Context {
	return context.Background()
}

func (s streamClient) SendMsg(m interface{}) error {
	return nil
}

func (s streamClient) RecvMsg(m interface{}) error {
	return nil
}

var streamClientFunc = func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return &streamClient{}, nil
}

// -----------------------------------------------------------------------------------------

var unaryServerInfo = &grpc.UnaryServerInfo{
	Server:     nil,
	FullMethod: "/ping",
}

var unaryServerHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
	return nil, nil
}

func newStreamServer(ctx context.Context) *streamServer {
	return &streamServer{
		ctx: ctx,
	}
}

type streamServer struct {
	ctx context.Context
}

func (s streamServer) SetHeader(md metadata.MD) error {
	return nil
}

func (s streamServer) SendHeader(md metadata.MD) error {
	return nil
}

func (s streamServer) SetTrailer(md metadata.MD) {}

func (s streamServer) Context() context.Context {
	return s.ctx
}

func (s streamServer) SendMsg(m interface{}) error {
	return nil
}

func (s streamServer) RecvMsg(m interface{}) error {
	return nil
}

var streamServerInfo = &grpc.StreamServerInfo{
	FullMethod:     "/test",
	IsClientStream: false,
	IsServerStream: false,
}

var streamServerHandler = func(srv interface{}, stream grpc.ServerStream) error {
	return nil
}
