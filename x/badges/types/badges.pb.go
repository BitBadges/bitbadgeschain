// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: badges/badges.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
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

// BitBadge defines a badge type. Think of this like the smart contract definition.
type BitBadge struct {
	// id defines the unique identifier of the Badge classification, similar to the contract address of ERC721
	// starts at 0 and increments by 1 each badge
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// uri object for the badge uri and subasset uris stored off chain. Stored in a special UriObject that attemtps to save space and avoid reused plaintext storage such as http:// and duplicate text for uri and subasset uris
	// data returned should corresponds to the Badge standard defined.
	Uri *UriObject `protobuf:"bytes,2,opt,name=uri,proto3" json:"uri,omitempty"`
	// these bytes can be used to store anything on-chain about the badge. This can be updatable or not depending on the permissions set.
	// Max 256 bytes allowed
	ArbitraryBytes string `protobuf:"bytes,3,opt,name=arbitraryBytes,proto3" json:"arbitraryBytes,omitempty"`
	// manager address of the class; can have special permissions; is used as the reserve address for the assets
	Manager uint64 `protobuf:"varint,4,opt,name=manager,proto3" json:"manager,omitempty"`
	//Store permissions packed in a uint where the bits correspond to permissions from left to right; leading zeroes are applied and any future additions will be appended to the right. See types/permissions.go
	Permissions uint64 `protobuf:"varint,5,opt,name=permissions,proto3" json:"permissions,omitempty"`
	//FreezeRanges defines what addresses are frozen or unfrozen. If permissions.FrozenByDefault is false, this is used for frozen addresses. If true, this is used for unfrozen addresses.
	FreezeRanges []*IdRange `protobuf:"bytes,10,rep,name=freezeRanges,proto3" json:"freezeRanges,omitempty"`
	// Starts at 0. Each subasset created will incrementally have an increasing ID #. Can't overflow.
	NextSubassetId uint64 `protobuf:"varint,12,opt,name=nextSubassetId,proto3" json:"nextSubassetId,omitempty"`
	//Subasset supplys are stored if the subasset supply != default. Balance => SubbadgeIdRange map
	SubassetSupplys []*BalanceObject `protobuf:"bytes,13,rep,name=subassetSupplys,proto3" json:"subassetSupplys,omitempty"`
	//Default subasset supply. If == 0, we assume == 1.
	DefaultSubassetSupply uint64 `protobuf:"varint,14,opt,name=defaultSubassetSupply,proto3" json:"defaultSubassetSupply,omitempty"`
	//Defines what standard this badge should implement. Must obey the rules of that standard.
	Standard uint64 `protobuf:"varint,15,opt,name=standard,proto3" json:"standard,omitempty"`
}

func (m *BitBadge) Reset()         { *m = BitBadge{} }
func (m *BitBadge) String() string { return proto.CompactTextString(m) }
func (*BitBadge) ProtoMessage()    {}
func (*BitBadge) Descriptor() ([]byte, []int) {
	return fileDescriptor_71eab594b779f631, []int{0}
}
func (m *BitBadge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BitBadge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BitBadge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BitBadge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BitBadge.Merge(m, src)
}
func (m *BitBadge) XXX_Size() int {
	return m.Size()
}
func (m *BitBadge) XXX_DiscardUnknown() {
	xxx_messageInfo_BitBadge.DiscardUnknown(m)
}

var xxx_messageInfo_BitBadge proto.InternalMessageInfo

func (m *BitBadge) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *BitBadge) GetUri() *UriObject {
	if m != nil {
		return m.Uri
	}
	return nil
}

func (m *BitBadge) GetArbitraryBytes() string {
	if m != nil {
		return m.ArbitraryBytes
	}
	return ""
}

func (m *BitBadge) GetManager() uint64 {
	if m != nil {
		return m.Manager
	}
	return 0
}

func (m *BitBadge) GetPermissions() uint64 {
	if m != nil {
		return m.Permissions
	}
	return 0
}

func (m *BitBadge) GetFreezeRanges() []*IdRange {
	if m != nil {
		return m.FreezeRanges
	}
	return nil
}

func (m *BitBadge) GetNextSubassetId() uint64 {
	if m != nil {
		return m.NextSubassetId
	}
	return 0
}

func (m *BitBadge) GetSubassetSupplys() []*BalanceObject {
	if m != nil {
		return m.SubassetSupplys
	}
	return nil
}

func (m *BitBadge) GetDefaultSubassetSupply() uint64 {
	if m != nil {
		return m.DefaultSubassetSupply
	}
	return 0
}

func (m *BitBadge) GetStandard() uint64 {
	if m != nil {
		return m.Standard
	}
	return 0
}

func init() {
	proto.RegisterType((*BitBadge)(nil), "bitbadges.bitbadgeschain.badges.BitBadge")
}

func init() { proto.RegisterFile("badges/badges.proto", fileDescriptor_71eab594b779f631) }

