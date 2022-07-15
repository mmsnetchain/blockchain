package payment

import (
	"math/big"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/token"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
)

const (
	Wallet_tx_class = config.Wallet_tx_type_token_payment
)

func TokenPay(srcAddress, addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string,
	txid string) (mining.TxItr, error) {

	txItr, err := mining.GetLongChain().GetBalance().BuildOtherTx(Wallet_tx_class, srcAddress,
		addr, amount, gas, frozenHeight, pwd, comment, txid)
	if err != nil {

	} else {

	}
	return txItr, err
}

func TokenPayFromSrcAddr(srcAddr, tokenSrcAddress, tokenAddr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string,
	txid []byte) (mining.TxItr, error) {

	var commentbs []byte
	if comment != "" {
		commentbs = []byte(comment)
	}

	tokenTotal, tokenTxItem := token.GetReadyPayToken(&txid, tokenSrcAddress, amount)
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
		Address:      *tokenAddr,
		FrozenHeight: frozenHeight,
	}
	tokenVouts = append(tokenVouts, &tokenVout)

	vins := make([]*mining.Vin, 0)
	chain := mining.GetLongChain()

	total, item := chain.GetBalance().BuildPayVinNew(srcAddr, gas)

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

	currentHeight := chain.GetCurrentBlock()
	var txin *TxToken
	for i := uint64(0); i < 10000; i++ {

		base := mining.TxBase{
			Type:       Wallet_tx_class,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
			Payload:    commentbs,
		}
		txin = &TxToken{
			TxBase:           base,
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
	mining.AddTx(txin)
	return txin, nil
}

func TokenPayMore(srcAddrStr, tokenSrcAddrStr crypto.AddressCoin, address []mining.PayNumber, gas uint64, pwd, comment string,
	txid []byte) (mining.TxItr, error) {

	amount := uint64(0)
	for _, one := range address {
		amount += one.Amount
	}

	tokenTotal, tokenTxItem := token.GetReadyPayToken(&txid, &tokenSrcAddrStr, amount)
	if tokenTotal < amount {
		return nil, config.ERROR_token_not_enough
	}
	tokenVins := make([]*mining.Vin, 0)

	puk, ok := keystore.GetPukByAddr(*tokenTxItem.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	tokenVin := mining.Vin{

		Puk: puk,
	}
	tokenVins = append(tokenVins, &tokenVin)

	tokenVouts := make([]*mining.Vout, 0, len(address))
	for _, one := range address {
		vout := mining.Vout{
			Value:        one.Amount,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		}
		tokenVouts = append(tokenVouts, &vout)
	}

	vins := make([]*mining.Vin, 0)
	chain := mining.GetLongChain()
	total, item := chain.GetBalance().BuildPayVinNew(&srcAddrStr, gas)

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

	commentbs := []byte{}
	if comment != "" {
		commentbs = []byte(comment)
	}

	currentHeight := chain.GetCurrentBlock()
	var txin *TxToken
	for i := uint64(0); i < 10000; i++ {

		base := mining.TxBase{
			Type:       Wallet_tx_class,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
			Payload:    commentbs,
		}
		txin = &TxToken{
			TxBase:           base,
			Token_txid:       txid,
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
	mining.AddTx(txin)

	return txin, nil
}
