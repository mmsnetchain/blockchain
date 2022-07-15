package tx_spaces_mining_in

import (
	"bytes"
	"mmschainnewaccount/cache/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/tx_spaces_mining"
	"mmschainnewaccount/config"
	"sync"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

func init() {
	ac := new(SpacesMiningInController)
	mining.RegisterTransaction(config.Wallet_tx_type_spaces_mining_in, ac)
}

type SpacesMiningInController struct {
}

func (this *SpacesMiningInController) Factory() interface{} {
	return new(Tx_SpacesMining)
}

func (this *SpacesMiningInController) CountBalance(balance *mining.TxItemManager, deposit *sync.Map, bhvo *mining.BlockHeadVO) {
	for _, txItr := range bhvo.Txs {

		if txItr.Class() != config.Wallet_tx_type_spaces_mining_in {
			continue
		}

		var depositIn *sync.Map
		v, ok := deposit.Load(config.Wallet_tx_type_spaces_mining_in)
		if ok {
			depositIn = v.(*sync.Map)
		} else {
			depositIn = new(sync.Map)
			deposit.Store(config.Wallet_tx_type_spaces_mining_in, depositIn)
		}

		for voutIndex, vout := range *txItr.GetVout() {

			if voutIndex == 0 {
				txItem := mining.TxItem{
					Addr: &vout.Address,

					Value:    vout.Value,
					Txid:     *txItr.GetHash(),
					OutIndex: uint64(voutIndex),
				}

				tsm := txItr.(*Tx_SpacesMining)

				spacesMining := tx_spaces_mining.SpacesMining{
					Deposit: vout.Value,
				}

				spacesMiningBS, _ := json.Marshal(spacesMining)

				db.Remove(append([]byte(config.DB_spaces_mining_addr), tsm.NetId...))
				db.Save(append([]byte(config.DB_spaces_mining_addr), tsm.NetId...), spacesMiningBS)

				_, ok := keystore.FindAddress(vout.Address)
				if !ok {
					continue
				}

				depositIn.Store(utils.Bytes2string(tsm.NetId), &txItem)

				continue
			}

		}

	}
}

func (this *SpacesMiningInController) CheckMultiplePayments(txItr mining.TxItr) error {
	return nil
}

func (this *SpacesMiningInController) SyncCount() {

}

func (this *SpacesMiningInController) RollbackBalance() {

}

func (this *SpacesMiningInController) BuildTx(balance *mining.TxItemManager, deposit *sync.Map, srcAddr,
	addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string, params ...interface{}) (mining.TxItr, error) {

	if amount < config.Mining_name_deposit_min {
		return nil, config.ERROR_name_deposit
	}

	var depositIn *sync.Map
	v, ok := deposit.Load(config.Wallet_tx_type_spaces_mining_in)
	if ok {
		depositIn = v.(*sync.Map)
	} else {
		depositIn = new(sync.Map)
		deposit.Store(config.Wallet_tx_type_spaces_mining_in, depositIn)
	}

	if len(params) < 1 {

		return nil, config.ERROR_params_not_enough
	}
	netidsMHash := params[1].(nodeStore.AddressNet)

	var commentBs []byte
	if comment != "" {
		commentBs = []byte(comment)
	}

	isReg := true

	chain := mining.GetLongChain()

	vins := make([]*mining.Vin, 0)
	total := uint64(0)

	var items []*mining.TxItem

	if total < amount+gas {

		total, items = chain.GetBalance().BuildPayVinNew(nil, amount+gas-total)
		if total < amount+gas {
			engine.Log.Error("11111111 %d %d %d", total, amount, gas)

			return nil, config.ERROR_not_enough
		}

		if len(items) > config.Mining_pay_vin_max {
			return nil, config.ERROR_pay_vin_too_much
		}

		for _, item := range items {
			puk, ok := keystore.GetPukByAddr(*item.Addr)
			if !ok {
				return nil, config.ERROR_public_key_not_exist
			}

			vin := mining.Vin{
				Txid: item.Txid,
				Vout: item.OutIndex,
				Puk:  puk,
			}
			vins = append(vins, &vin)
		}

	}

	if total < (amount + gas) {
		engine.Log.Error("222222222 %d %d %d", total, amount, gas)

		return nil, config.ERROR_not_enough
	}

	var dstAddr crypto.AddressCoin
	if addr == nil {
		dstAddr = keystore.GetCoinbase().Addr
	} else {
		dstAddr = *addr
	}

	vouts := make([]*mining.Vout, 0)

	vout := mining.Vout{
		Value:   amount,
		Address: dstAddr,
	}
	vouts = append(vouts, &vout)

	if total > amount+gas {
		newAddrCoin, err := keystore.GetNewAddr(pwd)
		if err != nil {

			return nil, config.ERROR_password_fail
		}
		vout := mining.Vout{
			Value:   total - amount - gas,
			Address: newAddrCoin,
		}
		vouts = append(vouts, &vout)
	}
	var class uint64
	if isReg {

		class = config.Wallet_tx_type_spaces_mining_in
	} else {

		class = config.Wallet_tx_type_account_cancel
	}

	currentHeight := chain.GetCurrentBlock()
	var txin *Tx_SpacesMining
	for i := uint64(0); i < 10000; i++ {

		base := mining.TxBase{
			Type:       class,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: block.Height + config.Wallet_tx_lockHeight + i,
			Payload:    commentBs,
		}
		txin = &Tx_SpacesMining{
			TxBase: base,
			NetId:  netidsMHash,
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
					sign := txin.GetSign(&prk, one.Txid, one.Vout, uint64(i))

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
