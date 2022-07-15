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

type TxNameIn struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	Account              []byte   `protobuf:"bytes,2,opt,name=Account,proto3" json:"Account,omitempty"`
	NetIds               [][]byte `protobuf:"bytes,3,rep,name=NetIds,proto3" json:"NetIds,omitempty"`
	NetIdsMerkleHash     []byte   `protobuf:"bytes,4,opt,name=NetIdsMerkleHash,proto3" json:"NetIdsMerkleHash,omitempty"`
	AddrCoins            [][]byte `protobuf:"bytes,5,rep,name=AddrCoins,proto3" json:"AddrCoins,omitempty"`
	AddrCoinsMerkleHash  []byte   `protobuf:"bytes,6,opt,name=AddrCoinsMerkleHash,proto3" json:"AddrCoinsMerkleHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxNameIn) Reset()         { *m = TxNameIn{} }
func (m *TxNameIn) String() string { return proto.CompactTextString(m) }
func (*TxNameIn) ProtoMessage()    {}
func (*TxNameIn) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d450409f40e2d34, []int{0}
}
func (m *TxNameIn) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxNameIn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxNameIn.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxNameIn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxNameIn.Merge(m, src)
}
func (m *TxNameIn) XXX_Size() int {
	return m.Size()
}
func (m *TxNameIn) XXX_DiscardUnknown() {
	xxx_messageInfo_TxNameIn.DiscardUnknown(m)
}

var xxx_messageInfo_TxNameIn proto.InternalMessageInfo

func (m *TxNameIn) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxNameIn) GetAccount() []byte {
	if m != nil {
		return m.Account
	}
	return nil
}

func (m *TxNameIn) GetNetIds() [][]byte {
	if m != nil {
		return m.NetIds
	}
	return nil
}

func (m *TxNameIn) GetNetIdsMerkleHash() []byte {
	if m != nil {
		return m.NetIdsMerkleHash
	}
	return nil
}

func (m *TxNameIn) GetAddrCoins() [][]byte {
	if m != nil {
		return m.AddrCoins
	}
	return nil
}

func (m *TxNameIn) GetAddrCoinsMerkleHash() []byte {
	if m != nil {
		return m.AddrCoinsMerkleHash
	}
	return nil
}

type TxNameOut struct {
	TxBase               *TxBase  `protobuf:"bytes,1,opt,name=TxBase,proto3" json:"TxBase,omitempty"`
	Account              []byte   `protobuf:"bytes,2,opt,name=Account,proto3" json:"Account,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxNameOut) Reset()         { *m = TxNameOut{} }
func (m *TxNameOut) String() string { return proto.CompactTextString(m) }
func (*TxNameOut) ProtoMessage()    {}
func (*TxNameOut) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d450409f40e2d34, []int{1}
}
func (m *TxNameOut) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxNameOut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxNameOut.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxNameOut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxNameOut.Merge(m, src)
}
func (m *TxNameOut) XXX_Size() int {
	return m.Size()
}
func (m *TxNameOut) XXX_DiscardUnknown() {
	xxx_messageInfo_TxNameOut.DiscardUnknown(m)
}

var xxx_messageInfo_TxNameOut proto.InternalMessageInfo

func (m *TxNameOut) GetTxBase() *TxBase {
	if m != nil {
		return m.TxBase
	}
	return nil
}

func (m *TxNameOut) GetAccount() []byte {
	if m != nil {
		return m.Account
	}
	return nil
}

func init() {
	proto.RegisterType((*TxNameIn)(nil), "go_protos.TxNameIn")
	proto.RegisterType((*TxNameOut)(nil), "go_protos.TxNameOut")
}

func init() { proto.RegisterFile("tx_name.proto", fileDescriptor_4d450409f40e2d34) }

var fileDescriptor_4d450409f40e2d34 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0xa9, 0x88, 0xcf,
	0x4b, 0xcc, 0x4d, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07, 0xb3,
	0x8a, 0xa5, 0x78, 0x4a, 0x2a, 0x92, 0x12, 0x8b, 0xa1, 0x12, 0x4a, 0x8f, 0x19, 0xb9, 0x38, 0x42,
	0x2a, 0xfc, 0x12, 0x73, 0x53, 0x3d, 0xf3, 0x84, 0x34, 0xb9, 0xd8, 0x42, 0x2a, 0x9c, 0x12, 0x8b,
	0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0x04, 0xf5, 0xe0, 0xda, 0xf4, 0x20, 0x12, 0x41,
	0x50, 0x05, 0x42, 0x12, 0x5c, 0xec, 0x8e, 0xc9, 0xc9, 0xf9, 0xa5, 0x79, 0x25, 0x12, 0x4c, 0x0a,
	0x8c, 0x1a, 0x3c, 0x41, 0x30, 0xae, 0x90, 0x18, 0x17, 0x9b, 0x5f, 0x6a, 0x89, 0x67, 0x4a, 0xb1,
	0x04, 0xb3, 0x02, 0xb3, 0x06, 0x4f, 0x10, 0x94, 0x27, 0xa4, 0xc5, 0x25, 0x00, 0x61, 0xf9, 0xa6,
	0x16, 0x65, 0xe7, 0xa4, 0x7a, 0x24, 0x16, 0x67, 0x48, 0xb0, 0x80, 0xb5, 0x62, 0x88, 0x0b, 0xc9,
	0x70, 0x71, 0x3a, 0xa6, 0xa4, 0x14, 0x39, 0xe7, 0x67, 0xe6, 0x15, 0x4b, 0xb0, 0x82, 0x8d, 0x41,
	0x08, 0x08, 0x19, 0x70, 0x09, 0xc3, 0x39, 0x48, 0x86, 0xb1, 0x81, 0x0d, 0xc3, 0x26, 0xa5, 0x14,
	0xc0, 0xc5, 0x09, 0xf1, 0xa4, 0x7f, 0x69, 0x09, 0x55, 0x7c, 0xe9, 0x24, 0x7b, 0xe2, 0x91, 0x1c,
	0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7, 0x10, 0xc5, 0xad,
	0xa7, 0x0f, 0x37, 0x27, 0x89, 0x0d, 0x4c, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xbc, 0xbe,
	0x0c, 0x7c, 0x87, 0x01, 0x00, 0x00,
}

