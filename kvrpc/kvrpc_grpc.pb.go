// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: kvrpc.proto

package kvrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ServiceKV_PostKV_FullMethodName   = "/kvrpc.ServiceKV/PostKV"
	ServiceKV_GetKV_FullMethodName    = "/kvrpc.ServiceKV/GetKV"
	ServiceKV_DeleteKV_FullMethodName = "/kvrpc.ServiceKV/DeleteKV"
)

// ServiceKVClient is the client API for ServiceKV service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The kv service definition.
type ServiceKVClient interface {
	PostKV(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*PostReply, error)
	GetKV(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
	DeleteKV(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error)
}

type serviceKVClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceKVClient(cc grpc.ClientConnInterface) ServiceKVClient {
	return &serviceKVClient{cc}
}

func (c *serviceKVClient) PostKV(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*PostReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostReply)
	err := c.cc.Invoke(ctx, ServiceKV_PostKV_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceKVClient) GetKV(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReply)
	err := c.cc.Invoke(ctx, ServiceKV_GetKV_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceKVClient) DeleteKV(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteReply)
	err := c.cc.Invoke(ctx, ServiceKV_DeleteKV_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceKVServer is the server API for ServiceKV service.
// All implementations must embed UnimplementedServiceKVServer
// for forward compatibility.
//
// The kv service definition.
type ServiceKVServer interface {
	PostKV(context.Context, *PostRequest) (*PostReply, error)
	GetKV(context.Context, *GetRequest) (*GetReply, error)
	DeleteKV(context.Context, *DeleteRequest) (*DeleteReply, error)
	mustEmbedUnimplementedServiceKVServer()
}

// UnimplementedServiceKVServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedServiceKVServer struct{}

func (UnimplementedServiceKVServer) PostKV(context.Context, *PostRequest) (*PostReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostKV not implemented")
}
func (UnimplementedServiceKVServer) GetKV(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKV not implemented")
}
func (UnimplementedServiceKVServer) DeleteKV(context.Context, *DeleteRequest) (*DeleteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKV not implemented")
}
func (UnimplementedServiceKVServer) mustEmbedUnimplementedServiceKVServer() {}
func (UnimplementedServiceKVServer) testEmbeddedByValue()                   {}

// UnsafeServiceKVServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceKVServer will
// result in compilation errors.
type UnsafeServiceKVServer interface {
	mustEmbedUnimplementedServiceKVServer()
}

func RegisterServiceKVServer(s grpc.ServiceRegistrar, srv ServiceKVServer) {
	// If the following call pancis, it indicates UnimplementedServiceKVServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ServiceKV_ServiceDesc, srv)
}

func _ServiceKV_PostKV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceKVServer).PostKV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceKV_PostKV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceKVServer).PostKV(ctx, req.(*PostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceKV_GetKV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceKVServer).GetKV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceKV_GetKV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceKVServer).GetKV(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceKV_DeleteKV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceKVServer).DeleteKV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceKV_DeleteKV_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceKVServer).DeleteKV(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServiceKV_ServiceDesc is the grpc.ServiceDesc for ServiceKV service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServiceKV_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kvrpc.ServiceKV",
	HandlerType: (*ServiceKVServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostKV",
			Handler:    _ServiceKV_PostKV_Handler,
		},
		{
			MethodName: "GetKV",
			Handler:    _ServiceKV_GetKV_Handler,
		},
		{
			MethodName: "DeleteKV",
			Handler:    _ServiceKV_DeleteKV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kvrpc.proto",
}