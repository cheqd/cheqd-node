// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: cheqd/resource/v2/query.proto

package resourcev2

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
	Query_Resource_FullMethodName                      = "/cheqd.resource.v2.Query/Resource"
	Query_ResourceMetadata_FullMethodName              = "/cheqd.resource.v2.Query/ResourceMetadata"
	Query_LatestResourceVersion_FullMethodName         = "/cheqd.resource.v2.Query/LatestResourceVersion"
	Query_LatestResourceVersionMetadata_FullMethodName = "/cheqd.resource.v2.Query/LatestResourceVersionMetadata"
	Query_CollectionResources_FullMethodName           = "/cheqd.resource.v2.Query/CollectionResources"
	Query_Params_FullMethodName                        = "/cheqd.resource.v2.Query/Params"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Query defines the gRPC querier service for the resource module
type QueryClient interface {
	// Fetch data/payload for a specific resource (without metadata)
	Resource(ctx context.Context, in *QueryResourceRequest, opts ...grpc.CallOption) (*QueryResourceResponse, error)
	// Fetch only metadata for a specific resource
	ResourceMetadata(ctx context.Context, in *QueryResourceMetadataRequest, opts ...grpc.CallOption) (*QueryResourceMetadataResponse, error)
	// Fetch latest version for a specific resource (without metadata)
	LatestResourceVersion(ctx context.Context, in *QueryLatestResourceVersionRequest, opts ...grpc.CallOption) (*QueryLatestResourceVersionResponse, error)
	// Fetch metadata of the latest version for a specific resource
	LatestResourceVersionMetadata(ctx context.Context, in *QueryLatestResourceVersionMetadataRequest, opts ...grpc.CallOption) (*QueryLatestResourceVersionMetadataResponse, error)
	// Fetch metadata for all resources in a collection
	CollectionResources(ctx context.Context, in *QueryCollectionResourcesRequest, opts ...grpc.CallOption) (*QueryCollectionResourcesResponse, error)
	// Params queries params of the resource module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Resource(ctx context.Context, in *QueryResourceRequest, opts ...grpc.CallOption) (*QueryResourceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryResourceResponse)
	err := c.cc.Invoke(ctx, Query_Resource_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ResourceMetadata(ctx context.Context, in *QueryResourceMetadataRequest, opts ...grpc.CallOption) (*QueryResourceMetadataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryResourceMetadataResponse)
	err := c.cc.Invoke(ctx, Query_ResourceMetadata_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) LatestResourceVersion(ctx context.Context, in *QueryLatestResourceVersionRequest, opts ...grpc.CallOption) (*QueryLatestResourceVersionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryLatestResourceVersionResponse)
	err := c.cc.Invoke(ctx, Query_LatestResourceVersion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) LatestResourceVersionMetadata(ctx context.Context, in *QueryLatestResourceVersionMetadataRequest, opts ...grpc.CallOption) (*QueryLatestResourceVersionMetadataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryLatestResourceVersionMetadataResponse)
	err := c.cc.Invoke(ctx, Query_LatestResourceVersionMetadata_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CollectionResources(ctx context.Context, in *QueryCollectionResourcesRequest, opts ...grpc.CallOption) (*QueryCollectionResourcesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryCollectionResourcesResponse)
	err := c.cc.Invoke(ctx, Query_CollectionResources_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
//
// Query defines the gRPC querier service for the resource module
type QueryServer interface {
	// Fetch data/payload for a specific resource (without metadata)
	Resource(context.Context, *QueryResourceRequest) (*QueryResourceResponse, error)
	// Fetch only metadata for a specific resource
	ResourceMetadata(context.Context, *QueryResourceMetadataRequest) (*QueryResourceMetadataResponse, error)
	// Fetch latest version for a specific resource (without metadata)
	LatestResourceVersion(context.Context, *QueryLatestResourceVersionRequest) (*QueryLatestResourceVersionResponse, error)
	// Fetch metadata of the latest version for a specific resource
	LatestResourceVersionMetadata(context.Context, *QueryLatestResourceVersionMetadataRequest) (*QueryLatestResourceVersionMetadataResponse, error)
	// Fetch metadata for all resources in a collection
	CollectionResources(context.Context, *QueryCollectionResourcesRequest) (*QueryCollectionResourcesResponse, error)
	// Params queries params of the resource module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) Resource(context.Context, *QueryResourceRequest) (*QueryResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Resource not implemented")
}
func (UnimplementedQueryServer) ResourceMetadata(context.Context, *QueryResourceMetadataRequest) (*QueryResourceMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResourceMetadata not implemented")
}
func (UnimplementedQueryServer) LatestResourceVersion(context.Context, *QueryLatestResourceVersionRequest) (*QueryLatestResourceVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LatestResourceVersion not implemented")
}
func (UnimplementedQueryServer) LatestResourceVersionMetadata(context.Context, *QueryLatestResourceVersionMetadataRequest) (*QueryLatestResourceVersionMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LatestResourceVersionMetadata not implemented")
}
func (UnimplementedQueryServer) CollectionResources(context.Context, *QueryCollectionResourcesRequest) (*QueryCollectionResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionResources not implemented")
}
func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}
func (UnimplementedQueryServer) testEmbeddedByValue()               {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	// If the following call pancis, it indicates UnimplementedQueryServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Resource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Resource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Resource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Resource(ctx, req.(*QueryResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ResourceMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryResourceMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ResourceMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ResourceMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ResourceMetadata(ctx, req.(*QueryResourceMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_LatestResourceVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryLatestResourceVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).LatestResourceVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_LatestResourceVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).LatestResourceVersion(ctx, req.(*QueryLatestResourceVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_LatestResourceVersionMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryLatestResourceVersionMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).LatestResourceVersionMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_LatestResourceVersionMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).LatestResourceVersionMetadata(ctx, req.(*QueryLatestResourceVersionMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CollectionResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCollectionResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CollectionResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_CollectionResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CollectionResources(ctx, req.(*QueryCollectionResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cheqd.resource.v2.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Resource",
			Handler:    _Query_Resource_Handler,
		},
		{
			MethodName: "ResourceMetadata",
			Handler:    _Query_ResourceMetadata_Handler,
		},
		{
			MethodName: "LatestResourceVersion",
			Handler:    _Query_LatestResourceVersion_Handler,
		},
		{
			MethodName: "LatestResourceVersionMetadata",
			Handler:    _Query_LatestResourceVersionMetadata_Handler,
		},
		{
			MethodName: "CollectionResources",
			Handler:    _Query_CollectionResources_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cheqd/resource/v2/query.proto",
}
