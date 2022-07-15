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

type BlockHead struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Height               uint64   `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	GroupHeight          uint64   `protobuf:"varint,3,opt,name=GroupHeight,proto3" json:"GroupHeight,omitempty"`
	GroupHeightGrowth    uint64   `protobuf:"varint,4,opt,name=GroupHeightGrowth,proto3" json:"GroupHeightGrowth,omitempty"`
	Previousblockhash    []byte   `protobuf:"bytes,5,opt,name=Previousblockhash,proto3" json:"Previousblockhash,omitempty"`
	Nextblockhash        []byte   `protobuf:"bytes,6,opt,name=Nextblockhash,proto3" json:"Nextblockhash,omitempty"`
	NTx                  uint64   `protobuf:"varint,7,opt,name=NTx,proto3" json:"NTx,omitempty"`
	MerkleRoot           []byte   `protobuf:"bytes,8,opt,name=MerkleRoot,proto3" json:"MerkleRoot,omitempty"`
	Tx                   [][]byte `protobuf:"bytes,9,rep,name=Tx,proto3" json:"Tx,omitempty"`
	Time                 int64    `protobuf:"varint,10,opt,name=Time,proto3" json:"Time,omitempty"`
	Witness              []byte   `protobuf:"bytes,11,opt,name=Witness,proto3" json:"Witness,omitempty"`
	Sign                 []byte   `protobuf:"bytes,12,opt,name=Sign,proto3" json:"Sign,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockHead) Reset()         { *m = BlockHead{} }
func (m *BlockHead) String() string { return proto.CompactTextString(m) }
func (*BlockHead) ProtoMessage()    {}
func (*BlockHead) Descriptor() ([]byte, []int) {
	return fileDescriptor_73a80c2de40f4050, []int{0}
}
func (m *BlockHead) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockHead) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockHead.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockHead) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHead.Merge(m, src)
}
func (m *BlockHead) XXX_Size() int {
	return m.Size()
}
func (m *BlockHead) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHead.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHead proto.InternalMessageInfo

func (m *BlockHead) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockHead) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHead) GetGroupHeight() uint64 {
	if m != nil {
		return m.GroupHeight
	}
	return 0
}

func (m *BlockHead) GetGroupHeightGrowth() uint64 {
	if m != nil {
		return m.GroupHeightGrowth
	}
	return 0
}

func (m *BlockHead) GetPreviousblockhash() []byte {
	if m != nil {
		return m.Previousblockhash
	}
	return nil
}

func (m *BlockHead) GetNextblockhash() []byte {
	if m != nil {
		return m.Nextblockhash
	}
	return nil
}

func (m *BlockHead) GetNTx() uint64 {
	if m != nil {
		return m.NTx
	}
	return 0
}

func (m *BlockHead) GetMerkleRoot() []byte {
	if m != nil {
		return m.MerkleRoot
	}
	return nil
}

func (m *BlockHead) GetTx() [][]byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *BlockHead) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *BlockHead) GetWitness() []byte {
	if m != nil {
		return m.Witness
	}
	return nil
}

func (m *BlockHead) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

func init() {
	proto.RegisterType((*BlockHead)(nil), "go_protos.BlockHead")
}

func init() { proto.RegisterFile("blockhead.proto", fileDescriptor_73a80c2de40f4050) }

