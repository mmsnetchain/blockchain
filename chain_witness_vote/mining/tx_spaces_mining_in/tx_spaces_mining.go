package tx_spaces_mining_in

import (
	"bytes"
	"encoding/binary"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"golang.org/x/crypto/ed25519"
)

type Tx_SpacesMining struct {
	mining.TxBase
	NetId nodeStore.AddressNet `json:"netid"`
}

type Tx_SpacesMining_VO struct {
	mining.TxBaseVO
	NetId string `json:"netid"`
}

func (this *Tx_SpacesMining) GetVOJSON() interface{} {
	return Tx_SpacesMining_VO{
		TxBaseVO: this.TxBase.ConversionVO(),
		NetId:    this.NetId.B58String(),
	}
}

func (this *Tx_SpacesMining) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_spaces_mining_in)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_SpacesMining) Json() (*[]byte, error) {
	bs, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *Tx_SpacesMining) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)

	buf.Write(this.NetId)

	*bs = buf.Bytes()
	return bs
}

func (this *Tx_SpacesMining) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.NetId...)

	sign := keystore.Sign(*key, *signDst)

	return &sign

}

func (this *Tx_SpacesMining) CheckSign() error {

	if len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}

	inTotal := uint64(0)
	for i, one := range this.Vin {

		sign := this.GetSignSerialize(nil, uint64(i))

		*sign = append(*sign, this.NetId...)

		puk := ed25519.PublicKey(one.Puk)
		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*sign))
		}
		if !ed25519.Verify(puk, *sign, one.Sign) {
			return config.ERROR_sign_fail
		}
	}

	outTotal := uint64(0)
	for _, one := range this.Vout {
		outTotal = outTotal + one.Value
	}

	if outTotal > inTotal {
		return config.ERROR_tx_fail
	}
	this.Gas = inTotal - outTotal

	return config.ERROR_tx_fail
}

func (this *Tx_SpacesMining) GetWitness() *crypto.AddressCoin {
	witness := crypto.BuildAddr(config.AddrPre, this.Vin[0].Puk)

	return &witness
}

func (this *Tx_SpacesMining) CheckRepeatedTx(txs ...mining.TxItr) bool {

	for _, one := range txs {
		if one.Class() != config.Wallet_tx_type_spaces_mining_in {
			continue
		}
		tsm := one.(*Tx_SpacesMining)
		if bytes.Equal(this.NetId, tsm.NetId) {
			return false
		}
	}
	return true
}
