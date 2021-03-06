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

type StorePeer struct {
	AddrNet              []byte   `protobuf:"bytes,1,opt,name=AddrNet,proto3" json:"AddrNet,omitempty"`
	SpaceNum             uint64   `protobuf:"varint,2,opt,name=SpaceNum,proto3" json:"SpaceNum,omitempty"`
	AddrCoin             []byte   `protobuf:"bytes,3,opt,name=AddrCoin,proto3" json:"AddrCoin,omitempty"`
	FlashTime            int64    `protobuf:"varint,4,opt,name=FlashTime,proto3" json:"FlashTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StorePeer) Reset()         { *m = StorePeer{} }
func (m *StorePeer) String() string { return proto.CompactTextString(m) }
func (*StorePeer) ProtoMessage()    {}
func (*StorePeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_cddbabfb23f01628, []int{0}
}
func (m *StorePeer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StorePeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StorePeer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StorePeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorePeer.Merge(m, src)
}
func (m *StorePeer) XXX_Size() int {
	return m.Size()
}
func (m *StorePeer) XXX_DiscardUnknown() {
	xxx_messageInfo_StorePeer.DiscardUnknown(m)
}

var xxx_messageInfo_StorePeer proto.InternalMessageInfo

func (m *StorePeer) GetAddrNet() []byte {
	if m != nil {
		return m.AddrNet
	}
	return nil
}

func (m *StorePeer) GetSpaceNum() uint64 {
	if m != nil {
		return m.SpaceNum
	}
	return 0
}

func (m *StorePeer) GetAddrCoin() []byte {
	if m != nil {
		return m.AddrCoin
	}
	return nil
}

func (m *StorePeer) GetFlashTime() int64 {
	if m != nil {
		return m.FlashTime
	}
	return 0
}

func init() {
	proto.RegisterType((*StorePeer)(nil), "go_protos.StorePeer")
}

func init() { proto.RegisterFile("storepeer.proto", fileDescriptor_cddbabfb23f01628) }

var fileDescriptor_cddbabfb23f01628 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x2e, 0xc9, 0x2f,
	0x4a, 0x2d, 0x48, 0x4d, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f,
	0x07, 0xb3, 0x8a, 0x95, 0xaa, 0xb9, 0x38, 0x83, 0x41, 0xb2, 0x01, 0xa9, 0xa9, 0x45, 0x42, 0x12,
	0x5c, 0xec, 0x8e, 0x29, 0x29, 0x45, 0x7e, 0xa9, 0x25, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41,
	0x30, 0xae, 0x90, 0x14, 0x17, 0x47, 0x70, 0x41, 0x62, 0x72, 0xaa, 0x5f, 0x69, 0xae, 0x04, 0x93,
	0x02, 0xa3, 0x06, 0x4b, 0x10, 0x9c, 0x0f, 0x92, 0x03, 0x29, 0x73, 0xce, 0xcf, 0xcc, 0x93, 0x60,
	0x06, 0x6b, 0x83, 0xf3, 0x85, 0x64, 0xb8, 0x38, 0xdd, 0x72, 0x12, 0x8b, 0x33, 0x42, 0x32, 0x73,
	0x53, 0x25, 0x58, 0x14, 0x18, 0x35, 0x98, 0x83, 0x10, 0x02, 0x4e, 0xb2, 0x27, 0x1e, 0xc9, 0x31,
	0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x8c, 0xc7, 0x72, 0x0c, 0x51, 0xdc, 0x7a,
	0xfa, 0x70, 0xb7, 0x25, 0xb1, 0x81, 0x69, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe2, 0x68,
	0x99, 0xf6, 0xc0, 0x00, 0x00, 0x00,
}

func (m *StorePeer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StorePeer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StorePeer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FlashTime != 0 {
		i = encodeVarintStorepeer(dAtA, i, uint64(m.FlashTime))
		i--
		dAtA[i] = 0x20
	}
	if len(m.AddrCoin) > 0 {
		i -= len(m.AddrCoin)
		copy(dAtA[i:], m.AddrCoin)
		i = encodeVarintStorepeer(dAtA, i, uint64(len(m.AddrCoin)))
		i--
		dAtA[i] = 0x1a
	}
	if m.SpaceNum != 0 {
		i = encodeVarintStorepeer(dAtA, i, uint64(m.SpaceNum))
		i--
		dAtA[i] = 0x10
	}
	if len(m.AddrNet) > 0 {
		i -= len(m.AddrNet)
		copy(dAtA[i:], m.AddrNet)
		i = encodeVarintStorepeer(dAtA, i, uint64(len(m.AddrNet)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintStorepeer(dAtA []byte, offset int, v uint64) int {
	offset -= sovStorepeer(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StorePeer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AddrNet)
	if l > 0 {
		n += 1 + l + sovStorepeer(uint64(l))
	}
	if m.SpaceNum != 0 {
		n += 1 + sovStorepeer(uint64(m.SpaceNum))
	}
	l = len(m.AddrCoin)
	if l > 0 {
		n += 1 + l + sovStorepeer(uint64(l))
	}
	if m.FlashTime != 0 {
		n += 1 + sovStorepeer(uint64(m.FlashTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovStorepeer(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStorepeer(x uint64) (n int) {
	return sovStorepeer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StorePeer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStorepeer
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
			return fmt.Errorf("proto: StorePeer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StorePeer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddrNet", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorepeer
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
				return ErrInvalidLengthStorepeer
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthStorepeer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddrNet = append(m.AddrNet[:0], dAtA[iNdEx:postIndex]...)
			if m.AddrNet == nil {
				m.AddrNet = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpaceNum", wireType)
			}
			m.SpaceNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorepeer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SpaceNum |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddrCoin", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorepeer
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
				return ErrInvalidLengthStorepeer
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthStorepeer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddrCoin = append(m.AddrCoin[:0], dAtA[iNdEx:postIndex]...)
			if m.AddrCoin == nil {
				m.AddrCoin = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FlashTime", wireType)
			}
			m.FlashTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStorepeer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FlashTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStorepeer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStorepeer
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthStorepeer
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
func skipStorepeer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStorepeer
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
					return 0, ErrIntOverflowStorepeer
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
					return 0, ErrIntOverflowStorepeer
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
				return 0, ErrInvalidLengthStorepeer
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStorepeer
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStorepeer
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStorepeer        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStorepeer          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStorepeer = fmt.Errorf("proto: unexpected end of group")
)
