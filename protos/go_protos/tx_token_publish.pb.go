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

type TxTokenPublish struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	TokenName            string   `protobuf:"bytes,2,opt,name=Token_name,json=TokenName,proto3" json:"Token_name,omitempty"`
	TokenSymbol          string   `protobuf:"bytes,3,opt,name=Token_symbol,json=TokenSymbol,proto3" json:"Token_symbol,omitempty"`
	TokenSupply          uint64   `protobuf:"varint,4,opt,name=Token_supply,json=TokenSupply,proto3" json:"Token_supply,omitempty"`
	Token_VoutTotal      uint64   `protobuf:"varint,5,opt,name=Token_Vout_total,json=TokenVoutTotal,proto3" json:"Token_Vout_total,omitempty"`
	Token_Vout           []*Vout  `protobuf:"bytes,6,rep,name=Token_Vout,json=TokenVout,proto3" json:"Token_Vout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxTokenPublish) Reset()         { *m = TxTokenPublish{} }
func (m *TxTokenPublish) String() string { return proto.CompactTextString(m) }
func (*TxTokenPublish) ProtoMessage()    {}
func (*TxTokenPublish) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b0b13739c865944, []int{0}
}
func (m *TxTokenPublish) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxTokenPublish) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxTokenPublish.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxTokenPublish) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxTokenPublish.Merge(m, src)
}
func (m *TxTokenPublish) XXX_Size() int {
	return m.Size()
}
func (m *TxTokenPublish) XXX_DiscardUnknown() {
	xxx_messageInfo_TxTokenPublish.DiscardUnknown(m)
}

var xxx_messageInfo_TxTokenPublish proto.InternalMessageInfo

func (m *TxTokenPublish) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxTokenPublish) GetTokenName() string {
	if m != nil {
		return m.TokenName
	}
	return ""
}

func (m *TxTokenPublish) GetTokenSymbol() string {
	if m != nil {
		return m.TokenSymbol
	}
	return ""
}

func (m *TxTokenPublish) GetTokenSupply() uint64 {
	if m != nil {
		return m.TokenSupply
	}
	return 0
}

func (m *TxTokenPublish) GetToken_VoutTotal() uint64 {
	if m != nil {
		return m.Token_VoutTotal
	}
	return 0
}

func (m *TxTokenPublish) GetToken_Vout() []*Vout {
	if m != nil {
		return m.Token_Vout
	}
	return nil
}

func init() {
	proto.RegisterType((*TxTokenPublish)(nil), "go_protos.TxTokenPublish")
}

func init() { proto.RegisterFile("tx_token_publish.proto", fileDescriptor_1b0b13739c865944) }

var fileDescriptor_1b0b13739c865944 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0xa9, 0x88, 0x2f,
	0xc9, 0xcf, 0x4e, 0xcd, 0x8b, 0x2f, 0x28, 0x4d, 0xca, 0xc9, 0x2c, 0xce, 0xd0, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a, 0xa5, 0x78, 0x4a, 0x2a, 0x92, 0x12,
	0x8b, 0x53, 0x21, 0x12, 0x4a, 0xbf, 0x18, 0xb9, 0xf8, 0x42, 0x2a, 0x42, 0x40, 0x5a, 0x02, 0x20,
	0x3a, 0x84, 0x34, 0xb9, 0xd8, 0x42, 0x2a, 0x9c, 0x12, 0x8b, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35,
	0xb8, 0x8d, 0x04, 0xf5, 0xe0, 0x9a, 0xf5, 0x20, 0x12, 0x41, 0x50, 0x05, 0x42, 0xb2, 0x5c, 0x5c,
	0x60, 0xad, 0xf1, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x9c, 0x60,
	0x11, 0xbf, 0xc4, 0xdc, 0x54, 0x21, 0x45, 0x2e, 0x1e, 0x88, 0x74, 0x71, 0x65, 0x6e, 0x52, 0x7e,
	0x8e, 0x04, 0x33, 0x58, 0x01, 0x37, 0x58, 0x2c, 0x18, 0x2c, 0x84, 0xa4, 0xa4, 0xb4, 0xa0, 0x20,
	0xa7, 0x52, 0x82, 0x45, 0x81, 0x51, 0x83, 0x05, 0xa6, 0x04, 0x2c, 0x24, 0xa4, 0xc1, 0x25, 0x00,
	0x51, 0x12, 0x96, 0x5f, 0x5a, 0x12, 0x5f, 0x92, 0x5f, 0x92, 0x98, 0x23, 0xc1, 0x0a, 0x56, 0xc6,
	0x07, 0x16, 0x07, 0x09, 0x87, 0x80, 0x44, 0x85, 0xf4, 0x60, 0xce, 0x01, 0x09, 0x49, 0xb0, 0x29,
	0x30, 0x6b, 0x70, 0x1b, 0xf1, 0x23, 0xb9, 0x1e, 0x24, 0x0c, 0x75, 0x1f, 0x88, 0xe9, 0x24, 0x7b,
	0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7,
	0x10, 0xc5, 0xad, 0xa7, 0x0f, 0x57, 0x9e, 0xc4, 0x06, 0xa6, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xda, 0xe1, 0x73, 0x5f, 0x55, 0x01, 0x00, 0x00,
}

func (m *TxTokenPublish) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxTokenPublish) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxTokenPublish) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Token_Vout) > 0 {
		for iNdEx := len(m.Token_Vout) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Token_Vout[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTxTokenPublish(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.Token_VoutTotal != 0 {
		i = encodeVarintTxTokenPublish(dAtA, i, uint64(m.Token_VoutTotal))
		i--
		dAtA[i] = 0x28
	}
	if m.TokenSupply != 0 {
		i = encodeVarintTxTokenPublish(dAtA, i, uint64(m.TokenSupply))
		i--
		dAtA[i] = 0x20
	}
	if len(m.TokenSymbol) > 0 {
		i -= len(m.TokenSymbol)
		copy(dAtA[i:], m.TokenSymbol)
		i = encodeVarintTxTokenPublish(dAtA, i, uint64(len(m.TokenSymbol)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TokenName) > 0 {
		i -= len(m.TokenName)
		copy(dAtA[i:], m.TokenName)
		i = encodeVarintTxTokenPublish(dAtA, i, uint64(len(m.TokenName)))
		i--
		dAtA[i] = 0x12
	}
	if m.TxBase != nil {
		{
			size, err := m.TxBase.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTxTokenPublish(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxTokenPublish(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxTokenPublish(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxTokenPublish) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxTokenPublish(uint64(l))
	}
	l = len(m.TokenName)
	if l > 0 {
		n += 1 + l + sovTxTokenPublish(uint64(l))
	}
	l = len(m.TokenSymbol)
	if l > 0 {
		n += 1 + l + sovTxTokenPublish(uint64(l))
	}
	if m.TokenSupply != 0 {
		n += 1 + sovTxTokenPublish(uint64(m.TokenSupply))
	}
	if m.Token_VoutTotal != 0 {
		n += 1 + sovTxTokenPublish(uint64(m.Token_VoutTotal))
	}
	if len(m.Token_Vout) > 0 {
		for _, e := range m.Token_Vout {
			l = e.Size()
			n += 1 + l + sovTxTokenPublish(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxTokenPublish(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxTokenPublish(x uint64) (n int) {
	return sovTxTokenPublish(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxTokenPublish) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxTokenPublish
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
			return fmt.Errorf("proto: TxTokenPublish: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxTokenPublish: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPublish
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
				return ErrInvalidLengthTxTokenPublish
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPublish
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TxBase == nil {
				m.TxBase = &TxBase{}
			}
			if err := m.TxBase.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPublish
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
				return ErrInvalidLengthTxTokenPublish
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPublish
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenSymbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPublish
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
				return ErrInvalidLengthTxTokenPublish
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPublish
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenSymbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenSupply", wireType)
			}
			m.TokenSupply = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPublish
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TokenSupply |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token_VoutTotal", wireType)
			}
			m.Token_VoutTotal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPublish
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Token_VoutTotal |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token_Vout", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPublish
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
				return ErrInvalidLengthTxTokenPublish
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPublish
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token_Vout = append(m.Token_Vout, &Vout{})
			if err := m.Token_Vout[len(m.Token_Vout)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxTokenPublish(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxTokenPublish
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxTokenPublish
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
func skipTxTokenPublish(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxTokenPublish
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
					return 0, ErrIntOverflowTxTokenPublish
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
					return 0, ErrIntOverflowTxTokenPublish
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
				return 0, ErrInvalidLengthTxTokenPublish
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxTokenPublish
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxTokenPublish
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxTokenPublish        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxTokenPublish          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxTokenPublish = fmt.Errorf("proto: unexpected end of group")
)
