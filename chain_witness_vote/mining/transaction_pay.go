package mining

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type Tx_Pay struct {
	TxBase
}

func (this *Tx_Pay) GetVOJSON() interface{} {
	return this.TxBase.ConversionVO()
}

func (this *Tx_Pay) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}

	bs := this.Serialize()
	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_pay)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_Pay) CheckSign() error {

	if len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}
	if this.Vout_total > config.Mining_pay_vout_max {
		return config.ERROR_pay_vout_too_much
	}
	if err := this.TxBase.CheckBase(); err != nil {
		return err
	}
	for _, vout := range this.Vout {
		if vout.Value <= 0 {
			return config.ERROR_amount_zero
		}
	}

	return nil
}

func (this *Tx_Pay) GetSpend() uint64 {
	spend := this.Gas
	for _, vout := range this.Vout {
		spend += vout.Value
	}
	return spend
}

func (this *Tx_Pay) CheckRepeatedTx(txs ...TxItr) bool {

	return true
}

func (this *Tx_Pay) CountTxItemsNew(height uint64) *TxItemCountMap {
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

func CreateTxPay(srcAddress, address *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string) (*Tx_Pay, error) {

	engine.Log.Info("start CreateTxPay")
	commentbs := []byte{}
	if comment != "" {
		commentbs = []byte(comment)
	}

	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()

	engine.Log.Info("currentHeight:%d", currentHeight)

	vins := make([]*Vin, 0)

	total, item := chain.Balance.BuildPayVinNew(srcAddress, amount+gas)

	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	puk, ok := keystore.GetPukByAddr(*item.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}
	nonce := chain.GetBalance().FindNonce(item.Addr)

	vin := Vin{

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
		Puk:   puk,
	}
	engine.Log.Info("nonce:%d", vin.Nonce.Uint64())
	vins = append(vins, &vin)

	vouts := make([]*Vout, 0)
	vout := Vout{
		Value:        amount,
		Address:      *address,
		FrozenHeight: frozenHeight,
	}
	vouts = append(vouts, &vout)

	var pay *Tx_Pay
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_pay,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
			Payload:    commentbs,
		}
		pay = &Tx_Pay{
			TxBase: base,
		}

		for i, one := range pay.Vin {
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

	engine.Log.Info("start CreateTxPay %s", hex.EncodeToString(*pay.GetHash()))

	chain.Balance.AddLockTx(pay)
	return pay, nil
}

func CreateTxsPay(srcAddr *crypto.AddressCoin, address []PayNumber, gas uint64, pwd, comment string) (*Tx_Pay, error) {

	commentbs := []byte{}
	if comment != "" {
		commentbs = []byte(comment)
	}

	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()
	amount := uint64(0)
	for _, one := range address {
		amount += one.Amount
	}

	vins := make([]*Vin, 0)

	total, item := chain.Balance.BuildPayVinNew(srcAddr, amount+gas)
	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	puk, ok := keystore.GetPukByAddr(*item.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := Vin{

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
		Puk:   puk,
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

	var pay *Tx_Pay
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_pay,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,

			Payload: commentbs,
		}
		pay = &Tx_Pay{
			TxBase: base,
		}

		for i, one := range pay.Vin {
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
	return pay, nil
}

type PayNumber struct {
	Address      crypto.AddressCoin
	Amount       uint64
	FrozenHeight uint64
}

func CreateTxsPayByPayload(address []PayNumber, gas uint64, pwd string, cs *CommunitySign) (*Tx_Pay, error) {
	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()
	amount := uint64(0)
	for _, one := range address {
		amount += one.Amount
	}

	vins := make([]*Vin, 0)

	total, item := chain.Balance.BuildPayVinNew(nil, amount+gas)
	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	puk, ok := keystore.GetPukByAddr(*item.Addr)
	if !ok {

		return nil, config.ERROR_public_key_not_exist
	}

	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := Vin{

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
		Puk:   puk,
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

	if total > amount+gas {
		vout := Vout{
			Value:   total - amount - gas,
			Address: keystore.GetAddr()[0].Addr,
		}
		vouts = append(vouts, &vout)
	}

	vouts = CleanZeroVouts(&vouts)
	vouts = MergeVouts(&vouts)

	var pay *Tx_Pay
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_pay,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
		}

		var txItr TxItr = &Tx_Pay{
			TxBase: base,
		}

		if cs != nil {
			addr := crypto.BuildAddr(config.AddrPre, cs.Puk)
			_, prk, _, err := keystore.GetKeyByAddr(addr, pwd)
			if err != nil {
				return nil, err
			}
			txItr = SignPayload(txItr, cs.Puk, prk, cs.StartHeight, cs.EndHeight)
		}

		pay = txItr.(*Tx_Pay)

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
	return pay, nil
}
