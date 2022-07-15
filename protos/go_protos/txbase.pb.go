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

type TxBase struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Type                 uint64   `protobuf:"varint,2,opt,name=Type,proto3" json:"Type,omitempty"`
	VinTotal             uint64   `protobuf:"varint,3,opt,name=Vin_total,json=VinTotal,proto3" json:"Vin_total,omitempty"`
	Vin                  []*Vin   `protobuf:"bytes,4,rep,name=Vin,proto3" json:"Vin,omitempty"`
	VoutTotal            uint64   `protobuf:"varint,5,opt,name=Vout_total,json=VoutTotal,proto3" json:"Vout_total,omitempty"`
	Vout                 []*Vout  `protobuf:"bytes,6,rep,name=Vout,proto3" json:"Vout,omitempty"`
	Gas                  uint64   `protobuf:"varint,7,opt,name=Gas,proto3" json:"Gas,omitempty"`
	LockHeight           uint64   `protobuf:"varint,8,opt,name=LockHeight,proto3" json:"LockHeight,omitempty"`
	Payload              []byte   `protobuf:"bytes,9,opt,name=Payload,proto3" json:"Payload,omitempty"`
	BlockHash            []byte   `protobuf:"bytes,10,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxBase) Reset()         { *m = TxBase{} }
func (m *TxBase) String() string { return proto.CompactTextString(m) }
func (*TxBase) ProtoMessage()    {}
func (*TxBase) Descriptor() ([]byte, []int) {
	return fileDescriptor_c44972010b73f2f4, []int{0}
}
func (m *TxBase) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxBase) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxBase.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxBase) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxBase.Merge(m, src)
}
func (m *TxBase) XXX_Size() int {
	return m.Size()
}
func (m *TxBase) XXX_DiscardUnknown() {
	xxx_messageInfo_TxBase.DiscardUnknown(m)
}

var xxx_messageInfo_TxBase proto.InternalMessageInfo

func (m *TxBase) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *TxBase) GetType() uint64 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *TxBase) GetVinTotal() uint64 {
	if m != nil {
		return m.VinTotal
	}
	return 0
}

func (m *TxBase) GetVin() []*Vin {
	if m != nil {
		return m.Vin
	}
	return nil
}

func (m *TxBase) GetVoutTotal() uint64 {
	if m != nil {
		return m.VoutTotal
	}
	return 0
}

func (m *TxBase) GetVout() []*Vout {
	if m != nil {
		return m.Vout
	}
	return nil
}

func (m *TxBase) GetGas() uint64 {
	if m != nil {
		return m.Gas
	}
	return 0
}

func (m *TxBase) GetLockHeight() uint64 {
	if m != nil {
		return m.LockHeight
	}
	return 0
}

func (m *TxBase) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *TxBase) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

type Vin struct {
	Puk                  []byte   `protobuf:"bytes,1,opt,name=Puk,proto3" json:"Puk,omitempty"`
	Sign                 []byte   `protobuf:"bytes,2,opt,name=Sign,proto3" json:"Sign,omitempty"`
	Nonce                []byte   `protobuf:"bytes,3,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vin) Reset()         { *m = Vin{} }
func (m *Vin) String() string { return proto.CompactTextString(m) }
func (*Vin) ProtoMessage()    {}
func (*Vin) Descriptor() ([]byte, []int) {
	return fileDescriptor_c44972010b73f2f4, []int{1}
}
func (m *Vin) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Vin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Vin.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Vin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vin.Merge(m, src)
}
func (m *Vin) XXX_Size() int {
	return m.Size()
}
func (m *Vin) XXX_DiscardUnknown() {
	xxx_messageInfo_Vin.DiscardUnknown(m)
}

var xxx_messageInfo_Vin proto.InternalMessageInfo

func (m *Vin) GetPuk() []byte {
	if m != nil {
		return m.Puk
	}
	return nil
}

func (m *Vin) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

func (m *Vin) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

type Vout struct {
	Value                uint64   `protobuf:"varint,1,opt,name=Value,proto3" json:"Value,omitempty"`
	Address              []byte   `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	FrozenHeight         uint64   `protobuf:"varint,3,opt,name=FrozenHeight,proto3" json:"FrozenHeight,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vout) Reset()         { *m = Vout{} }
func (m *Vout) String() string { return proto.CompactTextString(m) }
func (*Vout) ProtoMessage()    {}
func (*Vout) Descriptor() ([]byte, []int) {
	return fileDescriptor_c44972010b73f2f4, []int{2}
}
func (m *Vout) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Vout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Vout.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Vout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vout.Merge(m, src)
}
func (m *Vout) XXX_Size() int {
	return m.Size()
}
func (m *Vout) XXX_DiscardUnknown() {
	xxx_messageInfo_Vout.DiscardUnknown(m)
}

var xxx_messageInfo_Vout proto.InternalMessageInfo

