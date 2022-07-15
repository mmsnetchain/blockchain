package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/keystore/pubstore"
	"github.com/prestonTao/libp2parea/engine"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/token/payment"
	"mmschainnewaccount/config"
	"strconv"
	"sync"
)

const (
	wallet_db_path    = "D:/temp/peer/wallet/data"
	tx                = `{"hash":"0400000000000000f824c50fce25b94332f6ac2c96fe996ec61eabccc42a83d760867263ad78bc21","type":4,"vin_total":1,"vin":[{"txid":"0100000000000000eb95bf79cc5eb0345fdf14b38e0458c4412962d7c9c5d3dbd335eba5ee40441d","vout":0,"puk":"20d456c76f79de0f2d041092c2a4de5375b56afc948ace0c2050bf73f39ce3e9","sign":"ea69cc78c60768fc3f5cd4c02d530dc43a8d87604abe2e46944f4d38da4a9060a0aafb2666deac014f7dd6f75cbd9cde11e30fd4ceb0727c0bb6cf7f3970d703"}],"vout_total":2,"vout":[{"value":1000000000000000,"address":"ZHCQ3hzd66mYRGFU8HXLD2QguNnUpVh9wVVA4","frozen_height":0},{"value":298999999990000000,"address":"ZHCMmuVUkjAwpMz4fd1gEV55pEXXjGK1iPC4","frozen_height":0}],"gas":10000000,"lock_height":370901,"payload":"test","blockhash":"bb6298d4f0c9e8c0cdf406225f917305da0de020f8ff2ef2500ae8e1bd837b81"}`
	txIsVisualization = false
	seed              = "bExBxUCsLVSk5czubhZYm5mqrRE2tcCCZjWCQKPk9k2UoA5LJKCWbFsLKAsyiDNixD28Y9G2C6GLjTdFW2eDc9xwTd7pwdTWz3Xqm7Bt9z2UjtfdRyLt1dRxKtEvPa5pZmfjpEAKfAuyBnmB25HNQ5qcavHKRZYfozgLafSh8wHudDkgBdsfsbuy8dUDhHYdsJF8iTEsyyduWdEnApo6tyxyMrXUz7XHy3UvTLCQibkp8mr63ZTE7hM9GtCapkWhmEMmrEaJHyKRyeyeDNXJfKiyQHWf27P43ZTWCTdbN6Ja36Sp8XTZ5u4aL6Ux86MLXEUugERWt"
	pwd               = "123456789"
	lockHeight        = 1127226
)

func init() {
	tpc := new(payment.TokenPublishController)
	tpc.ActiveVoutIndex = new(sync.Map)
	mining.RegisterTransaction(config.Wallet_tx_type_account, tpc)

	err := db.InitDB(config.DB_path)
	if err != nil {
		fmt.Println("init db error:", err.Error())
		panic(err)
	}

}

func main() {
	example2()
}

func example1() {
	txItr := parseTx()
	if txItr.Class() == config.Wallet_tx_type_pay {
		createSign(txItr)
	} else if txItr.Class() == config.Wallet_tx_type_token_payment {
		txItr := parseTxTokenPay(tx)
		reSign(txItr, seed, pwd)
	} else {
		panic("")
	}
}
func example2() {
	txItr := parseTx()
	if txItr.Class() == config.Wallet_tx_type_pay {

		printVO(txItr)
	} else if txItr.Class() == config.Wallet_tx_type_token_payment {
		txItr := parseTxTokenPay(tx)

		printVO(txItr)
	} else {
		panic("")
	}

}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseTx() *mining.Tx_Pay {
	txbaseVO := new(mining.TxBaseVO)

	err := json.Unmarshal([]byte(tx), txbaseVO)
	panicError(err)

	txhash, err := hex.DecodeString(txbaseVO.Hash)
	panicError(err)

	vins := make([]*mining.Vin, 0)
	for _, one := range txbaseVO.Vin {
		txid, err := hex.DecodeString(one.Txid)
		panicError(err)
		puk, err := hex.DecodeString(one.Puk)
		panicError(err)
		vin := mining.Vin{
			Txid: txid,
			Vout: one.Vout,
			Puk:  puk,
		}
		vins = append(vins, &vin)
	}

	vouts := make([]*mining.Vout, 0)
	for _, one := range txbaseVO.Vout {

		vout := mining.Vout{
			Value:        one.Value,
			Address:      crypto.AddressFromB58String(one.Address),
			FrozenHeight: one.FrozenHeight,
		}
		vouts = append(vouts, &vout)
	}

	blockHash, err := hex.DecodeString(txbaseVO.BlockHash)
	panicError(err)

	txbase := mining.TxBase{
		Hash:       txhash,
		Type:       txbaseVO.Type,
		Vin_total:  txbaseVO.Vin_total,
		Vin:        vins,
		Vout_total: txbaseVO.Vout_total,
		Vout:       vouts,
		Gas:        txbaseVO.Gas,
		LockHeight: txbaseVO.LockHeight,

		Payload:   []byte(txbaseVO.Payload),
		BlockHash: blockHash,
	}
	txPay := mining.Tx_Pay{
		TxBase: txbase,
	}
	return &txPay
}

