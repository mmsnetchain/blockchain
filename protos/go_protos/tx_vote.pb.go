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

type TxVoteIn struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	Vote                 []byte   `protobuf:"bytes,2,opt,name=Vote,proto3" json:"Vote,omitempty"`
	VoteType             uint32   `protobuf:"varint,3,opt,name=VoteType,proto3" json:"VoteType,omitempty"`
	VoteAddr             []byte   `protobuf:"bytes,4,opt,name=VoteAddr,proto3" json:"VoteAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxVoteIn) Reset()         { *m = TxVoteIn{} }
func (m *TxVoteIn) String() string { return proto.CompactTextString(m) }
func (*TxVoteIn) ProtoMessage()    {}
func (*TxVoteIn) Descriptor() ([]byte, []int) {
	return fileDescriptor_82ebbf5ccea440ca, []int{0}
}
func (m *TxVoteIn) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxVoteIn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxVoteIn.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxVoteIn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxVoteIn.Merge(m, src)
}
func (m *TxVoteIn) XXX_Size() int {
	return m.Size()
}
func (m *TxVoteIn) XXX_DiscardUnknown() {
	xxx_messageInfo_TxVoteIn.DiscardUnknown(m)
}

var xxx_messageInfo_TxVoteIn proto.InternalMessageInfo

func (m *TxVoteIn) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxVoteIn) GetVote() []byte {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *TxVoteIn) GetVoteType() uint32 {
	if m != nil {
		return m.VoteType
	}
	return 0
}

func (m *TxVoteIn) GetVoteAddr() []byte {
	if m != nil {
		return m.VoteAddr
	}
	return nil
}

type TxVoteOut struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	Vote                 []byte   `protobuf:"bytes,2,opt,name=Vote,proto3" json:"Vote,omitempty"`
	VoteType             uint32   `protobuf:"varint,3,opt,name=VoteType,proto3" json:"VoteType,omitempty"`
	VoteAddr             []byte   `protobuf:"bytes,4,opt,name=VoteAddr,proto3" json:"VoteAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxVoteOut) Reset()         { *m = TxVoteOut{} }
func (m *TxVoteOut) String() string { return proto.CompactTextString(m) }
func (*TxVoteOut) ProtoMessage()    {}
func (*TxVoteOut) Descriptor() ([]byte, []int) {
	return fileDescriptor_82ebbf5ccea440ca, []int{1}
}
func (m *TxVoteOut) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxVoteOut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxVoteOut.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxVoteOut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxVoteOut.Merge(m, src)
}
func (m *TxVoteOut) XXX_Size() int {
	return m.Size()
}
func (m *TxVoteOut) XXX_DiscardUnknown() {
	xxx_messageInfo_TxVoteOut.DiscardUnknown(m)
}

var xxx_messageInfo_TxVoteOut proto.InternalMessageInfo

func (m *TxVoteOut) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxVoteOut) GetVote() []byte {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *TxVoteOut) GetVoteType() uint32 {
	if m != nil {
		return m.VoteType
	}
	return 0
}

func (m *TxVoteOut) GetVoteAddr() []byte {
	if m != nil {
		return m.VoteAddr
	}
	return nil
}

func init() {
	proto.RegisterType((*TxVoteIn)(nil), "go_protos.TxVoteIn")
	proto.RegisterType((*TxVoteOut)(nil), "go_protos.TxVoteOut")
}

func init() { proto.RegisterFile("tx_vote.proto", fileDescriptor_82ebbf5ccea440ca) }

var fileDescriptor_82ebbf5ccea440ca = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0xa9, 0x88, 0x2f,
	0xcb, 0x2f, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3,
	0x8a, 0xa5, 0x78, 0x4a, 0x2a, 0x92, 0x12, 0x8b, 0xa1, 0x12, 0x4a, 0x8d, 0x8c, 0x5c, 0x1c, 0x21,
	0x15, 0x61, 0xf9, 0x25, 0xa9, 0x9e, 0x79, 0x42, 0x9a, 0x5c, 0x6c, 0x21, 0x15, 0x4e, 0x89, 0xc5,
	0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x82, 0x7a, 0x70, 0x6d, 0x7a, 0x10, 0x89, 0x20,
	0xa8, 0x02, 0x21, 0x21, 0x2e, 0x16, 0x90, 0x26, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x30,
	0x5b, 0x48, 0x8a, 0x8b, 0x03, 0x44, 0x87, 0x54, 0x16, 0xa4, 0x4a, 0x30, 0x2b, 0x30, 0x6a, 0xf0,
	0x06, 0xc1, 0xf9, 0x30, 0x39, 0xc7, 0x94, 0x94, 0x22, 0x09, 0x16, 0xb0, 0x1e, 0x38, 0x5f, 0xa9,
	0x89, 0x91, 0x8b, 0x13, 0xe2, 0x06, 0xff, 0xd2, 0x92, 0x01, 0x72, 0x84, 0x93, 0xec, 0x89, 0x47,
	0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe3, 0xb1, 0x1c, 0x43, 0x14,
	0xb7, 0x9e, 0x3e, 0xdc, 0xe6, 0x24, 0x36, 0x30, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x3f,
	0xc9, 0x91, 0x08, 0x58, 0x01, 0x00, 0x00,
}

