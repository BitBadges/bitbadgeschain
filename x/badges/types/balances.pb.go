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

//indexed by badgeid-subassetid-uniqueaccountnumber (26 bytes)
type BadgeBalanceInfo struct {
	IdsForBalances []*SubbadgeRange   `protobuf:"bytes,1,rep,name=idsForBalances,proto3" json:"idsForBalances,omitempty"`
	Balances       []uint64           `protobuf:"varint,2,rep,packed,name=balances,proto3" json:"balances,omitempty"`
	PendingNonce   uint64             `protobuf:"varint,3,opt,name=pending_nonce,json=pendingNonce,proto3" json:"pending_nonce,omitempty"`
	Pending        []*PendingTransfer `protobuf:"bytes,4,rep,name=pending,proto3" json:"pending,omitempty"`
	Approvals      []*Approval        `protobuf:"bytes,5,rep,name=approvals,proto3" json:"approvals,omitempty"`
	UserFlags      uint64             `protobuf:"varint,6,opt,name=user_flags,json=userFlags,proto3" json:"user_flags,omitempty"`
}

func (m *BadgeBalanceInfo) Reset()         { *m = BadgeBalanceInfo{} }
func (m *BadgeBalanceInfo) String() string { return proto.CompactTextString(m) }
func (*BadgeBalanceInfo) ProtoMessage()    {}
func (*BadgeBalanceInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_233d29a167e739f0, []int{0}
}
func (m *BadgeBalanceInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BadgeBalanceInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BadgeBalanceInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BadgeBalanceInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BadgeBalanceInfo.Merge(m, src)
}
func (m *BadgeBalanceInfo) XXX_Size() int {
	return m.Size()
}
func (m *BadgeBalanceInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_BadgeBalanceInfo.DiscardUnknown(m)
}

var xxx_messageInfo_BadgeBalanceInfo proto.InternalMessageInfo

func (m *BadgeBalanceInfo) GetIdsForBalances() []*SubbadgeRange {
	if m != nil {
		return m.IdsForBalances
	}
	return nil
}

func (m *BadgeBalanceInfo) GetBalances() []uint64 {
	if m != nil {
		return m.Balances
	}
	return nil
}

func (m *BadgeBalanceInfo) GetPendingNonce() uint64 {
	if m != nil {
		return m.PendingNonce
	}
	return 0
}

func (m *BadgeBalanceInfo) GetPending() []*PendingTransfer {
	if m != nil {
		return m.Pending
	}
	return nil
}

func (m *BadgeBalanceInfo) GetApprovals() []*Approval {
	if m != nil {
		return m.Approvals
	}
	return nil
}

func (m *BadgeBalanceInfo) GetUserFlags() uint64 {
	if m != nil {
		return m.UserFlags
	}
	return 0
}

type Approval struct {
	AddressNum uint64 `protobuf:"varint,1,opt,name=addressNum,proto3" json:"addressNum,omitempty"`
	Amount     uint64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	SubbadgeId uint64 `protobuf:"varint,3,opt,name=subbadgeId,proto3" json:"subbadgeId,omitempty"`
}

func (m *Approval) Reset()         { *m = Approval{} }
func (m *Approval) String() string { return proto.CompactTextString(m) }
func (*Approval) ProtoMessage()    {}
func (*Approval) Descriptor() ([]byte, []int) {
	return fileDescriptor_233d29a167e739f0, []int{1}
}
func (m *Approval) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Approval) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Approval.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Approval) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Approval.Merge(m, src)
}
func (m *Approval) XXX_Size() int {
	return m.Size()
}
func (m *Approval) XXX_DiscardUnknown() {
	xxx_messageInfo_Approval.DiscardUnknown(m)
}

var xxx_messageInfo_Approval proto.InternalMessageInfo

func (m *Approval) GetAddressNum() uint64 {
	if m != nil {
		return m.AddressNum
	}
	return 0
}

func (m *Approval) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *Approval) GetSubbadgeId() uint64 {
	if m != nil {
		return m.SubbadgeId
	}
	return 0
}

//Pending transfers will not be saved after accept / reject
type PendingTransfer struct {
	SubbadgeId        uint64 `protobuf:"varint,1,opt,name=subbadgeId,proto3" json:"subbadgeId,omitempty"`
	ThisPendingNonce  uint64 `protobuf:"varint,2,opt,name=this_pending_nonce,json=thisPendingNonce,proto3" json:"this_pending_nonce,omitempty"`
	OtherPendingNonce uint64 `protobuf:"varint,3,opt,name=other_pending_nonce,json=otherPendingNonce,proto3" json:"other_pending_nonce,omitempty"`
	Amount            uint64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
	SendRequest       bool   `protobuf:"varint,5,opt,name=send_request,json=sendRequest,proto3" json:"send_request,omitempty"`
	To                uint64 `protobuf:"varint,6,opt,name=to,proto3" json:"to,omitempty"`
	From              uint64 `protobuf:"varint,7,opt,name=from,proto3" json:"from,omitempty"`
	ApprovedBy        uint64 `protobuf:"varint,9,opt,name=approved_by,json=approvedBy,proto3" json:"approved_by,omitempty"`
}

