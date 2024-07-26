// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: badges/tx.proto

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
	Msg_UpdateParams_FullMethodName              = "/badges.Msg/UpdateParams"
	Msg_UniversalUpdateCollection_FullMethodName = "/badges.Msg/UniversalUpdateCollection"
	Msg_CreateAddressLists_FullMethodName        = "/badges.Msg/CreateAddressLists"
	Msg_TransferBadges_FullMethodName            = "/badges.Msg/TransferBadges"
	Msg_UpdateUserApprovals_FullMethodName       = "/badges.Msg/UpdateUserApprovals"
	Msg_DeleteCollection_FullMethodName          = "/badges.Msg/DeleteCollection"
	Msg_UpdateCollection_FullMethodName          = "/badges.Msg/UpdateCollection"
	Msg_CreateCollection_FullMethodName          = "/badges.Msg/CreateCollection"
	Msg_GlobalArchive_FullMethodName             = "/badges.Msg/GlobalArchive"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	UniversalUpdateCollection(ctx context.Context, in *MsgUniversalUpdateCollection, opts ...grpc.CallOption) (*MsgUniversalUpdateCollectionResponse, error)
	CreateAddressLists(ctx context.Context, in *MsgCreateAddressLists, opts ...grpc.CallOption) (*MsgCreateAddressListsResponse, error)
	TransferBadges(ctx context.Context, in *MsgTransferBadges, opts ...grpc.CallOption) (*MsgTransferBadgesResponse, error)
	UpdateUserApprovals(ctx context.Context, in *MsgUpdateUserApprovals, opts ...grpc.CallOption) (*MsgUpdateUserApprovalsResponse, error)
	DeleteCollection(ctx context.Context, in *MsgDeleteCollection, opts ...grpc.CallOption) (*MsgDeleteCollectionResponse, error)
	UpdateCollection(ctx context.Context, in *MsgUpdateCollection, opts ...grpc.CallOption) (*MsgUpdateCollectionResponse, error)
	CreateCollection(ctx context.Context, in *MsgCreateCollection, opts ...grpc.CallOption) (*MsgCreateCollectionResponse, error)
	GlobalArchive(ctx context.Context, in *MsgGlobalArchive, opts ...grpc.CallOption) (*MsgGlobalArchiveResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UniversalUpdateCollection(ctx context.Context, in *MsgUniversalUpdateCollection, opts ...grpc.CallOption) (*MsgUniversalUpdateCollectionResponse, error) {
	out := new(MsgUniversalUpdateCollectionResponse)
	err := c.cc.Invoke(ctx, Msg_UniversalUpdateCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateAddressLists(ctx context.Context, in *MsgCreateAddressLists, opts ...grpc.CallOption) (*MsgCreateAddressListsResponse, error) {
	out := new(MsgCreateAddressListsResponse)
	err := c.cc.Invoke(ctx, Msg_CreateAddressLists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) TransferBadges(ctx context.Context, in *MsgTransferBadges, opts ...grpc.CallOption) (*MsgTransferBadgesResponse, error) {
	out := new(MsgTransferBadgesResponse)
	err := c.cc.Invoke(ctx, Msg_TransferBadges_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateUserApprovals(ctx context.Context, in *MsgUpdateUserApprovals, opts ...grpc.CallOption) (*MsgUpdateUserApprovalsResponse, error) {
	out := new(MsgUpdateUserApprovalsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateUserApprovals_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeleteCollection(ctx context.Context, in *MsgDeleteCollection, opts ...grpc.CallOption) (*MsgDeleteCollectionResponse, error) {
	out := new(MsgDeleteCollectionResponse)
	err := c.cc.Invoke(ctx, Msg_DeleteCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateCollection(ctx context.Context, in *MsgUpdateCollection, opts ...grpc.CallOption) (*MsgUpdateCollectionResponse, error) {
	out := new(MsgUpdateCollectionResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateCollection(ctx context.Context, in *MsgCreateCollection, opts ...grpc.CallOption) (*MsgCreateCollectionResponse, error) {
	out := new(MsgCreateCollectionResponse)
	err := c.cc.Invoke(ctx, Msg_CreateCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) GlobalArchive(ctx context.Context, in *MsgGlobalArchive, opts ...grpc.CallOption) (*MsgGlobalArchiveResponse, error) {
	out := new(MsgGlobalArchiveResponse)
	err := c.cc.Invoke(ctx, Msg_GlobalArchive_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	UniversalUpdateCollection(context.Context, *MsgUniversalUpdateCollection) (*MsgUniversalUpdateCollectionResponse, error)
	CreateAddressLists(context.Context, *MsgCreateAddressLists) (*MsgCreateAddressListsResponse, error)
	TransferBadges(context.Context, *MsgTransferBadges) (*MsgTransferBadgesResponse, error)
	UpdateUserApprovals(context.Context, *MsgUpdateUserApprovals) (*MsgUpdateUserApprovalsResponse, error)
	DeleteCollection(context.Context, *MsgDeleteCollection) (*MsgDeleteCollectionResponse, error)
	UpdateCollection(context.Context, *MsgUpdateCollection) (*MsgUpdateCollectionResponse, error)
	CreateCollection(context.Context, *MsgCreateCollection) (*MsgCreateCollectionResponse, error)
	GlobalArchive(context.Context, *MsgGlobalArchive) (*MsgGlobalArchiveResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) UniversalUpdateCollection(context.Context, *MsgUniversalUpdateCollection) (*MsgUniversalUpdateCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UniversalUpdateCollection not implemented")
}
func (UnimplementedMsgServer) CreateAddressLists(context.Context, *MsgCreateAddressLists) (*MsgCreateAddressListsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAddressLists not implemented")
}
func (UnimplementedMsgServer) TransferBadges(context.Context, *MsgTransferBadges) (*MsgTransferBadgesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferBadges not implemented")
}
func (UnimplementedMsgServer) UpdateUserApprovals(context.Context, *MsgUpdateUserApprovals) (*MsgUpdateUserApprovalsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserApprovals not implemented")
}
func (UnimplementedMsgServer) DeleteCollection(context.Context, *MsgDeleteCollection) (*MsgDeleteCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCollection not implemented")
}
func (UnimplementedMsgServer) UpdateCollection(context.Context, *MsgUpdateCollection) (*MsgUpdateCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCollection not implemented")
}
func (UnimplementedMsgServer) CreateCollection(context.Context, *MsgCreateCollection) (*MsgCreateCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCollection not implemented")
}
func (UnimplementedMsgServer) GlobalArchive(context.Context, *MsgGlobalArchive) (*MsgGlobalArchiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GlobalArchive not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UniversalUpdateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUniversalUpdateCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UniversalUpdateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UniversalUpdateCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UniversalUpdateCollection(ctx, req.(*MsgUniversalUpdateCollection))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateAddressLists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateAddressLists)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateAddressLists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateAddressLists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateAddressLists(ctx, req.(*MsgCreateAddressLists))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_TransferBadges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgTransferBadges)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).TransferBadges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_TransferBadges_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).TransferBadges(ctx, req.(*MsgTransferBadges))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateUserApprovals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateUserApprovals)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateUserApprovals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateUserApprovals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateUserApprovals(ctx, req.(*MsgUpdateUserApprovals))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeleteCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeleteCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeleteCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_DeleteCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeleteCollection(ctx, req.(*MsgDeleteCollection))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateCollection(ctx, req.(*MsgUpdateCollection))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateCollection(ctx, req.(*MsgCreateCollection))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_GlobalArchive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgGlobalArchive)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).GlobalArchive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_GlobalArchive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).GlobalArchive(ctx, req.(*MsgGlobalArchive))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "badges.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "UniversalUpdateCollection",
			Handler:    _Msg_UniversalUpdateCollection_Handler,
		},
		{
			MethodName: "CreateAddressLists",
			Handler:    _Msg_CreateAddressLists_Handler,
		},
		{
			MethodName: "TransferBadges",
			Handler:    _Msg_TransferBadges_Handler,
		},
		{
			MethodName: "UpdateUserApprovals",
			Handler:    _Msg_UpdateUserApprovals_Handler,
		},
		{
			MethodName: "DeleteCollection",
			Handler:    _Msg_DeleteCollection_Handler,
		},
		{
			MethodName: "UpdateCollection",
			Handler:    _Msg_UpdateCollection_Handler,
		},
		{
			MethodName: "CreateCollection",
			Handler:    _Msg_CreateCollection_Handler,
		},
		{
			MethodName: "GlobalArchive",
			Handler:    _Msg_GlobalArchive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "badges/tx.proto",
}