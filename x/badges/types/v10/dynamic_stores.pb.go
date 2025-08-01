// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: badges/v10/dynamic_stores.proto

package v10

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

// A DynamicStore is a flexible storage object that can store arbitrary data.
// It is identified by a unique ID assigned by the blockchain, which is a uint64 that increments.
// Dynamic stores are created by users and can only be updated or deleted by their creator.
// They provide a way to store custom data on-chain with proper access control.
type DynamicStore struct {
	// The unique identifier for this dynamic store. This is assigned by the blockchain.
	StoreId Uint `protobuf:"bytes,1,opt,name=storeId,proto3,customtype=Uint" json:"storeId"`
	// The address of the creator of this dynamic store.
	CreatedBy string `protobuf:"bytes,2,opt,name=createdBy,proto3" json:"createdBy,omitempty"`
	// The default value for uninitialized addresses.
	DefaultValue bool `protobuf:"varint,3,opt,name=defaultValue,proto3" json:"defaultValue,omitempty"`
}

func (m *DynamicStore) Reset()         { *m = DynamicStore{} }
func (m *DynamicStore) String() string { return proto.CompactTextString(m) }
func (*DynamicStore) ProtoMessage()    {}
func (*DynamicStore) Descriptor() ([]byte, []int) {
	return fileDescriptor_576760bfd9c270d3, []int{0}
}
func (m *DynamicStore) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DynamicStore) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DynamicStore.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DynamicStore) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynamicStore.Merge(m, src)
}
func (m *DynamicStore) XXX_Size() int {
	return m.Size()
}
func (m *DynamicStore) XXX_DiscardUnknown() {
	xxx_messageInfo_DynamicStore.DiscardUnknown(m)
}

var xxx_messageInfo_DynamicStore proto.InternalMessageInfo

func (m *DynamicStore) GetCreatedBy() string {
	if m != nil {
		return m.CreatedBy
	}
	return ""
}

func (m *DynamicStore) GetDefaultValue() bool {
	if m != nil {
		return m.DefaultValue
	}
	return false
}

