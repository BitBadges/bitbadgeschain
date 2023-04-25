// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: wasmx/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryWasmxParamsRequest is the request type for the Query/WasmxParams RPC method.
type QueryWasmxParamsRequest struct {
}

func (m *QueryWasmxParamsRequest) Reset()         { *m = QueryWasmxParamsRequest{} }
func (m *QueryWasmxParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryWasmxParamsRequest) ProtoMessage()    {}
func (*QueryWasmxParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ba433635f5517a5, []int{0}
}
func (m *QueryWasmxParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryWasmxParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryWasmxParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryWasmxParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryWasmxParamsRequest.Merge(m, src)
}
func (m *QueryWasmxParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryWasmxParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryWasmxParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryWasmxParamsRequest proto.InternalMessageInfo

// QueryWasmxParamsRequest is the response type for the Query/WasmxParams RPC method.
type QueryWasmxParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryWasmxParamsResponse) Reset()         { *m = QueryWasmxParamsResponse{} }
func (m *QueryWasmxParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryWasmxParamsResponse) ProtoMessage()    {}
func (*QueryWasmxParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ba433635f5517a5, []int{1}
}
func (m *QueryWasmxParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryWasmxParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryWasmxParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryWasmxParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryWasmxParamsResponse.Merge(m, src)
}
func (m *QueryWasmxParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryWasmxParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryWasmxParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryWasmxParamsResponse proto.InternalMessageInfo

func (m *QueryWasmxParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// QueryModuleStateRequest is the request type for the Query/WasmxModuleState RPC method.
type QueryModuleStateRequest struct {
}

func (m *QueryModuleStateRequest) Reset()         { *m = QueryModuleStateRequest{} }
func (m *QueryModuleStateRequest) String() string { return proto.CompactTextString(m) }
func (*QueryModuleStateRequest) ProtoMessage()    {}
func (*QueryModuleStateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ba433635f5517a5, []int{2}
}
func (m *QueryModuleStateRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryModuleStateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryModuleStateRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryModuleStateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryModuleStateRequest.Merge(m, src)
}
func (m *QueryModuleStateRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryModuleStateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryModuleStateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryModuleStateRequest proto.InternalMessageInfo

// QueryModuleStateResponse is the response type for the Query/WasmxModuleState RPC method.
type QueryModuleStateResponse struct {
	State *GenesisState `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
}

func (m *QueryModuleStateResponse) Reset()         { *m = QueryModuleStateResponse{} }
func (m *QueryModuleStateResponse) String() string { return proto.CompactTextString(m) }
func (*QueryModuleStateResponse) ProtoMessage()    {}
func (*QueryModuleStateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ba433635f5517a5, []int{3}
}
func (m *QueryModuleStateResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryModuleStateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryModuleStateResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryModuleStateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryModuleStateResponse.Merge(m, src)
}
func (m *QueryModuleStateResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryModuleStateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryModuleStateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryModuleStateResponse proto.InternalMessageInfo

func (m *QueryModuleStateResponse) GetState() *GenesisState {
	if m != nil {
		return m.State
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryWasmxParamsRequest)(nil), "bitbadges.bitbadgeschain.wasmx.QueryWasmxParamsRequest")
	proto.RegisterType((*QueryWasmxParamsResponse)(nil), "bitbadges.bitbadgeschain.wasmx.QueryWasmxParamsResponse")
	proto.RegisterType((*QueryModuleStateRequest)(nil), "bitbadges.bitbadgeschain.wasmx.QueryModuleStateRequest")
	proto.RegisterType((*QueryModuleStateResponse)(nil), "bitbadges.bitbadgeschain.wasmx.QueryModuleStateResponse")
}

func init() { proto.RegisterFile("wasmx/query.proto", fileDescriptor_4ba433635f5517a5) }

var fileDescriptor_4ba433635f5517a5 = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x4f, 0x2c, 0xce,
	0xad, 0xd0, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x4b,
	0xca, 0x2c, 0x49, 0x4a, 0x4c, 0x49, 0x4f, 0x2d, 0xd6, 0x83, 0xb3, 0x92, 0x33, 0x12, 0x33, 0xf3,
	0xf4, 0xc0, 0x6a, 0xa5, 0x64, 0xd2, 0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0x13, 0x0b, 0x32, 0xf5,
	0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12, 0x4b, 0x32, 0xf3, 0xf3, 0x8a, 0x21, 0xba, 0xa5, 0xa0, 0x06,
	0x82, 0x49, 0xa8, 0x90, 0x30, 0x44, 0x28, 0x3d, 0x35, 0x2f, 0xb5, 0x38, 0x13, 0xa6, 0x4e, 0x24,
	0x3d, 0x3f, 0x3d, 0x1f, 0xcc, 0xd4, 0x07, 0xb1, 0x20, 0xa2, 0x4a, 0x92, 0x5c, 0xe2, 0x81, 0x20,
	0xa7, 0x84, 0x83, 0x74, 0x04, 0x24, 0x16, 0x25, 0xe6, 0x16, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16,
	0x97, 0x28, 0x25, 0x70, 0x49, 0x60, 0x4a, 0x15, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0xb9, 0x70,
	0xb1, 0x15, 0x80, 0x45, 0x24, 0x18, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0xd4, 0xf4, 0xf0, 0xfb, 0x41,
	0x0f, 0xa2, 0xdf, 0x89, 0xe5, 0xc4, 0x3d, 0x79, 0x86, 0x20, 0xa8, 0x5e, 0xb8, 0xe5, 0xbe, 0xf9,
	0x29, 0xa5, 0x39, 0xa9, 0xc1, 0x25, 0x89, 0x25, 0xa9, 0x30, 0xcb, 0xe3, 0xa0, 0x96, 0xa3, 0x48,
	0x41, 0x2d, 0x77, 0xe2, 0x62, 0x2d, 0x06, 0x09, 0x40, 0xed, 0xd6, 0x21, 0x64, 0xb7, 0x3b, 0x24,
	0x1c, 0x20, 0x86, 0x40, 0xb4, 0x1a, 0xfd, 0x63, 0xe2, 0x62, 0x05, 0x5b, 0x20, 0xb4, 0x99, 0x91,
	0x8b, 0x1b, 0xc9, 0x8b, 0x42, 0xe6, 0x84, 0x8c, 0xc3, 0x11, 0x5e, 0x52, 0x16, 0xa4, 0x6b, 0x84,
	0x78, 0x48, 0xc9, 0xb0, 0xe9, 0xf2, 0x93, 0xc9, 0x4c, 0xda, 0x42, 0x9a, 0xfa, 0x70, 0x7d, 0xfa,
	0xa8, 0x26, 0x40, 0xa2, 0x57, 0xbf, 0xcc, 0x50, 0x1f, 0x12, 0x74, 0x42, 0xfb, 0x18, 0xb9, 0x04,
	0xc0, 0x46, 0x21, 0x05, 0x10, 0x91, 0x4e, 0xc7, 0x0c, 0x6d, 0x22, 0x9d, 0x8e, 0x25, 0x2e, 0x94,
	0xcc, 0xc1, 0x4e, 0x37, 0x14, 0xd2, 0x27, 0xc2, 0xe9, 0xb9, 0x60, 0xfd, 0xf1, 0xe0, 0x08, 0x70,
	0xf2, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c,
	0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xc3, 0xf4, 0xcc, 0x92,
	0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xdc, 0x86, 0x42, 0x13, 0xbc, 0x7e, 0x49, 0x65, 0x41,
	0x6a, 0x71, 0x12, 0x1b, 0x38, 0x31, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xd7, 0x42, 0x71,
	0x1b, 0x5d, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Retrieves wasmx params
	WasmxParams(ctx context.Context, in *QueryWasmxParamsRequest, opts ...grpc.CallOption) (*QueryWasmxParamsResponse, error)
	// Retrieves the entire wasmx module's state
	WasmxModuleState(ctx context.Context, in *QueryModuleStateRequest, opts ...grpc.CallOption) (*QueryModuleStateResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) WasmxParams(ctx context.Context, in *QueryWasmxParamsRequest, opts ...grpc.CallOption) (*QueryWasmxParamsResponse, error) {
	out := new(QueryWasmxParamsResponse)
	err := c.cc.Invoke(ctx, "/bitbadges.bitbadgeschain.wasmx.Query/WasmxParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) WasmxModuleState(ctx context.Context, in *QueryModuleStateRequest, opts ...grpc.CallOption) (*QueryModuleStateResponse, error) {
	out := new(QueryModuleStateResponse)
	err := c.cc.Invoke(ctx, "/bitbadges.bitbadgeschain.wasmx.Query/WasmxModuleState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Retrieves wasmx params
	WasmxParams(context.Context, *QueryWasmxParamsRequest) (*QueryWasmxParamsResponse, error)
	// Retrieves the entire wasmx module's state
	WasmxModuleState(context.Context, *QueryModuleStateRequest) (*QueryModuleStateResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) WasmxParams(ctx context.Context, req *QueryWasmxParamsRequest) (*QueryWasmxParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WasmxParams not implemented")
}
func (*UnimplementedQueryServer) WasmxModuleState(ctx context.Context, req *QueryModuleStateRequest) (*QueryModuleStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WasmxModuleState not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_WasmxParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryWasmxParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).WasmxParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bitbadges.bitbadgeschain.wasmx.Query/WasmxParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).WasmxParams(ctx, req.(*QueryWasmxParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_WasmxModuleState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryModuleStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).WasmxModuleState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bitbadges.bitbadgeschain.wasmx.Query/WasmxModuleState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).WasmxModuleState(ctx, req.(*QueryModuleStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bitbadges.bitbadgeschain.wasmx.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WasmxParams",
			Handler:    _Query_WasmxParams_Handler,
		},
		{
			MethodName: "WasmxModuleState",
			Handler:    _Query_WasmxModuleState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wasmx/query.proto",
}

func (m *QueryWasmxParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryWasmxParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryWasmxParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryWasmxParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryWasmxParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryWasmxParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryModuleStateRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryModuleStateRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryModuleStateRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryModuleStateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryModuleStateResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryModuleStateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.State != nil {
		{
			size, err := m.State.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryWasmxParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryWasmxParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryModuleStateRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryModuleStateResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.State != nil {
		l = m.State.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryWasmxParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryWasmxParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryWasmxParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryWasmxParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryWasmxParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryWasmxParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryModuleStateRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryModuleStateRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryModuleStateRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryModuleStateResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryModuleStateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryModuleStateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.State == nil {
				m.State = &GenesisState{}
			}
			if err := m.State.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)