package publish

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"golang.org/x/crypto/ed25519"
)

type TxToken struct {
	mining.TxBase
	Token_name       string        `json:"token_name"`
	Token_symbol     string        `json:"token_symbol"`
	Token_supply     uint64        `json:"token_supply"`
	Token_Vout_total uint64        `json:"token_vout_total"`
	Token_Vout       []mining.Vout `json:"token_vout"`
}

type TxToken_VO struct {
	mining.TxBaseVO
	Token_name       string           `json:"token_name"`
	Token_symbol     string           `json:"token_symbol"`
	Token_supply     uint64           `json:"token_supply"`
	Token_Vout_total uint64           `json:"token_vout_total"`
	Token_Vout       []*mining.VoutVO `json:"token_vout"`
}

func (this *TxToken) GetVOJSON() interface{} {

	vouts := make([]*mining.VoutVO, 0)
	for _, one := range this.Token_Vout {
		vouts = append(vouts, one.ConversionVO())
	}

	return TxToken_VO{
		TxBaseVO:         this.TxBase.ConversionVO(),
		Token_name:       this.Token_name,
		Token_symbol:     this.Token_symbol,
		Token_supply:     this.Token_supply,
		Token_Vout_total: this.Token_Vout_total,
		Token_Vout:       vouts,
	}
}

func (this *TxToken) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, Wallet_tx_class)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *TxToken) Proto() (*[]byte, error) {
	vins := make([]*go_protos.Vin, 0)
	for _, one := range this.Vin {
		vins = append(vins, &go_protos.Vin{

			Puk:   one.Puk,
			Sign:  one.Sign,
			Nonce: one.Nonce.Bytes(),
		})
	}
	vouts := make([]*go_protos.Vout, 0)
	for _, one := range this.Vout {
		vouts = append(vouts, &go_protos.Vout{
			Value:        one.Value,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		})
	}
	txBase := go_protos.TxBase{
		Hash:       this.Hash,
		Type:       this.Type,
		VinTotal:   this.Vin_total,
		Vin:        vins,
		VoutTotal:  this.Vout_total,
		Vout:       vouts,
		Gas:        this.Gas,
		LockHeight: this.LockHeight,
		Payload:    this.Payload,
		BlockHash:  this.BlockHash,
	}

	tokenVouts := make([]*go_protos.Vout, 0)
	for _, one := range this.Token_Vout {
		tokenVouts = append(tokenVouts, &go_protos.Vout{
			Value:        one.Value,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		})
	}

	txPay := go_protos.TxTokenPublish{
		TxBase:          &txBase,
		TokenName:       this.Token_name,
		TokenSymbol:     this.Token_symbol,
		TokenSupply:     this.Token_supply,
		Token_VoutTotal: this.Token_Vout_total,
		Token_Vout:      tokenVouts,
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *TxToken) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)

	buf.Write([]byte(this.Token_name))
	buf.Write([]byte(this.Token_symbol))
	buf.Write(utils.Uint64ToBytes(this.Token_supply))

	buf.Write(utils.Uint64ToBytes(this.Token_Vout_total))
	if this.Token_Vout != nil {
		for _, one := range this.Token_Vout {
			bs := one.Serialize()
			buf.Write(*bs)
		}
	}

	*bs = buf.Bytes()
	return bs
}

func (this *TxToken) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Token_name...)
	*signDst = append(*signDst, this.Token_symbol...)
	*signDst = append(*signDst, utils.Uint64ToBytes(this.Token_supply)...)

	*signDst = append(*signDst, utils.Uint64ToBytes(this.Token_Vout_total)...)
	if this.Token_Vout != nil {
		for _, one := range this.Token_Vout {
			*signDst = append(*signDst, *one.Serialize()...)
		}
	}

	sign := keystore.Sign(*key, *signDst)

	return &sign

}

func (this *TxToken) GetTokenVoutSignSerialize(voutIndex uint64) *[]byte {
	bufVout := bytes.NewBuffer(nil)

	bufVout.Write(utils.Uint64ToBytes(voutIndex))
	vout := this.Token_Vout[voutIndex]
	bs := vout.Serialize()
	bufVout.Write(*vout.Serialize())
	*bs = bufVout.Bytes()
	return bs
}

func (this *TxToken) CheckSign() error {

	if len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}
	if this.Vout_total != 0 {
		return config.ERROR_pay_vout_too_much
	}

	for i, one := range this.Vin {

		sign := this.GetSignSerialize(nil, uint64(i))

		*sign = append(*sign, this.Token_name...)
		*sign = append(*sign, this.Token_symbol...)
		*sign = append(*sign, utils.Uint64ToBytes(this.Token_supply)...)

		*sign = append(*sign, utils.Uint64ToBytes(this.Token_Vout_total)...)
		if this.Token_Vout != nil {
			for _, one := range this.Token_Vout {
				*sign = append(*sign, *one.Serialize()...)
			}
		}

		puk := ed25519.PublicKey(one.Puk)
		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*sign))
		}
		if !ed25519.Verify(puk, *sign, one.Sign) {
			return config.ERROR_sign_fail
		}
	}

	return nil

}

func (this *TxToken) GetWitness() *crypto.AddressCoin {
	witness := crypto.BuildAddr(config.AddrPre, this.Vin[0].Puk)

	return &witness
}

func (this *TxToken) GetSpend() uint64 {
	return this.Gas
}

func (this *TxToken) CheckRepeatedTx(txs ...mining.TxItr) bool {

	for _, one := range txs {
		if one.Class() != Wallet_tx_class {
			continue
		}

	}
	return true
}

func (this *TxToken) CountTxItemsNew(height uint64) *mining.TxItemCountMap {
	itemCount := mining.TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vin)),
		Nonce:    make(map[string]big.Int),
	}

	totalValue := this.Gas

	from := this.Vin[0].GetPukToAddr()
	itemCount.Nonce[utils.Bytes2string(*from)] = this.Vin[0].Nonce

	frozenMap := make(map[uint64]int64, 0)
	frozenMap[0] = (0 - int64(totalValue))
	itemCount.AddItems[utils.Bytes2string(*from)] = &frozenMap

	return &itemCount
}
