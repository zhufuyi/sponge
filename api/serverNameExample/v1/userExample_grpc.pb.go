// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: api/serverNameExample/v1/userExample.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserExampleClient is the client API for UserExample service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserExampleClient interface {
	// create userExample
	Create(ctx context.Context, in *CreateUserExampleRequest, opts ...grpc.CallOption) (*CreateUserExampleReply, error)
	// delete userExample by id
	DeleteByID(ctx context.Context, in *DeleteUserExampleByIDRequest, opts ...grpc.CallOption) (*DeleteUserExampleByIDReply, error)
	// delete userExample by batch id
	DeleteByIDs(ctx context.Context, in *DeleteUserExampleByIDsRequest, opts ...grpc.CallOption) (*DeleteUserExampleByIDsReply, error)
	// update userExample by id
	UpdateByID(ctx context.Context, in *UpdateUserExampleByIDRequest, opts ...grpc.CallOption) (*UpdateUserExampleByIDReply, error)
	// get userExample by id
	GetByID(ctx context.Context, in *GetUserExampleByIDRequest, opts ...grpc.CallOption) (*GetUserExampleByIDReply, error)
	// get userExample by condition
	GetByCondition(ctx context.Context, in *GetUserExampleByConditionRequest, opts ...grpc.CallOption) (*GetUserExampleByConditionReply, error)
	// list of userExample by batch id
	ListByIDs(ctx context.Context, in *ListUserExampleByIDsRequest, opts ...grpc.CallOption) (*ListUserExampleByIDsReply, error)
	// list userExample by last id
	ListByLastID(ctx context.Context, in *ListUserExampleByLastIDRequest, opts ...grpc.CallOption) (*ListUserExampleByLastIDReply, error)
	// list of userExample by query parameters
	List(ctx context.Context, in *ListUserExampleRequest, opts ...grpc.CallOption) (*ListUserExampleReply, error)
}

type userExampleClient struct {
	cc grpc.ClientConnInterface
}

func NewUserExampleClient(cc grpc.ClientConnInterface) UserExampleClient {
	return &userExampleClient{cc}
}

