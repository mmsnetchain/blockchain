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

type TxVoteReward struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	StartHeight          uint64   `protobuf:"varint,2,opt,name=StartHeight,proto3" json:"StartHeight,omitempty"`
	EndHeight            uint64   `protobuf:"varint,3,opt,name=EndHeight,proto3" json:"EndHeight,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxVoteReward) Reset()         { *m = TxVoteReward{} }
func (m *TxVoteReward) String() string { return proto.CompactTextString(m) }
func (*TxVoteReward) ProtoMessage()    {}
func (*TxVoteReward) Descriptor() ([]byte, []int) {
	return fileDescriptor_593fbc538261b21b, []int{0}
}
func (m *TxVoteReward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxVoteReward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxVoteReward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxVoteReward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxVoteReward.Merge(m, src)
}
func (m *TxVoteReward) XXX_Size() int {
	return m.Size()
}
func (m *TxVoteReward) XXX_DiscardUnknown() {
	xxx_messageInfo_TxVoteReward.DiscardUnknown(m)
}

var xxx_messageInfo_TxVoteReward proto.InternalMessageInfo

func (m *TxVoteReward) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxVoteReward) GetStartHeight() uint64 {
	if m != nil {
		return m.StartHeight
	}
	return 0
}

func (m *TxVoteReward) GetEndHeight() uint64 {
	if m != nil {
		return m.EndHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*TxVoteReward)(nil), "go_protos.TxVoteReward")
}

func init() { proto.RegisterFile("tx_vote_reward.proto", fileDescriptor_593fbc538261b21b) }

var fileDescriptor_593fbc538261b21b = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xa9, 0x88, 0x2f,
	0xcb, 0x2f, 0x49, 0x8d, 0x2f, 0x4a, 0x2d, 0x4f, 0x2c, 0x4a, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a, 0xa5, 0x78, 0x4a, 0x2a, 0x92, 0x12, 0x8b, 0x53,
	0x21, 0x12, 0x4a, 0x95, 0x5c, 0x3c, 0x21, 0x15, 0x61, 0xf9, 0x25, 0xa9, 0x41, 0x60, 0xe5, 0x42,
	0x9a, 0x5c, 0x6c, 0x21, 0x15, 0x4e, 0x89, 0xc5, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46,
	0x82, 0x7a, 0x70, 0x9d, 0x7a, 0x10, 0x89, 0x20, 0xa8, 0x02, 0x21, 0x05, 0x2e, 0xee, 0xe0, 0x92,
	0xc4, 0xa2, 0x12, 0x8f, 0xd4, 0xcc, 0xf4, 0x8c, 0x12, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x96, 0x20,
	0x64, 0x21, 0x21, 0x19, 0x2e, 0x4e, 0xd7, 0xbc, 0x14, 0xa8, 0x3c, 0x33, 0x58, 0x1e, 0x21, 0xe0,
	0x24, 0x7b, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0xce, 0x78,
	0x2c, 0xc7, 0x10, 0xc5, 0xad, 0xa7, 0x0f, 0xb7, 0x2d, 0x89, 0x0d, 0x4c, 0x1b, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x47, 0x13, 0x38, 0x92, 0xd1, 0x00, 0x00, 0x00,
}

func (m *TxVoteReward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxVoteReward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxVoteReward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.EndHeight != 0 {
		i = encodeVarintTxVoteReward(dAtA, i, uint64(m.EndHeight))
		i--
		dAtA[i] = 0x18
	}
	if m.StartHeight != 0 {
		i = encodeVarintTxVoteReward(dAtA, i, uint64(m.StartHeight))
		i--
		dAtA[i] = 0x10
	}
	if m.TxBase != nil {
		{
			size, err := m.TxBase.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTxVoteReward(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxVoteReward(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxVoteReward(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxVoteReward) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxVoteReward(uint64(l))
	}
	if m.StartHeight != 0 {
		n += 1 + sovTxVoteReward(uint64(m.StartHeight))
	}
	if m.EndHeight != 0 {
		n += 1 + sovTxVoteReward(uint64(m.EndHeight))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxVoteReward(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxVoteReward(x uint64) (n int) {
	return sovTxVoteReward(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxVoteReward) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxVoteReward
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
			return fmt.Errorf("proto: TxVoteReward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxVoteReward: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVoteReward
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
				return ErrInvalidLengthTxVoteReward
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxVoteReward
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartHeight", wireType)
			}
			m.StartHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVoteReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndHeight", wireType)
			}
			m.EndHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVoteReward
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTxVoteReward(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxVoteReward
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxVoteReward
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
func skipTxVoteReward(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxVoteReward
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
					return 0, ErrIntOverflowTxVoteReward
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
					return 0, ErrIntOverflowTxVoteReward
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
				return 0, ErrInvalidLengthTxVoteReward
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxVoteReward
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxVoteReward
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxVoteReward        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxVoteReward          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxVoteReward = fmt.Errorf("proto: unexpected end of group")
)
