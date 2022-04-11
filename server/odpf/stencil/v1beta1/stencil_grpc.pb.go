// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package stencilv1beta1

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

// StencilServiceClient is the client API for StencilService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StencilServiceClient interface {
	ListNamespaces(ctx context.Context, in *ListNamespacesRequest, opts ...grpc.CallOption) (*ListNamespacesResponse, error)
	GetNamespace(ctx context.Context, in *GetNamespaceRequest, opts ...grpc.CallOption) (*GetNamespaceResponse, error)
	CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, opts ...grpc.CallOption) (*CreateNamespaceResponse, error)
	UpdateNamespace(ctx context.Context, in *UpdateNamespaceRequest, opts ...grpc.CallOption) (*UpdateNamespaceResponse, error)
	DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, opts ...grpc.CallOption) (*DeleteNamespaceResponse, error)
	ListSchemas(ctx context.Context, in *ListSchemasRequest, opts ...grpc.CallOption) (*ListSchemasResponse, error)
	CreateSchema(ctx context.Context, in *CreateSchemaRequest, opts ...grpc.CallOption) (*CreateSchemaResponse, error)
	CheckCompatibility(ctx context.Context, in *CheckCompatibilityRequest, opts ...grpc.CallOption) (*CheckCompatibilityResponse, error)
	GetSchemaMetadata(ctx context.Context, in *GetSchemaMetadataRequest, opts ...grpc.CallOption) (*GetSchemaMetadataResponse, error)
	UpdateSchemaMetadata(ctx context.Context, in *UpdateSchemaMetadataRequest, opts ...grpc.CallOption) (*UpdateSchemaMetadataResponse, error)
	GetLatestSchema(ctx context.Context, in *GetLatestSchemaRequest, opts ...grpc.CallOption) (*GetLatestSchemaResponse, error)
	DeleteSchema(ctx context.Context, in *DeleteSchemaRequest, opts ...grpc.CallOption) (*DeleteSchemaResponse, error)
	GetSchema(ctx context.Context, in *GetSchemaRequest, opts ...grpc.CallOption) (*GetSchemaResponse, error)
	ListVersions(ctx context.Context, in *ListVersionsRequest, opts ...grpc.CallOption) (*ListVersionsResponse, error)
	DeleteVersion(ctx context.Context, in *DeleteVersionRequest, opts ...grpc.CallOption) (*DeleteVersionResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type stencilServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStencilServiceClient(cc grpc.ClientConnInterface) StencilServiceClient {
	return &stencilServiceClient{cc}
}

func (c *stencilServiceClient) ListNamespaces(ctx context.Context, in *ListNamespacesRequest, opts ...grpc.CallOption) (*ListNamespacesResponse, error) {
	out := new(ListNamespacesResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/ListNamespaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) GetNamespace(ctx context.Context, in *GetNamespaceRequest, opts ...grpc.CallOption) (*GetNamespaceResponse, error) {
	out := new(GetNamespaceResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/GetNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, opts ...grpc.CallOption) (*CreateNamespaceResponse, error) {
	out := new(CreateNamespaceResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/CreateNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) UpdateNamespace(ctx context.Context, in *UpdateNamespaceRequest, opts ...grpc.CallOption) (*UpdateNamespaceResponse, error) {
	out := new(UpdateNamespaceResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/UpdateNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, opts ...grpc.CallOption) (*DeleteNamespaceResponse, error) {
	out := new(DeleteNamespaceResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/DeleteNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) ListSchemas(ctx context.Context, in *ListSchemasRequest, opts ...grpc.CallOption) (*ListSchemasResponse, error) {
	out := new(ListSchemasResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/ListSchemas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) CreateSchema(ctx context.Context, in *CreateSchemaRequest, opts ...grpc.CallOption) (*CreateSchemaResponse, error) {
	out := new(CreateSchemaResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/CreateSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) CheckCompatibility(ctx context.Context, in *CheckCompatibilityRequest, opts ...grpc.CallOption) (*CheckCompatibilityResponse, error) {
	out := new(CheckCompatibilityResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/CheckCompatibility", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) GetSchemaMetadata(ctx context.Context, in *GetSchemaMetadataRequest, opts ...grpc.CallOption) (*GetSchemaMetadataResponse, error) {
	out := new(GetSchemaMetadataResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/GetSchemaMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) UpdateSchemaMetadata(ctx context.Context, in *UpdateSchemaMetadataRequest, opts ...grpc.CallOption) (*UpdateSchemaMetadataResponse, error) {
	out := new(UpdateSchemaMetadataResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/UpdateSchemaMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) GetLatestSchema(ctx context.Context, in *GetLatestSchemaRequest, opts ...grpc.CallOption) (*GetLatestSchemaResponse, error) {
	out := new(GetLatestSchemaResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/GetLatestSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) DeleteSchema(ctx context.Context, in *DeleteSchemaRequest, opts ...grpc.CallOption) (*DeleteSchemaResponse, error) {
	out := new(DeleteSchemaResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/DeleteSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) GetSchema(ctx context.Context, in *GetSchemaRequest, opts ...grpc.CallOption) (*GetSchemaResponse, error) {
	out := new(GetSchemaResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/GetSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) ListVersions(ctx context.Context, in *ListVersionsRequest, opts ...grpc.CallOption) (*ListVersionsResponse, error) {
	out := new(ListVersionsResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/ListVersions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) DeleteVersion(ctx context.Context, in *DeleteVersionRequest, opts ...grpc.CallOption) (*DeleteVersionResponse, error) {
	out := new(DeleteVersionResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/DeleteVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stencilServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/odpf.stencil.v1beta1.StencilService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StencilServiceServer is the server API for StencilService service.
// All implementations must embed UnimplementedStencilServiceServer
// for forward compatibility
type StencilServiceServer interface {
	ListNamespaces(context.Context, *ListNamespacesRequest) (*ListNamespacesResponse, error)
	GetNamespace(context.Context, *GetNamespaceRequest) (*GetNamespaceResponse, error)
	CreateNamespace(context.Context, *CreateNamespaceRequest) (*CreateNamespaceResponse, error)
	UpdateNamespace(context.Context, *UpdateNamespaceRequest) (*UpdateNamespaceResponse, error)
	DeleteNamespace(context.Context, *DeleteNamespaceRequest) (*DeleteNamespaceResponse, error)
	ListSchemas(context.Context, *ListSchemasRequest) (*ListSchemasResponse, error)
	CreateSchema(context.Context, *CreateSchemaRequest) (*CreateSchemaResponse, error)
	CheckCompatibility(context.Context, *CheckCompatibilityRequest) (*CheckCompatibilityResponse, error)
	GetSchemaMetadata(context.Context, *GetSchemaMetadataRequest) (*GetSchemaMetadataResponse, error)
	UpdateSchemaMetadata(context.Context, *UpdateSchemaMetadataRequest) (*UpdateSchemaMetadataResponse, error)
	GetLatestSchema(context.Context, *GetLatestSchemaRequest) (*GetLatestSchemaResponse, error)
	DeleteSchema(context.Context, *DeleteSchemaRequest) (*DeleteSchemaResponse, error)
	GetSchema(context.Context, *GetSchemaRequest) (*GetSchemaResponse, error)
	ListVersions(context.Context, *ListVersionsRequest) (*ListVersionsResponse, error)
	DeleteVersion(context.Context, *DeleteVersionRequest) (*DeleteVersionResponse, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	mustEmbedUnimplementedStencilServiceServer()
}

// UnimplementedStencilServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStencilServiceServer struct {
}

func (UnimplementedStencilServiceServer) ListNamespaces(context.Context, *ListNamespacesRequest) (*ListNamespacesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNamespaces not implemented")
}
func (UnimplementedStencilServiceServer) GetNamespace(context.Context, *GetNamespaceRequest) (*GetNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNamespace not implemented")
}
func (UnimplementedStencilServiceServer) CreateNamespace(context.Context, *CreateNamespaceRequest) (*CreateNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNamespace not implemented")
}
func (UnimplementedStencilServiceServer) UpdateNamespace(context.Context, *UpdateNamespaceRequest) (*UpdateNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNamespace not implemented")
}
func (UnimplementedStencilServiceServer) DeleteNamespace(context.Context, *DeleteNamespaceRequest) (*DeleteNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNamespace not implemented")
}
func (UnimplementedStencilServiceServer) ListSchemas(context.Context, *ListSchemasRequest) (*ListSchemasResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSchemas not implemented")
}
func (UnimplementedStencilServiceServer) CreateSchema(context.Context, *CreateSchemaRequest) (*CreateSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSchema not implemented")
}
func (UnimplementedStencilServiceServer) CheckCompatibility(context.Context, *CheckCompatibilityRequest) (*CheckCompatibilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckCompatibility not implemented")
}
func (UnimplementedStencilServiceServer) GetSchemaMetadata(context.Context, *GetSchemaMetadataRequest) (*GetSchemaMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchemaMetadata not implemented")
}
func (UnimplementedStencilServiceServer) UpdateSchemaMetadata(context.Context, *UpdateSchemaMetadataRequest) (*UpdateSchemaMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSchemaMetadata not implemented")
}
func (UnimplementedStencilServiceServer) GetLatestSchema(context.Context, *GetLatestSchemaRequest) (*GetLatestSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLatestSchema not implemented")
}
func (UnimplementedStencilServiceServer) DeleteSchema(context.Context, *DeleteSchemaRequest) (*DeleteSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSchema not implemented")
}
func (UnimplementedStencilServiceServer) GetSchema(context.Context, *GetSchemaRequest) (*GetSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchema not implemented")
}
func (UnimplementedStencilServiceServer) ListVersions(context.Context, *ListVersionsRequest) (*ListVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVersions not implemented")
}
func (UnimplementedStencilServiceServer) DeleteVersion(context.Context, *DeleteVersionRequest) (*DeleteVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVersion not implemented")
}
func (UnimplementedStencilServiceServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedStencilServiceServer) mustEmbedUnimplementedStencilServiceServer() {}

// UnsafeStencilServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StencilServiceServer will
// result in compilation errors.
type UnsafeStencilServiceServer interface {
	mustEmbedUnimplementedStencilServiceServer()
}

func RegisterStencilServiceServer(s grpc.ServiceRegistrar, srv StencilServiceServer) {
	s.RegisterService(&StencilService_ServiceDesc, srv)
}

func _StencilService_ListNamespaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNamespacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).ListNamespaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/ListNamespaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).ListNamespaces(ctx, req.(*ListNamespacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_GetNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).GetNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/GetNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).GetNamespace(ctx, req.(*GetNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_CreateNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).CreateNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/CreateNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).CreateNamespace(ctx, req.(*CreateNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_UpdateNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).UpdateNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/UpdateNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).UpdateNamespace(ctx, req.(*UpdateNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_DeleteNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).DeleteNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/DeleteNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).DeleteNamespace(ctx, req.(*DeleteNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_ListSchemas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSchemasRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).ListSchemas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/ListSchemas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).ListSchemas(ctx, req.(*ListSchemasRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_CreateSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).CreateSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/CreateSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).CreateSchema(ctx, req.(*CreateSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_CheckCompatibility_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckCompatibilityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).CheckCompatibility(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/CheckCompatibility",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).CheckCompatibility(ctx, req.(*CheckCompatibilityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_GetSchemaMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).GetSchemaMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/GetSchemaMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).GetSchemaMetadata(ctx, req.(*GetSchemaMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_UpdateSchemaMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSchemaMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).UpdateSchemaMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/UpdateSchemaMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).UpdateSchemaMetadata(ctx, req.(*UpdateSchemaMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_GetLatestSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLatestSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).GetLatestSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/GetLatestSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).GetLatestSchema(ctx, req.(*GetLatestSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_DeleteSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).DeleteSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/DeleteSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).DeleteSchema(ctx, req.(*DeleteSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_GetSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).GetSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/GetSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).GetSchema(ctx, req.(*GetSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_ListVersions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVersionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).ListVersions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/ListVersions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).ListVersions(ctx, req.(*ListVersionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_DeleteVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).DeleteVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/DeleteVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).DeleteVersion(ctx, req.(*DeleteVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StencilService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StencilServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.stencil.v1beta1.StencilService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StencilServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StencilService_ServiceDesc is the grpc.ServiceDesc for StencilService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StencilService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "odpf.stencil.v1beta1.StencilService",
	HandlerType: (*StencilServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListNamespaces",
			Handler:    _StencilService_ListNamespaces_Handler,
		},
		{
			MethodName: "GetNamespace",
			Handler:    _StencilService_GetNamespace_Handler,
		},
		{
			MethodName: "CreateNamespace",
			Handler:    _StencilService_CreateNamespace_Handler,
		},
		{
			MethodName: "UpdateNamespace",
			Handler:    _StencilService_UpdateNamespace_Handler,
		},
		{
			MethodName: "DeleteNamespace",
			Handler:    _StencilService_DeleteNamespace_Handler,
		},
		{
			MethodName: "ListSchemas",
			Handler:    _StencilService_ListSchemas_Handler,
		},
		{
			MethodName: "CreateSchema",
			Handler:    _StencilService_CreateSchema_Handler,
		},
		{
			MethodName: "CheckCompatibility",
			Handler:    _StencilService_CheckCompatibility_Handler,
		},
		{
			MethodName: "GetSchemaMetadata",
			Handler:    _StencilService_GetSchemaMetadata_Handler,
		},
		{
			MethodName: "UpdateSchemaMetadata",
			Handler:    _StencilService_UpdateSchemaMetadata_Handler,
		},
		{
			MethodName: "GetLatestSchema",
			Handler:    _StencilService_GetLatestSchema_Handler,
		},
		{
			MethodName: "DeleteSchema",
			Handler:    _StencilService_DeleteSchema_Handler,
		},
		{
			MethodName: "GetSchema",
			Handler:    _StencilService_GetSchema_Handler,
		},
		{
			MethodName: "ListVersions",
			Handler:    _StencilService_ListVersions_Handler,
		},
		{
			MethodName: "DeleteVersion",
			Handler:    _StencilService_DeleteVersion_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _StencilService_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "odpf/stencil/v1beta1/stencil.proto",
}