func (c *userExampleClient) Create(ctx context.Context, in *CreateUserExampleRequest, opts ...grpc.CallOption) (*CreateUserExampleReply, error) {
	out := new(CreateUserExampleReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) DeleteByID(ctx context.Context, in *DeleteUserExampleByIDRequest, opts ...grpc.CallOption) (*DeleteUserExampleByIDReply, error) {
	out := new(DeleteUserExampleByIDReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/DeleteByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) DeleteByIDs(ctx context.Context, in *DeleteUserExampleByIDsRequest, opts ...grpc.CallOption) (*DeleteUserExampleByIDsReply, error) {
	out := new(DeleteUserExampleByIDsReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/DeleteByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) UpdateByID(ctx context.Context, in *UpdateUserExampleByIDRequest, opts ...grpc.CallOption) (*UpdateUserExampleByIDReply, error) {
	out := new(UpdateUserExampleByIDReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/UpdateByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) GetByID(ctx context.Context, in *GetUserExampleByIDRequest, opts ...grpc.CallOption) (*GetUserExampleByIDReply, error) {
	out := new(GetUserExampleByIDReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/GetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) GetByCondition(ctx context.Context, in *GetUserExampleByConditionRequest, opts ...grpc.CallOption) (*GetUserExampleByConditionReply, error) {
	out := new(GetUserExampleByConditionReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/GetByCondition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) ListByIDs(ctx context.Context, in *ListUserExampleByIDsRequest, opts ...grpc.CallOption) (*ListUserExampleByIDsReply, error) {
	out := new(ListUserExampleByIDsReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/ListByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) ListByLastID(ctx context.Context, in *ListUserExampleByLastIDRequest, opts ...grpc.CallOption) (*ListUserExampleByLastIDReply, error) {
	out := new(ListUserExampleByLastIDReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/ListByLastID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userExampleClient) List(ctx context.Context, in *ListUserExampleRequest, opts ...grpc.CallOption) (*ListUserExampleReply, error) {
	out := new(ListUserExampleReply)
	err := c.cc.Invoke(ctx, "/api.serverNameExample.v1.userExample/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserExampleServer is the server API for UserExample service.
// All implementations must embed UnimplementedUserExampleServer
// for forward compatibility
type UserExampleServer interface {
	// create userExample
	Create(context.Context, *CreateUserExampleRequest) (*CreateUserExampleReply, error)
	// delete userExample by id
	DeleteByID(context.Context, *DeleteUserExampleByIDRequest) (*DeleteUserExampleByIDReply, error)
	// delete userExample by batch id
	DeleteByIDs(context.Context, *DeleteUserExampleByIDsRequest) (*DeleteUserExampleByIDsReply, error)
	// update userExample by id
	UpdateByID(context.Context, *UpdateUserExampleByIDRequest) (*UpdateUserExampleByIDReply, error)
	// get userExample by id
	GetByID(context.Context, *GetUserExampleByIDRequest) (*GetUserExampleByIDReply, error)
	// get userExample by condition
	GetByCondition(context.Context, *GetUserExampleByConditionRequest) (*GetUserExampleByConditionReply, error)
	// list of userExample by batch id
	ListByIDs(context.Context, *ListUserExampleByIDsRequest) (*ListUserExampleByIDsReply, error)
	// list userExample by last id
	ListByLastID(context.Context, *ListUserExampleByLastIDRequest) (*ListUserExampleByLastIDReply, error)
	// list of userExample by query parameters
	List(context.Context, *ListUserExampleRequest) (*ListUserExampleReply, error)
	mustEmbedUnimplementedUserExampleServer()
}

// UnimplementedUserExampleServer must be embedded to have forward compatible implementations.
type UnimplementedUserExampleServer struct {
}

func (UnimplementedUserExampleServer) Create(context.Context, *CreateUserExampleRequest) (*CreateUserExampleReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserExampleServer) DeleteByID(context.Context, *DeleteUserExampleByIDRequest) (*DeleteUserExampleByIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteByID not implemented")
}
func (UnimplementedUserExampleServer) DeleteByIDs(context.Context, *DeleteUserExampleByIDsRequest) (*DeleteUserExampleByIDsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteByIDs not implemented")
}
func (UnimplementedUserExampleServer) UpdateByID(context.Context, *UpdateUserExampleByIDRequest) (*UpdateUserExampleByIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateByID not implemented")
}
func (UnimplementedUserExampleServer) GetByID(context.Context, *GetUserExampleByIDRequest) (*GetUserExampleByIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedUserExampleServer) GetByCondition(context.Context, *GetUserExampleByConditionRequest) (*GetUserExampleByConditionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByCondition not implemented")
}
func (UnimplementedUserExampleServer) ListByIDs(context.Context, *ListUserExampleByIDsRequest) (*ListUserExampleByIDsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListByIDs not implemented")
}
func (UnimplementedUserExampleServer) ListByLastID(context.Context, *ListUserExampleByLastIDRequest) (*ListUserExampleByLastIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListByLastID not implemented")
}
func (UnimplementedUserExampleServer) List(context.Context, *ListUserExampleRequest) (*ListUserExampleReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedUserExampleServer) mustEmbedUnimplementedUserExampleServer() {}

// UnsafeUserExampleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserExampleServer will
// result in compilation errors.
type UnsafeUserExampleServer interface {
	mustEmbedUnimplementedUserExampleServer()
}

func RegisterUserExampleServer(s grpc.ServiceRegistrar, srv UserExampleServer) {
	s.RegisterService(&UserExample_ServiceDesc, srv)
}

func _UserExample_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).Create(ctx, req.(*CreateUserExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_DeleteByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserExampleByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).DeleteByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/DeleteByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).DeleteByID(ctx, req.(*DeleteUserExampleByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_DeleteByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserExampleByIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).DeleteByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/DeleteByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).DeleteByIDs(ctx, req.(*DeleteUserExampleByIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_UpdateByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserExampleByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).UpdateByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/UpdateByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).UpdateByID(ctx, req.(*UpdateUserExampleByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserExampleByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/GetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).GetByID(ctx, req.(*GetUserExampleByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_GetByCondition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserExampleByConditionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).GetByCondition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/GetByCondition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).GetByCondition(ctx, req.(*GetUserExampleByConditionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_ListByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserExampleByIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).ListByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/ListByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).ListByIDs(ctx, req.(*ListUserExampleByIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_ListByLastID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserExampleByLastIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).ListByLastID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/ListByLastID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).ListByLastID(ctx, req.(*ListUserExampleByLastIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserExample_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserExampleServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.serverNameExample.v1.userExample/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserExampleServer).List(ctx, req.(*ListUserExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserExample_ServiceDesc is the grpc.ServiceDesc for UserExample service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserExample_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.serverNameExample.v1.userExample",
	HandlerType: (*UserExampleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserExample_Create_Handler,
		},
		{
			MethodName: "DeleteByID",
			Handler:    _UserExample_DeleteByID_Handler,
		},
		{
			MethodName: "DeleteByIDs",
			Handler:    _UserExample_DeleteByIDs_Handler,
		},
		{
			MethodName: "UpdateByID",
			Handler:    _UserExample_UpdateByID_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _UserExample_GetByID_Handler,
		},
		{
			MethodName: "GetByCondition",
			Handler:    _UserExample_GetByCondition_Handler,
		},
		{
			MethodName: "ListByIDs",
			Handler:    _UserExample_ListByIDs_Handler,
		},
		{
			MethodName: "ListByLastID",
			Handler:    _UserExample_ListByLastID_Handler,
		},
		{
			MethodName: "List",
			Handler:    _UserExample_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/serverNameExample/v1/userExample.proto",
}
