package interceptor

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"

	"github.com/zhufuyi/sponge/pkg/utils"
)

func newUnaryRPCServer(unaryServerInterceptors ...grpc.UnaryServerInterceptor) string {
	return newRPCServer(unaryServerInterceptors, nil)
}

func newStreamRPCServer(streamServerInterceptors ...grpc.StreamServerInterceptor) string {
	return newRPCServer(nil, streamServerInterceptors)
}

func newRPCServer(unaryServerInterceptors []grpc.UnaryServerInterceptor, streamServerInterceptors []grpc.StreamServerInterceptor) string {
	serverAddr, _ := utils.GetLocalHTTPAddrPairs()

	list, err := net.Listen("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

	options1 := grpc_middleware.WithUnaryServerChain(unaryServerInterceptors...)
	options2 := grpc_middleware.WithStreamServerChain(streamServerInterceptors...)
	server := grpc.NewServer(options1, options2)

	RegisterGreeterServer(server, &greeterServer{})

	go func() {
		err = server.Serve(list)
		if err != nil {
			panic(err)
		}
	}()

	return serverAddr
}

func newUnaryRPCClient(addr string, unaryClientInterceptors ...grpc.UnaryClientInterceptor) GreeterClient {
	return newRPCClient(addr, unaryClientInterceptors, nil)
}

func newStreamRPCClient(addr string, streamClientInterceptors ...grpc.StreamClientInterceptor) GreeterClient {
	return newRPCClient(addr, nil, streamClientInterceptors)
}

func newRPCClient(addr string, unaryClientInterceptors []grpc.UnaryClientInterceptor, streamClientInterceptors []grpc.StreamClientInterceptor) GreeterClient {
	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	option1 := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryClientInterceptors...))
	option2 := grpc.WithChainStreamInterceptor(grpc_middleware.ChainStreamClient(streamClientInterceptors...))
	options = append(options, option1, option2)

	addr = "127.0.0.1" + addr
	conn, err := grpc.Dial(addr, options...)
	if err != nil {
		panic(err)
	}

	return NewGreeterClient(conn)
}

func sayHelloMethod(client GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &HelloRequest{Name: "foo"})
	if err != nil {
		return err
	}

	fmt.Println("resp:", resp.Message)
	return nil
}

func discussHelloMethod(client GreeterClient) error {
	stream, err := client.DiscussHello(context.Background())
	if err != nil {
		return err
	}

	names := []string{"foo1", "foo2"}
	var resp *HelloReply
	for _, name := range names {
		err = stream.Send(&HelloRequest{Name: name})
		if err != nil {
			return err
		}

		resp, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fmt.Println("client receive:", resp.Message)
	}

	time.Sleep(10 * time.Millisecond)
	err = stream.CloseSend()
	if err != nil {
		return err
	}

	return nil
}

type greeterServer struct {
	UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "hello " + r.Name}, nil
}

func (g *greeterServer) DiscussHello(stream Greeter_DiscussHelloServer) error {
	recValues := []string{}
	sendValues := []string{}

	defer func() {
		fmt.Println("\nserver receive: ", recValues)
		fmt.Println("server send    : ", sendValues)
	}()

	var resp *HelloRequest
	var err error
	for {
		resp, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		recValues = append(recValues, resp.Name)

		sendMsg := "hello " + resp.Name
		err = stream.Send(&HelloReply{Message: sendMsg})
		if err != nil {
			return err
		}
		sendValues = append(sendValues, sendMsg)
	}
}

// -----------------------------------hello.pb.go-------------------------------------------

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hello_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_hello_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_hello_proto_rawDescGZIP(), []int{0}
}

func (x *HelloRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type HelloReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *HelloReply) Reset() {
	*x = HelloReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hello_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloReply) ProtoMessage() {}

func (x *HelloReply) ProtoReflect() protoreflect.Message {
	mi := &file_hello_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloReply.ProtoReflect.Descriptor instead.
func (*HelloReply) Descriptor() ([]byte, []int) {
	return file_hello_proto_rawDescGZIP(), []int{1}
}

func (x *HelloReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_hello_proto protoreflect.FileDescriptor

var file_hello_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x26, 0x0a, 0x0a, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x7d, 0x0a, 0x07, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x08, 0x53,
	0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x12, 0x3c, 0x0a, 0x0c, 0x44, 0x69, 0x73, 0x63, 0x75, 0x73, 0x73, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42,
	0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hello_proto_rawDescOnce sync.Once
	file_hello_proto_rawDescData = file_hello_proto_rawDesc
)

func file_hello_proto_rawDescGZIP() []byte {
	file_hello_proto_rawDescOnce.Do(func() {
		file_hello_proto_rawDescData = protoimpl.X.CompressGZIP(file_hello_proto_rawDescData)
	})
	return file_hello_proto_rawDescData
}

var file_hello_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_hello_proto_goTypes = []interface{}{
	(*HelloRequest)(nil), // 0: proto.HelloRequest
	(*HelloReply)(nil),   // 1: proto.HelloReply
}
var file_hello_proto_depIdxs = []int32{
	0, // 0: proto.Greeter.SayHello:input_type -> proto.HelloRequest
	0, // 1: proto.Greeter.DiscussHello:input_type -> proto.HelloRequest
	1, // 2: proto.Greeter.SayHello:output_type -> proto.HelloReply
	1, // 3: proto.Greeter.DiscussHello:output_type -> proto.HelloReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_hello_proto_init() }
func file_hello_proto_init() {
	if File_hello_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hello_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hello_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hello_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_hello_proto_goTypes,
		DependencyIndexes: file_hello_proto_depIdxs,
		MessageInfos:      file_hello_proto_msgTypes,
	}.Build()
	File_hello_proto = out.File
	file_hello_proto_rawDesc = nil
	file_hello_proto_goTypes = nil
	file_hello_proto_depIdxs = nil
}

// -----------------------------------hello_grpc.pb.go-------------------------------------
const _ = grpc.SupportPackageIsVersion7

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	DiscussHello(ctx context.Context, opts ...grpc.CallOption) (Greeter_DiscussHelloClient, error)
}

type greeterClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/proto.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) DiscussHello(ctx context.Context, opts ...grpc.CallOption) (Greeter_DiscussHelloClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeter_ServiceDesc.Streams[0], "/proto.Greeter/DiscussHello", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterDiscussHelloClient{stream}
	return x, nil
}

type Greeter_DiscussHelloClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type greeterDiscussHelloClient struct {
	grpc.ClientStream
}

func (x *greeterDiscussHelloClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greeterDiscussHelloClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	DiscussHello(Greeter_DiscussHelloServer) error
	mustEmbedUnimplementedGreeterServer()
}

// UnimplementedGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

func (UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedGreeterServer) DiscussHello(Greeter_DiscussHelloServer) error {
	return status.Errorf(codes.Unimplemented, "method DiscussHello not implemented")
}
func (UnimplementedGreeterServer) mustEmbedUnimplementedGreeterServer() {}

// UnsafeGreeterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreeterServer will
// result in compilation errors.
type UnsafeGreeterServer interface {
	mustEmbedUnimplementedGreeterServer()
}

func RegisterGreeterServer(s grpc.ServiceRegistrar, srv GreeterServer) {
	s.RegisterService(&Greeter_ServiceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_DiscussHello_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreeterServer).DiscussHello(&greeterDiscussHelloServer{stream})
}

type Greeter_DiscussHelloServer interface {
	Send(*HelloReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greeterDiscussHelloServer struct {
	grpc.ServerStream
}

func (x *greeterDiscussHelloServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greeterDiscussHelloServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Greeter_ServiceDesc is the grpc.ServiceDesc for Greeter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Greeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DiscussHello",
			Handler:       _Greeter_DiscussHello_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello.proto",
}

// ------------------------------------------------------------------------------------------

func TestUnaryClientRequestID(t *testing.T) {
	addr := newUnaryRPCServer()
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr, UnaryClientRequestID())
	_ = sayHelloMethod(cli)
}

func TestUnaryServerRequestID(t *testing.T) {
	addr := newUnaryRPCServer(UnaryServerRequestID())
	time.Sleep(time.Millisecond * 200)
	cli := newUnaryRPCClient(addr)
	_ = sayHelloMethod(cli)
}

func TestStreamClientRequestID(t *testing.T) {
	addr := newStreamRPCServer()
	time.Sleep(time.Millisecond * 200)
	cli := newStreamRPCClient(addr, StreamClientRequestID())
	_ = discussHelloMethod(cli)
	time.Sleep(time.Millisecond)
}

func TestStreamServerRequestID(t *testing.T) {
	addr := newStreamRPCServer(StreamServerRequestID())
	time.Sleep(time.Millisecond * 200)
	cli := newStreamRPCClient(addr)
	_ = discussHelloMethod(cli)
	time.Sleep(time.Millisecond)
}

func TestCtxRequestID(t *testing.T) {
	_ = ClientCtxRequestID(context.Background())
	field := CtxRequestIDField(context.Background())
	assert.NotNil(t, field)
	field = ClientCtxRequestIDField(context.Background())
	assert.NotNil(t, field)

	ctx := WrapServerCtx(context.Background())
	assert.NotNil(t, ctx)
	ctx = WrapServerCtx(context.Background(), KV{Key: "foo", Val: "bar"})
	assert.NotNil(t, ctx)
	_ = ServerCtxRequestID(context.Background())
	field = ServerCtxRequestIDField(context.Background())
	assert.NotNil(t, field)
}

func TestSetContextRequestIDKey(t *testing.T) {
	SetContextRequestIDKey("my_request_id")
	SetContextRequestIDKey("foo_bar") // invalid key, sync.Once
	t.Log(ContextRequestIDKey)
	SetContextRequestIDKey("xx") // invalid key
}
