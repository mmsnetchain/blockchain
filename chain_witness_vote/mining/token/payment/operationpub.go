package payment

import (
	"errors"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/keystore/pubstore"
)

func CreateTokenPayPub(pwd, seed string, txid string, items, tokenTxItems []*mining.TxItem,
	addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, comment string) (mining.TxItr, error) {
	keystore, err := pubstore.GetPubStore(pwd, seed)
	if err != nil {
		return nil, err
	}

	if len(tokenTxItems) == 0 {
		return nil, errors.New("token items empty")
	}

	var commentbs []byte
	if comment != "" {
		commentbs = []byte(comment)
	}

	var tokenTotal uint64
	tokenVins := make([]*mining.Vin, 0)
	for _, item := range tokenTxItems {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := mining.Vin{

			Puk: puk,
		}
		tokenVins = append(tokenVins, &vin)
		tokenTotal = tokenTotal + item.Value
	}
	if tokenTotal < amount {
		return nil, config.ERROR_token_not_enough
	}

	tokenVouts := make([]*mining.Vout, 0)

	tokenVout := mining.Vout{
		Value:        amount,
		Address:      *addr,
		FrozenHeight: frozenHeight,
	}
	tokenVouts = append(tokenVouts, &tokenVout)

	if tokenTotal > amount {
		tokenVout := mining.Vout{
			Value:   tokenTotal - amount,
			Address: *tokenTxItems[0].Addr,
		}
		tokenVouts = append(tokenVouts, &tokenVout)
	}

	vins := make([]*mining.Vin, 0)
	chain := mining.GetLongChain()

	var total uint64
	for _, item := range items {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := mining.Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)
		total = total + item.Value
	}

	vouts := make([]*mining.Vout, 0)

	if total > gas {
		vout := mining.Vout{
			Value:   total - gas,
			Address: *items[0].Addr,
		}
		vouts = append(vouts, &vout)
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
	return txin, nil
}

func CreateTokenPayPubs(pwd, seed string, txid string, items, tokenTxItems []*mining.TxItem,
	address []mining.PayNumber, gas uint64, comment string) (mining.TxItr, error) {
	keystore, err := pubstore.GetPubStore(pwd, seed)
	if err != nil {
		return nil, err
	}

	if len(tokenTxItems) == 0 {
		return nil, errors.New("token items empty")
	}

	var amount uint64
	for _, one := range address {
		amount += one.Amount
	}

	var tokenTotal uint64
	tokenVins := make([]*mining.Vin, 0)
	for _, item := range tokenTxItems {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := mining.Vin{

			Puk: puk,
		}
		tokenTotal = tokenTotal + item.Value
		tokenVins = append(tokenVins, &vin)
	}
	if tokenTotal < amount {

		return nil, config.ERROR_not_enough
	}

	tokenVouts := make([]*mining.Vout, 0)
	for _, one := range address {
		vout := mining.Vout{
			Value:        one.Amount,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		}
		tokenVouts = append(tokenVouts, &vout)
	}

	if tokenTotal > amount {
		tokenVout := mining.Vout{
			Value:   tokenTotal - amount,
			Address: *tokenTxItems[0].Addr,
		}
		tokenVouts = append(tokenVouts, &tokenVout)
	}

	vins := make([]*mining.Vin, 0)
	chain := mining.GetLongChain()

	var total uint64
	for _, item := range items {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := mining.Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)
		total = total + item.Value
	}
	if total < gas {

		return nil, config.ERROR_not_enough
	}

	vouts := make([]*mining.Vout, 0)

	if total > gas {
		vout := mining.Vout{
			Value:   total - gas,
			Address: *items[0].Addr,
		}
		vouts = append(vouts, &vout)
	}

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