var fileDescriptor_73a80c2de40f4050 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x51, 0xcd, 0x4e, 0x83, 0x40,
	0x10, 0x76, 0xa1, 0x52, 0x19, 0xf0, 0x6f, 0x0e, 0x66, 0x2e, 0x12, 0x62, 0x3c, 0x70, 0x30, 0xf5,
	0xe0, 0x1b, 0xf4, 0x52, 0x2e, 0x36, 0x06, 0x49, 0x4c, 0xbc, 0x18, 0x6a, 0x37, 0x40, 0x5a, 0xbb,
	0x0d, 0x6c, 0x95, 0x47, 0xf1, 0x91, 0x3c, 0xfa, 0x06, 0x1a, 0x7c, 0x11, 0xb3, 0xd3, 0xaa, 0x98,
	0x9e, 0xf6, 0xfb, 0xdb, 0xdd, 0x2f, 0x33, 0x70, 0x38, 0x99, 0xab, 0xc7, 0x59, 0x21, 0xb3, 0xe9,
	0x60, 0x59, 0x29, 0xad, 0xd0, 0xcd, 0xd5, 0x03, 0xa3, 0xfa, 0xec, 0xc3, 0x02, 0x77, 0x68, 0xec,
	0x58, 0x66, 0x53, 0x44, 0xe8, 0xc5, 0x59, 0x5d, 0x90, 0x08, 0x45, 0xe4, 0x27, 0x8c, 0xf1, 0x04,
	0x9c, 0x58, 0x96, 0x79, 0xa1, 0xc9, 0x0a, 0x45, 0xd4, 0x4b, 0x36, 0x0c, 0x43, 0xf0, 0x46, 0x95,
	0x5a, 0x2d, 0x37, 0xa6, 0xcd, 0x66, 0x57, 0xc2, 0x0b, 0x38, 0xee, 0xd0, 0x51, 0xa5, 0x5e, 0x74,
	0x41, 0x3d, 0xce, 0x6d, 0x1b, 0x26, 0x7d, 0x53, 0xc9, 0xe7, 0x52, 0xad, 0xea, 0x75, 0x5f, 0x53,
	0x64, 0x97, 0x8b, 0x6c, 0x1b, 0x78, 0x0e, 0xfb, 0x63, 0xd9, 0xe8, 0xbf, 0xa4, 0xc3, 0xc9, 0xff,
	0x22, 0x1e, 0x81, 0x3d, 0x4e, 0x1b, 0xea, 0xf3, 0x9f, 0x06, 0x62, 0x00, 0x70, 0x2d, 0xab, 0xd9,
	0x5c, 0x26, 0x4a, 0x69, 0xda, 0xe3, 0x4b, 0x1d, 0x05, 0x0f, 0xc0, 0x4a, 0x1b, 0x72, 0x43, 0x3b,
	0xf2, 0x13, 0x2b, 0x6d, 0xcc, 0x44, 0xd2, 0xf2, 0x49, 0x12, 0x84, 0x22, 0xb2, 0x13, 0xc6, 0x48,
	0xd0, 0xbf, 0x2b, 0xf5, 0x42, 0xd6, 0x35, 0x79, 0xfc, 0xc0, 0x0f, 0x35, 0xe9, 0xdb, 0x32, 0x5f,
	0x90, 0xbf, 0x9e, 0x9f, 0xc1, 0xc3, 0xd3, 0xb7, 0x36, 0x10, 0xef, 0x6d, 0x20, 0x3e, 0xdb, 0x40,
	0xbc, 0x7e, 0x05, 0x3b, 0xf7, 0xde, 0xe0, 0xf2, 0x77, 0x01, 0x13, 0x87, 0xcf, 0xab, 0xef, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xd6, 0x60, 0xd7, 0x36, 0xa5, 0x01, 0x00, 0x00,
}

