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

type Tx_deposit_in struct {
	TxBase
	Puk []byte `json:"puk"`
}

type Tx_deposit_in_VO struct {
	TxBaseVO
	Puk string `json:"puk"`
}

func (this *Tx_deposit_in) GetVOJSON() interface{} {
	return Tx_deposit_in_VO{
		TxBaseVO: this.TxBase.ConversionVO(),
		Puk:      hex.EncodeToString(this.Puk),
	}
}

func (this *Tx_deposit_in) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_deposit_in)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_deposit_in) Proto() (*[]byte, error) {
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

	txPay := go_protos.TxDepositIn{
		TxBase: &txBase,
		Puk:    this.Puk,
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *Tx_deposit_in) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)
	buf.Write(this.Puk)
	*bs = buf.Bytes()
	return bs
}

func (this *Tx_deposit_in) GetWaitSign(vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Puk...)

	return signDst

}

func (this *Tx_deposit_in) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Puk...)

	sign := keystore.Sign(*key, *signDst)

	return &sign

}

func (this *Tx_deposit_in) CheckSign() error {

	if len(this.Vin) != 1 {
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

	for i, one := range this.Vin {

		signDst := this.GetSignSerialize(nil, uint64(i))

		*signDst = append(*signDst, this.Puk...)

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

func (this *Tx_deposit_in) GetWitness() *crypto.AddressCoin {
	witness := crypto.BuildAddr(config.AddrPre, this.Vin[0].Puk)

	return &witness
}

func (this *Tx_deposit_in) GetSpend() uint64 {
	return this.Vout[0].Value + this.Gas
}

func (this *Tx_deposit_in) CheckRepeatedTx(txs ...TxItr) bool {

	addrSelf := this.Vout[0].Address
	for _, one := range txs {
		if one.Class() != config.Wallet_tx_type_deposit_in {
			continue
		}
		addrOne := (*one.GetVout())[0].Address
		if bytes.Equal(addrSelf, addrOne) {
			return false
		} else {

			value := GetDepositWitnessAddr(&addrOne)
			if value > 0 {
				return false
			}
		}
	}
	return true
}

func (this *Tx_deposit_in) CountTxItemsNew(height uint64) *TxItemCountMap {
	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vin)),
		Nonce:    make(map[string]big.Int),
	}
	if height == 1 {

		return &itemCount
	}

	totalValue := this.Gas + this.Vout[0].Value

	from := this.Vin[0].GetPukToAddr()
	itemCount.Nonce[utils.Bytes2string(*from)] = this.Vin[0].Nonce

	frozenMap := make(map[uint64]int64, 0)
	frozenMap[0] = (0 - int64(totalValue))
	itemCount.AddItems[utils.Bytes2string(*from)] = &frozenMap
	return &itemCount
}

func (this *Tx_deposit_in) CountTxHistory(height uint64) {

	hiOut := HistoryItem{
		IsIn:    false,
		Type:    this.Class(),
		InAddr:  make([]*crypto.AddressCoin, 0),
		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
	}

	addrCoin := make(map[string]bool)
	for _, vin := range this.Vin {
		addrInfo, isSelf := keystore.FindPuk(vin.Puk)

		if !isSelf {
			continue
		}
		if _, ok := addrCoin[utils.Bytes2string(addrInfo.Addr)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(addrInfo.Addr)] = false
		}
		hiOut.InAddr = append(hiOut.InAddr, &addrInfo.Addr)
	}

	addrCoin = make(map[string]bool)
	for voutIndex, vout := range this.Vout {
		if voutIndex != 0 {
			continue
		}
		hiOut.OutAddr = append(hiOut.OutAddr, &vout.Address)
		hiOut.Value += vout.Value
		_, ok := keystore.FindAddress(vout.Address)
		if !ok {
			continue
		}

		if _, ok := addrCoin[utils.Bytes2string(vout.Address)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(vout.Address)] = false
		}

	}
	if len(hiOut.InAddr) > 0 {
		balanceHistoryManager.Add(hiOut)
	}

}

type Tx_deposit_out struct {
	TxBase
}

func (this *Tx_deposit_out) GetVOJSON() interface{} {
	return this.TxBase.ConversionVO()
}

