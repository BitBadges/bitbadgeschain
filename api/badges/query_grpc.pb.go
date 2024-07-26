// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: badges/query.proto

package badges

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

const (
	Query_Params_FullMethodName              = "/badges.Query/Params"
	Query_GetCollection_FullMethodName       = "/badges.Query/GetCollection"
	Query_GetAddressList_FullMethodName      = "/badges.Query/GetAddressList"
	Query_GetApprovalTracker_FullMethodName  = "/badges.Query/GetApprovalTracker"
	Query_GetChallengeTracker_FullMethodName = "/badges.Query/GetChallengeTracker"
	Query_GetBalance_FullMethodName          = "/badges.Query/GetBalance"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a badge collection by ID.
	GetCollection(ctx context.Context, in *QueryGetCollectionRequest, opts ...grpc.CallOption) (*QueryGetCollectionResponse, error)
	// Queries an address list by ID.
	GetAddressList(ctx context.Context, in *QueryGetAddressListRequest, opts ...grpc.CallOption) (*QueryGetAddressListResponse, error)
	// Queries an approvals tracker by ID.
	GetApprovalTracker(ctx context.Context, in *QueryGetApprovalTrackerRequest, opts ...grpc.CallOption) (*QueryGetApprovalTrackerResponse, error)
	// Queries the number of times a given leaf has been used for a given merkle challenge.
	GetChallengeTracker(ctx context.Context, in *QueryGetChallengeTrackerRequest, opts ...grpc.CallOption) (*QueryGetChallengeTrackerResponse, error)
	// Queries an addresses balance for a badge collection, specified by its ID.
	GetBalance(ctx context.Context, in *QueryGetBalanceRequest, opts ...grpc.CallOption) (*QueryGetBalanceResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetCollection(ctx context.Context, in *QueryGetCollectionRequest, opts ...grpc.CallOption) (*QueryGetCollectionResponse, error) {
	out := new(QueryGetCollectionResponse)
	err := c.cc.Invoke(ctx, Query_GetCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetAddressList(ctx context.Context, in *QueryGetAddressListRequest, opts ...grpc.CallOption) (*QueryGetAddressListResponse, error) {
	out := new(QueryGetAddressListResponse)
	err := c.cc.Invoke(ctx, Query_GetAddressList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetApprovalTracker(ctx context.Context, in *QueryGetApprovalTrackerRequest, opts ...grpc.CallOption) (*QueryGetApprovalTrackerResponse, error) {
	out := new(QueryGetApprovalTrackerResponse)
	err := c.cc.Invoke(ctx, Query_GetApprovalTracker_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetChallengeTracker(ctx context.Context, in *QueryGetChallengeTrackerRequest, opts ...grpc.CallOption) (*QueryGetChallengeTrackerResponse, error) {
	out := new(QueryGetChallengeTrackerResponse)
	err := c.cc.Invoke(ctx, Query_GetChallengeTracker_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetBalance(ctx context.Context, in *QueryGetBalanceRequest, opts ...grpc.CallOption) (*QueryGetBalanceResponse, error) {
	out := new(QueryGetBalanceResponse)
	err := c.cc.Invoke(ctx, Query_GetBalance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a badge collection by ID.
	GetCollection(context.Context, *QueryGetCollectionRequest) (*QueryGetCollectionResponse, error)
	// Queries an address list by ID.
	GetAddressList(context.Context, *QueryGetAddressListRequest) (*QueryGetAddressListResponse, error)
	// Queries an approvals tracker by ID.
	GetApprovalTracker(context.Context, *QueryGetApprovalTrackerRequest) (*QueryGetApprovalTrackerResponse, error)
	// Queries the number of times a given leaf has been used for a given merkle challenge.
	GetChallengeTracker(context.Context, *QueryGetChallengeTrackerRequest) (*QueryGetChallengeTrackerResponse, error)
	// Queries an addresses balance for a badge collection, specified by its ID.
	GetBalance(context.Context, *QueryGetBalanceRequest) (*QueryGetBalanceResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) GetCollection(context.Context, *QueryGetCollectionRequest) (*QueryGetCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollection not implemented")
}
func (UnimplementedQueryServer) GetAddressList(context.Context, *QueryGetAddressListRequest) (*QueryGetAddressListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressList not implemented")
}
func (UnimplementedQueryServer) GetApprovalTracker(context.Context, *QueryGetApprovalTrackerRequest) (*QueryGetApprovalTrackerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApprovalTracker not implemented")
}
func (UnimplementedQueryServer) GetChallengeTracker(context.Context, *QueryGetChallengeTrackerRequest) (*QueryGetChallengeTrackerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChallengeTracker not implemented")
}
func (UnimplementedQueryServer) GetBalance(context.Context, *QueryGetBalanceRequest) (*QueryGetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
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

func _Query_GetCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetCollection(ctx, req.(*QueryGetCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetAddressList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetAddressListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetAddressList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetAddressList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetAddressList(ctx, req.(*QueryGetAddressListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetApprovalTracker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetApprovalTrackerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetApprovalTracker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetApprovalTracker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetApprovalTracker(ctx, req.(*QueryGetApprovalTrackerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetChallengeTracker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetChallengeTrackerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetChallengeTracker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetChallengeTracker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetChallengeTracker(ctx, req.(*QueryGetChallengeTrackerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetBalance(ctx, req.(*QueryGetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "badges.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "GetCollection",
			Handler:    _Query_GetCollection_Handler,
		},
		{
			MethodName: "GetAddressList",
			Handler:    _Query_GetAddressList_Handler,
		},
		{
			MethodName: "GetApprovalTracker",
			Handler:    _Query_GetApprovalTracker_Handler,
		},
		{
			MethodName: "GetChallengeTracker",
			Handler:    _Query_GetChallengeTracker_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _Query_GetBalance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "badges/query.proto",
}