func (m *PendingTransfer) Reset()         { *m = PendingTransfer{} }
func (m *PendingTransfer) String() string { return proto.CompactTextString(m) }
func (*PendingTransfer) ProtoMessage()    {}
func (*PendingTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_233d29a167e739f0, []int{2}
}
func (m *PendingTransfer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingTransfer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingTransfer.Merge(m, src)
}
func (m *PendingTransfer) XXX_Size() int {
	return m.Size()
}
func (m *PendingTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_PendingTransfer proto.InternalMessageInfo

func (m *PendingTransfer) GetSubbadgeId() uint64 {
	if m != nil {
		return m.SubbadgeId
	}
	return 0
}

func (m *PendingTransfer) GetThisPendingNonce() uint64 {
	if m != nil {
		return m.ThisPendingNonce
	}
	return 0
}

func (m *PendingTransfer) GetOtherPendingNonce() uint64 {
	if m != nil {
		return m.OtherPendingNonce
	}
	return 0
}

func (m *PendingTransfer) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *PendingTransfer) GetSendRequest() bool {
	if m != nil {
		return m.SendRequest
	}
	return false
}

func (m *PendingTransfer) GetTo() uint64 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *PendingTransfer) GetFrom() uint64 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *PendingTransfer) GetApprovedBy() uint64 {
	if m != nil {
		return m.ApprovedBy
	}
	return 0
}

func init() {
	proto.RegisterType((*BadgeBalanceInfo)(nil), "trevormil.bitbadgeschain.badges.BadgeBalanceInfo")
	proto.RegisterType((*Approval)(nil), "trevormil.bitbadgeschain.badges.Approval")
	proto.RegisterType((*PendingTransfer)(nil), "trevormil.bitbadgeschain.badges.PendingTransfer")
}

func init() { proto.RegisterFile("badges/balances.proto", fileDescriptor_233d29a167e739f0) }

