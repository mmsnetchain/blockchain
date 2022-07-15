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

type ChainInfo struct {
	StartBlockHash       []byte   `protobuf:"bytes,1,opt,name=StartBlockHash,proto3" json:"StartBlockHash,omitempty"`
	HightBlock           uint64   `protobuf:"varint,2,opt,name=HightBlock,proto3" json:"HightBlock,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChainInfo) Reset()         { *m = ChainInfo{} }
func (m *ChainInfo) String() string { return proto.CompactTextString(m) }
func (*ChainInfo) ProtoMessage()    {}
func (*ChainInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06c819fdb62c369, []int{0}
}
func (m *ChainInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainInfo.Merge(m, src)
}
func (m *ChainInfo) XXX_Size() int {
	return m.Size()
}
func (m *ChainInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChainInfo proto.InternalMessageInfo

func (m *ChainInfo) GetStartBlockHash() []byte {
	if m != nil {
		return m.StartBlockHash
	}
	return nil
}

func (m *ChainInfo) GetHightBlock() uint64 {
	if m != nil {
		return m.HightBlock
	}
	return 0
}

func init() {
	proto.RegisterType((*ChainInfo)(nil), "go_protos.ChainInfo")
}

func init() { proto.RegisterFile("chaininfo.proto", fileDescriptor_c06c819fdb62c369) }

var fileDescriptor_c06c819fdb62c369 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0xce, 0x48, 0xcc,
	0xcc, 0xcb, 0xcc, 0x4b, 0xcb, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f,
	0x07, 0xb3, 0x8a, 0x95, 0x82, 0xb9, 0x38, 0x9d, 0x41, 0xb2, 0x9e, 0x79, 0x69, 0xf9, 0x42, 0x6a,
	0x5c, 0x7c, 0xc1, 0x25, 0x89, 0x45, 0x25, 0x4e, 0x39, 0xf9, 0xc9, 0xd9, 0x1e, 0x89, 0xc5, 0x19,
	0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x68, 0xa2, 0x42, 0x72, 0x5c, 0x5c, 0x1e, 0x99, 0xe9,
	0x19, 0x10, 0x11, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x24, 0x11, 0x27, 0xd9, 0x13, 0x8f,
	0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc6, 0x63, 0x39, 0x86, 0x28,
	0x6e, 0x3d, 0x7d, 0xb8, 0x9d, 0x49, 0x6c, 0x60, 0xda, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x3a,
	0xab, 0xfe, 0xc9, 0x98, 0x00, 0x00, 0x00,
}

func (m *ChainInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.HightBlock != 0 {
		i = encodeVarintChaininfo(dAtA, i, uint64(m.HightBlock))
		i--
		dAtA[i] = 0x10
	}
	if len(m.StartBlockHash) > 0 {
		i -= len(m.StartBlockHash)
		copy(dAtA[i:], m.StartBlockHash)
		i = encodeVarintChaininfo(dAtA, i, uint64(len(m.StartBlockHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintChaininfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovChaininfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChainInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.StartBlockHash)
	if l > 0 {
		n += 1 + l + sovChaininfo(uint64(l))
	}
	if m.HightBlock != 0 {
		n += 1 + sovChaininfo(uint64(m.HightBlock))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovChaininfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChaininfo(x uint64) (n int) {
	return sovChaininfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChainInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChaininfo
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
			return fmt.Errorf("proto: ChainInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartBlockHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChaininfo
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
				return ErrInvalidLengthChaininfo
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthChaininfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StartBlockHash = append(m.StartBlockHash[:0], dAtA[iNdEx:postIndex]...)
			if m.StartBlockHash == nil {
				m.StartBlockHash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HightBlock", wireType)
			}
			m.HightBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChaininfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HightBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipChaininfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthChaininfo
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthChaininfo
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
func skipChaininfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChaininfo
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
					return 0, ErrIntOverflowChaininfo
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
					return 0, ErrIntOverflowChaininfo
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
				return 0, ErrInvalidLengthChaininfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChaininfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChaininfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChaininfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChaininfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChaininfo = fmt.Errorf("proto: unexpected end of group")
)
