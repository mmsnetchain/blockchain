package tx_name_out

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

type Tx_account struct {
	mining.TxBase
	Account []byte `json:"a"`
}

type Tx_account_VO struct {
	mining.TxBaseVO
	Account string `json:"a"`
}

func (this *Tx_account) GetVOJSON() interface{} {
	return Tx_account_VO{
		TxBaseVO: this.TxBase.ConversionVO(),
		Account:  string(this.Account),
	}
}

func (this *Tx_account) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_account_cancel)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_account) Proto() (*[]byte, error) {
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

	txPay := go_protos.TxNameOut{
		TxBase:  &txBase,
		Account: this.Account,
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *Tx_account) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)
	buf.Write(this.Account)
	*bs = buf.Bytes()
	return bs
}

func (this *Tx_account) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Account...)

	sign := keystore.Sign(*key, *signDst)

	return &sign
}

func (this *Tx_account) CheckSign() error {
	if this.Vin == nil || len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}
	if this.Vout_total != 1 {
		return config.ERROR_pay_vout_too_much
	}
	if this.Vout[0].Value <= 0 {
		return config.ERROR_amount_zero
	}

	oldTxClass := uint64(config.Wallet_tx_type_mining)

	for i, one := range this.Vin {

		sign := this.GetSignSerialize(nil, uint64(i))

		*sign = append(*sign, this.Account...)

		puk := ed25519.PublicKey(one.Puk)
		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*sign))
		}
		if !ed25519.Verify(puk, *sign, one.Sign) {
			return config.ERROR_sign_fail
		}

	}
	if this.Type == config.Wallet_tx_type_account_cancel && oldTxClass == config.Wallet_tx_type_account {
		return nil
	}

	return config.ERROR_tx_fail
}

func (this *Tx_account) GetWitness() *crypto.AddressCoin {
	witness := crypto.BuildAddr(config.AddrPre, this.Vin[0].Puk)

	return &witness
}

func (this *Tx_account) GetSpend() uint64 {
	return this.Gas
}

func (this *Tx_account) CheckRepeatedTx(txs ...mining.TxItr) bool {

	for _, one := range txs {
		if one.Class() != config.Wallet_tx_type_account_cancel {
			continue
		}
		ta := one.(*Tx_account)
		if bytes.Equal(ta.Account, this.Account) {
			return false
		}
	}
	return true
}

func (this *Tx_account) CountTxItemsNew(height uint64) *mining.TxItemCountMap {
	itemCount := mining.TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vout)+len(this.Vin)),
		Nonce:    make(map[string]big.Int),
	}

	totalValue := this.Gas

	from := this.Vin[0].GetPukToAddr()
	itemCount.Nonce[utils.Bytes2string(*from)] = this.Vin[0].Nonce
	frozenMap, ok := itemCount.AddItems[utils.Bytes2string(*from)]
	if ok {
		oldValue, ok := (*frozenMap)[0]
		if ok {
			oldValue -= int64(totalValue)
			(*frozenMap)[0] = oldValue
		} else {
			(*frozenMap)[0] = (0 - int64(totalValue))
		}
	} else {
		frozenMap := make(map[uint64]int64, 0)
		frozenMap[0] = (0 - int64(totalValue))
		itemCount.AddItems[utils.Bytes2string(*from)] = &frozenMap
	}
	return &itemCount
}
