// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: badges/balances.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

// uintRange is a range of IDs from some start to some end (inclusive).
//
// uintRanges are one of the core types used in the BitBadgesChain module.
// They are used for evrything from badge IDs to time ranges to min / max balance amounts.
type UintRange struct {
	Start Uint `protobuf:"bytes,1,opt,name=start,proto3,customtype=Uint" json:"start"`
	End   Uint `protobuf:"bytes,2,opt,name=end,proto3,customtype=Uint" json:"end"`
}

func (m *UintRange) Reset()         { *m = UintRange{} }
func (m *UintRange) String() string { return proto.CompactTextString(m) }
func (*UintRange) ProtoMessage()    {}
func (*UintRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_233d29a167e739f0, []int{0}
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

// Balance represents the balance of a badge for a specific user.
// The user amounts xAmount of a badge for the badgeID specified for the time ranges specified.
//
// Ex: User A owns x10 of badge IDs 1-10 from 1/1/2020 to 1/1/2021.
//
// If times or badgeIDs have len > 1, then the user owns all badge IDs specified for all time ranges specified.
type Balance struct {
	Amount         Uint       `protobuf:"bytes,1,opt,name=amount,proto3,customtype=Uint" json:"amount"`
	OwnershipTimes []*UintRange `protobuf:"bytes,2,rep,name=ownershipTimes,proto3" json:"ownershipTimes,omitempty"`
	BadgeIds       []*UintRange `protobuf:"bytes,3,rep,name=badgeIds,proto3" json:"badgeIds,omitempty"`
}

func (m *Balance) Reset()         { *m = Balance{} }
func (m *Balance) String() string { return proto.CompactTextString(m) }
func (*Balance) ProtoMessage()    {}
func (*Balance) Descriptor() ([]byte, []int) {
	return fileDescriptor_233d29a167e739f0, []int{1}
}
func (m *Balance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Balance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Balance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Balance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Balance.Merge(m, src)
}
func (m *Balance) XXX_Size() int {
	return m.Size()
}
func (m *Balance) XXX_DiscardUnknown() {
	xxx_messageInfo_Balance.DiscardUnknown(m)
}

var xxx_messageInfo_Balance proto.InternalMessageInfo

func (m *Balance) GetOwnershipTimes() []*UintRange {
	if m != nil {
		return m.OwnershipTimes
	}
	return nil
}

func (m *Balance) GetBadgeIds() []*UintRange {
	if m != nil {
		return m.BadgeIds
	}
	return nil
}

// InheritedBalances are a powerful feature of the BitBadges module.
// They allow a colllection to inherit the balances from another collection.
// Ex: Badges from Collection A inherits the balances from badges from Collection B.
//
// The badgeIds specified will inherit the balances from the parent collection and badges specified.
// If the total number of parent badges == 1, then all the badgeIds will inherit the balance from that parent badge.
// Otherwise, the total number of parent badges must equal the total number of badgeIds specified.
// By total number, we mean the sum of the number of badgeIds in each UintRange.
type InheritedBalance struct {
	BadgeIds           []*UintRange `protobuf:"bytes,1,rep,name=badgeIds,proto3" json:"badgeIds,omitempty"`
	ParentCollectionId Uint       `protobuf:"bytes,2,opt,name=parentCollectionId,proto3,customtype=Uint" json:"parentCollectionId"`
	ParentBadgeIds     []*UintRange `protobuf:"bytes,3,rep,name=parentBadgeIds,proto3" json:"parentBadgeIds,omitempty"`
}

func (m *InheritedBalance) Reset()         { *m = InheritedBalance{} }
func (m *InheritedBalance) String() string { return proto.CompactTextString(m) }
func (*InheritedBalance) ProtoMessage()    {}
func (*InheritedBalance) Descriptor() ([]byte, []int) {
	return fileDescriptor_233d29a167e739f0, []int{2}
}
func (m *InheritedBalance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InheritedBalance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InheritedBalance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InheritedBalance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InheritedBalance.Merge(m, src)
}
func (m *InheritedBalance) XXX_Size() int {
	return m.Size()
}
func (m *InheritedBalance) XXX_DiscardUnknown() {
	xxx_messageInfo_InheritedBalance.DiscardUnknown(m)
}

var xxx_messageInfo_InheritedBalance proto.InternalMessageInfo

func (m *InheritedBalance) GetBadgeIds() []*UintRange {
	if m != nil {
		return m.BadgeIds
	}
	return nil
}

func (m *InheritedBalance) GetParentBadgeIds() []*UintRange {
	if m != nil {
		return m.ParentBadgeIds
	}
	return nil
}

func init() {
	proto.RegisterType((*UintRange)(nil), "bitbadges.bitbadgeschain.badges.UintRange")
	proto.RegisterType((*Balance)(nil), "bitbadges.bitbadgeschain.badges.Balance")
	proto.RegisterType((*InheritedBalance)(nil), "bitbadges.bitbadgeschain.badges.InheritedBalance")
}

func init() { proto.RegisterFile("badges/balances.proto", fileDescriptor_233d29a167e739f0) }

var fileDescriptor_233d29a167e739f0 = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x41, 0x4b, 0xfb, 0x30,
	0x18, 0xc6, 0x9b, 0xed, 0xff, 0xdf, 0x34, 0x8a, 0x48, 0x54, 0x28, 0x3b, 0x64, 0xa3, 0x78, 0xd8,
	0x29, 0x81, 0x79, 0xf5, 0x54, 0xbd, 0x14, 0x14, 0x64, 0xe8, 0xc5, 0x5b, 0xda, 0x86, 0x36, 0xb0,
	0x26, 0x25, 0xc9, 0x50, 0xbf, 0x85, 0xdf, 0xca, 0x1d, 0x77, 0x14, 0x0f, 0x43, 0xb6, 0x8b, 0x1f,
	0x43, 0xb6, 0xd4, 0xa1, 0xa3, 0x20, 0xee, 0xf6, 0x92, 0xe7, 0x79, 0x7e, 0xbc, 0xcf, 0x4b, 0xe0,
	0x49, 0xcc, 0xd2, 0x8c, 0x1b, 0x1a, 0xb3, 0x11, 0x93, 0x09, 0x37, 0xa4, 0xd4, 0xca, 0x2a, 0xd4,
	0x8d, 0x85, 0x75, 0x0a, 0x59, 0x4f, 0x49, 0xce, 0x84, 0x24, 0x6e, 0xee, 0x1c, 0x67, 0x2a, 0x53,
	0x2b, 0x2f, 0x5d, 0x4e, 0x2e, 0xd6, 0x39, 0xaa, 0x68, 0x25, 0xd3, 0xac, 0xa8, 0x58, 0xc1, 0x35,
	0x6c, 0x47, 0xe9, 0x90, 0xc9, 0x8c, 0xa3, 0x00, 0xfe, 0x37, 0x96, 0x69, 0xeb, 0x83, 0x1e, 0xe8,
	0xef, 0x86, 0xfb, 0x93, 0x59, 0xd7, 0x7b, 0x9b, 0x75, 0xff, 0xdd, 0x09, 0x69, 0x87, 0x4e, 0x42,
	0x18, 0x36, 0xb9, 0x4c, 0xfd, 0x46, 0x8d, 0x63, 0x29, 0x04, 0x2f, 0x00, 0xb6, 0x43, 0xb7, 0x2d,
	0x3a, 0x85, 0x2d, 0x56, 0xa8, 0xb1, 0xac, 0x07, 0x56, 0x1a, 0xba, 0x81, 0x07, 0xea, 0x41, 0x72,
	0x6d, 0x72, 0x51, 0xde, 0x8a, 0x82, 0x1b, 0xbf, 0xd1, 0x6b, 0xf6, 0xf7, 0x06, 0x7d, 0xf2, 0x4b,
	0x4b, 0x52, 0xed, 0x3d, 0xdc, 0xc8, 0xa3, 0x4b, 0xb8, 0xb3, 0x72, 0x44, 0xa9, 0xf1, 0x9b, 0x7f,
	0x64, 0xad, 0x93, 0xc1, 0x07, 0x80, 0x87, 0x91, 0xcc, 0xb9, 0x16, 0x96, 0xa7, 0x5f, 0x95, 0xbe,
	0xa3, 0xc1, 0xb6, 0x68, 0x74, 0x0e, 0x51, 0xc9, 0x34, 0x97, 0xf6, 0x42, 0x8d, 0x46, 0x3c, 0xb1,
	0x42, 0xc9, 0xa8, 0xfe, 0xa6, 0x35, 0xbe, 0xe5, 0xc1, 0xdc, 0x6b, 0xb8, 0x6d, 0xc9, 0x8d, 0x7c,
	0x78, 0x35, 0x99, 0x63, 0x30, 0x9d, 0x63, 0xf0, 0x3e, 0xc7, 0xe0, 0x79, 0x81, 0xbd, 0xe9, 0x02,
	0x7b, 0xaf, 0x0b, 0xec, 0xdd, 0x0f, 0x32, 0x61, 0xf3, 0x71, 0x4c, 0x12, 0x55, 0xd0, 0x35, 0x93,
	0xfe, 0xa4, 0xd3, 0x47, 0x5a, 0xbd, 0xdb, 0xa7, 0x92, 0x9b, 0xb8, 0xb5, 0xfa, 0x58, 0x67, 0x9f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x67, 0xf0, 0xca, 0x22, 0xbd, 0x02, 0x00, 0x00,
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

func (m *Balance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Balance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Balance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BadgeIds) > 0 {
		for iNdEx := len(m.BadgeIds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.BadgeIds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBalances(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.OwnershipTimes) > 0 {
		for iNdEx := len(m.OwnershipTimes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.OwnershipTimes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBalances(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBalances(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *InheritedBalance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InheritedBalance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InheritedBalance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ParentBadgeIds) > 0 {
		for iNdEx := len(m.ParentBadgeIds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ParentBadgeIds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBalances(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size := m.ParentCollectionId.Size()
		i -= size
		if _, err := m.ParentCollectionId.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintBalances(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.BadgeIds) > 0 {
		for iNdEx := len(m.BadgeIds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.BadgeIds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBalances(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
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

func (m *Balance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Amount.Size()
	n += 1 + l + sovBalances(uint64(l))
	if len(m.OwnershipTimes) > 0 {
		for _, e := range m.OwnershipTimes {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
	if len(m.BadgeIds) > 0 {
		for _, e := range m.BadgeIds {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
	return n
}

func (m *InheritedBalance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.BadgeIds) > 0 {
		for _, e := range m.BadgeIds {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
	l = m.ParentCollectionId.Size()
	n += 1 + l + sovBalances(uint64(l))
	if len(m.ParentBadgeIds) > 0 {
		for _, e := range m.ParentBadgeIds {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
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
func (m *Balance) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Balance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Balance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OwnershipTimes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
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
				return ErrInvalidLengthBalances
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBalances
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OwnershipTimes = append(m.OwnershipTimes, &UintRange{})
			if err := m.OwnershipTimes[len(m.OwnershipTimes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BadgeIds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
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
				return ErrInvalidLengthBalances
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBalances
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BadgeIds = append(m.BadgeIds, &UintRange{})
			if err := m.BadgeIds[len(m.BadgeIds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *InheritedBalance) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: InheritedBalance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InheritedBalance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BadgeIds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
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
				return ErrInvalidLengthBalances
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBalances
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BadgeIds = append(m.BadgeIds, &UintRange{})
			if err := m.BadgeIds[len(m.BadgeIds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentCollectionId", wireType)
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
			if err := m.ParentCollectionId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentBadgeIds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
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
				return ErrInvalidLengthBalances
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBalances
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ParentBadgeIds = append(m.ParentBadgeIds, &UintRange{})
			if err := m.ParentBadgeIds[len(m.ParentBadgeIds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