func (m *Vout) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Vout) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Vout) GetFrozenHeight() uint64 {
	if m != nil {
		return m.FrozenHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*TxBase)(nil), "go_protos.TxBase")
	proto.RegisterType((*Vin)(nil), "go_protos.Vin")
	proto.RegisterType((*Vout)(nil), "go_protos.Vout")
}

func init() { proto.RegisterFile("txbase.proto", fileDescriptor_c44972010b73f2f4) }

var fileDescriptor_c44972010b73f2f4 = []byte{

	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x50, 0x4d, 0x4e, 0xc2, 0x40,
	0x18, 0xb5, 0xb4, 0xfc, 0xf4, 0xa3, 0x51, 0x33, 0x71, 0x31, 0x89, 0xd2, 0x34, 0x75, 0xc3, 0xaa,
	0x26, 0x7a, 0x02, 0x58, 0x28, 0x0b, 0x63, 0xc8, 0x48, 0xba, 0x60, 0x43, 0x06, 0x98, 0x40, 0x43,
	0x33, 0x43, 0x98, 0x69, 0x02, 0x9e, 0xc4, 0x8d, 0xf7, 0x71, 0xe9, 0x11, 0x0c, 0x5e, 0xc4, 0xcc,
	0xd7, 0x82, 0xb8, 0xe2, 0xfd, 0xf0, 0xbd, 0xe9, 0x7b, 0x10, 0x98, 0xed, 0x94, 0x6b, 0x91, 0xac,
	0x37, 0xca, 0x28, 0xe2, 0x2f, 0xd4, 0x04, 0x91, 0x8e, 0x3f, 0x6a, 0xd0, 0x18, 0x6d, 0xfb, 0x5c,
	0x0b, 0x42, 0xc0, 0x1b, 0x70, 0xbd, 0xa4, 0x4e, 0xe4, 0x74, 0x03, 0x86, 0xd8, 0x6a, 0xa3, 0xdd,
	0x5a, 0xd0, 0x5a, 0xe4, 0x74, 0x3d, 0x86, 0x98, 0x5c, 0x83, 0x9f, 0x66, 0x72, 0x62, 0x94, 0xe1,
	0x39, 0x75, 0xd1, 0x68, 0xa5, 0x99, 0x1c, 0x59, 0x4e, 0x22, 0x70, 0xd3, 0x4c, 0x52, 0x2f, 0x72,
	0xbb, 0xed, 0xfb, 0xf3, 0xe4, 0xf8, 0x50, 0x92, 0x66, 0x92, 0x59, 0x8b, 0x74, 0x00, 0x52, 0x55,
	0x98, 0xea, 0xbe, 0x8e, 0xf7, 0xbe, 0x55, 0xca, 0x80, 0x5b, 0xf0, 0x2c, 0xa1, 0x0d, 0x4c, 0xb8,
	0x38, 0x4d, 0x50, 0x85, 0x61, 0x68, 0x92, 0x4b, 0x70, 0x9f, 0xb8, 0xa6, 0x4d, 0x3c, 0xb6, 0x90,
	0x84, 0x00, 0xcf, 0x6a, 0xb6, 0x1a, 0x88, 0x6c, 0xb1, 0x34, 0xb4, 0x85, 0xc6, 0x89, 0x42, 0x28,
	0x34, 0x87, 0x7c, 0x97, 0x2b, 0x3e, 0xa7, 0x3e, 0xf6, 0x3b, 0x50, 0x72, 0x03, 0x7e, 0x3f, 0xb7,
	0x7f, 0xb4, 0xdd, 0x01, 0xbd, 0x3f, 0x21, 0xee, 0x61, 0x1f, 0xfb, 0xe0, 0xb0, 0x58, 0x55, 0xd3,
	0x58, 0x68, 0x97, 0x79, 0xcd, 0x16, 0x12, 0x97, 0x09, 0x18, 0x62, 0x72, 0x05, 0xf5, 0x17, 0x25,
	0x67, 0x02, 0x57, 0x09, 0x58, 0x49, 0xe2, 0x71, 0xd9, 0xc8, 0xba, 0x29, 0xcf, 0x0b, 0x81, 0x29,
	0x1e, 0x2b, 0x89, 0xfd, 0xb0, 0xde, 0x7c, 0xbe, 0x11, 0x5a, 0x57, 0x51, 0x07, 0x4a, 0x62, 0x08,
	0x1e, 0x37, 0xea, 0x4d, 0xc8, 0xaa, 0x54, 0x39, 0xf5, 0x3f, 0xad, 0xdf, 0xf9, 0xdc, 0x87, 0xce,
	0xd7, 0x3e, 0x74, 0xbe, 0xf7, 0xa1, 0xf3, 0xfe, 0x13, 0x9e, 0x8d, 0xdb, 0xc9, 0xdd, 0x71, 0xb2,
	0x69, 0x03, 0x7f, 0x1f, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0x27, 0xd6, 0x41, 0x1e, 0xff, 0x01,
	0x00, 0x00,
}