var fileDescriptor_71eab594b779f631 = []byte{
	// 384 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xc1, 0x6a, 0xe3, 0x30,
	0x10, 0x8d, 0xe3, 0xec, 0x6e, 0x56, 0xc9, 0x26, 0xac, 0x96, 0x05, 0xad, 0x0f, 0x5e, 0xb3, 0x87,
	0xc5, 0xf4, 0x60, 0x43, 0xda, 0x63, 0x4f, 0xbe, 0x05, 0x02, 0x05, 0x87, 0x42, 0xe9, 0x4d, 0xb2,
	0x14, 0x47, 0xc5, 0x91, 0x8d, 0x24, 0x43, 0xdc, 0xaf, 0xe8, 0x67, 0xf5, 0x98, 0x63, 0x4f, 0xa5,
	0x24, 0x3f, 0x52, 0x22, 0x3b, 0xa6, 0x09, 0x85, 0x9c, 0x3c, 0xf3, 0xde, 0xbc, 0x79, 0xcf, 0x62,
	0xc0, 0x2f, 0x82, 0x69, 0xca, 0x54, 0x58, 0x7f, 0x82, 0x42, 0xe6, 0x3a, 0x87, 0x7f, 0x09, 0xd7,
	0x0d, 0xd0, 0x56, 0xc9, 0x12, 0x73, 0x11, 0xd4, 0xb5, 0xf3, 0x27, 0xcd, 0xf3, 0x34, 0x63, 0xa1,
	0x19, 0x27, 0xe5, 0x22, 0xc4, 0xa2, 0xaa, 0xb5, 0xce, 0x61, 0xa1, 0xc4, 0xa2, 0x5d, 0xe8, 0xfc,
	0x6c, 0xc0, 0x52, 0xf2, 0x06, 0xfa, 0xf7, 0x6a, 0x83, 0x7e, 0xc4, 0x75, 0xb4, 0x27, 0xe0, 0x08,
	0x74, 0x39, 0x45, 0x96, 0x67, 0xf9, 0xbd, 0xb8, 0xcb, 0x29, 0xbc, 0x06, 0x76, 0x29, 0x39, 0xea,
	0x7a, 0x96, 0x3f, 0x98, 0x5c, 0x04, 0x67, 0xe2, 0x04, 0xb7, 0x92, 0xdf, 0x90, 0x07, 0x96, 0xe8,
	0x78, 0x2f, 0x83, 0xff, 0xc1, 0x08, 0x4b, 0xc2, 0xb5, 0xc4, 0xb2, 0x8a, 0x2a, 0xcd, 0x14, 0xb2,
	0x3d, 0xcb, 0xff, 0x1e, 0x9f, 0xa0, 0x10, 0x81, 0x6f, 0x2b, 0x2c, 0x70, 0xca, 0x24, 0xea, 0x19,
	0xeb, 0x43, 0x0b, 0x3d, 0x30, 0x28, 0x98, 0x5c, 0x71, 0xa5, 0x78, 0x2e, 0x14, 0xfa, 0x62, 0xd8,
	0x8f, 0x10, 0x9c, 0x81, 0xe1, 0x42, 0x32, 0xf6, 0xc8, 0x62, 0xf3, 0x9f, 0x08, 0x78, 0xb6, 0x3f,
	0x98, 0xf8, 0x67, 0xa3, 0x4e, 0xa9, 0x11, 0xc4, 0x47, 0xea, 0x7d, 0x62, 0xc1, 0xd6, 0x7a, 0x5e,
	0x12, 0xac, 0x14, 0xd3, 0x53, 0x8a, 0x86, 0xc6, 0xf2, 0x04, 0x85, 0x77, 0x60, 0xac, 0x9a, 0x6e,
	0x5e, 0x16, 0x45, 0x56, 0x29, 0xf4, 0xc3, 0x18, 0x07, 0x67, 0x8d, 0x23, 0x9c, 0x61, 0x91, 0xb0,
	0xe6, 0x9d, 0x4e, 0xd7, 0xc0, 0x2b, 0xf0, 0x9b, 0xb2, 0x05, 0x2e, 0xb3, 0xd6, 0xae, 0x66, 0xd0,
	0xc8, 0x04, 0xf9, 0x9c, 0x84, 0x0e, 0xe8, 0x2b, 0x8d, 0x05, 0xc5, 0x92, 0xa2, 0xb1, 0x19, 0x6c,
	0xfb, 0x68, 0xf6, 0xbc, 0x75, 0xad, 0xcd, 0xd6, 0xb5, 0xde, 0xb6, 0xae, 0xf5, 0xb4, 0x73, 0x3b,
	0x9b, 0x9d, 0xdb, 0x79, 0xd9, 0xb9, 0x9d, 0xfb, 0x49, 0xca, 0xf5, 0xb2, 0x24, 0x41, 0x92, 0xaf,
	0xc2, 0x36, 0x6c, 0x78, 0x1c, 0x3b, 0x5c, 0x37, 0x27, 0x19, 0xea, 0xaa, 0x60, 0x8a, 0x7c, 0x35,
	0x57, 0x73, 0xf9, 0x1e, 0x00, 0x00, 0xff, 0xff, 0x94, 0x3c, 0x4b, 0x04, 0xb0, 0x02, 0x00, 0x00,
}

