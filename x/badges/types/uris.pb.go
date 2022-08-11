// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: badges/uris.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type UriObject struct {
	Scheme uint64 `protobuf:"varint,1,opt,name=scheme,proto3" json:"scheme,omitempty"`
	//We can add more URI backwards compatible space savers in the future here like .com or .io
	Uri                    []byte   `protobuf:"bytes,2,opt,name=uri,proto3" json:"uri,omitempty"`
	IdxRangeToRemove       *IdRange `protobuf:"bytes,4,opt,name=idxRangeToRemove,proto3" json:"idxRangeToRemove,omitempty"`
	InsertSubassetBytesIdx uint64   `protobuf:"varint,5,opt,name=insertSubassetBytesIdx,proto3" json:"insertSubassetBytesIdx,omitempty"`
	BytesToInsert          []byte   `protobuf:"bytes,6,opt,name=bytesToInsert,proto3" json:"bytesToInsert,omitempty"`
	InsertIdIdx            uint64   `protobuf:"varint,7,opt,name=insertIdIdx,proto3" json:"insertIdIdx,omitempty"`
}

func (m *UriObject) Reset()         { *m = UriObject{} }
func (m *UriObject) String() string { return proto.CompactTextString(m) }
func (*UriObject) ProtoMessage()    {}
func (*UriObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9aa3e50e83f5b03, []int{0}
}
func (m *UriObject) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UriObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UriObject.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UriObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UriObject.Merge(m, src)
}
func (m *UriObject) XXX_Size() int {
	return m.Size()
}
func (m *UriObject) XXX_DiscardUnknown() {
	xxx_messageInfo_UriObject.DiscardUnknown(m)
}

var xxx_messageInfo_UriObject proto.InternalMessageInfo

func (m *UriObject) GetScheme() uint64 {
	if m != nil {
		return m.Scheme
	}
	return 0
}

func (m *UriObject) GetUri() []byte {
	if m != nil {
		return m.Uri
	}
	return nil
}

func (m *UriObject) GetIdxRangeToRemove() *IdRange {
	if m != nil {
		return m.IdxRangeToRemove
	}
	return nil
}

func (m *UriObject) GetInsertSubassetBytesIdx() uint64 {
	if m != nil {
		return m.InsertSubassetBytesIdx
	}
	return 0
}

func (m *UriObject) GetBytesToInsert() []byte {
	if m != nil {
		return m.BytesToInsert
	}
	return nil
}

func (m *UriObject) GetInsertIdIdx() uint64 {
	if m != nil {
		return m.InsertIdIdx
	}
	return 0
}

func init() {
	proto.RegisterType((*UriObject)(nil), "trevormil.bitbadgeschain.badges.UriObject")
}

func init() { proto.RegisterFile("badges/uris.proto", fileDescriptor_f9aa3e50e83f5b03) }

var fileDescriptor_f9aa3e50e83f5b03 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0x97, 0x39, 0x27, 0x66, 0x0a, 0x33, 0xc2, 0x28, 0x1e, 0x62, 0x11, 0x0f, 0x3d, 0xa5,
	0x30, 0xc1, 0x1f, 0xb0, 0x5b, 0x41, 0x10, 0x6a, 0xbd, 0x78, 0x6b, 0xda, 0x8f, 0x36, 0x62, 0x9b,
	0x91, 0xa4, 0xa3, 0xfb, 0x0d, 0x5e, 0xfc, 0x59, 0x1e, 0x77, 0xf4, 0x28, 0xed, 0x1f, 0x91, 0xa6,
	0x45, 0x1c, 0x22, 0xde, 0xde, 0x7c, 0x5f, 0xde, 0x87, 0xf7, 0x7b, 0xf1, 0x19, 0x8f, 0xd3, 0x0c,
	0xb4, 0x5f, 0x29, 0xa1, 0xd9, 0x5a, 0x49, 0x23, 0xc9, 0xa5, 0x51, 0xb0, 0x91, 0xaa, 0x10, 0x2f,
	0x8c, 0x0b, 0xd3, 0xef, 0x93, 0x3c, 0x16, 0x25, 0xeb, 0xf5, 0xc5, 0xf9, 0xe0, 0x51, 0x71, 0x99,
	0xc1, 0xe0, 0xba, 0x7a, 0x1d, 0xe3, 0xe3, 0x47, 0x25, 0xee, 0xf9, 0x33, 0x24, 0x86, 0x2c, 0xf0,
	0x54, 0x27, 0x39, 0x14, 0xe0, 0x20, 0x17, 0x79, 0x93, 0x70, 0x78, 0x91, 0x39, 0x3e, 0xa8, 0x94,
	0x70, 0xc6, 0x2e, 0xf2, 0x4e, 0xc2, 0x4e, 0x92, 0x08, 0xcf, 0x45, 0x5a, 0x87, 0x1d, 0x2a, 0x92,
	0x21, 0x14, 0x72, 0x03, 0xce, 0xc4, 0x45, 0xde, 0x6c, 0xe9, 0xb1, 0x7f, 0x82, 0xb0, 0x20, 0xb5,
	0xbe, 0xf0, 0x17, 0x81, 0xdc, 0xe2, 0x85, 0x28, 0x35, 0x28, 0xf3, 0x50, 0xf1, 0x58, 0x6b, 0x30,
	0xab, 0xad, 0x01, 0x1d, 0xa4, 0xb5, 0x73, 0x68, 0xf3, 0xfc, 0xb1, 0x25, 0xd7, 0xf8, 0x94, 0x77,
	0x3a, 0x92, 0x81, 0xfd, 0xe0, 0x4c, 0x6d, 0xd2, 0xfd, 0x21, 0x71, 0xf1, 0xac, 0xf7, 0x07, 0x69,
	0x87, 0x3c, 0xb2, 0xc8, 0x9f, 0xa3, 0xd5, 0xdd, 0x7b, 0x43, 0xd1, 0xae, 0xa1, 0xe8, 0xb3, 0xa1,
	0xe8, 0xad, 0xa5, 0xa3, 0x5d, 0x4b, 0x47, 0x1f, 0x2d, 0x1d, 0x3d, 0x2d, 0x33, 0x61, 0xf2, 0x8a,
	0xb3, 0x44, 0x16, 0xfe, 0xf7, 0x7d, 0xfe, 0xfe, 0x7d, 0x7e, 0xed, 0x0f, 0x15, 0x9b, 0xed, 0x1a,
	0x34, 0x9f, 0xda, 0x8a, 0x6f, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x63, 0x66, 0xd0, 0x9a, 0xad,
	0x01, 0x00, 0x00,
}

