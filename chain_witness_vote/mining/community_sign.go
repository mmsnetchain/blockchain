package mining

import (
	"crypto/ed25519"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"

	gogoproto "github.com/gogo/protobuf/proto"
)

const (
	SIGN_TYPE_community_reward = 1
)

type CommunitySign struct {
	Type        uint64
	StartHeight uint64
	EndHeight   uint64
	Rand        uint64
	Puk         []byte
	Sign        []byte
}

func (this *CommunitySign) Protobuf() []byte {
	csp := go_protos.CommunitySign{
		Type:        this.Type,
		StartHeight: this.StartHeight,
		EndHeight:   this.EndHeight,
		Rand:        this.Rand,
		Puk:         this.Puk,
		Sign:        this.Sign,
	}
	bs, _ := csp.Marshal()

	return bs
}

func ParseCommunitySign(bs []byte) (*CommunitySign, error) {
	csp := new(go_protos.CommunitySign)
	err := gogoproto.Unmarshal(bs, csp)
	if err != nil {
		return nil, err
	}

	cs := CommunitySign{
		Type:        csp.Type,
		StartHeight: csp.StartHeight,
		EndHeight:   csp.EndHeight,
		Rand:        csp.Rand,
		Puk:         csp.Puk,
		Sign:        csp.Sign,
	}
	return &cs, nil

}

func NewCommunitySign(puk []byte, startHeight, endHeight uint64) *CommunitySign {
	max := utils.BytesToUint64([]byte{255 - 128, 255, 255, 255, 255, 255, 255, 255})
	r := utils.GetRandNum(int64(max))
	return &CommunitySign{
		Type:        SIGN_TYPE_community_reward,
		StartHeight: startHeight,
		EndHeight:   endHeight,
		Rand:        uint64(r),
		Puk:         puk,
		Sign:        nil,
	}
}

func SignPayload(txItr TxItr, puk []byte, prk ed25519.PrivateKey, startHeight, endHeight uint64) TxItr {
	cs := NewCommunitySign(puk, startHeight, endHeight)
	txItr.SetPayload(cs.Protobuf())

	for i, _ := range *txItr.GetVin() {
		txItr.SetSign(uint64(i), nil)
	}
	signDst := txItr.Serialize()
	sign := keystore.Sign(prk, *signDst)
	cs.Sign = sign
	txItr.SetPayload(cs.Protobuf())
	return txItr
}

func CheckPayload(txItr TxItr) (crypto.AddressCoin, bool, *CommunitySign) {

	bs := txItr.GetPayload()
	if bs == nil || len(bs) <= 0 {
		return nil, false, nil
	}
	cs, err := ParseCommunitySign(bs)
	if err != nil {
		return nil, false, nil
	}
	if cs.Puk == nil || len(cs.Puk) <= 0 || cs.Sign == nil || len(cs.Sign) <= 0 {
		return nil, false, nil
	}
	addr := crypto.BuildAddr(config.AddrPre, cs.Puk)
	if cs.Type != SIGN_TYPE_community_reward {
		return nil, false, nil
	}
	signtmp := cs.Sign
	cs.Sign = nil
	txItr.SetPayload(cs.Protobuf())

	signs := make([][]byte, 0)

	for i, _ := range *txItr.GetVin() {

		signs = append(signs, (*txItr.GetVin())[i].Sign)
		txItr.SetSign(uint64(i), nil)
	}
	signDst := txItr.Serialize()
	cs.Sign = signtmp

	txItr.SetPayload(cs.Protobuf())
	for i, _ := range signs {
		txItr.SetSign(uint64(i), signs[i])
	}
	if !ed25519.Verify(cs.Puk, *signDst, cs.Sign) {
		return addr, false, nil
	}

	return addr, true, cs
}
