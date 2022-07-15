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

type BlockForWitness struct {
	GroupHeight          uint64   `protobuf:"varint,1,opt,name=GroupHeight,proto3" json:"GroupHeight,omitempty"`
	Addr                 []byte   `protobuf:"bytes,2,opt,name=Addr,proto3" json:"Addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockForWitness) Reset()         { *m = BlockForWitness{} }
func (m *BlockForWitness) String() string { return proto.CompactTextString(m) }
func (*BlockForWitness) ProtoMessage()    {}
func (*BlockForWitness) Descriptor() ([]byte, []int) {
	return fileDescriptor_59b39cd754ad3c42, []int{0}
}
func (m *BlockForWitness) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockForWitness) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockForWitness.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockForWitness) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockForWitness.Merge(m, src)
}
func (m *BlockForWitness) XXX_Size() int {
	return m.Size()
}
func (m *BlockForWitness) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockForWitness.DiscardUnknown(m)
}

var xxx_messageInfo_BlockForWitness proto.InternalMessageInfo

func (m *BlockForWitness) GetGroupHeight() uint64 {
	if m != nil {
		return m.GroupHeight
	}
	return 0
}

func (m *BlockForWitness) GetAddr() []byte {
	if m != nil {
		return m.Addr
	}
	return nil
}

func init() {
	proto.RegisterType((*BlockForWitness)(nil), "go_protos.BlockForWitness")
}

func init() { proto.RegisterFile("blockforwitness.proto", fileDescriptor_59b39cd754ad3c42) }

var fileDescriptor_59b39cd754ad3c42 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4d, 0xca, 0xc9, 0x4f,
	0xce, 0x4e, 0xcb, 0x2f, 0x2a, 0xcf, 0x2c, 0xc9, 0x4b, 0x2d, 0x2e, 0xd6, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a, 0x95, 0xdc, 0xb9, 0xf8, 0x9d, 0x40, 0x6a,
	0xdc, 0xf2, 0x8b, 0xc2, 0x21, 0x6a, 0x84, 0x14, 0xb8, 0xb8, 0xdd, 0x8b, 0xf2, 0x4b, 0x0b, 0x3c,
	0x52, 0x33, 0xd3, 0x33, 0x4a, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0x90, 0x85, 0x84, 0x84,
	0xb8, 0x58, 0x1c, 0x53, 0x52, 0x8a, 0x24, 0x98, 0x14, 0x18, 0x35, 0x78, 0x82, 0xc0, 0x6c, 0x27,
	0xd9, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc6, 0x63,
	0x39, 0x86, 0x28, 0x6e, 0x3d, 0x7d, 0xb8, 0x3d, 0x49, 0x6c, 0x60, 0xda, 0x18, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x68, 0x53, 0x11, 0x05, 0x92, 0x00, 0x00, 0x00,
}

func (m *BlockForWitness) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockForWitness) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockForWitness) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Addr) > 0 {
		i -= len(m.Addr)
		copy(dAtA[i:], m.Addr)
		i = encodeVarintBlockforwitness(dAtA, i, uint64(len(m.Addr)))
		i--
		dAtA[i] = 0x12
	}
	if m.GroupHeight != 0 {
		i = encodeVarintBlockforwitness(dAtA, i, uint64(m.GroupHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintBlockforwitness(dAtA []byte, offset int, v uint64) int {
	offset -= sovBlockforwitness(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BlockForWitness) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupHeight != 0 {
		n += 1 + sovBlockforwitness(uint64(m.GroupHeight))
	}
	l = len(m.Addr)
	if l > 0 {
		n += 1 + l + sovBlockforwitness(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBlockforwitness(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBlockforwitness(x uint64) (n int) {
	return sovBlockforwitness(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BlockForWitness) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBlockforwitness
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
			return fmt.Errorf("proto: BlockForWitness: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockForWitness: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupHeight", wireType)
			}
			m.GroupHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockforwitness
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockforwitness
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
				return ErrInvalidLengthBlockforwitness
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockforwitness
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addr = append(m.Addr[:0], dAtA[iNdEx:postIndex]...)
			if m.Addr == nil {
				m.Addr = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBlockforwitness(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBlockforwitness
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBlockforwitness
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
func skipBlockforwitness(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBlockforwitness
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
					return 0, ErrIntOverflowBlockforwitness
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
					return 0, ErrIntOverflowBlockforwitness
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
				return 0, ErrInvalidLengthBlockforwitness
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBlockforwitness
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBlockforwitness
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBlockforwitness        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBlockforwitness          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBlockforwitness = fmt.Errorf("proto: unexpected end of group")
)
