package go_protos

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const _ = proto.ProtoPackageIsVersion3

type HistoryItem struct {
	GenerateId           []byte   `protobuf:"bytes,1,opt,name=GenerateId,proto3" json:"GenerateId,omitempty"`
	IsIn                 bool     `protobuf:"varint,2,opt,name=IsIn,proto3" json:"IsIn,omitempty"`
	Type                 uint64   `protobuf:"varint,3,opt,name=Type,proto3" json:"Type,omitempty"`
	InAddr               [][]byte `protobuf:"bytes,4,rep,name=InAddr,proto3" json:"InAddr,omitempty"`
	OutAddr              [][]byte `protobuf:"bytes,5,rep,name=OutAddr,proto3" json:"OutAddr,omitempty"`
	Value                uint64   `protobuf:"varint,6,opt,name=Value,proto3" json:"Value,omitempty"`
	Txid                 []byte   `protobuf:"bytes,7,opt,name=Txid,proto3" json:"Txid,omitempty"`
	Height               uint64   `protobuf:"varint,8,opt,name=Height,proto3" json:"Height,omitempty"`
	Payload              []byte   `protobuf:"bytes,9,opt,name=Payload,proto3" json:"Payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HistoryItem) Reset()         { *m = HistoryItem{} }
func (m *HistoryItem) String() string { return proto.CompactTextString(m) }
func (*HistoryItem) ProtoMessage()    {}
func (*HistoryItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_2530718a19f54266, []int{0}
}
func (m *HistoryItem) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoryItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoryItem.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoryItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoryItem.Merge(m, src)
}
func (m *HistoryItem) XXX_Size() int {
	return m.Size()
}
func (m *HistoryItem) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoryItem.DiscardUnknown(m)
}

var xxx_messageInfo_HistoryItem proto.InternalMessageInfo

func (m *HistoryItem) GetGenerateId() []byte {
	if m != nil {
		return m.GenerateId
	}
	return nil
}

func (m *HistoryItem) GetIsIn() bool {
	if m != nil {
		return m.IsIn
	}
	return false
}

func (m *HistoryItem) GetType() uint64 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *HistoryItem) GetInAddr() [][]byte {
	if m != nil {
		return m.InAddr
	}
	return nil
}

func (m *HistoryItem) GetOutAddr() [][]byte {
	if m != nil {
		return m.OutAddr
	}
	return nil
}

func (m *HistoryItem) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *HistoryItem) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *HistoryItem) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *HistoryItem) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*HistoryItem)(nil), "go_protos.HistoryItem")
}

func init() { proto.RegisterFile("history_item.proto", fileDescriptor_2530718a19f54266) }

var fileDescriptor_2530718a19f54266 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xc8, 0x2c, 0x2e,
	0xc9, 0x2f, 0xaa, 0x8c, 0xcf, 0x2c, 0x49, 0xcd, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a, 0x95, 0x1e, 0x33, 0x72, 0x71, 0x7b, 0x40, 0x54, 0x78, 0x96,
	0xa4, 0xe6, 0x0a, 0xc9, 0x71, 0x71, 0xb9, 0xa7, 0xe6, 0xa5, 0x16, 0x25, 0x96, 0xa4, 0x7a, 0xa6,
	0x48, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x21, 0x89, 0x08, 0x09, 0x71, 0xb1, 0x78, 0x16, 0x7b,
	0xe6, 0x49, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x04, 0x81, 0xd9, 0x20, 0xb1, 0x90, 0xca, 0x82, 0x54,
	0x09, 0x66, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x30, 0x5b, 0x48, 0x8c, 0x8b, 0xcd, 0x33, 0xcf, 0x31,
	0x25, 0xa5, 0x48, 0x82, 0x45, 0x81, 0x59, 0x83, 0x27, 0x08, 0xca, 0x13, 0x92, 0xe0, 0x62, 0xf7,
	0x2f, 0x2d, 0x01, 0x4b, 0xb0, 0x82, 0x25, 0x60, 0x5c, 0x21, 0x11, 0x2e, 0xd6, 0xb0, 0xc4, 0x9c,
	0xd2, 0x54, 0x09, 0x36, 0xb0, 0x31, 0x10, 0x0e, 0xd8, 0xec, 0x8a, 0xcc, 0x14, 0x09, 0x76, 0xb0,
	0x4b, 0xc0, 0x6c, 0x90, 0xd9, 0x1e, 0xa9, 0x99, 0xe9, 0x19, 0x25, 0x12, 0x1c, 0x60, 0xa5, 0x50,
	0x1e, 0xc8, 0xec, 0x80, 0xc4, 0xca, 0x9c, 0xfc, 0xc4, 0x14, 0x09, 0x4e, 0xb0, 0x72, 0x18, 0xd7,
	0x49, 0xf6, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf1,
	0x58, 0x8e, 0x21, 0x8a, 0x5b, 0x4f, 0x1f, 0x1e, 0x08, 0x49, 0x6c, 0x60, 0xda, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0x6d, 0x00, 0x9a, 0xd2, 0x2c, 0x01, 0x00, 0x00,
}

func (m *HistoryItem) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoryItem) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoryItem) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintHistoryItem(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x4a
	}
	if m.Height != 0 {
		i = encodeVarintHistoryItem(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x40
	}
	if len(m.Txid) > 0 {
		i -= len(m.Txid)
		copy(dAtA[i:], m.Txid)
		i = encodeVarintHistoryItem(dAtA, i, uint64(len(m.Txid)))
		i--
		dAtA[i] = 0x3a
	}
	if m.Value != 0 {
		i = encodeVarintHistoryItem(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x30
	}
	if len(m.OutAddr) > 0 {
		for iNdEx := len(m.OutAddr) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.OutAddr[iNdEx])
			copy(dAtA[i:], m.OutAddr[iNdEx])
			i = encodeVarintHistoryItem(dAtA, i, uint64(len(m.OutAddr[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.InAddr) > 0 {
		for iNdEx := len(m.InAddr) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.InAddr[iNdEx])
			copy(dAtA[i:], m.InAddr[iNdEx])
			i = encodeVarintHistoryItem(dAtA, i, uint64(len(m.InAddr[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Type != 0 {
		i = encodeVarintHistoryItem(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x18
	}
	if m.IsIn {
		i--
		if m.IsIn {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.GenerateId) > 0 {
		i -= len(m.GenerateId)
		copy(dAtA[i:], m.GenerateId)
		i = encodeVarintHistoryItem(dAtA, i, uint64(len(m.GenerateId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintHistoryItem(dAtA []byte, offset int, v uint64) int {
	offset -= sovHistoryItem(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *HistoryItem) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GenerateId)
	if l > 0 {
		n += 1 + l + sovHistoryItem(uint64(l))
	}
	if m.IsIn {
		n += 2
	}
	if m.Type != 0 {
		n += 1 + sovHistoryItem(uint64(m.Type))
	}
	if len(m.InAddr) > 0 {
		for _, b := range m.InAddr {
			l = len(b)
			n += 1 + l + sovHistoryItem(uint64(l))
		}
	}
	if len(m.OutAddr) > 0 {
		for _, b := range m.OutAddr {
			l = len(b)
			n += 1 + l + sovHistoryItem(uint64(l))
		}
	}
	if m.Value != 0 {
		n += 1 + sovHistoryItem(uint64(m.Value))
	}
	l = len(m.Txid)
	if l > 0 {
		n += 1 + l + sovHistoryItem(uint64(l))
	}
	if m.Height != 0 {
		n += 1 + sovHistoryItem(uint64(m.Height))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovHistoryItem(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovHistoryItem(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozHistoryItem(x uint64) (n int) {
	return sovHistoryItem(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HistoryItem) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHistoryItem
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
			return fmt.Errorf("proto: HistoryItem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoryItem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenerateId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
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
				return ErrInvalidLengthHistoryItem
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GenerateId = append(m.GenerateId[:0], dAtA[iNdEx:postIndex]...)
			if m.GenerateId == nil {
				m.GenerateId = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsIn", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
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
			m.IsIn = bool(v != 0)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InAddr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
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
				return ErrInvalidLengthHistoryItem
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InAddr = append(m.InAddr, make([]byte, postIndex-iNdEx))
			copy(m.InAddr[len(m.InAddr)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutAddr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
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
				return ErrInvalidLengthHistoryItem
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OutAddr = append(m.OutAddr, make([]byte, postIndex-iNdEx))
			copy(m.OutAddr[len(m.OutAddr)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
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
				return ErrInvalidLengthHistoryItem
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Txid = append(m.Txid[:0], dAtA[iNdEx:postIndex]...)
			if m.Txid == nil {
				m.Txid = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHistoryItem
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
				return ErrInvalidLengthHistoryItem
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHistoryItem(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHistoryItem
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipHistoryItem(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHistoryItem
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
					return 0, ErrIntOverflowHistoryItem
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
					return 0, ErrIntOverflowHistoryItem
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
				return 0, ErrInvalidLengthHistoryItem
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupHistoryItem
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthHistoryItem
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthHistoryItem        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHistoryItem          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupHistoryItem = fmt.Errorf("proto: unexpected end of group")
)
