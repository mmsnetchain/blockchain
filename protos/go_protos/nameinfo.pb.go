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

type Nameinfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Txid                 []byte   `protobuf:"bytes,2,opt,name=Txid,proto3" json:"Txid,omitempty"`
	NetIds               [][]byte `protobuf:"bytes,3,rep,name=NetIds,proto3" json:"NetIds,omitempty"`
	AddrCoins            [][]byte `protobuf:"bytes,4,rep,name=AddrCoins,proto3" json:"AddrCoins,omitempty"`
	Height               uint64   `protobuf:"varint,5,opt,name=Height,proto3" json:"Height,omitempty"`
	NameOfValidity       uint64   `protobuf:"varint,6,opt,name=NameOfValidity,proto3" json:"NameOfValidity,omitempty"`
	Deposit              uint64   `protobuf:"varint,7,opt,name=Deposit,proto3" json:"Deposit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Nameinfo) Reset()         { *m = Nameinfo{} }
func (m *Nameinfo) String() string { return proto.CompactTextString(m) }
func (*Nameinfo) ProtoMessage()    {}
func (*Nameinfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_998d3ff6fede2178, []int{0}
}
func (m *Nameinfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Nameinfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Nameinfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Nameinfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nameinfo.Merge(m, src)
}
func (m *Nameinfo) XXX_Size() int {
	return m.Size()
}
func (m *Nameinfo) XXX_DiscardUnknown() {
	xxx_messageInfo_Nameinfo.DiscardUnknown(m)
}

var xxx_messageInfo_Nameinfo proto.InternalMessageInfo

func (m *Nameinfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Nameinfo) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *Nameinfo) GetNetIds() [][]byte {
	if m != nil {
		return m.NetIds
	}
	return nil
}

func (m *Nameinfo) GetAddrCoins() [][]byte {
	if m != nil {
		return m.AddrCoins
	}
	return nil
}

func (m *Nameinfo) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Nameinfo) GetNameOfValidity() uint64 {
	if m != nil {
		return m.NameOfValidity
	}
	return 0
}

func (m *Nameinfo) GetDeposit() uint64 {
	if m != nil {
		return m.Deposit
	}
	return 0
}

func init() {
	proto.RegisterType((*Nameinfo)(nil), "go_protos.Nameinfo")
}

func init() { proto.RegisterFile("nameinfo.proto", fileDescriptor_998d3ff6fede2178) }

var fileDescriptor_998d3ff6fede2178 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0x4b, 0xcc, 0x4d,
	0xcd, 0xcc, 0x4b, 0xcb, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xcf, 0x8f, 0x07,
	0xb3, 0x8a, 0x95, 0x0e, 0x31, 0x72, 0x71, 0xf8, 0x41, 0x65, 0x85, 0x84, 0xb8, 0x58, 0x40, 0x6c,
	0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x30, 0x1b, 0x24, 0x16, 0x52, 0x91, 0x99, 0x22, 0xc1,
	0xa4, 0xc0, 0xa8, 0xc1, 0x13, 0x04, 0x66, 0x0b, 0x89, 0x71, 0xb1, 0xf9, 0xa5, 0x96, 0x78, 0xa6,
	0x14, 0x4b, 0x30, 0x2b, 0x30, 0x6b, 0xf0, 0x04, 0x41, 0x79, 0x42, 0x32, 0x5c, 0x9c, 0x8e, 0x29,
	0x29, 0x45, 0xce, 0xf9, 0x99, 0x79, 0xc5, 0x12, 0x2c, 0x60, 0x29, 0x84, 0x00, 0x48, 0x97, 0x47,
	0x6a, 0x66, 0x7a, 0x46, 0x89, 0x04, 0xab, 0x02, 0xa3, 0x06, 0x4b, 0x10, 0x94, 0x27, 0xa4, 0xc6,
	0xc5, 0x07, 0xb2, 0xc9, 0x3f, 0x2d, 0x2c, 0x31, 0x27, 0x33, 0x25, 0xb3, 0xa4, 0x52, 0x82, 0x0d,
	0x2c, 0x8f, 0x26, 0x2a, 0x24, 0xc1, 0xc5, 0xee, 0x92, 0x5a, 0x90, 0x5f, 0x9c, 0x59, 0x22, 0xc1,
	0x0e, 0x56, 0x00, 0xe3, 0x3a, 0xc9, 0x9e, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83,
	0x47, 0x72, 0x8c, 0x33, 0x1e, 0xcb, 0x31, 0x44, 0x71, 0xeb, 0xe9, 0xc3, 0xfd, 0x98, 0xc4, 0x06,
	0xa6, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x70, 0x8b, 0xcb, 0x0f, 0x07, 0x01, 0x00, 0x00,
}

func (m *Nameinfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Nameinfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Nameinfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Deposit != 0 {
		i = encodeVarintNameinfo(dAtA, i, uint64(m.Deposit))
		i--
		dAtA[i] = 0x38
	}
	if m.NameOfValidity != 0 {
		i = encodeVarintNameinfo(dAtA, i, uint64(m.NameOfValidity))
		i--
		dAtA[i] = 0x30
	}
	if m.Height != 0 {
		i = encodeVarintNameinfo(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x28
	}
	if len(m.AddrCoins) > 0 {
		for iNdEx := len(m.AddrCoins) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AddrCoins[iNdEx])
			copy(dAtA[i:], m.AddrCoins[iNdEx])
			i = encodeVarintNameinfo(dAtA, i, uint64(len(m.AddrCoins[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.NetIds) > 0 {
		for iNdEx := len(m.NetIds) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.NetIds[iNdEx])
			copy(dAtA[i:], m.NetIds[iNdEx])
			i = encodeVarintNameinfo(dAtA, i, uint64(len(m.NetIds[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Txid) > 0 {
		i -= len(m.Txid)
		copy(dAtA[i:], m.Txid)
		i = encodeVarintNameinfo(dAtA, i, uint64(len(m.Txid)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintNameinfo(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintNameinfo(dAtA []byte, offset int, v uint64) int {
	offset -= sovNameinfo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Nameinfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovNameinfo(uint64(l))
	}
	l = len(m.Txid)
	if l > 0 {
		n += 1 + l + sovNameinfo(uint64(l))
	}
	if len(m.NetIds) > 0 {
		for _, b := range m.NetIds {
			l = len(b)
			n += 1 + l + sovNameinfo(uint64(l))
		}
	}
	if len(m.AddrCoins) > 0 {
		for _, b := range m.AddrCoins {
			l = len(b)
			n += 1 + l + sovNameinfo(uint64(l))
		}
	}
	if m.Height != 0 {
		n += 1 + sovNameinfo(uint64(m.Height))
	}
	if m.NameOfValidity != 0 {
		n += 1 + sovNameinfo(uint64(m.NameOfValidity))
	}
	if m.Deposit != 0 {
		n += 1 + sovNameinfo(uint64(m.Deposit))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovNameinfo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNameinfo(x uint64) (n int) {
	return sovNameinfo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Nameinfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNameinfo
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
			return fmt.Errorf("proto: Nameinfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Nameinfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthNameinfo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNameinfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Txid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
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
				return ErrInvalidLengthNameinfo
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthNameinfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Txid = append(m.Txid[:0], dAtA[iNdEx:postIndex]...)
			if m.Txid == nil {
				m.Txid = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetIds", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
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
				return ErrInvalidLengthNameinfo
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthNameinfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetIds = append(m.NetIds, make([]byte, postIndex-iNdEx))
			copy(m.NetIds[len(m.NetIds)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddrCoins", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
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
				return ErrInvalidLengthNameinfo
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthNameinfo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddrCoins = append(m.AddrCoins, make([]byte, postIndex-iNdEx))
			copy(m.AddrCoins[len(m.AddrCoins)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
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
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NameOfValidity", wireType)
			}
			m.NameOfValidity = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NameOfValidity |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposit", wireType)
			}
			m.Deposit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNameinfo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Deposit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipNameinfo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthNameinfo
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthNameinfo
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
func skipNameinfo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNameinfo
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
					return 0, ErrIntOverflowNameinfo
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
					return 0, ErrIntOverflowNameinfo
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
				return 0, ErrInvalidLengthNameinfo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNameinfo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNameinfo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNameinfo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNameinfo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNameinfo = fmt.Errorf("proto: unexpected end of group")
)