// A DynamicStoreValue stores a 0/1 flag for a specific address in a dynamic store.
// This allows the creator to set boolean values per address.
type DynamicStoreValue struct {
	// The unique identifier for this dynamic store.
	StoreId Uint `protobuf:"bytes,1,opt,name=storeId,proto3,customtype=Uint" json:"storeId"`
	// The address for which this value is stored.
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// The boolean value (0 or 1).
	Value bool `protobuf:"varint,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *DynamicStoreValue) Reset()         { *m = DynamicStoreValue{} }
func (m *DynamicStoreValue) String() string { return proto.CompactTextString(m) }
func (*DynamicStoreValue) ProtoMessage()    {}
func (*DynamicStoreValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_576760bfd9c270d3, []int{1}
}
func (m *DynamicStoreValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DynamicStoreValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DynamicStoreValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DynamicStoreValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DynamicStoreValue.Merge(m, src)
}
func (m *DynamicStoreValue) XXX_Size() int {
	return m.Size()
}
func (m *DynamicStoreValue) XXX_DiscardUnknown() {
	xxx_messageInfo_DynamicStoreValue.DiscardUnknown(m)
}

var xxx_messageInfo_DynamicStoreValue proto.InternalMessageInfo

func (m *DynamicStoreValue) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *DynamicStoreValue) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

func init() {
	proto.RegisterType((*DynamicStore)(nil), "badges.v10.DynamicStore")
	proto.RegisterType((*DynamicStoreValue)(nil), "badges.v10.DynamicStoreValue")
}

func init() { proto.RegisterFile("badges/v10/dynamic_stores.proto", fileDescriptor_576760bfd9c270d3) }

var fileDescriptor_576760bfd9c270d3 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x4a, 0x4c, 0x49,
	0x4f, 0x2d, 0xd6, 0x2f, 0x33, 0x34, 0xd0, 0x4f, 0xa9, 0xcc, 0x4b, 0xcc, 0xcd, 0x4c, 0x8e, 0x2f,
	0x2e, 0xc9, 0x2f, 0x4a, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x82, 0x28, 0xd0,
	0x2b, 0x33, 0x34, 0x90, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x0b, 0xeb, 0x83, 0x58, 0x10, 0x15,
	0x4a, 0x15, 0x5c, 0x3c, 0x2e, 0x10, 0x9d, 0xc1, 0x20, 0x8d, 0x42, 0x6a, 0x5c, 0xec, 0x60, 0x13,
	0x3c, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x9d, 0x78, 0x4e, 0xdc, 0x93, 0x67, 0xb8, 0x75,
	0x4f, 0x9e, 0x25, 0x34, 0x33, 0xaf, 0x24, 0x08, 0x26, 0x29, 0x24, 0xc3, 0xc5, 0x99, 0x5c, 0x94,
	0x9a, 0x58, 0x92, 0x9a, 0xe2, 0x54, 0x29, 0xc1, 0x04, 0x52, 0x19, 0x84, 0x10, 0x10, 0x52, 0xe2,
	0xe2, 0x49, 0x49, 0x4d, 0x4b, 0x2c, 0xcd, 0x29, 0x09, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x56,
	0x60, 0xd4, 0xe0, 0x08, 0x42, 0x11, 0x53, 0xca, 0xe6, 0x12, 0x44, 0xb6, 0x19, 0x2c, 0x48, 0xb4,
	0xf5, 0x12, 0x5c, 0xec, 0x89, 0x29, 0x29, 0x45, 0xa9, 0xc5, 0xc5, 0x50, 0xcb, 0x61, 0x5c, 0x21,
	0x11, 0x2e, 0xd6, 0x32, 0x24, 0x3b, 0x21, 0x1c, 0xa7, 0x80, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c,
	0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e,
	0x3c, 0x96, 0x63, 0x88, 0x32, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5,
	0x4f, 0xca, 0x2c, 0x81, 0x86, 0x28, 0x9c, 0x95, 0x9c, 0x91, 0x98, 0x99, 0xa7, 0x5f, 0xa1, 0x0f,
	0x15, 0x2f, 0xa9, 0x2c, 0x80, 0x84, 0x77, 0x12, 0x1b, 0x38, 0xfc, 0x8c, 0x01, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x63, 0x62, 0x01, 0x80, 0x84, 0x01, 0x00, 0x00,
}

func (m *DynamicStore) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DynamicStore) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DynamicStore) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.DefaultValue {
		i--
		if m.DefaultValue {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.CreatedBy) > 0 {
		i -= len(m.CreatedBy)
		copy(dAtA[i:], m.CreatedBy)
		i = encodeVarintDynamicStores(dAtA, i, uint64(len(m.CreatedBy)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.StoreId.Size()
		i -= size
		if _, err := m.StoreId.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDynamicStores(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *DynamicStoreValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DynamicStoreValue) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DynamicStoreValue) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Value {
		i--
		if m.Value {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintDynamicStores(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.StoreId.Size()
		i -= size
		if _, err := m.StoreId.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDynamicStores(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintDynamicStores(dAtA []byte, offset int, v uint64) int {
	offset -= sovDynamicStores(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DynamicStore) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.StoreId.Size()
	n += 1 + l + sovDynamicStores(uint64(l))
	l = len(m.CreatedBy)
	if l > 0 {
		n += 1 + l + sovDynamicStores(uint64(l))
	}
	if m.DefaultValue {
		n += 2
	}
	return n
}

func (m *DynamicStoreValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.StoreId.Size()
	n += 1 + l + sovDynamicStores(uint64(l))
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovDynamicStores(uint64(l))
	}
	if m.Value {
		n += 2
	}
	return n
}

func sovDynamicStores(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDynamicStores(x uint64) (n int) {
	return sovDynamicStores(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DynamicStore) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDynamicStores
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
			return fmt.Errorf("proto: DynamicStore: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DynamicStore: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicStores
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
				return ErrInvalidLengthDynamicStores
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDynamicStores
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StoreId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedBy", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicStores
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
				return ErrInvalidLengthDynamicStores
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDynamicStores
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CreatedBy = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultValue", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicStores
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
			m.DefaultValue = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipDynamicStores(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDynamicStores
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
func (m *DynamicStoreValue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDynamicStores
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
			return fmt.Errorf("proto: DynamicStoreValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DynamicStoreValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicStores
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
				return ErrInvalidLengthDynamicStores
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDynamicStores
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StoreId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicStores
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
				return ErrInvalidLengthDynamicStores
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDynamicStores
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDynamicStores
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
			m.Value = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipDynamicStores(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDynamicStores
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
func skipDynamicStores(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDynamicStores
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
					return 0, ErrIntOverflowDynamicStores
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
					return 0, ErrIntOverflowDynamicStores
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
				return 0, ErrInvalidLengthDynamicStores
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDynamicStores
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDynamicStores
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDynamicStores        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDynamicStores          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDynamicStores = fmt.Errorf("proto: unexpected end of group")
)
