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

//A URI object defines a uri and subasset uri for a badge and its subbadges. Designed to save storage and avoid reused text and common patterns.
type UriObject struct {
	DecodeScheme uint64 `protobuf:"varint,1,opt,name=decodeScheme,proto3" json:"decodeScheme,omitempty"`
	Scheme       uint64 `protobuf:"varint,2,opt,name=scheme,proto3" json:"scheme,omitempty"`
	Uri          string `protobuf:"bytes,3,opt,name=uri,proto3" json:"uri,omitempty"`
	//The four fields below are used to convert the uri from above to the subasset URI.
	IdxRangeToRemove       *IdRange `protobuf:"bytes,4,opt,name=idxRangeToRemove,proto3" json:"idxRangeToRemove,omitempty"`
	InsertSubassetBytesIdx uint64   `protobuf:"varint,5,opt,name=insertSubassetBytesIdx,proto3" json:"insertSubassetBytesIdx,omitempty"`
	BytesToInsert          string   `protobuf:"bytes,6,opt,name=bytesToInsert,proto3" json:"bytesToInsert,omitempty"`
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

func (m *UriObject) GetDecodeScheme() uint64 {
	if m != nil {
		return m.DecodeScheme
	}
	return 0
}

func (m *UriObject) GetScheme() uint64 {
	if m != nil {
		return m.Scheme
	}
	return 0
}

func (m *UriObject) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
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

func (m *UriObject) GetBytesToInsert() string {
	if m != nil {
		return m.BytesToInsert
	}
	return ""
}

func (m *UriObject) GetInsertIdIdx() uint64 {
	if m != nil {
		return m.InsertIdIdx
	}
	return 0
}

func init() {
	proto.RegisterType((*UriObject)(nil), "bitbadges.bitbadgeschain.badges.UriObject")
}

func init() { proto.RegisterFile("badges/uris.proto", fileDescriptor_f9aa3e50e83f5b03) }

var fileDescriptor_f9aa3e50e83f5b03 = []byte{
	// 301 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0x41, 0x6a, 0xeb, 0x30,
	0x14, 0x8c, 0x92, 0xfc, 0x7c, 0xa2, 0xb4, 0x90, 0xaa, 0x10, 0x4c, 0x17, 0xaa, 0x09, 0x5d, 0x78,
	0x25, 0x43, 0x0a, 0x3d, 0x40, 0x76, 0x86, 0x42, 0xc1, 0x49, 0x37, 0xdd, 0x59, 0xd6, 0x23, 0x51,
	0x21, 0x56, 0x90, 0xe4, 0xe2, 0xdc, 0xa2, 0x37, 0xe9, 0x35, 0xba, 0xcc, 0xb2, 0xcb, 0x62, 0x5f,
	0xa4, 0x58, 0x36, 0xa1, 0xa6, 0x94, 0xee, 0x46, 0x33, 0x6f, 0x46, 0xf3, 0x78, 0xf8, 0x82, 0x27,
	0x62, 0x03, 0x26, 0xcc, 0xb5, 0x34, 0x6c, 0xaf, 0x95, 0x55, 0xe4, 0x9a, 0x4b, 0xdb, 0xb0, 0xec,
	0x84, 0xd2, 0x6d, 0x22, 0x33, 0xd6, 0xe0, 0xab, 0xcb, 0xd6, 0xa3, 0x93, 0xac, 0x1e, 0x72, 0xae,
	0xf9, 0x5b, 0x1f, 0x8f, 0x1f, 0xb5, 0x7c, 0xe0, 0xcf, 0x90, 0x5a, 0x32, 0xc7, 0x67, 0x02, 0x52,
	0x25, 0x60, 0x95, 0x6e, 0x61, 0x07, 0x1e, 0xf2, 0x51, 0x30, 0x8c, 0x3b, 0x1c, 0x99, 0xe1, 0x91,
	0x69, 0xd4, 0xbe, 0x53, 0xdb, 0x17, 0x99, 0xe2, 0x41, 0xae, 0xa5, 0x37, 0xf0, 0x51, 0x30, 0x8e,
	0x6b, 0x48, 0xd6, 0x78, 0x2a, 0x45, 0x11, 0xd7, 0xdf, 0xad, 0x55, 0x0c, 0x3b, 0xf5, 0x02, 0xde,
	0xd0, 0x47, 0xc1, 0x64, 0x11, 0xb0, 0x3f, 0xca, 0xb2, 0x48, 0x38, 0x5f, 0xfc, 0x23, 0x81, 0xdc,
	0xe1, 0x99, 0xcc, 0x0c, 0x68, 0xbb, 0xca, 0x79, 0x62, 0x0c, 0xd8, 0xe5, 0xc1, 0x82, 0x89, 0x44,
	0xe1, 0xfd, 0x73, 0x7d, 0x7e, 0x51, 0xc9, 0x0d, 0x3e, 0xe7, 0x35, 0x5e, 0xab, 0xc8, 0x0d, 0x78,
	0x23, 0xd7, 0xb4, 0x4b, 0x12, 0x1f, 0x4f, 0x1a, 0x7f, 0x24, 0xea, 0xc8, 0xff, 0x2e, 0xf2, 0x3b,
	0xb5, 0xbc, 0x7f, 0x2f, 0x29, 0x3a, 0x96, 0x14, 0x7d, 0x96, 0x14, 0xbd, 0x56, 0xb4, 0x77, 0xac,
	0x68, 0xef, 0xa3, 0xa2, 0xbd, 0xa7, 0xc5, 0x46, 0xda, 0x6d, 0xce, 0x59, 0xaa, 0x76, 0xe1, 0x69,
	0xab, 0xb0, 0xbb, 0x5f, 0x58, 0x84, 0x2d, 0x6f, 0x0f, 0x7b, 0x30, 0x7c, 0xe4, 0xce, 0x70, 0xfb,
	0x15, 0x00, 0x00, 0xff, 0xff, 0xf8, 0xd6, 0x14, 0xb2, 0xd1, 0x01, 0x00, 0x00,
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
		dAtA[i] = 0x1a
	}
	if m.Scheme != 0 {
		i = encodeVarintUris(dAtA, i, uint64(m.Scheme))
		i--
		dAtA[i] = 0x10
	}
	if m.DecodeScheme != 0 {
		i = encodeVarintUris(dAtA, i, uint64(m.DecodeScheme))
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
	if m.DecodeScheme != 0 {
		n += 1 + sovUris(uint64(m.DecodeScheme))
	}
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
				return fmt.Errorf("proto: wrong wireType = %d for field DecodeScheme", wireType)
			}
			m.DecodeScheme = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DecodeScheme |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uri", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
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
				return ErrInvalidLengthUris
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUris
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uri = string(dAtA[iNdEx:postIndex])
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
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUris
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
				return ErrInvalidLengthUris
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUris
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BytesToInsert = string(dAtA[iNdEx:postIndex])
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