func (m *UriObject) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UriObject) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UriObject) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.InsertIdIdx != 0 {
		i = encodeVarintUris(dAtA, i, uint64(m.InsertIdIdx))
		i--
		dAtA[i] = 0x38
	}
	if len(m.BytesToInsert) > 0 {
		i -= len(m.BytesToInsert)
		copy(dAtA[i:], m.BytesToInsert)
		i = encodeVarintUris(dAtA, i, uint64(len(m.BytesToInsert)))
		i--
		dAtA[i] = 0x32
	}
	if m.InsertSubassetBytesIdx != 0 {
		i = encodeVarintUris(dAtA, i, uint64(m.InsertSubassetBytesIdx))
		i--
		dAtA[i] = 0x28
	}
	if m.IdxRangeToRemove != nil {
		{
			size, err := m.IdxRangeToRemove.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintUris(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.Uri) > 0 {
		i -= len(m.Uri)
		copy(dAtA[i:], m.Uri)
		i = encodeVarintUris(dAtA, i, uint64(len(m.Uri)))
		i--
		dAtA[i] = 0x12
	}
	if m.Scheme != 0 {
		i = encodeVarintUris(dAtA, i, uint64(m.Scheme))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintUris(dAtA []byte, offset int, v uint64) int {
	offset -= sovUris(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UriObject) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Scheme != 0 {
		n += 1 + sovUris(uint64(m.Scheme))
	}
	l = len(m.Uri)
	if l > 0 {
		n += 1 + l + sovUris(uint64(l))
	}
	if m.IdxRangeToRemove != nil {
		l = m.IdxRangeToRemove.Size()
		n += 1 + l + sovUris(uint64(l))
	}
	if m.InsertSubassetBytesIdx != 0 {
		n += 1 + sovUris(uint64(m.InsertSubassetBytesIdx))
	}
	l = len(m.BytesToInsert)
	if l > 0 {
		n += 1 + l + sovUris(uint64(l))
	}
	if m.InsertIdIdx != 0 {
		n += 1 + sovUris(uint64(m.InsertIdIdx))
	}
	return n
}

func sovUris(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozUris(x uint64) (n int) {
	return sovUris(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UriObject) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUris
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
			return fmt.Errorf("proto: UriObject: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UriObject: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Scheme", wireType)
			}
			m.Scheme = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Scheme |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uri", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthUris
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthUris
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uri = append(m.Uri[:0], dAtA[iNdEx:postIndex]...)
			if m.Uri == nil {
				m.Uri = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IdxRangeToRemove", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
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
				return ErrInvalidLengthUris
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthUris
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.IdxRangeToRemove == nil {
				m.IdxRangeToRemove = &IdRange{}
			}
			if err := m.IdxRangeToRemove.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InsertSubassetBytesIdx", wireType)
			}
			m.InsertSubassetBytesIdx = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InsertSubassetBytesIdx |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BytesToInsert", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthUris
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthUris
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BytesToInsert = append(m.BytesToInsert[:0], dAtA[iNdEx:postIndex]...)
			if m.BytesToInsert == nil {
				m.BytesToInsert = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InsertIdIdx", wireType)
			}
			m.InsertIdIdx = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InsertIdIdx |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipUris(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUris
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
func skipUris(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUris
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
					return 0, ErrIntOverflowUris
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
					return 0, ErrIntOverflowUris
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
				return 0, ErrInvalidLengthUris
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUris
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUris
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUris        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUris          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUris = fmt.Errorf("proto: unexpected end of group")
)
