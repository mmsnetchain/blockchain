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

type TxTokenPay struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	Token_Txid           []byte   `protobuf:"bytes,2,opt,name=Token_Txid,json=TokenTxid,proto3" json:"Token_Txid,omitempty"`
	Token_VinTotal       uint64   `protobuf:"varint,3,opt,name=Token_Vin_total,json=TokenVinTotal,proto3" json:"Token_Vin_total,omitempty"`
	Token_Vin            []*Vin   `protobuf:"bytes,4,rep,name=Token_Vin,json=TokenVin,proto3" json:"Token_Vin,omitempty"`
	Token_VoutTotal      uint64   `protobuf:"varint,5,opt,name=Token_Vout_total,json=TokenVoutTotal,proto3" json:"Token_Vout_total,omitempty"`
	Token_Vout           []*Vout  `protobuf:"bytes,6,rep,name=Token_Vout,json=TokenVout,proto3" json:"Token_Vout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxTokenPay) Reset()         { *m = TxTokenPay{} }
func (m *TxTokenPay) String() string { return proto.CompactTextString(m) }
func (*TxTokenPay) ProtoMessage()    {}
func (*TxTokenPay) Descriptor() ([]byte, []int) {
	return fileDescriptor_802412a3d36ae489, []int{0}
}
func (m *TxTokenPay) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxTokenPay) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxTokenPay.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxTokenPay) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxTokenPay.Merge(m, src)
}
func (m *TxTokenPay) XXX_Size() int {
	return m.Size()
}
func (m *TxTokenPay) XXX_DiscardUnknown() {
	xxx_messageInfo_TxTokenPay.DiscardUnknown(m)
}

var xxx_messageInfo_TxTokenPay proto.InternalMessageInfo

func (m *TxTokenPay) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxTokenPay) GetToken_Txid() []byte {
	if m != nil {
		return m.Token_Txid
	}
	return nil
}

func (m *TxTokenPay) GetToken_VinTotal() uint64 {
	if m != nil {
		return m.Token_VinTotal
	}
	return 0
}

func (m *TxTokenPay) GetToken_Vin() []*Vin {
	if m != nil {
		return m.Token_Vin
	}
	return nil
}

func (m *TxTokenPay) GetToken_VoutTotal() uint64 {
	if m != nil {
		return m.Token_VoutTotal
	}
	return 0
}

func (m *TxTokenPay) GetToken_Vout() []*Vout {
	if m != nil {
		return m.Token_Vout
	}
	return nil
}

func init() {
	proto.RegisterType((*TxTokenPay)(nil), "go_protos.TxTokenPay")
}

func init() { proto.RegisterFile("tx_token_pay.proto", fileDescriptor_802412a3d36ae489) }

var fileDescriptor_802412a3d36ae489 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0xa9, 0x88, 0x2f,
	0xc9, 0xcf, 0x4e, 0xcd, 0x8b, 0x2f, 0x48, 0xac, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a, 0xa5, 0x78, 0x4a, 0x2a, 0x92, 0x12, 0x8b, 0x53, 0x21, 0x12,
	0x4a, 0xad, 0x4c, 0x5c, 0x5c, 0x21, 0x15, 0x21, 0x20, 0xe5, 0x01, 0x89, 0x95, 0x42, 0x9a, 0x5c,
	0x6c, 0x21, 0x15, 0x4e, 0x89, 0xc5, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x82, 0x7a,
	0x70, 0x8d, 0x7a, 0x10, 0x89, 0x20, 0xa8, 0x02, 0x21, 0x59, 0x2e, 0x2e, 0xb0, 0xb6, 0xf8, 0x90,
	0x8a, 0xcc, 0x14, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x4e, 0xb0, 0x08, 0x48, 0x40, 0x48,
	0x8d, 0x8b, 0x1f, 0x22, 0x1d, 0x96, 0x99, 0x17, 0x5f, 0x92, 0x5f, 0x92, 0x98, 0x23, 0xc1, 0xac,
	0xc0, 0xa8, 0xc1, 0x12, 0xc4, 0x0b, 0x16, 0x0e, 0xcb, 0xcc, 0x0b, 0x01, 0x09, 0x0a, 0x69, 0x73,
	0x71, 0xc2, 0xd5, 0x49, 0xb0, 0x28, 0x30, 0x6b, 0x70, 0x1b, 0xf1, 0x21, 0x59, 0x1a, 0x96, 0x99,
	0x17, 0xc4, 0x01, 0xd3, 0x21, 0xa4, 0xc1, 0x25, 0x00, 0x55, 0x9c, 0x5f, 0x5a, 0x02, 0x35, 0x95,
	0x15, 0x6c, 0x2a, 0x1f, 0x44, 0x4d, 0x7e, 0x69, 0x09, 0xc4, 0x58, 0x3d, 0x98, 0xeb, 0x40, 0x42,
	0x12, 0x6c, 0x60, 0x73, 0xf9, 0x91, 0xcd, 0xcd, 0x2f, 0x2d, 0x81, 0x3a, 0x17, 0xc4, 0x74, 0x92,
	0x3d, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x67, 0x3c, 0x96,
	0x63, 0x88, 0xe2, 0xd6, 0xd3, 0x87, 0x2b, 0x4f, 0x62, 0x03, 0xd3, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xd5, 0x54, 0x5f, 0xe0, 0x5c, 0x01, 0x00, 0x00,
}

func (m *TxTokenPay) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxTokenPay) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxTokenPay) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
				i = encodeVarintTxTokenPay(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.Token_VoutTotal != 0 {
		i = encodeVarintTxTokenPay(dAtA, i, uint64(m.Token_VoutTotal))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Token_Vin) > 0 {
		for iNdEx := len(m.Token_Vin) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Token_Vin[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTxTokenPay(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Token_VinTotal != 0 {
		i = encodeVarintTxTokenPay(dAtA, i, uint64(m.Token_VinTotal))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Token_Txid) > 0 {
		i -= len(m.Token_Txid)
		copy(dAtA[i:], m.Token_Txid)
		i = encodeVarintTxTokenPay(dAtA, i, uint64(len(m.Token_Txid)))
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
			i = encodeVarintTxTokenPay(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxTokenPay(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxTokenPay(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxTokenPay) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxTokenPay(uint64(l))
	}
	l = len(m.Token_Txid)
	if l > 0 {
		n += 1 + l + sovTxTokenPay(uint64(l))
	}
	if m.Token_VinTotal != 0 {
		n += 1 + sovTxTokenPay(uint64(m.Token_VinTotal))
	}
	if len(m.Token_Vin) > 0 {
		for _, e := range m.Token_Vin {
			l = e.Size()
			n += 1 + l + sovTxTokenPay(uint64(l))
		}
	}
	if m.Token_VoutTotal != 0 {
		n += 1 + sovTxTokenPay(uint64(m.Token_VoutTotal))
	}
	if len(m.Token_Vout) > 0 {
		for _, e := range m.Token_Vout {
			l = e.Size()
			n += 1 + l + sovTxTokenPay(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxTokenPay(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxTokenPay(x uint64) (n int) {
	return sovTxTokenPay(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxTokenPay) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxTokenPay
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
			return fmt.Errorf("proto: TxTokenPay: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxTokenPay: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPay
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
				return ErrInvalidLengthTxTokenPay
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPay
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
				return fmt.Errorf("proto: wrong wireType = %d for field Token_Txid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPay
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
				return ErrInvalidLengthTxTokenPay
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPay
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token_Txid = append(m.Token_Txid[:0], dAtA[iNdEx:postIndex]...)
			if m.Token_Txid == nil {
				m.Token_Txid = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token_VinTotal", wireType)
			}
			m.Token_VinTotal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPay
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Token_VinTotal |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token_Vin", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPay
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
				return ErrInvalidLengthTxTokenPay
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPay
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token_Vin = append(m.Token_Vin, &Vin{})
			if err := m.Token_Vin[len(m.Token_Vin)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token_VoutTotal", wireType)
			}
			m.Token_VoutTotal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxTokenPay
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
					return ErrIntOverflowTxTokenPay
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
				return ErrInvalidLengthTxTokenPay
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxTokenPay
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
			skippy, err := skipTxTokenPay(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxTokenPay
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxTokenPay
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
func skipTxTokenPay(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxTokenPay
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
					return 0, ErrIntOverflowTxTokenPay
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
					return 0, ErrIntOverflowTxTokenPay
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
				return 0, ErrInvalidLengthTxTokenPay
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxTokenPay
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxTokenPay
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxTokenPay        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxTokenPay          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxTokenPay = fmt.Errorf("proto: unexpected end of group")
)