var fileDescriptor_233d29a167e739f0 = []byte{
	// 492 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xd1, 0x6e, 0xd3, 0x3c,
	0x14, 0xc7, 0x9b, 0x34, 0xeb, 0xda, 0xd3, 0x7d, 0xdb, 0x3e, 0x0f, 0x90, 0x55, 0x89, 0xac, 0x94,
	0x9b, 0x22, 0xa1, 0x04, 0x8d, 0x27, 0xa0, 0x17, 0x43, 0x43, 0x68, 0x9a, 0x02, 0xe2, 0x82, 0x9b,
	0xc8, 0x69, 0xdc, 0x34, 0x52, 0x63, 0x07, 0xdb, 0x99, 0xd6, 0xa7, 0x80, 0xc7, 0xe2, 0x72, 0x97,
	0xdc, 0x81, 0xda, 0x17, 0x41, 0xb1, 0xdd, 0xad, 0x8d, 0x90, 0x7a, 0x95, 0xe3, 0x9f, 0xcf, 0xff,
	0xe4, 0x9c, 0xff, 0x91, 0xe1, 0x69, 0x42, 0xd2, 0x8c, 0xca, 0x30, 0x21, 0x0b, 0xc2, 0xa6, 0x54,
	0x06, 0xa5, 0xe0, 0x8a, 0xa3, 0x73, 0x25, 0xe8, 0x2d, 0x17, 0x45, 0xbe, 0x08, 0x92, 0x5c, 0x99,
	0x9c, 0xe9, 0x9c, 0xe4, 0x2c, 0x30, 0xf1, 0xe0, 0x49, 0xc6, 0x33, 0xae, 0x73, 0xc3, 0x3a, 0x32,
	0xb2, 0xc1, 0x99, 0xad, 0x56, 0x12, 0x41, 0x0a, 0xd9, 0x80, 0xe6, 0x63, 0xe1, 0x89, 0x85, 0xea,
	0xce, 0x80, 0xd1, 0x6f, 0x17, 0x4e, 0x27, 0x35, 0x9b, 0x98, 0x4e, 0xae, 0xd8, 0x8c, 0xa3, 0x2f,
	0x70, 0x9c, 0xa7, 0xf2, 0x92, 0x0b, 0x0b, 0x25, 0x76, 0x86, 0xed, 0x71, 0xff, 0x22, 0x08, 0xf6,
	0xf4, 0x17, 0x7c, 0xaa, 0x12, 0x1d, 0x45, 0x84, 0x65, 0x34, 0x6a, 0x54, 0x41, 0x03, 0xe8, 0x6e,
	0x06, 0xc6, 0xee, 0xb0, 0x3d, 0xf6, 0xa2, 0x87, 0x33, 0x7a, 0x09, 0xff, 0x95, 0x94, 0xa5, 0x39,
	0xcb, 0x62, 0xc6, 0xd9, 0x94, 0xe2, 0xf6, 0xd0, 0x19, 0x7b, 0xd1, 0x91, 0x85, 0xd7, 0x35, 0x43,
	0x1f, 0xe0, 0xd0, 0x9e, 0xb1, 0xa7, 0x3b, 0x7a, 0xb3, 0xb7, 0xa3, 0x1b, 0x93, 0xff, 0x59, 0x10,
	0x26, 0x67, 0x54, 0x44, 0x9b, 0x02, 0xe8, 0x3d, 0xf4, 0x48, 0x59, 0x0a, 0x7e, 0x4b, 0x16, 0x12,
	0x1f, 0xe8, 0x6a, 0xaf, 0xf6, 0x56, 0x7b, 0x67, 0x15, 0xd1, 0xa3, 0x16, 0x3d, 0x07, 0xa8, 0x24,
	0x15, 0xf1, 0x6c, 0x41, 0x32, 0x89, 0x3b, 0xba, 0xed, 0x5e, 0x4d, 0x2e, 0x6b, 0x30, 0x4a, 0xa0,
	0xbb, 0x51, 0x21, 0x1f, 0x80, 0xa4, 0xa9, 0xa0, 0x52, 0x5e, 0x57, 0x05, 0x76, 0x74, 0xea, 0x16,
	0x41, 0xcf, 0xa0, 0x43, 0x0a, 0x5e, 0x31, 0x85, 0x5d, 0x7d, 0x67, 0x4f, 0xb5, 0x4e, 0x5a, 0x67,
	0xaf, 0x52, 0xeb, 0xcc, 0x16, 0x19, 0x7d, 0x77, 0xe1, 0xa4, 0x31, 0x68, 0x43, 0xe3, 0x34, 0x35,
	0xe8, 0x35, 0x20, 0x35, 0xcf, 0x65, 0xbc, 0xeb, 0xba, 0xf9, 0xef, 0x69, 0x7d, 0x73, 0xb3, 0xed,
	0x7c, 0x00, 0x67, 0x5c, 0xcd, 0xa9, 0x88, 0xff, 0xb5, 0xa4, 0xff, 0xf5, 0xd5, 0x4e, 0xfe, 0xe3,
	0x24, 0xde, 0xce, 0x24, 0x2f, 0xe0, 0x48, 0x52, 0x96, 0xc6, 0x82, 0x7e, 0xab, 0xa8, 0x54, 0xf8,
	0x60, 0xe8, 0x8c, 0xbb, 0x51, 0xbf, 0x66, 0x91, 0x41, 0xe8, 0x18, 0x5c, 0xc5, 0xad, 0x8f, 0xae,
	0xe2, 0x08, 0x81, 0x37, 0x13, 0xbc, 0xc0, 0x87, 0x9a, 0xe8, 0x18, 0x9d, 0x43, 0xdf, 0x2c, 0x80,
	0xa6, 0x71, 0xb2, 0xc4, 0x3d, 0xeb, 0xa4, 0x45, 0x93, 0xe5, 0xe4, 0xe3, 0xcf, 0x95, 0xef, 0xdc,
	0xaf, 0x7c, 0xe7, 0xcf, 0xca, 0x77, 0x7e, 0xac, 0xfd, 0xd6, 0xfd, 0xda, 0x6f, 0xfd, 0x5a, 0xfb,
	0xad, 0xaf, 0x17, 0x59, 0xae, 0xe6, 0x55, 0x12, 0x4c, 0x79, 0x11, 0x3e, 0xac, 0x3b, 0xdc, 0x5d,
	0x77, 0x78, 0x17, 0x6e, 0x1e, 0xca, 0xb2, 0xa4, 0x32, 0xe9, 0xe8, 0xc7, 0xf2, 0xf6, 0x6f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x36, 0x09, 0xb2, 0xb8, 0xb7, 0x03, 0x00, 0x00,
}

