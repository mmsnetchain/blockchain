package mining

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type Tx_reward struct {
	TxBase
}

func (this *Tx_reward) GetVOJSON() interface{} {
	return this.TxBase.ConversionVO()
}

func (this *Tx_reward) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	sign := keystore.Sign(*key, *signDst)

	return &sign
}

func (this *Tx_reward) CheckSign() error {

	if this.Vin == nil || len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) != 0 {

		return config.ERROR_pay_nonce_is_nil
	}

	one := this.Vin[0]
	signDst := this.GetSignSerialize(nil, uint64(0))

	puk := ed25519.PublicKey(one.Puk)

	if config.Wallet_print_serialize_hex {
		engine.Log.Info("sign serialize:%s", hex.EncodeToString(*signDst))
	}
	if !ed25519.Verify(puk, *signDst, one.Sign) {
		return config.ERROR_sign_fail
	}

	outTotal := uint64(0)
	for _, one := range this.Vout {
		outTotal = outTotal + one.Value
	}

	return nil
}

func (this *Tx_reward) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {

		return
	}

	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_mining)
	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)

}

func (this *Tx_reward) GetSpend() uint64 {
	return 0
}

func (this *Tx_reward) CheckRepeatedTx(txs ...TxItr) bool {
	return true
}

func (this *Tx_reward) CountTxItemsNew(height uint64) *TxItemCountMap {
	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vout)),
		Nonce:    make(map[string]big.Int),
	}

	for _, vout := range this.Vout {

		frozenMap, ok := itemCount.AddItems[utils.Bytes2string(vout.Address)]
		if ok {
			oldValue, ok := (*frozenMap)[vout.FrozenHeight]
			if ok {
				oldValue += int64(vout.Value)
				(*frozenMap)[vout.FrozenHeight] = oldValue
			} else {
				(*frozenMap)[vout.FrozenHeight] = int64(vout.Value)
			}
		} else {
			frozenMap := make(map[uint64]int64, 0)
			frozenMap[vout.FrozenHeight] = int64(vout.Value)
			itemCount.AddItems[utils.Bytes2string(vout.Address)] = &frozenMap
		}
	}

	return &itemCount
}

func (this *Tx_reward) CountTxHistory(height uint64) {

	hiIn := HistoryItem{
		IsIn: true,
		Type: this.Class(),

		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
	}

	addrCoin := make(map[string]bool)

	for _, vout := range this.Vout {

		_, ok := keystore.FindAddress(vout.Address)
		if !ok {
			continue
		}
		hiIn.Value += vout.Value
		if _, ok := addrCoin[utils.Bytes2string(vout.Address)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(vout.Address)] = false
		}
		hiIn.OutAddr = append(hiIn.OutAddr, &vout.Address)
	}

	if len(hiIn.OutAddr) > 0 && hiIn.Value > 0 {
		balanceHistoryManager.Add(hiIn)
	}
}
