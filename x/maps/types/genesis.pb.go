// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: maps/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// GenesisState defines the maps module's genesis state.
type GenesisState struct {
	Params             Params        `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	PortId             string        `protobuf:"bytes,2,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	Maps               []*Map        `protobuf:"bytes,3,rep,name=maps,proto3" json:"maps,omitempty"`
	FullKeys           []string      `protobuf:"bytes,4,rep,name=full_keys,json=fullKeys,proto3" json:"full_keys,omitempty"`
	Values             []*ValueStore `protobuf:"bytes,5,rep,name=values,proto3" json:"values,omitempty"`
	DuplicatesFullKeys []string      `protobuf:"bytes,6,rep,name=duplicates_full_keys,json=duplicatesFullKeys,proto3" json:"duplicates_full_keys,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_a965bd191a7e837c, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetPortId() string {
	if m != nil {
		return m.PortId
	}
	return ""
}

func (m *GenesisState) GetMaps() []*Map {
	if m != nil {
		return m.Maps
	}
	return nil
}

func (m *GenesisState) GetFullKeys() []string {
	if m != nil {
		return m.FullKeys
	}
	return nil
}

func (m *GenesisState) GetValues() []*ValueStore {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *GenesisState) GetDuplicatesFullKeys() []string {
	if m != nil {
		return m.DuplicatesFullKeys
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "maps.GenesisState")
}

func init() { proto.RegisterFile("maps/genesis.proto", fileDescriptor_a965bd191a7e837c) }

var fileDescriptor_a965bd191a7e837c = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0x41, 0x4b, 0xf3, 0x30,
	0x18, 0xc7, 0x9b, 0x77, 0x7b, 0xab, 0xcd, 0x26, 0x68, 0x18, 0x58, 0x3a, 0x8c, 0xc5, 0x53, 0xf1,
	0xd0, 0xca, 0xc4, 0x2f, 0xb0, 0x83, 0x22, 0x22, 0x48, 0x07, 0x1e, 0xbc, 0x94, 0x74, 0x8d, 0xb5,
	0xd8, 0x2d, 0xa1, 0x49, 0x65, 0xfd, 0x16, 0x7e, 0xac, 0x1d, 0x77, 0xf4, 0x24, 0xd2, 0x1e, 0xfd,
	0x12, 0x92, 0xa4, 0xb0, 0x5b, 0x9e, 0xdf, 0xf3, 0xff, 0xff, 0x08, 0x0f, 0x44, 0x2b, 0xc2, 0x45,
	0x94, 0xd3, 0x35, 0x15, 0x85, 0x08, 0x79, 0xc5, 0x24, 0x43, 0x43, 0xc5, 0xbc, 0x49, 0xce, 0x72,
	0xa6, 0x41, 0xa4, 0x5e, 0x66, 0xe7, 0x9d, 0xe8, 0x3c, 0x27, 0x15, 0x59, 0xf5, 0x71, 0xef, 0x48,
	0x23, 0xb9, 0x31, 0xe3, 0xc5, 0x2f, 0x80, 0xe3, 0x3b, 0xe3, 0x5b, 0x48, 0x22, 0x29, 0xba, 0x84,
	0xb6, 0xc9, 0xbb, 0xc0, 0x07, 0xc1, 0x68, 0x36, 0x0e, 0x55, 0x21, 0x7c, 0xd2, 0x6c, 0x3e, 0xdc,
	0x7e, 0x9f, 0x5b, 0x71, 0x9f, 0x40, 0xa7, 0xf0, 0x80, 0xb3, 0x4a, 0x26, 0x45, 0xe6, 0xfe, 0xf3,
	0x41, 0xe0, 0xc4, 0xb6, 0x1a, 0xef, 0x33, 0x74, 0x06, 0xf5, 0xaf, 0xdc, 0x81, 0x3f, 0x08, 0x46,
	0x33, 0xc7, 0x28, 0x1e, 0x09, 0x8f, 0x35, 0x46, 0x53, 0xe8, 0xbc, 0xd6, 0x65, 0x99, 0xbc, 0xd3,
	0x46, 0xb8, 0x43, 0x7f, 0x10, 0x38, 0xf1, 0xa1, 0x02, 0x0f, 0xb4, 0x11, 0x28, 0x80, 0xf6, 0x07,
	0x29, 0x6b, 0x2a, 0xdc, 0xff, 0xba, 0x7d, 0x6c, 0xda, 0xcf, 0x8a, 0x2d, 0x24, 0xab, 0x68, 0xdc,
	0xef, 0xd1, 0x15, 0x9c, 0x64, 0x35, 0x2f, 0x8b, 0x25, 0x91, 0x54, 0x24, 0x7b, 0xa3, 0xad, 0x8d,
	0x68, 0xbf, 0xbb, 0xed, 0xdd, 0xf3, 0x9b, 0x6d, 0x8b, 0xc1, 0xae, 0xc5, 0xe0, 0xa7, 0xc5, 0xe0,
	0xb3, 0xc3, 0xd6, 0xae, 0xc3, 0xd6, 0x57, 0x87, 0xad, 0x97, 0x69, 0x5a, 0xc8, 0x94, 0x64, 0x39,
	0x15, 0xcb, 0x37, 0x52, 0xac, 0xa3, 0x4d, 0x64, 0xee, 0xd4, 0x70, 0x2a, 0x52, 0x5b, 0xdf, 0xea,
	0xfa, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xc5, 0x8b, 0xb2, 0x7f, 0x01, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DuplicatesFullKeys) > 0 {
		for iNdEx := len(m.DuplicatesFullKeys) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.DuplicatesFullKeys[iNdEx])
			copy(dAtA[i:], m.DuplicatesFullKeys[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.DuplicatesFullKeys[iNdEx])))
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Values) > 0 {
		for iNdEx := len(m.Values) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Values[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.FullKeys) > 0 {
		for iNdEx := len(m.FullKeys) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.FullKeys[iNdEx])
			copy(dAtA[i:], m.FullKeys[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.FullKeys[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Maps) > 0 {
		for iNdEx := len(m.Maps) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Maps[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.PortId) > 0 {
		i -= len(m.PortId)
		copy(dAtA[i:], m.PortId)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.PortId)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = len(m.PortId)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Maps) > 0 {
		for _, e := range m.Maps {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.FullKeys) > 0 {
		for _, s := range m.FullKeys {
			l = len(s)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Values) > 0 {
		for _, e := range m.Values {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.DuplicatesFullKeys) > 0 {
		for _, s := range m.DuplicatesFullKeys {
			l = len(s)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PortId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PortId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Maps", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Maps = append(m.Maps, &Map{})
			if err := m.Maps[len(m.Maps)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FullKeys", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FullKeys = append(m.FullKeys, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Values", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Values = append(m.Values, &ValueStore{})
			if err := m.Values[len(m.Values)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DuplicatesFullKeys", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DuplicatesFullKeys = append(m.DuplicatesFullKeys, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