func (m *BlockHead) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockHead) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockHead) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Sign) > 0 {
		i -= len(m.Sign)
		copy(dAtA[i:], m.Sign)
		i = encodeVarintBlockhead(dAtA, i, uint64(len(m.Sign)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.Witness) > 0 {
		i -= len(m.Witness)
		copy(dAtA[i:], m.Witness)
		i = encodeVarintBlockhead(dAtA, i, uint64(len(m.Witness)))
		i--
		dAtA[i] = 0x5a
	}
	if m.Time != 0 {
		i = encodeVarintBlockhead(dAtA, i, uint64(m.Time))
		i--
		dAtA[i] = 0x50
	}
	if len(m.Tx) > 0 {
		for iNdEx := len(m.Tx) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Tx[iNdEx])
			copy(dAtA[i:], m.Tx[iNdEx])
			i = encodeVarintBlockhead(dAtA, i, uint64(len(m.Tx[iNdEx])))
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.MerkleRoot) > 0 {
		i -= len(m.MerkleRoot)
		copy(dAtA[i:], m.MerkleRoot)
		i = encodeVarintBlockhead(dAtA, i, uint64(len(m.MerkleRoot)))
		i--
		dAtA[i] = 0x42
	}
	if m.NTx != 0 {
		i = encodeVarintBlockhead(dAtA, i, uint64(m.NTx))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Nextblockhash) > 0 {
		i -= len(m.Nextblockhash)
		copy(dAtA[i:], m.Nextblockhash)
		i = encodeVarintBlockhead(dAtA, i, uint64(len(m.Nextblockhash)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Previousblockhash) > 0 {
		i -= len(m.Previousblockhash)
		copy(dAtA[i:], m.Previousblockhash)
		i = encodeVarintBlockhead(dAtA, i, uint64(len(m.Previousblockhash)))
		i--
		dAtA[i] = 0x2a
	}
	if m.GroupHeightGrowth != 0 {
		i = encodeVarintBlockhead(dAtA, i, uint64(m.GroupHeightGrowth))
		i--
		dAtA[i] = 0x20
	}
	if m.GroupHeight != 0 {
		i = encodeVarintBlockhead(dAtA, i, uint64(m.GroupHeight))
		i--
		dAtA[i] = 0x18
	}
	if m.Height != 0 {
		i = encodeVarintBlockhead(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintBlockhead(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBlockhead(dAtA []byte, offset int, v uint64) int {
	offset -= sovBlockhead(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BlockHead) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovBlockhead(uint64(l))
	}
	if m.Height != 0 {
		n += 1 + sovBlockhead(uint64(m.Height))
	}
	if m.GroupHeight != 0 {
		n += 1 + sovBlockhead(uint64(m.GroupHeight))
	}
	if m.GroupHeightGrowth != 0 {
		n += 1 + sovBlockhead(uint64(m.GroupHeightGrowth))
	}
	l = len(m.Previousblockhash)
	if l > 0 {
		n += 1 + l + sovBlockhead(uint64(l))
	}
	l = len(m.Nextblockhash)
	if l > 0 {
		n += 1 + l + sovBlockhead(uint64(l))
	}
	if m.NTx != 0 {
		n += 1 + sovBlockhead(uint64(m.NTx))
	}
	l = len(m.MerkleRoot)
	if l > 0 {
		n += 1 + l + sovBlockhead(uint64(l))
	}
	if len(m.Tx) > 0 {
		for _, b := range m.Tx {
			l = len(b)
			n += 1 + l + sovBlockhead(uint64(l))
		}
	}
	if m.Time != 0 {
		n += 1 + sovBlockhead(uint64(m.Time))
	}
	l = len(m.Witness)
	if l > 0 {
		n += 1 + l + sovBlockhead(uint64(l))
	}
	l = len(m.Sign)
	if l > 0 {
		n += 1 + l + sovBlockhead(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBlockhead(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBlockhead(x uint64) (n int) {
	return sovBlockhead(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BlockHead) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBlockhead
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
			return fmt.Errorf("proto: BlockHead: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockHead: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupHeight", wireType)
			}
			m.GroupHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupHeightGrowth", wireType)
			}
			m.GroupHeightGrowth = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupHeightGrowth |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Previousblockhash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Previousblockhash = append(m.Previousblockhash[:0], dAtA[iNdEx:postIndex]...)
			if m.Previousblockhash == nil {
				m.Previousblockhash = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nextblockhash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nextblockhash = append(m.Nextblockhash[:0], dAtA[iNdEx:postIndex]...)
			if m.Nextblockhash == nil {
				m.Nextblockhash = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NTx", wireType)
			}
			m.NTx = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NTx |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MerkleRoot", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MerkleRoot = append(m.MerkleRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.MerkleRoot == nil {
				m.MerkleRoot = []byte{}
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tx", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tx = append(m.Tx, make([]byte, postIndex-iNdEx))
			copy(m.Tx[len(m.Tx)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			m.Time = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Time |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Witness", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Witness = append(m.Witness[:0], dAtA[iNdEx:postIndex]...)
			if m.Witness == nil {
				m.Witness = []byte{}
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sign", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBlockhead
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
				return ErrInvalidLengthBlockhead
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBlockhead
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sign = append(m.Sign[:0], dAtA[iNdEx:postIndex]...)
			if m.Sign == nil {
				m.Sign = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBlockhead(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBlockhead
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBlockhead
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
func skipBlockhead(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBlockhead
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
					return 0, ErrIntOverflowBlockhead
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
					return 0, ErrIntOverflowBlockhead
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
				return 0, ErrInvalidLengthBlockhead
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBlockhead
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBlockhead
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBlockhead        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBlockhead          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBlockhead = fmt.Errorf("proto: unexpected end of group")
)