func (m *TxBase) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxBase) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxBase) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.BlockHash) > 0 {
		i -= len(m.BlockHash)
		copy(dAtA[i:], m.BlockHash)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.BlockHash)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x4a
	}
	if m.LockHeight != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.LockHeight))
		i--
		dAtA[i] = 0x40
	}
	if m.Gas != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.Gas))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Vout) > 0 {
		for iNdEx := len(m.Vout) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Vout[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTxbase(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.VoutTotal != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.VoutTotal))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Vin) > 0 {
		for iNdEx := len(m.Vin) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Vin[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTxbase(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.VinTotal != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.VinTotal))
		i--
		dAtA[i] = 0x18
	}
	if m.Type != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Vin) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Vin) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Vin) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Nonce) > 0 {
		i -= len(m.Nonce)
		copy(dAtA[i:], m.Nonce)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.Nonce)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Sign) > 0 {
		i -= len(m.Sign)
		copy(dAtA[i:], m.Sign)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.Sign)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Puk) > 0 {
		i -= len(m.Puk)
		copy(dAtA[i:], m.Puk)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.Puk)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Vout) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Vout) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Vout) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FrozenHeight != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.FrozenHeight))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTxbase(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.Value != 0 {
		i = encodeVarintTxbase(dAtA, i, uint64(m.Value))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTxbase(dAtA []byte, offset int, v uint64) int {
	offset -= sovTxbase(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxBase) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovTxbase(uint64(m.Type))
	}
	if m.VinTotal != 0 {
		n += 1 + sovTxbase(uint64(m.VinTotal))
	}
	if len(m.Vin) > 0 {
		for _, e := range m.Vin {
			l = e.Size()
			n += 1 + l + sovTxbase(uint64(l))
		}
	}
	if m.VoutTotal != 0 {
		n += 1 + sovTxbase(uint64(m.VoutTotal))
	}
	if len(m.Vout) > 0 {
		for _, e := range m.Vout {
			l = e.Size()
			n += 1 + l + sovTxbase(uint64(l))
		}
	}
	if m.Gas != 0 {
		n += 1 + sovTxbase(uint64(m.Gas))
	}
	if m.LockHeight != 0 {
		n += 1 + sovTxbase(uint64(m.LockHeight))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	l = len(m.BlockHash)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Vin) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Puk)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	l = len(m.Sign)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	l = len(m.Nonce)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Vout) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Value != 0 {
		n += 1 + sovTxbase(uint64(m.Value))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTxbase(uint64(l))
	}
	if m.FrozenHeight != 0 {
		n += 1 + sovTxbase(uint64(m.FrozenHeight))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTxbase(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTxbase(x uint64) (n int) {
	return sovTxbase(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxBase) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxbase
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
			return fmt.Errorf("proto: TxBase: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxBase: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
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
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VinTotal", wireType)
			}
			m.VinTotal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VinTotal |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vin", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vin = append(m.Vin, &Vin{})
			if err := m.Vin[len(m.Vin)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoutTotal", wireType)
			}
			m.VoutTotal = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VoutTotal |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vout", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vout = append(m.Vout, &Vout{})
			if err := m.Vout[len(m.Vout)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gas", wireType)
			}
			m.Gas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Gas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockHeight", wireType)
			}
			m.LockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BlockHash = append(m.BlockHash[:0], dAtA[iNdEx:postIndex]...)
			if m.BlockHash == nil {
				m.BlockHash = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxbase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxbase
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxbase
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
func (m *Vin) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxbase
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
			return fmt.Errorf("proto: Vin: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vin: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Puk", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Puk = append(m.Puk[:0], dAtA[iNdEx:postIndex]...)
			if m.Puk == nil {
				m.Puk = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sign", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sign = append(m.Sign[:0], dAtA[iNdEx:postIndex]...)
			if m.Sign == nil {
				m.Sign = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nonce = append(m.Nonce[:0], dAtA[iNdEx:postIndex]...)
			if m.Nonce == nil {
				m.Nonce = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTxbase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxbase
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxbase
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
func (m *Vout) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTxbase
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
			return fmt.Errorf("proto: Vout: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vout: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			m.Value = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Value |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
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
				return ErrInvalidLengthTxbase
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTxbase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FrozenHeight", wireType)
			}
			m.FrozenHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTxbase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FrozenHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTxbase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTxbase
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTxbase
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
func skipTxbase(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTxbase
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
					return 0, ErrIntOverflowTxbase
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
					return 0, ErrIntOverflowTxbase
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
				return 0, ErrInvalidLengthTxbase
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTxbase
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTxbase
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTxbase        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTxbase          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTxbase = fmt.Errorf("proto: unexpected end of group")
)
