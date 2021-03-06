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

type TxPay struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxPay) Reset()         { *m = TxPay{} }
func (m *TxPay) String() string { return proto.CompactTextString(m) }
func (*TxPay) ProtoMessage()    {}
func (*TxPay) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d938a31d2830631, []int{0}
}
func (m *TxPay) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxPay) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxPay.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxPay) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxPay.Merge(m, src)
}
func (m *TxPay) XXX_Size() int {
	return m.Size()
}
func (m *TxPay) XXX_DiscardUnknown() {
	xxx_messageInfo_TxPay.DiscardUnknown(m)
}

var xxx_messageInfo_TxPay proto.InternalMessageInfo

func (m *TxPay) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func init() {
	proto.RegisterType((*TxPay)(nil), "go_protos.TxPay")
}

func init() { proto.RegisterFile("tx_pay.proto", fileDescriptor_8d938a31d2830631) }

var fileDescriptor_8d938a31d2830631 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0xa9, 0x88, 0x2f,
	0x48, 0xac, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a,
	0xa5, 0x78, 0x4a, 0x2a, 0x92, 0x12, 0x8b, 0x53, 0x21, 0x12, 0x4a, 0x46, 0x5c, 0xac, 0x21, 0x15,
	0x01, 0x89, 0x95, 0x42, 0x9a, 0x5c, 0x6c, 0x21, 0x15, 0x4e, 0x89, 0xc5, 0xa9, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0xdc, 0x46, 0x82, 0x7a, 0x70, 0x2d, 0x7a, 0x10, 0x89, 0x20, 0xa8, 0x02, 0x27, 0xd9,
	0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc6, 0x63, 0x39,
	0x86, 0x28, 0x6e, 0x3d, 0x7d, 0xb8, 0xea, 0x24, 0x36, 0x30, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff,
	0xff, 0xed, 0x01, 0x4a, 0x7e, 0x82, 0x00, 0x00, 0x00,
}

func (m *TxPay) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxPay) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxPay) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.TxBase != nil {
		{
			size, err := m.TxBase.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTxPay(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxPay(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxPay(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxPay) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxPay(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxPay(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxPay(x uint64) (n int) {
	return sovTxPay(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxPay) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxPay
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
			return fmt.Errorf("proto: TxPay: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxPay: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxPay
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
				return ErrInvalidLengthTxPay
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxPay
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
		default:
			iNdEx = preIndex
			skippy, err := skipTxPay(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxPay
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxPay
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
func skipTxPay(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxPay
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
					return 0, ErrIntOverflowTxPay
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
					return 0, ErrIntOverflowTxPay
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
				return 0, ErrInvalidLengthTxPay
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxPay
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxPay
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxPay        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxPay          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxPay = fmt.Errorf("proto: unexpected end of group")
)
