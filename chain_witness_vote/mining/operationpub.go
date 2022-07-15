package mining

import (
	"errors"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/keystore/pubstore"
)

func CreateTxPayPub(pwd, seed string, items []*TxItem, address *crypto.AddressCoin, amount, gas, frozenHeight uint64, comment string) (*Tx_Pay, error) {
	keystore, err := pubstore.GetPubStore(pwd, seed)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("items empty")
	}
	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()

	vins := make([]*Vin, 0)

	var total uint64
	for _, item := range items {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := &Vin{

			Puk: puk,
		}
		vins = append(vins, vin)
		total = total + uint64(item.Value)
	}
	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	vouts := make([]*Vout, 0)
	vout := &Vout{
		Value:        amount,
		Address:      *address,
		FrozenHeight: frozenHeight,
	}
	vouts = append(vouts, vout)

	if total > amount+gas {
		vout := &Vout{
			Value:   total - amount - gas,
			Address: *items[0].Addr,
		}
		vouts = append(vouts, vout)
	}

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

			Payload: []byte(comment),
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

func CreateTxsPayPub(pwd, seed string, items []*TxItem, address []PayNumber, gas uint64, comment string) (*Tx_Pay, error) {
	keystore, err := pubstore.GetPubStore(pwd, seed)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("items empty")
	}
	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()
	amount := uint64(0)
	for _, one := range address {
		amount += one.Amount
	}

	vins := make([]*Vin, 0)

	var total uint64
	for _, item := range items {
		puk, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {
			return nil, config.ERROR_public_key_not_exist
		}

		vin := Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)
		total = total + uint64(item.Value)
	}
	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

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
			Address: *items[0].Addr,
		}
		vouts = append(vouts, &vout)
	}

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

			Payload: []byte(comment),
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
