// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: maps/tx.proto

package maps

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
	Msg_UpdateParams_FullMethodName = "/maps.Msg/UpdateParams"
	Msg_CreateMap_FullMethodName    = "/maps.Msg/CreateMap"
	Msg_UpdateMap_FullMethodName    = "/maps.Msg/UpdateMap"
	Msg_DeleteMap_FullMethodName    = "/maps.Msg/DeleteMap"
	Msg_SetValue_FullMethodName     = "/maps.Msg/SetValue"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	CreateMap(ctx context.Context, in *MsgCreateMap, opts ...grpc.CallOption) (*MsgCreateMapResponse, error)
	UpdateMap(ctx context.Context, in *MsgUpdateMap, opts ...grpc.CallOption) (*MsgUpdateMapResponse, error)
	DeleteMap(ctx context.Context, in *MsgDeleteMap, opts ...grpc.CallOption) (*MsgDeleteMapResponse, error)
	SetValue(ctx context.Context, in *MsgSetValue, opts ...grpc.CallOption) (*MsgSetValueResponse, error)
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

func (c *msgClient) CreateMap(ctx context.Context, in *MsgCreateMap, opts ...grpc.CallOption) (*MsgCreateMapResponse, error) {
	out := new(MsgCreateMapResponse)
	err := c.cc.Invoke(ctx, Msg_CreateMap_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateMap(ctx context.Context, in *MsgUpdateMap, opts ...grpc.CallOption) (*MsgUpdateMapResponse, error) {
	out := new(MsgUpdateMapResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateMap_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeleteMap(ctx context.Context, in *MsgDeleteMap, opts ...grpc.CallOption) (*MsgDeleteMapResponse, error) {
	out := new(MsgDeleteMapResponse)
	err := c.cc.Invoke(ctx, Msg_DeleteMap_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetValue(ctx context.Context, in *MsgSetValue, opts ...grpc.CallOption) (*MsgSetValueResponse, error) {
	out := new(MsgSetValueResponse)
	err := c.cc.Invoke(ctx, Msg_SetValue_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	CreateMap(context.Context, *MsgCreateMap) (*MsgCreateMapResponse, error)
	UpdateMap(context.Context, *MsgUpdateMap) (*MsgUpdateMapResponse, error)
	DeleteMap(context.Context, *MsgDeleteMap) (*MsgDeleteMapResponse, error)
	SetValue(context.Context, *MsgSetValue) (*MsgSetValueResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) CreateMap(context.Context, *MsgCreateMap) (*MsgCreateMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMap not implemented")
}
func (UnimplementedMsgServer) UpdateMap(context.Context, *MsgUpdateMap) (*MsgUpdateMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMap not implemented")
}
func (UnimplementedMsgServer) DeleteMap(context.Context, *MsgDeleteMap) (*MsgDeleteMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMap not implemented")
}
func (UnimplementedMsgServer) SetValue(context.Context, *MsgSetValue) (*MsgSetValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetValue not implemented")
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

func _Msg_CreateMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateMap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateMap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateMap(ctx, req.(*MsgCreateMap))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateMap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateMap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateMap(ctx, req.(*MsgUpdateMap))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeleteMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeleteMap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeleteMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_DeleteMap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeleteMap(ctx, req.(*MsgDeleteMap))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetValue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetValue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetValue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetValue(ctx, req.(*MsgSetValue))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "maps.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "CreateMap",
			Handler:    _Msg_CreateMap_Handler,
		},
		{
			MethodName: "UpdateMap",
			Handler:    _Msg_UpdateMap_Handler,
		},
		{
			MethodName: "DeleteMap",
			Handler:    _Msg_DeleteMap_Handler,
		},
		{
			MethodName: "SetValue",
			Handler:    _Msg_SetValue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "maps/tx.proto",
}