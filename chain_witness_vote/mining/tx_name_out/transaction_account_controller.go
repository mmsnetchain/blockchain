package tx_name_out

import (
	"bytes"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
)

func init() {
	ac := new(AccountController)
	mining.RegisterTransaction(config.Wallet_tx_type_account_cancel, ac)
}

type AccountController struct {
}

func (this *AccountController) Factory() interface{} {
	return new(Tx_account)
}

func (this *AccountController) ParseProto(bs *[]byte) (interface{}, error) {
	if bs == nil {
		return nil, nil
	}
	txProto := new(go_protos.TxNameOut)
	err := proto.Unmarshal(*bs, txProto)
	if err != nil {
		return nil, err
	}
	vins := make([]*mining.Vin, 0)
	for _, one := range txProto.TxBase.Vin {
		nonce := new(big.Int).SetBytes(one.Nonce)
		vins = append(vins, &mining.Vin{

			Puk:   one.Puk,
			Sign:  one.Sign,
			Nonce: *nonce,
		})
	}
	vouts := make([]*mining.Vout, 0)
	for _, one := range txProto.TxBase.Vout {
		vouts = append(vouts, &mining.Vout{
			Value:        one.Value,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		})
	}
	txBase := mining.TxBase{}
	txBase.Hash = txProto.TxBase.Hash
	txBase.Type = txProto.TxBase.Type
	txBase.Vin_total = txProto.TxBase.VinTotal
	txBase.Vin = vins
	txBase.Vout_total = txProto.TxBase.VoutTotal
	txBase.Vout = vouts
	txBase.Gas = txProto.TxBase.Gas
	txBase.LockHeight = txProto.TxBase.LockHeight
	txBase.Payload = txProto.TxBase.Payload
	txBase.BlockHash = txProto.TxBase.BlockHash
	tx := &Tx_account{
		TxBase:  txBase,
		Account: txProto.Account,
	}
	return tx, nil
}

func (this *AccountController) CountBalance(deposit *sync.Map, bhvo *mining.BlockHeadVO) {

	for _, txItr := range bhvo.Txs {

		if txItr.Class() != config.Wallet_tx_type_account_cancel {
			continue
		}

		var depositIn *sync.Map
		v, ok := deposit.Load(config.Wallet_tx_type_account)
		if ok {
			depositIn = v.(*sync.Map)
		} else {
			depositIn = new(sync.Map)
			deposit.Store(config.Wallet_tx_type_account, depositIn)
		}

		nameOut := txItr.(*Tx_account)
		vin := (*txItr.GetVin())[0]

		db.LevelTempDB.Remove(append([]byte(config.Name), nameOut.Account...))

		_, ok = keystore.FindAddress(*vin.GetPukToAddr())
		if !ok {
			continue
		}

		depositIn.Delete(string(nameOut.Account))

		name.DelName(nameOut.Account)

	}
}

func (this *AccountController) CheckMultiplePayments(txItr mining.TxItr) error {
	return nil
}

func (this *AccountController) SyncCount() {

}

func (this *AccountController) RollbackBalance() {

}

func (this *AccountController) BuildTx(deposit *sync.Map, srcAddr,
	addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string, params ...interface{}) (mining.TxItr, error) {
	var depositIn *sync.Map
	v, ok := deposit.Load(config.Wallet_tx_type_account)
	if ok {
		depositIn = v.(*sync.Map)
	} else {
		depositIn = new(sync.Map)
		deposit.Store(config.Wallet_tx_type_account, depositIn)
	}

	if len(params) < 1 {

		return nil, config.ERROR_params_not_enough
	}
	nameStr := params[0].(string)

	var commentBs []byte
	if comment != "" {
		commentBs = []byte(comment)
	}

	nameInTxid := name.FindName(nameStr)
	if nameInTxid == nil {
		return nil, config.ERROR_name_not_exist
	}

	itemItr, ok := depositIn.Load(nameStr)
	if !ok {

		return nil, config.ERROR_deposit_not_exist
	}

	depositItem := itemItr.(*mining.TxItem)

	chain := mining.GetLongChain()
	pukBs, ok := keystore.GetPukByAddr(*depositItem.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	total, item := chain.GetBalance().BuildPayVinNew(depositItem.Addr, gas)
	if total < gas {

		return nil, config.ERROR_not_enough
	}
	vins := make([]*mining.Vin, 0)
	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := mining.Vin{

		Puk: pukBs,

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}

	vins = append(vins, &vin)

	var dstAddr crypto.AddressCoin
	if addr == nil {

		dstAddr = keystore.GetCoinbase().Addr
	} else {

		dstAddr = *addr
	}

	vouts := make([]*mining.Vout, 0)

	vout := &mining.Vout{
		Value:   depositItem.Value,
		Address: dstAddr,
	}
	vouts = append(vouts, vout)

	currentHeight := chain.GetCurrentBlock()

	var txin *Tx_account
	for i := uint64(0); i < 10000; i++ {
		base := mining.TxBase{
			Type:       config.Wallet_tx_type_account_cancel,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
			Payload:    commentBs,
		}
		txin = &Tx_account{
			TxBase:  base,
			Account: []byte(nameStr),
		}

		for i, one := range txin.Vin {
			for _, key := range keystore.GetAddr() {
				puk, ok := keystore.GetPukByAddr(key.Addr)
				if !ok {

					return nil, config.ERROR_public_key_not_exist
				}

				if bytes.Equal(puk, one.Puk) {

					_, prk, _, err := keystore.GetKeyByAddr(key.Addr, pwd)

					if err != nil {
						return nil, err
					}
					sign := txin.GetSign(&prk, uint64(i))

					txin.Vin[i].Sign = *sign
				}
			}
		}

		txin.BuildHash()
		if txin.CheckHashExist() {
			txin = nil
			continue
		} else {
			break
		}
	}

	chain.GetBalance().AddLockTx(txin)
	return txin, nil
}
