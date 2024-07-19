// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: wasmx/tx.proto

package wasmx

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
	Msg_UpdateParams_FullMethodName              = "/wasmx.Msg/UpdateParams"
	Msg_ExecuteContractCompat_FullMethodName     = "/wasmx.Msg/ExecuteContractCompat"
	Msg_InstantiateContractCompat_FullMethodName = "/wasmx.Msg/InstantiateContractCompat"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	ExecuteContractCompat(ctx context.Context, in *MsgExecuteContractCompat, opts ...grpc.CallOption) (*MsgExecuteContractCompatResponse, error)
	InstantiateContractCompat(ctx context.Context, in *MsgInstantiateContractCompat, opts ...grpc.CallOption) (*MsgInstantiateContractCompatResponse, error)
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

func (c *msgClient) ExecuteContractCompat(ctx context.Context, in *MsgExecuteContractCompat, opts ...grpc.CallOption) (*MsgExecuteContractCompatResponse, error) {
	out := new(MsgExecuteContractCompatResponse)
	err := c.cc.Invoke(ctx, Msg_ExecuteContractCompat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) InstantiateContractCompat(ctx context.Context, in *MsgInstantiateContractCompat, opts ...grpc.CallOption) (*MsgInstantiateContractCompatResponse, error) {
	out := new(MsgInstantiateContractCompatResponse)
	err := c.cc.Invoke(ctx, Msg_InstantiateContractCompat_FullMethodName, in, out, opts...)
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
	ExecuteContractCompat(context.Context, *MsgExecuteContractCompat) (*MsgExecuteContractCompatResponse, error)
	InstantiateContractCompat(context.Context, *MsgInstantiateContractCompat) (*MsgInstantiateContractCompatResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) ExecuteContractCompat(context.Context, *MsgExecuteContractCompat) (*MsgExecuteContractCompatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteContractCompat not implemented")
}
func (UnimplementedMsgServer) InstantiateContractCompat(context.Context, *MsgInstantiateContractCompat) (*MsgInstantiateContractCompatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstantiateContractCompat not implemented")
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

func _Msg_ExecuteContractCompat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExecuteContractCompat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ExecuteContractCompat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_ExecuteContractCompat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ExecuteContractCompat(ctx, req.(*MsgExecuteContractCompat))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_InstantiateContractCompat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgInstantiateContractCompat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).InstantiateContractCompat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_InstantiateContractCompat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).InstantiateContractCompat(ctx, req.(*MsgInstantiateContractCompat))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wasmx.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "ExecuteContractCompat",
			Handler:    _Msg_ExecuteContractCompat_Handler,
		},
		{
			MethodName: "InstantiateContractCompat",
			Handler:    _Msg_InstantiateContractCompat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wasmx/tx.proto",
}
