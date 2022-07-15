package payment

import (
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/token"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/utils"
)

func init() {
	tpc := new(TokenPublishController)
	tpc.ActiveVoutIndex = new(sync.Map)
	mining.RegisterTransaction(Wallet_tx_class, tpc)
}

type TokenPublishController struct {
	ActiveVoutIndex *sync.Map
}

func (this *TokenPublishController) Factory() interface{} {
	return new(TxToken)
}

func (this *TokenPublishController) ParseProto(bs *[]byte) (interface{}, error) {
	if bs == nil {
		return nil, nil
	}
	txProto := new(go_protos.TxTokenPay)
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

	tokenVins := make([]*mining.Vin, 0)
	for _, one := range txProto.Token_Vin {
		nonce := new(big.Int).SetBytes(one.Nonce)
		tokenVins = append(tokenVins, &mining.Vin{

			Puk:   one.Puk,
			Sign:  one.Sign,
			Nonce: *nonce,
		})
	}

	tokenVouts := make([]*mining.Vout, 0)
	for _, one := range txProto.Token_Vout {
		tokenVouts = append(tokenVouts, &mining.Vout{
			Value:        one.Value,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		})
	}
	tx := &TxToken{
		TxBase:           txBase,
		Token_Vin_total:  txProto.Token_VinTotal,
		Token_Vin:        tokenVins,
		Token_Vout_total: txProto.Token_VoutTotal,
		Token_Vout:       tokenVouts,
	}
	return tx, nil
}

func (this *TokenPublishController) CountBalance(deposit *sync.Map, bhvo *mining.BlockHeadVO) {

	txItemCounts := make(map[string]*map[string]int64, 0)

	for _, txItr := range bhvo.Txs {
		if txItr.Class() != Wallet_tx_class {
			continue
		}

		txToken := txItr.(*TxToken)

		from := txToken.Token_Vin[0].GetPukToAddr()
		txid := utils.Bytes2string(txToken.Token_txid)

		tokenMap, ok := txItemCounts[txid]
		if !ok {
			newm := make(map[string]int64, 0)
			tokenMap = &newm
			txItemCounts[txid] = &newm
		}

		totalValue := uint64(0)
		for _, vout := range txToken.Token_Vout {
			totalValue += vout.Value
			value := (*tokenMap)[utils.Bytes2string(vout.Address)]
			value += int64(vout.Value)
			(*tokenMap)[utils.Bytes2string(vout.Address)] = value
		}

		value, _ := (*tokenMap)[utils.Bytes2string(*from)]
		value -= int64(totalValue)
		(*tokenMap)[utils.Bytes2string(*from)] = value

	}

	token.CountToken(&txItemCounts)

	token.UnfrozenToken(bhvo.BH.Height - 1)

}

func (this *TokenPublishController) CheckMultiplePayments(txItr mining.TxItr) error {

	return nil
}

func (this *TokenPublishController) SyncCount() {

}

func (this *TokenPublishController) RollbackBalance() {

}

func (this *TokenPublishController) BuildTx(deposit *sync.Map, srcAddr, addr *crypto.AddressCoin,
	amount, gas, frozenHeight uint64, pwd, comment string, params ...interface{}) (mining.TxItr, error) {

	if len(params) < 1 {

		return nil, config.ERROR_params_not_enough
	}

	txid := params[0].(string)
	txidBs, err := hex.DecodeString(txid)
	if err != nil {
		return nil, config.ERROR_params_fail
	}

	var commentbs []byte
	if comment != "" {
		commentbs = []byte(comment)
	}

	tokenTotal, tokenTxItem := token.GetReadyPayToken(&txidBs, srcAddr, amount)
	if tokenTotal < amount {
		return nil, config.ERROR_token_not_enough
	}
	tokenVins := make([]*mining.Vin, 0)

	puk, ok := keystore.GetPukByAddr(*tokenTxItem.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}
	tokenvin := mining.Vin{

		Puk: puk,
	}
	tokenVins = append(tokenVins, &tokenvin)

	tokenVouts := make([]*mining.Vout, 0)

	tokenVout := mining.Vout{
		Value:        amount,
		Address:      *addr,
		FrozenHeight: frozenHeight,
	}
	tokenVouts = append(tokenVouts, &tokenVout)

	vins := make([]*mining.Vin, 0)
	chain := mining.GetLongChain()
	total, item := chain.GetBalance().BuildPayVinNew(nil, gas)
	if total < gas {

		return nil, config.ERROR_not_enough
	}

	puk, ok = keystore.GetPukByAddr(*item.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := mining.Vin{

		Puk: puk,

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}
	vins = append(vins, &vin)

	vouts := make([]*mining.Vout, 0)

	_, block := chain.GetLastBlock()
	var txin *TxToken
	for i := uint64(0); i < 10000; i++ {

		base := mining.TxBase{
			Type:       Wallet_tx_class,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: block.Height + config.Wallet_tx_lockHeight + i,
			Payload:    commentbs,
		}
		txin = &TxToken{
			TxBase:           base,
			Token_txid:       txidBs,
			Token_Vin_total:  uint64(len(tokenVins)),
			Token_Vin:        tokenVins,
			Token_Vout_total: uint64(len(tokenVouts)),
			Token_Vout:       tokenVouts,
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

	chain.GetBalance().AddLockTx(txin)
	return txin, nil
}
