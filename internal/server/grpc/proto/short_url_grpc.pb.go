// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: grpc/short_url.proto

package proto

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

// ShortURLClient is the client API for ShortURL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortURLClient interface {
	CreateShortURL(ctx context.Context, in *ShortURLRequest, opts ...grpc.CallOption) (*ShortURLResponse, error)
	CreateBatchShortURL(ctx context.Context, in *BatchShortURLRequests, opts ...grpc.CallOption) (*BatchShortURLResponses, error)
	GetShortURL(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetUserURLs(ctx context.Context, in *GetUserURLRequest, opts ...grpc.CallOption) (*GetUserURLResponses, error)
	DeleteUserURLs(ctx context.Context, in *DeleteUserURLRequest, opts ...grpc.CallOption) (*DeleteUserURLResponse, error)
}

type shortURLClient struct {
	cc grpc.ClientConnInterface
}

func NewShortURLClient(cc grpc.ClientConnInterface) ShortURLClient {
	return &shortURLClient{cc}
}

func (c *shortURLClient) CreateShortURL(ctx context.Context, in *ShortURLRequest, opts ...grpc.CallOption) (*ShortURLResponse, error) {
	out := new(ShortURLResponse)
	err := c.cc.Invoke(ctx, "/grpc.ShortURL/CreateShortURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLClient) CreateBatchShortURL(ctx context.Context, in *BatchShortURLRequests, opts ...grpc.CallOption) (*BatchShortURLResponses, error) {
	out := new(BatchShortURLResponses)
	err := c.cc.Invoke(ctx, "/grpc.ShortURL/CreateBatchShortURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLClient) GetShortURL(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/grpc.ShortURL/GetShortURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLClient) GetUserURLs(ctx context.Context, in *GetUserURLRequest, opts ...grpc.CallOption) (*GetUserURLResponses, error) {
	out := new(GetUserURLResponses)
	err := c.cc.Invoke(ctx, "/grpc.ShortURL/GetUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLClient) DeleteUserURLs(ctx context.Context, in *DeleteUserURLRequest, opts ...grpc.CallOption) (*DeleteUserURLResponse, error) {
	out := new(DeleteUserURLResponse)
	err := c.cc.Invoke(ctx, "/grpc.ShortURL/DeleteUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortURLServer is the server API for ShortURL service.
// All implementations must embed UnimplementedShortURLServer
// for forward compatibility
type ShortURLServer interface {
	CreateShortURL(context.Context, *ShortURLRequest) (*ShortURLResponse, error)
	CreateBatchShortURL(context.Context, *BatchShortURLRequests) (*BatchShortURLResponses, error)
	GetShortURL(context.Context, *GetRequest) (*GetResponse, error)
	GetUserURLs(context.Context, *GetUserURLRequest) (*GetUserURLResponses, error)
	DeleteUserURLs(context.Context, *DeleteUserURLRequest) (*DeleteUserURLResponse, error)
	mustEmbedUnimplementedShortURLServer()
}

// UnimplementedShortURLServer must be embedded to have forward compatible implementations.
type UnimplementedShortURLServer struct {
}

func (UnimplementedShortURLServer) CreateShortURL(context.Context, *ShortURLRequest) (*ShortURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortURL not implemented")
}
func (UnimplementedShortURLServer) CreateBatchShortURL(context.Context, *BatchShortURLRequests) (*BatchShortURLResponses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBatchShortURL not implemented")
}
func (UnimplementedShortURLServer) GetShortURL(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortURL not implemented")
}
func (UnimplementedShortURLServer) GetUserURLs(context.Context, *GetUserURLRequest) (*GetUserURLResponses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
}
func (UnimplementedShortURLServer) DeleteUserURLs(context.Context, *DeleteUserURLRequest) (*DeleteUserURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserURLs not implemented")
}
func (UnimplementedShortURLServer) mustEmbedUnimplementedShortURLServer() {}

// UnsafeShortURLServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortURLServer will
// result in compilation errors.
type UnsafeShortURLServer interface {
	mustEmbedUnimplementedShortURLServer()
}

func RegisterShortURLServer(s grpc.ServiceRegistrar, srv ShortURLServer) {
	s.RegisterService(&ShortURL_ServiceDesc, srv)
}

func _ShortURL_CreateShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServer).CreateShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ShortURL/CreateShortURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServer).CreateShortURL(ctx, req.(*ShortURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURL_CreateBatchShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchShortURLRequests)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServer).CreateBatchShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ShortURL/CreateBatchShortURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServer).CreateBatchShortURL(ctx, req.(*BatchShortURLRequests))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURL_GetShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServer).GetShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ShortURL/GetShortURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServer).GetShortURL(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURL_GetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServer).GetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ShortURL/GetUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServer).GetUserURLs(ctx, req.(*GetUserURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURL_DeleteUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServer).DeleteUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.ShortURL/DeleteUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServer).DeleteUserURLs(ctx, req.(*DeleteUserURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortURL_ServiceDesc is the grpc.ServiceDesc for ShortURL service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortURL_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.ShortURL",
	HandlerType: (*ShortURLServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortURL",
			Handler:    _ShortURL_CreateShortURL_Handler,
		},
		{
			MethodName: "CreateBatchShortURL",
			Handler:    _ShortURL_CreateBatchShortURL_Handler,
		},
		{
			MethodName: "GetShortURL",
			Handler:    _ShortURL_GetShortURL_Handler,
		},
		{
			MethodName: "GetUserURLs",
			Handler:    _ShortURL_GetUserURLs_Handler,
		},
		{
			MethodName: "DeleteUserURLs",
			Handler:    _ShortURL_DeleteUserURLs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/short_url.proto",
}