func (m *TxVoteIn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxVoteIn) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxVoteIn) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.VoteAddr) > 0 {
		i -= len(m.VoteAddr)
		copy(dAtA[i:], m.VoteAddr)
		i = encodeVarintTxVote(dAtA, i, uint64(len(m.VoteAddr)))
		i--
		dAtA[i] = 0x22
	}
	if m.VoteType != 0 {
		i = encodeVarintTxVote(dAtA, i, uint64(m.VoteType))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Vote) > 0 {
		i -= len(m.Vote)
		copy(dAtA[i:], m.Vote)
		i = encodeVarintTxVote(dAtA, i, uint64(len(m.Vote)))
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
			i = encodeVarintTxVote(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TxVoteOut) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxVoteOut) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxVoteOut) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.VoteAddr) > 0 {
		i -= len(m.VoteAddr)
		copy(dAtA[i:], m.VoteAddr)
		i = encodeVarintTxVote(dAtA, i, uint64(len(m.VoteAddr)))
		i--
		dAtA[i] = 0x22
	}
	if m.VoteType != 0 {
		i = encodeVarintTxVote(dAtA, i, uint64(m.VoteType))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Vote) > 0 {
		i -= len(m.Vote)
		copy(dAtA[i:], m.Vote)
		i = encodeVarintTxVote(dAtA, i, uint64(len(m.Vote)))
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
			i = encodeVarintTxVote(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxVote(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxVote(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxVoteIn) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxVote(uint64(l))
	}
	l = len(m.Vote)
	if l > 0 {
		n += 1 + l + sovTxVote(uint64(l))
	}
	if m.VoteType != 0 {
		n += 1 + sovTxVote(uint64(m.VoteType))
	}
	l = len(m.VoteAddr)
	if l > 0 {
		n += 1 + l + sovTxVote(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TxVoteOut) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxVote(uint64(l))
	}
	l = len(m.Vote)
	if l > 0 {
		n += 1 + l + sovTxVote(uint64(l))
	}
	if m.VoteType != 0 {
		n += 1 + sovTxVote(uint64(m.VoteType))
	}
	l = len(m.VoteAddr)
	if l > 0 {
		n += 1 + l + sovTxVote(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxVote(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxVote(x uint64) (n int) {
	return sovTxVote(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxVoteIn) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxVote
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
			return fmt.Errorf("proto: TxVoteIn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxVoteIn: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
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
				return ErrInvalidLengthTxVote
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxVote
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
				return fmt.Errorf("proto: wrong wireType = %d for field Vote", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
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
				return ErrInvalidLengthTxVote
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vote = append(m.Vote[:0], dAtA[iNdEx:postIndex]...)
			if m.Vote == nil {
				m.Vote = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteType", wireType)
			}
			m.VoteType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VoteType |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteAddr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
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
				return ErrInvalidLengthTxVote
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VoteAddr = append(m.VoteAddr[:0], dAtA[iNdEx:postIndex]...)
			if m.VoteAddr == nil {
				m.VoteAddr = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxVote
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
func (m *TxVoteOut) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxVote
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
			return fmt.Errorf("proto: TxVoteOut: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxVoteOut: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
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
				return ErrInvalidLengthTxVote
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxVote
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
				return fmt.Errorf("proto: wrong wireType = %d for field Vote", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
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
				return ErrInvalidLengthTxVote
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vote = append(m.Vote[:0], dAtA[iNdEx:postIndex]...)
			if m.Vote == nil {
				m.Vote = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteType", wireType)
			}
			m.VoteType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VoteType |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteAddr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxVote
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
				return ErrInvalidLengthTxVote
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VoteAddr = append(m.VoteAddr[:0], dAtA[iNdEx:postIndex]...)
			if m.VoteAddr == nil {
				m.VoteAddr = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxVote
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
func skipTxVote(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxVote
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
					return 0, ErrIntOverflowTxVote
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
					return 0, ErrIntOverflowTxVote
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
				return 0, ErrInvalidLengthTxVote
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxVote
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxVote
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxVote        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxVote          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxVote = fmt.Errorf("proto: unexpected end of group")
)