func (m *BitBadge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BitBadge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BitBadge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Standard != 0 {
		i = encodeVarintBadges(dAtA, i, uint64(m.Standard))
		i--
		dAtA[i] = 0x78
	}
	if m.DefaultSubassetSupply != 0 {
		i = encodeVarintBadges(dAtA, i, uint64(m.DefaultSubassetSupply))
		i--
		dAtA[i] = 0x70
	}
	if len(m.SubassetSupplys) > 0 {
		for iNdEx := len(m.SubassetSupplys) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SubassetSupplys[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBadges(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x6a
		}
	}
	if m.NextSubassetId != 0 {
		i = encodeVarintBadges(dAtA, i, uint64(m.NextSubassetId))
		i--
		dAtA[i] = 0x60
	}
	if len(m.FreezeRanges) > 0 {
		for iNdEx := len(m.FreezeRanges) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FreezeRanges[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBadges(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x52
		}
	}
	if m.Permissions != 0 {
		i = encodeVarintBadges(dAtA, i, uint64(m.Permissions))
		i--
		dAtA[i] = 0x28
	}
	if m.Manager != 0 {
		i = encodeVarintBadges(dAtA, i, uint64(m.Manager))
		i--
		dAtA[i] = 0x20
	}
	if len(m.ArbitraryBytes) > 0 {
		i -= len(m.ArbitraryBytes)
		copy(dAtA[i:], m.ArbitraryBytes)
		i = encodeVarintBadges(dAtA, i, uint64(len(m.ArbitraryBytes)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Uri != nil {
		{
			size, err := m.Uri.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBadges(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintBadges(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintBadges(dAtA []byte, offset int, v uint64) int {
	offset -= sovBadges(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BitBadge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovBadges(uint64(m.Id))
	}
	if m.Uri != nil {
		l = m.Uri.Size()
		n += 1 + l + sovBadges(uint64(l))
	}
	l = len(m.ArbitraryBytes)
	if l > 0 {
		n += 1 + l + sovBadges(uint64(l))
	}
	if m.Manager != 0 {
		n += 1 + sovBadges(uint64(m.Manager))
	}
	if m.Permissions != 0 {
		n += 1 + sovBadges(uint64(m.Permissions))
	}
	if len(m.FreezeRanges) > 0 {
		for _, e := range m.FreezeRanges {
			l = e.Size()
			n += 1 + l + sovBadges(uint64(l))
		}
	}
	if m.NextSubassetId != 0 {
		n += 1 + sovBadges(uint64(m.NextSubassetId))
	}
	if len(m.SubassetSupplys) > 0 {
		for _, e := range m.SubassetSupplys {
			l = e.Size()
			n += 1 + l + sovBadges(uint64(l))
		}
	}
	if m.DefaultSubassetSupply != 0 {
		n += 1 + sovBadges(uint64(m.DefaultSubassetSupply))
	}
	if m.Standard != 0 {
		n += 1 + sovBadges(uint64(m.Standard))
	}
	return n
}

func sovBadges(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBadges(x uint64) (n int) {
	return sovBadges(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BitBadge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBadges
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
			return fmt.Errorf("proto: BitBadge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BitBadge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uri", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
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
				return ErrInvalidLengthBadges
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBadges
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Uri == nil {
				m.Uri = &UriObject{}
			}
			if err := m.Uri.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ArbitraryBytes", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
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
				return ErrInvalidLengthBadges
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBadges
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ArbitraryBytes = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Manager", wireType)
			}
			m.Manager = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Manager |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Permissions", wireType)
			}
			m.Permissions = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Permissions |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FreezeRanges", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
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
				return ErrInvalidLengthBadges
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBadges
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FreezeRanges = append(m.FreezeRanges, &IdRange{})
			if err := m.FreezeRanges[len(m.FreezeRanges)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextSubassetId", wireType)
			}
			m.NextSubassetId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextSubassetId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubassetSupplys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
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
				return ErrInvalidLengthBadges
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBadges
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubassetSupplys = append(m.SubassetSupplys, &BalanceObject{})
			if err := m.SubassetSupplys[len(m.SubassetSupplys)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultSubassetSupply", wireType)
			}
			m.DefaultSubassetSupply = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DefaultSubassetSupply |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Standard", wireType)
			}
			m.Standard = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBadges
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Standard |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBadges(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBadges
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
func skipBadges(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBadges
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
					return 0, ErrIntOverflowBadges
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
					return 0, ErrIntOverflowBadges
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
				return 0, ErrInvalidLengthBadges
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBadges
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBadges
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBadges        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBadges          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBadges = fmt.Errorf("proto: unexpected end of group")
)