func (this *Tx_deposit_out) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_deposit_out)
	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_deposit_out) CheckSign() error {
	if len(this.Vin) != 1 {
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
	return this.TxBase.CheckBase()
}

func (this *Tx_deposit_out) GetSpend() uint64 {
	return this.Gas
}

func (this *Tx_deposit_out) CheckRepeatedTx(txs ...TxItr) bool {
	witnessAddr := this.Vin[0].GetPukToAddr()

	for _, one := range txs {
		if one.Class() != config.Wallet_tx_type_deposit_out {
			continue
		}
		addrOne := (*one.GetVin())[0].GetPukToAddr()
		if bytes.Equal(*witnessAddr, *addrOne) {

			return false
		} else {

			value := GetDepositWitnessAddr(addrOne)
			if value <= 0 {
				return false
			}
		}
	}
	return true
}

func (this *Tx_deposit_out) CountTxItemsNew(height uint64) *TxItemCountMap {
	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vout)+len(this.Vin)),
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

	from := this.Vin[0].GetPukToAddr()
	itemCount.Nonce[utils.Bytes2string(*from)] = this.Vin[0].Nonce
	frozenMap, ok := itemCount.AddItems[utils.Bytes2string(*from)]
	if ok {
		oldValue, ok := (*frozenMap)[0]
		if ok {
			oldValue -= int64(this.Gas)
			(*frozenMap)[0] = oldValue
		} else {
			(*frozenMap)[0] = (0 - int64(this.Gas))
		}
	} else {
		frozenMap := make(map[uint64]int64, 0)
		frozenMap[0] = (0 - int64(this.Gas))
		itemCount.AddItems[utils.Bytes2string(*from)] = &frozenMap
	}
	return &itemCount
}

func (this *Tx_deposit_out) CountTxHistory(height uint64) {

	hiIn := HistoryItem{
		IsIn:    true,
		Type:    this.Class(),
		InAddr:  make([]*crypto.AddressCoin, 0),
		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
	}

	addrCoin := make(map[string]bool)
	for _, vin := range this.Vin {
		addrInfo, isSelf := keystore.FindPuk(vin.Puk)
		hiIn.InAddr = append(hiIn.InAddr, &addrInfo.Addr)
		if !isSelf {
			continue
		}
		if _, ok := addrCoin[utils.Bytes2string(addrInfo.Addr)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(addrInfo.Addr)] = false
		}

	}

	addrCoin = make(map[string]bool)
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

	if len(hiIn.OutAddr) > 0 {
		balanceHistoryManager.Add(hiIn)
	}
}

func CreateTxDepositIn(amount, gas uint64, pwd, payload string) (*Tx_deposit_in, error) {

	if amount != config.Mining_deposit {

		return nil, config.ERROR_deposit_witness
	}
	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()

	key := keystore.GetCoinbase()

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

		Puk: puk,

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}
	vins = append(vins, &vin)

	vouts := make([]*Vout, 0)

	vout := Vout{
		Value:   amount,
		Address: key.Addr,
	}
	vouts = append(vouts, &vout)

	var txin *Tx_deposit_in
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_deposit_in,
			Vin_total:  uint64(len(vouts)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,

			Payload: []byte(payload),
		}
		txin = &Tx_deposit_in{
			TxBase: base,
			Puk:    key.Puk,
		}

		for i, one := range txin.Vin {
			_, prk, err := keystore.GetKeyByPuk(one.Puk, pwd)
			if err != nil {
				return nil, err
			}

			sign := txin.GetSign(&prk, uint64(i))

			txin.Vin[i].Sign = *sign

		}
		txin.BuildHash()
		if txin.CheckHashExist() {
			txin = nil
			continue
		} else {
			break
		}
	}

	chain.Balance.AddLockTx(txin)
	return txin, nil
}

func CreateTxDepositOut(addr string, amount, gas uint64, pwd string) (*Tx_deposit_out, error) {

	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()
	item := chain.Balance.GetDepositIn()
	if item == nil {

		return nil, config.ERROR_deposit_not_exist
	}
	amount = item.Value

	vins := make([]*Vin, 0)

	totalAll, item := chain.Balance.BuildPayVinNew(item.Addr, gas)

	if totalAll < gas {

		return nil, config.ERROR_not_enough
	}

	puk, ok := keystore.GetPukByAddr(*item.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := Vin{

		Puk: puk,

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}
	vins = append(vins, &vin)

	var dstAddr crypto.AddressCoin
	if addr == "" {

		dstAddr = keystore.GetAddr()[0].Addr
	} else {

		dstAddr = crypto.AddressFromB58String(addr)

	}

	vouts := make([]*Vout, 0)

	vout := Vout{
		Value:   amount,
		Address: dstAddr,
	}
	vouts = append(vouts, &vout)

	var txin *Tx_deposit_out
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_deposit_out,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
		}
		txin = &Tx_deposit_out{
			TxBase: base,
		}

		for i, one := range txin.Vin {
			_, prk, err := keystore.GetKeyByPuk(one.Puk, pwd)
			if err != nil {
				return nil, err
			}

			sign := txin.GetSign(&prk, uint64(i))

			txin.Vin[i].Sign = *sign

		}
		txin.BuildHash()
		if txin.CheckHashExist() {
			txin = nil
			continue
		} else {
			break
		}
	}
	chain.Balance.AddLockTx(txin)
	return txin, nil
}