func createSign(txItr mining.TxItr) {

	txPay := txItr.(*mining.Tx_Pay)

	items := make([]*mining.TxItem, 0)

	for _, vin := range *txPay.GetVin() {

		bs, err := db.Find(vin.Txid)
		if err != nil {
			fmt.Println(":", err.Error())
			return
		}

		txItrOne, err := mining.ParseTxBase(mining.ParseTxClass(vin.Txid), bs)
		if err != nil {
			fmt.Println(":", err.Error())
			return
		}

		vouts := txItrOne.GetVout()
		voutOne := (*vouts)[vin.Vout]

		txItem := mining.TxItem{
			Addr:     &voutOne.Address,
			Value:    voutOne.Value,
			Txid:     vin.Txid,
			OutIndex: vin.Vout,
		}

		items = append(items, &txItem)

	}

	voutOne := (*txPay.GetVout())[0]

	txPay, err := CreateTxPayPub(txItr, pwd, seed, items, &voutOne.Address, voutOne.Value, 0, 0, "")
	if err != nil {
		fmt.Println(":", err.Error())
		return
	}
	fmt.Println("4444444444444444444444")
	txPayBs, err := json.Marshal(txPay.GetVOJSON())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(txPayBs))

	fmt.Println(":")

	txBs, _ := txPay.Json()
	fmt.Println(strconv.Quote(string(*txBs)))

	fmt.Println("end!!!!!!")
}

func CreateTxPayPub(txItr mining.TxItr, pwd, seed string, items []*mining.TxItem, address *crypto.AddressCoin,
	amount, gas, frozenHeight uint64, comment string) (*mining.Tx_Pay, error) {
	txBase := txItr.(*mining.Tx_Pay)

	keystore, err := pubstore.GetPubStore(pwd, seed)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("items empty")
	}

	vins := make([]*mining.Vin, 0)

	var total uint64
	for _, item := range items {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := &mining.Vin{
			Txid: item.Txid,
			Vout: item.OutIndex,
			Puk:  puk,
		}
		vins = append(vins, vin)
		total = total + item.Value
	}
	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	vouts := make([]*mining.Vout, 0)
	vout := &mining.Vout{
		Value:        amount,
		Address:      *address,
		FrozenHeight: frozenHeight,
	}
	vouts = append(vouts, vout)

	if total > amount+gas {
		vout := &mining.Vout{
			Value:   total - amount - gas,
			Address: *items[0].Addr,
		}
		vouts = append(vouts, vout)
	}

	var pay *mining.Tx_Pay
	for i := uint64(0); i < 10000; i++ {

		base := mining.TxBase{
			Type:       config.Wallet_tx_type_pay,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        txBase.Gas,

			LockHeight: lockHeight,

			Payload: txBase.Payload,
		}
		pay = &mining.Tx_Pay{
			TxBase: base,
		}

		for i, one := range pay.Vin {
			_, prk, err := keystore.GetKeyByPuk(one.Puk, pwd)
			if err != nil {
				return nil, err
			}
			sign := pay.GetSign(&prk, one.Txid, one.Vout, uint64(i))
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

	return pay, nil
}

func reSign(txItr mining.TxItr, seed, pwd string) {
	keystore, err := pubstore.GetPubStore(pwd, seed)
	if err != nil {
		return
	}
	pay := txItr.(*payment.TxToken)
	pay.BlockHash = nil
	pay.Hash = nil

	for i, one := range *pay.GetVin() {
		_, prk, err := keystore.GetKeyByPuk(one.Puk, pwd)
		panicError(err)

		engine.Log.Info("key %s", hex.EncodeToString(prk))
		sign := txItr.GetSign(&prk, one.Txid, one.Vout, uint64(i))
		txItr.SetSign(uint64(i), *sign)

	}

	txItr.BuildHash()

	if txItr.CheckHashExist() {
		panic("hash")
	}

	txPayBs, err := json.Marshal(txItr.GetVOJSON())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(txPayBs))

	fmt.Println(":")

	txBs, _ := txItr.Json()
	fmt.Println(strconv.Quote(string(*txBs)))

	fmt.Println("end!!!!!!")

}

