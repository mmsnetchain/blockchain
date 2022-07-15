package mining

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type Tx_Vote_Reward struct {
	TxBase
	StartHeight uint64
	EndHeight   uint64
}

type Tx_Vote_Reward_VO struct {
	TxBaseVO
	StartHeight uint64 `json:"StartHeight"`
	EndHeight   uint64 `json:"EndHeight"`
}

func (this *Tx_Vote_Reward) GetVOJSON() interface{} {
	return Tx_Vote_Reward_VO{
		TxBaseVO:    this.TxBase.ConversionVO(),
		StartHeight: this.StartHeight,
		EndHeight:   this.EndHeight,
	}
}

func (this *Tx_Vote_Reward) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()
	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_voting_reward)
	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_Vote_Reward) Proto() (*[]byte, error) {
	vins := make([]*go_protos.Vin, 0, len(this.Vin))
	for _, one := range this.Vin {
		vins = append(vins, &go_protos.Vin{

			Puk:   one.Puk,
			Sign:  one.Sign,
			Nonce: one.Nonce.Bytes(),
		})
	}
	vouts := make([]*go_protos.Vout, 0, len(this.Vout))
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

	txPay := go_protos.TxVoteReward{
		TxBase:      &txBase,
		StartHeight: this.StartHeight,
		EndHeight:   this.EndHeight,
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *Tx_Vote_Reward) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)
	buf.Write(utils.Uint64ToBytes(this.StartHeight))
	buf.Write(utils.Uint64ToBytes(this.EndHeight))
	*bs = buf.Bytes()
	return bs
}

func (this *Tx_Vote_Reward) GetWaitSign(vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)
	*signDst = append(*signDst, utils.Uint64ToBytes(this.StartHeight)...)
	*signDst = append(*signDst, utils.Uint64ToBytes(this.EndHeight)...)

	return signDst

}

func (this *Tx_Vote_Reward) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, utils.Uint64ToBytes(this.StartHeight)...)
	*signDst = append(*signDst, utils.Uint64ToBytes(this.EndHeight)...)

	sign := keystore.Sign(*key, *signDst)

	return &sign

}

func (this *Tx_Vote_Reward) CheckSign() error {
	if this.Vin == nil || len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}
	if this.Vout_total > config.Mining_pay_vout_max {
		return config.ERROR_pay_vout_too_much
	}

	for _, vout := range this.Vout {
		if vout.Value <= 0 {
			return config.ERROR_amount_zero
		}
	}

	var puk *[]byte

	for i, one := range this.Vin {
		if i == 0 {
			puk = &one.Puk
		} else {
			if !bytes.Equal(*puk, one.Puk) {
				return config.ERROR_vote_reward_addr_disunity
			}
		}

		signDst := this.GetSignSerialize(nil, uint64(i))

		*signDst = append(*signDst, utils.Uint64ToBytes(this.StartHeight)...)
		*signDst = append(*signDst, utils.Uint64ToBytes(this.EndHeight)...)

		puk := ed25519.PublicKey(one.Puk)

		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*signDst))
		}
		if !ed25519.Verify(puk, *signDst, one.Sign) {
			return config.ERROR_sign_fail
		}

	}

	return nil
}

func (this *Tx_Vote_Reward) GetSpend() uint64 {
	spend := this.Gas
	for _, vout := range this.Vout {
		spend += vout.Value
	}
	return spend
}

func (this *Tx_Vote_Reward) CheckRepeatedTx(txs ...TxItr) bool {

	for _, txOne := range txs {
		if txOne.Class() != config.Wallet_tx_type_voting_reward {
			continue
		}
		if bytes.Equal(this.Vin[0].Puk, (*txOne.GetVin())[0].Puk) {
			return false
		}
	}

	return true
}

func (this *Tx_Vote_Reward) CountTxItemsNew(height uint64) *TxItemCountMap {
	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vout)+len(this.Vin)),
		Nonce:    make(map[string]big.Int),
	}

	totalValue := this.Gas
	for _, vout := range this.Vout {
		totalValue += vout.Value
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

func CreateTxVoteReward(addr *crypto.AddressCoin, address []PayNumber, gas uint64, pwd string, startHeight, endHeight uint64) (*Tx_Vote_Reward, error) {
	engine.Log.Info("CreateTxVoteReward start")
	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()
	amount := uint64(0)
	for _, one := range address {
		amount += one.Amount
	}

	vins := make([]*Vin, 0)

	total := GetCommunityVoteRewardFrozen(addr)
	_, lockValue := chain.GetBalance().FindLockTotalByAddr(addr)
	total -= lockValue

	if total < gas+amount {

		return nil, config.ERROR_not_enough
	}

	puk, ok := keystore.GetPukByAddr(*addr)
	if !ok {

		return nil, config.ERROR_public_key_not_exist
	}

	nonce := chain.GetBalance().FindNonce(addr)
	vin := Vin{
		Puk:   puk,
		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}
	vins = append(vins, &vin)

	vouts := make([]*Vout, 0)
	for _, one := range address {
		vout := Vout{
			Value:        one.Amount,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		}
		vouts = append(vouts, &vout)
	}

	vouts = CleanZeroVouts(&vouts)
	vouts = MergeVouts(&vouts)

	var pay *Tx_Vote_Reward
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_voting_reward,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
		}

		var txItr TxItr = &Tx_Vote_Reward{
			TxBase:      base,
			StartHeight: startHeight,
			EndHeight:   endHeight,
		}

		pay = txItr.(*Tx_Vote_Reward)

		for i, one := range *pay.GetVin() {
			_, prk, err := keystore.GetKeyByPuk(one.Puk, pwd)
			if err != nil {
				return nil, err
			}

			sign := pay.GetSign(&prk, uint64(i))

			pay.Vin[i].Sign = *sign
		}
		pay.BuildHash()
		if pay.CheckHashExist() {
			pay = nil
			continue
		} else {
			break
		}
	}

	chain.Balance.AddLockTx(pay)
	engine.Log.Info("CreateTxVoteReward end hash:%s", hex.EncodeToString(*pay.GetHash()))
	return pay, nil
}