func (m *BadgeBalanceInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BadgeBalanceInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BadgeBalanceInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.UserFlags != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.UserFlags))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Approvals) > 0 {
		for iNdEx := len(m.Approvals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Approvals[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBalances(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Pending) > 0 {
		for iNdEx := len(m.Pending) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pending[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBalances(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.PendingNonce != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.PendingNonce))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Balances) > 0 {
		dAtA2 := make([]byte, len(m.Balances)*10)
		var j1 int
		for _, num := range m.Balances {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintBalances(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x12
	}
	if len(m.IdsForBalances) > 0 {
		for iNdEx := len(m.IdsForBalances) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.IdsForBalances[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Approval) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Approval) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Approval) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SubbadgeId != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.SubbadgeId))
		i--
		dAtA[i] = 0x18
	}
	if m.Amount != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x10
	}
	if m.AddressNum != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.AddressNum))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PendingTransfer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingTransfer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingTransfer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ApprovedBy != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.ApprovedBy))
		i--
		dAtA[i] = 0x48
	}
	if m.From != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.From))
		i--
		dAtA[i] = 0x38
	}
	if m.To != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.To))
		i--
		dAtA[i] = 0x30
	}
	if m.SendRequest {
		i--
		if m.SendRequest {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.Amount != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x20
	}
	if m.OtherPendingNonce != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.OtherPendingNonce))
		i--
		dAtA[i] = 0x18
	}
	if m.ThisPendingNonce != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.ThisPendingNonce))
		i--
		dAtA[i] = 0x10
	}
	if m.SubbadgeId != 0 {
		i = encodeVarintBalances(dAtA, i, uint64(m.SubbadgeId))
		i--
		dAtA[i] = 0x8
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
func (m *BadgeBalanceInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.IdsForBalances) > 0 {
		for _, e := range m.IdsForBalances {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
	if len(m.Balances) > 0 {
		l = 0
		for _, e := range m.Balances {
			l += sovBalances(uint64(e))
		}
		n += 1 + sovBalances(uint64(l)) + l
	}
	if m.PendingNonce != 0 {
		n += 1 + sovBalances(uint64(m.PendingNonce))
	}
	if len(m.Pending) > 0 {
		for _, e := range m.Pending {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
	if len(m.Approvals) > 0 {
		for _, e := range m.Approvals {
			l = e.Size()
			n += 1 + l + sovBalances(uint64(l))
		}
	}
	if m.UserFlags != 0 {
		n += 1 + sovBalances(uint64(m.UserFlags))
	}
	return n
}

func (m *Approval) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AddressNum != 0 {
		n += 1 + sovBalances(uint64(m.AddressNum))
	}
	if m.Amount != 0 {
		n += 1 + sovBalances(uint64(m.Amount))
	}
	if m.SubbadgeId != 0 {
		n += 1 + sovBalances(uint64(m.SubbadgeId))
	}
	return n
}

func (m *PendingTransfer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubbadgeId != 0 {
		n += 1 + sovBalances(uint64(m.SubbadgeId))
	}
	if m.ThisPendingNonce != 0 {
		n += 1 + sovBalances(uint64(m.ThisPendingNonce))
	}
	if m.OtherPendingNonce != 0 {
		n += 1 + sovBalances(uint64(m.OtherPendingNonce))
	}
	if m.Amount != 0 {
		n += 1 + sovBalances(uint64(m.Amount))
	}
	if m.SendRequest {
		n += 2
	}
	if m.To != 0 {
		n += 1 + sovBalances(uint64(m.To))
	}
	if m.From != 0 {
		n += 1 + sovBalances(uint64(m.From))
	}
	if m.ApprovedBy != 0 {
		n += 1 + sovBalances(uint64(m.ApprovedBy))
	}
	return n
}

func sovBalances(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBalances(x uint64) (n int) {
	return sovBalances(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BadgeBalanceInfo) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BadgeBalanceInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BadgeBalanceInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IdsForBalances", wireType)
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
			m.IdsForBalances = append(m.IdsForBalances, &SubbadgeRange{})
			if err := m.IdsForBalances[len(m.IdsForBalances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBalances
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Balances = append(m.Balances, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBalances
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthBalances
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthBalances
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Balances) == 0 {
					m.Balances = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBalances
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Balances = append(m.Balances, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Balances", wireType)
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PendingNonce", wireType)
			}
			m.PendingNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PendingNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pending", wireType)
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
			m.Pending = append(m.Pending, &PendingTransfer{})
			if err := m.Pending[len(m.Pending)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Approvals", wireType)
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
			m.Approvals = append(m.Approvals, &Approval{})
			if err := m.Approvals[len(m.Approvals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserFlags", wireType)
			}
			m.UserFlags = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UserFlags |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *Approval) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Approval: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Approval: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressNum", wireType)
			}
			m.AddressNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AddressNum |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubbadgeId", wireType)
			}
			m.SubbadgeId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SubbadgeId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *PendingTransfer) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PendingTransfer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingTransfer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubbadgeId", wireType)
			}
			m.SubbadgeId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SubbadgeId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ThisPendingNonce", wireType)
			}
			m.ThisPendingNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ThisPendingNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OtherPendingNonce", wireType)
			}
			m.OtherPendingNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OtherPendingNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SendRequest", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SendRequest = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
			}
			m.To = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.To |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			m.From = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.From |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApprovedBy", wireType)
			}
			m.ApprovedBy = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBalances
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ApprovedBy |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