func (m *TxNameIn) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxNameIn) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxNameIn) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.AddrCoinsMerkleHash) > 0 {
		i -= len(m.AddrCoinsMerkleHash)
		copy(dAtA[i:], m.AddrCoinsMerkleHash)
		i = encodeVarintTxName(dAtA, i, uint64(len(m.AddrCoinsMerkleHash)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.AddrCoins) > 0 {
		for iNdEx := len(m.AddrCoins) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AddrCoins[iNdEx])
			copy(dAtA[i:], m.AddrCoins[iNdEx])
			i = encodeVarintTxName(dAtA, i, uint64(len(m.AddrCoins[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.NetIdsMerkleHash) > 0 {
		i -= len(m.NetIdsMerkleHash)
		copy(dAtA[i:], m.NetIdsMerkleHash)
		i = encodeVarintTxName(dAtA, i, uint64(len(m.NetIdsMerkleHash)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.NetIds) > 0 {
		for iNdEx := len(m.NetIds) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.NetIds[iNdEx])
			copy(dAtA[i:], m.NetIds[iNdEx])
			i = encodeVarintTxName(dAtA, i, uint64(len(m.NetIds[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarintTxName(dAtA, i, uint64(len(m.Account)))
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
			i = encodeVarintTxName(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TxNameOut) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxNameOut) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxNameOut) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Account) > 0 {
		i -= len(m.Account)
		copy(dAtA[i:], m.Account)
		i = encodeVarintTxName(dAtA, i, uint64(len(m.Account)))
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
			i = encodeVarintTxName(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxName(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxName(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxNameIn) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxName(uint64(l))
	}
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sovTxName(uint64(l))
	}
	if len(m.NetIds) > 0 {
		for _, b := range m.NetIds {
			l = len(b)
			n += 1 + l + sovTxName(uint64(l))
		}
	}
	l = len(m.NetIdsMerkleHash)
	if l > 0 {
		n += 1 + l + sovTxName(uint64(l))
	}
	if len(m.AddrCoins) > 0 {
		for _, b := range m.AddrCoins {
			l = len(b)
			n += 1 + l + sovTxName(uint64(l))
		}
	}
	l = len(m.AddrCoinsMerkleHash)
	if l > 0 {
		n += 1 + l + sovTxName(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TxNameOut) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxBase != nil {
		l = m.TxBase.Size()
		n += 1 + l + sovTxName(uint64(l))
	}
	l = len(m.Account)
	if l > 0 {
		n += 1 + l + sovTxName(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxName(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxName(x uint64) (n int) {
	return sovTxName(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxNameIn) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxName
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
			return fmt.Errorf("proto: TxNameIn: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxNameIn: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
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
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Account = append(m.Account[:0], dAtA[iNdEx:postIndex]...)
			if m.Account == nil {
				m.Account = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetIds", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetIds = append(m.NetIds, make([]byte, postIndex-iNdEx))
			copy(m.NetIds[len(m.NetIds)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetIdsMerkleHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetIdsMerkleHash = append(m.NetIdsMerkleHash[:0], dAtA[iNdEx:postIndex]...)
			if m.NetIdsMerkleHash == nil {
				m.NetIdsMerkleHash = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddrCoins", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddrCoins = append(m.AddrCoins, make([]byte, postIndex-iNdEx))
			copy(m.AddrCoins[len(m.AddrCoins)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddrCoinsMerkleHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddrCoinsMerkleHash = append(m.AddrCoinsMerkleHash[:0], dAtA[iNdEx:postIndex]...)
			if m.AddrCoinsMerkleHash == nil {
				m.AddrCoinsMerkleHash = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxName(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxName
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxName
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
func (m *TxNameOut) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxName
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
			return fmt.Errorf("proto: TxNameOut: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxNameOut: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBase", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
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
				return fmt.Errorf("proto: wrong wireType = %d for field Account", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxName
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
				return ErrInvalidLengthTxName
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxName
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Account = append(m.Account[:0], dAtA[iNdEx:postIndex]...)
			if m.Account == nil {
				m.Account = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxName(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxName
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxName
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
func skipTxName(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxName
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
					return 0, ErrIntOverflowTxName
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
					return 0, ErrIntOverflowTxName
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
				return 0, ErrInvalidLengthTxName
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxName
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxName
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxName        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxName          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxName = fmt.Errorf("proto: unexpected end of group")
)
