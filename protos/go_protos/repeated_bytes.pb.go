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

type RepeatedBytes struct {
	Bss                  [][]byte `protobuf:"bytes,1,rep,name=bss,proto3" json:"bss,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RepeatedBytes) Reset()         { *m = RepeatedBytes{} }
func (m *RepeatedBytes) String() string { return proto.CompactTextString(m) }
func (*RepeatedBytes) ProtoMessage()    {}
func (*RepeatedBytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee7e1a7b81ad15a0, []int{0}
}
func (m *RepeatedBytes) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RepeatedBytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RepeatedBytes.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RepeatedBytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RepeatedBytes.Merge(m, src)
}
func (m *RepeatedBytes) XXX_Size() int {
	return m.Size()
}
func (m *RepeatedBytes) XXX_DiscardUnknown() {
	xxx_messageInfo_RepeatedBytes.DiscardUnknown(m)
}

var xxx_messageInfo_RepeatedBytes proto.InternalMessageInfo

func (m *RepeatedBytes) GetBss() [][]byte {
	if m != nil {
		return m.Bss
	}
	return nil
}

func init() {
	proto.RegisterType((*RepeatedBytes)(nil), "go_protos.RepeatedBytes")
}

func init() { proto.RegisterFile("repeated_bytes.proto", fileDescriptor_ee7e1a7b81ad15a0) }

var fileDescriptor_ee7e1a7b81ad15a0 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x4a, 0x2d, 0x48,
	0x4d, 0x2c, 0x49, 0x4d, 0x89, 0x4f, 0xaa, 0x2c, 0x49, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3, 0x8a, 0x95, 0x14, 0xb9, 0x78, 0x83, 0xa0, 0x4a, 0x9c,
	0x40, 0x2a, 0x84, 0x04, 0xb8, 0x98, 0x93, 0x8a, 0x8b, 0x25, 0x18, 0x15, 0x98, 0x35, 0x78, 0x82,
	0x40, 0x4c, 0x27, 0xd9, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e,
	0x71, 0xc6, 0x63, 0x39, 0x86, 0x28, 0x6e, 0x3d, 0x7d, 0xb8, 0x09, 0x49, 0x6c, 0x60, 0xda, 0x18,
	0x10, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xd2, 0x22, 0x23, 0x6b, 0x00, 0x00, 0x00,
}

func (m *RepeatedBytes) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RepeatedBytes) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RepeatedBytes) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Bss) > 0 {
		for iNdEx := len(m.Bss) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Bss[iNdEx])
			copy(dAtA[i:], m.Bss[iNdEx])
			i = encodeVarintRepeatedBytes(dAtA, i, uint64(len(m.Bss[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintRepeatedBytes(dAtA []byte, offset int, v uint64) int {
	offset -= sovRepeatedBytes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RepeatedBytes) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Bss) > 0 {
		for _, b := range m.Bss {
			l = len(b)
			n += 1 + l + sovRepeatedBytes(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovRepeatedBytes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRepeatedBytes(x uint64) (n int) {
	return sovRepeatedBytes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RepeatedBytes) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRepeatedBytes
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
			return fmt.Errorf("proto: RepeatedBytes: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RepeatedBytes: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bss", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRepeatedBytes
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
				return ErrInvalidLengthRepeatedBytes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthRepeatedBytes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bss = append(m.Bss, make([]byte, postIndex-iNdEx))
			copy(m.Bss[len(m.Bss)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRepeatedBytes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRepeatedBytes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRepeatedBytes
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
func skipRepeatedBytes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRepeatedBytes
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
					return 0, ErrIntOverflowRepeatedBytes
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
					return 0, ErrIntOverflowRepeatedBytes
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
				return 0, ErrInvalidLengthRepeatedBytes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRepeatedBytes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRepeatedBytes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRepeatedBytes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRepeatedBytes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRepeatedBytes = fmt.Errorf("proto: unexpected end of group")
)
