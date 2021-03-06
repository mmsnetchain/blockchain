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

type RepeatedBlockHeadAndTxs struct {
	Bhat                 []*BlockHeadAndTxs `protobuf:"bytes,1,rep,name=bhat,proto3" json:"bhat,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *RepeatedBlockHeadAndTxs) Reset()         { *m = RepeatedBlockHeadAndTxs{} }
func (m *RepeatedBlockHeadAndTxs) String() string { return proto.CompactTextString(m) }
func (*RepeatedBlockHeadAndTxs) ProtoMessage()    {}
func (*RepeatedBlockHeadAndTxs) Descriptor() ([]byte, []int) {
	return fileDescriptor_900106ab601f1fc6, []int{0}
}
func (m *RepeatedBlockHeadAndTxs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RepeatedBlockHeadAndTxs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RepeatedBlockHeadAndTxs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RepeatedBlockHeadAndTxs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RepeatedBlockHeadAndTxs.Merge(m, src)
}
func (m *RepeatedBlockHeadAndTxs) XXX_Size() int {
	return m.Size()
}
func (m *RepeatedBlockHeadAndTxs) XXX_DiscardUnknown() {
	xxx_messageInfo_RepeatedBlockHeadAndTxs.DiscardUnknown(m)
}

var xxx_messageInfo_RepeatedBlockHeadAndTxs proto.InternalMessageInfo

func (m *RepeatedBlockHeadAndTxs) GetBhat() []*BlockHeadAndTxs {
	if m != nil {
		return m.Bhat
	}
	return nil
}

func init() {
	proto.RegisterType((*RepeatedBlockHeadAndTxs)(nil), "go_protos.RepeatedBlockHeadAndTxs")
}

func init() { proto.RegisterFile("blockheadandtxs_repeated.proto", fileDescriptor_900106ab601f1fc6) }

var fileDescriptor_900106ab601f1fc6 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xca, 0xc9, 0x4f,
	0xce, 0xce, 0x48, 0x4d, 0x4c, 0x49, 0xcc, 0x4b, 0x29, 0xa9, 0x28, 0x8e, 0x2f, 0x4a, 0x2d, 0x48,
	0x4d, 0x2c, 0x49, 0x4d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07,
	0xb3, 0x8a, 0xa5, 0x44, 0xd1, 0x94, 0x42, 0x54, 0x28, 0x79, 0x72, 0x89, 0x07, 0x41, 0xf5, 0x38,
	0x81, 0x14, 0x78, 0xa4, 0x26, 0xa6, 0x38, 0xe6, 0xa5, 0x84, 0x54, 0x14, 0x0b, 0xe9, 0x71, 0xb1,
	0x24, 0x65, 0x24, 0x96, 0x48, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0x49, 0xe9, 0xc1, 0xcd, 0xd2,
	0x43, 0x53, 0x19, 0x04, 0x56, 0xe7, 0x24, 0x7b, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c,
	0x0f, 0x1e, 0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7, 0x10, 0xc5, 0xad, 0xa7, 0x0f, 0xd7, 0x94, 0xc4,
	0x06, 0xa6, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1f, 0x03, 0xc6, 0xe9, 0xb4, 0x00, 0x00,
	0x00,
}

func (m *RepeatedBlockHeadAndTxs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RepeatedBlockHeadAndTxs) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RepeatedBlockHeadAndTxs) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Bhat) > 0 {
		for iNdEx := len(m.Bhat) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Bhat[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBlockheadandtxsRepeated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintBlockheadandtxsRepeated(dAtA []byte, offset int, v uint64) int {
	offset -= sovBlockheadandtxsRepeated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RepeatedBlockHeadAndTxs) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Bhat) > 0 {
		for _, e := range m.Bhat {
			l = e.Size()
			n += 1 + l + sovBlockheadandtxsRepeated(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBlockheadandtxsRepeated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBlockheadandtxsRepeated(x uint64) (n int) {
	return sovBlockheadandtxsRepeated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RepeatedBlockHeadAndTxs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBlockheadandtxsRepeated
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
			return fmt.Errorf("proto: RepeatedBlockHeadAndTxs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RepeatedBlockHeadAndTxs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bhat", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockheadandtxsRepeated
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
				return ErrInvalidLengthBlockheadandtxsRepeated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBlockheadandtxsRepeated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bhat = append(m.Bhat, &BlockHeadAndTxs{})
			if err := m.Bhat[len(m.Bhat)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBlockheadandtxsRepeated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBlockheadandtxsRepeated
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBlockheadandtxsRepeated
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
func skipBlockheadandtxsRepeated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBlockheadandtxsRepeated
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
					return 0, ErrIntOverflowBlockheadandtxsRepeated
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
					return 0, ErrIntOverflowBlockheadandtxsRepeated
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
				return 0, ErrInvalidLengthBlockheadandtxsRepeated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBlockheadandtxsRepeated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBlockheadandtxsRepeated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBlockheadandtxsRepeated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBlockheadandtxsRepeated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBlockheadandtxsRepeated = fmt.Errorf("proto: unexpected end of group")
)
