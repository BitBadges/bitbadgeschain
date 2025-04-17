// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: maps/balances.proto

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

// The UintRange is a range of IDs from some start to some end (inclusive).
//
// uintRanges are one of the core types used in the BitBadgesChain module.
// They are used for everything from badge IDs to time ranges to min/max balance amounts.
//
// See the BitBadges documentation for more information.
type UintRange struct {
	// The starting value of the range (inclusive).
	Start Uint `protobuf:"bytes,1,opt,name=start,proto3,customtype=Uint" json:"start"`
	// The ending value of the range (inclusive).
	End Uint `protobuf:"bytes,2,opt,name=end,proto3,customtype=Uint" json:"end"`
}

func (m *UintRange) Reset()         { *m = UintRange{} }
func (m *UintRange) String() string { return proto.CompactTextString(m) }
func (*UintRange) ProtoMessage()    {}
func (*UintRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_f712f54dd13b7176, []int{0}
}
func (m *UintRange) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UintRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UintRange.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UintRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UintRange.Merge(m, src)
}
func (m *UintRange) XXX_Size() int {
	return m.Size()
}
func (m *UintRange) XXX_DiscardUnknown() {
	xxx_messageInfo_UintRange.DiscardUnknown(m)
}

var xxx_messageInfo_UintRange proto.InternalMessageInfo

func init() {
	proto.RegisterType((*UintRange)(nil), "maps.UintRange")
}

func init() { proto.RegisterFile("maps/balances.proto", fileDescriptor_f712f54dd13b7176) }

var fileDescriptor_f712f54dd13b7176 = []byte{
	// 200 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xce, 0x4d, 0x2c, 0x28,
	0xd6, 0x4f, 0x4a, 0xcc, 0x49, 0xcc, 0x4b, 0x4e, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x62, 0x01, 0x09, 0x4a, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0x05, 0xf4, 0x41, 0x2c, 0x88, 0x9c,
	0x94, 0x20, 0x58, 0x43, 0x41, 0x62, 0x51, 0x62, 0x2e, 0x54, 0xb9, 0x92, 0x3f, 0x17, 0x67, 0x68,
	0x66, 0x5e, 0x49, 0x50, 0x62, 0x5e, 0x7a, 0xaa, 0x90, 0x12, 0x17, 0x6b, 0x71, 0x49, 0x62, 0x51,
	0x89, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0xa7, 0x13, 0xcf, 0x89, 0x7b, 0xf2, 0x0c, 0xb7, 0xee, 0xc9,
	0xb3, 0x80, 0x55, 0x40, 0xa4, 0x84, 0xe4, 0xb8, 0x98, 0x53, 0xf3, 0x52, 0x24, 0x98, 0xb0, 0xa8,
	0x00, 0x49, 0x38, 0x79, 0x9d, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72,
	0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x41,
	0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x52, 0x66, 0x49, 0x52, 0x62,
	0x4a, 0x7a, 0x6a, 0x31, 0x82, 0x95, 0x9c, 0x91, 0x98, 0x99, 0xa7, 0x5f, 0xa1, 0x0f, 0x76, 0x63,
	0x49, 0x65, 0x41, 0x6a, 0x71, 0x12, 0x1b, 0xd8, 0x8d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x9f, 0x38, 0x7a, 0xb8, 0xe9, 0x00, 0x00, 0x00,
}

func (m *UintRange) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UintRange) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UintRange) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.End.Size()
		i -= size
		if _, err := m.End.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBalances(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Start.Size()
		i -= size
		if _, err := m.Start.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBalances(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintBalances(dAtA []byte, offset int, v uint64) int {
	offset -= sovBalances(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UintRange) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Start.Size()
	n += 1 + l + sovBalances(uint64(l))
	l = m.End.Size()
	n += 1 + l + sovBalances(uint64(l))
	return n
}

func sovBalances(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBalances(x uint64) (n int) {
	return sovBalances(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UintRange) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBalances
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
			return fmt.Errorf("proto: UintRange: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UintRange: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
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
				return ErrInvalidLengthBalances
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBalances
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Start.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field End", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
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
				return ErrInvalidLengthBalances
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBalances
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.End.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBalances(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBalances
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
func skipBalances(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBalances
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
					return 0, ErrIntOverflowBalances
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
					return 0, ErrIntOverflowBalances
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
				return 0, ErrInvalidLengthBalances
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBalances
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBalances
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBalances        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBalances          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBalances = fmt.Errorf("proto: unexpected end of group")
)