func parseTxTokenPay(txStr string) *payment.TxToken {
	txbaseVO := new(payment.TxToken_VO)

	err := json.Unmarshal([]byte(txStr), txbaseVO)
	panicError(err)

	txhash, err := hex.DecodeString(txbaseVO.Hash)
	panicError(err)

	vins := make([]*mining.Vin, 0)
	for _, one := range txbaseVO.Vin {
		txid, err := hex.DecodeString(one.Txid)
		panicError(err)
		puk, err := hex.DecodeString(one.Puk)
		panicError(err)
		vin := mining.Vin{
			Txid: txid,
			Vout: one.Vout,
			Puk:  puk,
		}
		vins = append(vins, &vin)
	}

	vouts := make([]*mining.Vout, 0)
	for _, one := range txbaseVO.Vout {

		vout := mining.Vout{
			Value:        one.Value,
			Address:      crypto.AddressFromB58String(one.Address),
			FrozenHeight: one.FrozenHeight,
		}
		vouts = append(vouts, &vout)
	}

	blockHash, err := hex.DecodeString(txbaseVO.BlockHash)
	panicError(err)

	txbase := mining.TxBase{
		Hash:       txhash,
		Type:       txbaseVO.Type,
		Vin_total:  txbaseVO.Vin_total,
		Vin:        vins,
		Vout_total: txbaseVO.Vout_total,
		Vout:       vouts,
		Gas:        txbaseVO.Gas,
		LockHeight: txbaseVO.LockHeight,

		Payload:   []byte(txbaseVO.Payload),
		BlockHash: blockHash,
	}

	tokenVin := make([]*mining.Vin, 0)
	for _, one := range txbaseVO.Token_Vin {
		txid, err := hex.DecodeString(one.Txid)
		panicError(err)
		puk, err := hex.DecodeString(one.Puk)
		panicError(err)
		vin := mining.Vin{
			Txid: txid,
			Vout: one.Vout,
			Puk:  puk,
		}
		tokenVin = append(tokenVin, &vin)
	}

	tokenVouts := make([]*mining.Vout, 0)
	for _, one := range txbaseVO.Token_Vout {

		vout := mining.Vout{
			Value:        one.Value,
			Address:      crypto.AddressFromB58String(one.Address),
			FrozenHeight: one.FrozenHeight,
		}
		tokenVouts = append(tokenVouts, &vout)
	}

	txPay := payment.TxToken{
		TxBase:           txbase,
		Token_Vin_total:  txbaseVO.Token_Vin_total,
		Token_Vin:        tokenVin,
		Token_Vout_total: txbaseVO.Token_Vout_total,
		Token_Vout:       tokenVouts,
	}
	return &txPay
}

func printVO(txItr mining.TxItr) {
	txPayBs, err := json.Marshal(txItr.GetVOJSON())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(txPayBs))

	fmt.Println(":")

	txBs, _ := txItr.Json()
	fmt.Println(strconv.Quote(string(*txBs)))

	fmt.Println("end!!!!!!")
}
