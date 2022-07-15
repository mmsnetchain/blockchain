package tx_name_in

import (
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
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

func init() {
	ac := new(AccountController)
	mining.RegisterTransaction(config.Wallet_tx_type_account, ac)
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
	txProto := new(go_protos.TxNameIn)
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

	addrNet := make([]nodeStore.AddressNet, 0)
	for i, _ := range txProto.NetIds {
		one := txProto.NetIds[i]
		addrNet = append(addrNet, one)
	}
	addrCoins := make([]crypto.AddressCoin, 0)
	for i, _ := range txProto.AddrCoins {
		one := txProto.AddrCoins[i]
		addrCoins = append(addrCoins, one)
	}
	tx := &Tx_account{
		TxBase:              txBase,
		Account:             txProto.Account,
		NetIds:              addrNet,
		NetIdsMerkleHash:    txProto.NetIdsMerkleHash,
		AddrCoins:           addrCoins,
		AddrCoinsMerkleHash: txProto.AddrCoinsMerkleHash,
	}
	return tx, nil
}

func (this *AccountController) CountBalance(deposit *sync.Map, bhvo *mining.BlockHeadVO) {
	for _, txItr := range bhvo.Txs {

		if txItr.Class() != config.Wallet_tx_type_account {
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

		for voutIndex, vout := range *txItr.GetVout() {

			if voutIndex == 0 {
				txItem := mining.TxItem{
					Addr: &vout.Address,

					Value: vout.Value,
				}

				txAcc := txItr.(*Tx_account)

				nameObj := name.Nameinfo{
					Name:      string(txAcc.Account),
					Txid:      *txItr.GetHash(),
					NetIds:    txAcc.NetIds,
					AddrCoins: txAcc.AddrCoins,
					Height:    bhvo.BH.Height,
					Deposit:   vout.Value,
				}

				nameinfoBS, _ := nameObj.Proto()

				db.LevelTempDB.Remove(append([]byte(config.Name), txAcc.Account...))
				db.LevelTempDB.Save(append([]byte(config.Name), txAcc.Account...), &nameinfoBS)

				_, ok := keystore.FindAddress(vout.Address)
				if !ok {

					depositIn.Delete(string(txAcc.Account))
					name.DelName(txAcc.Account)
					continue
				}

				depositIn.Store(string(txAcc.Account), &txItem)

				name.AddName(nameObj)
				continue
			}

		}

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

	if amount < config.Mining_name_deposit_min {
		return nil, config.ERROR_name_deposit
	}

	var depositIn *sync.Map
	v, ok := deposit.Load(config.Wallet_tx_type_account)
	if ok {
		depositIn = v.(*sync.Map)
	} else {
		depositIn = new(sync.Map)
		deposit.Store(config.Wallet_tx_type_account, depositIn)
	}

	if len(params) < 3 {

		return nil, config.ERROR_params_not_enough
	}
	nameStr := params[0].(string)
	netidsMHash := params[1].([]nodeStore.AddressNet)
	addrCoins := params[2].([]crypto.AddressCoin)

	var commentBs []byte
	if comment != "" {
		commentBs = []byte(comment)
	}

	isReg := true

	netids := make([][]byte, 0)
	for _, one := range netidsMHash {
		netids = append(netids, one)
	}

	addrCoinBs := make([][]byte, 0)
	for _, one := range addrCoins {
		addrCoinBs = append(addrCoinBs, one)
	}

	isHave := false
	isOvertime := false
	{

		nameinfo := name.FindNameToNet(nameStr)
		if nameinfo != nil {
			isHave = true
			isOvertime = nameinfo.CheckIsOvertime(mining.GetHighestBlock())
		}
	}

	chain := mining.GetLongChain()

	nameinfo := name.FindName(nameStr)

	vins := make([]*mining.Vin, 0)

	vouts := make([]*mining.Vout, 0)
	var item *mining.TxItem
	total := uint64(0)
	if (isReg && isHave && isOvertime && nameinfo == nil) || (isReg && !isHave) {

		total, item = chain.GetBalance().BuildPayVinNew(srcAddr, amount+gas)
		if total < amount+gas {

			return nil, config.ERROR_not_enough
		}
		if addr == nil || len(*addr) <= 0 {
			addr = item.Addr
		}
		pukBs, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {

			return nil, config.ERROR_public_key_not_exist
		}
		nonce := chain.GetBalance().FindNonce(item.Addr)
		vin := mining.Vin{
			Puk:   pukBs,
			Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
		}
		vins = append(vins, &vin)

		vout := mining.Vout{
			Value:   amount,
			Address: *addr,
		}
		vouts = append(vouts, &vout)
	} else if isReg && isHave && nameinfo != nil {

		itemItr, ok := depositIn.Load(nameStr)
		if !ok {

			return nil, config.ERROR_deposit_not_exist
		}
		item = itemItr.(*mining.TxItem)

		total, _ = chain.GetBalance().BuildPayVinNew(item.Addr, amount+gas-item.Value)
		if total+item.Value < amount+gas {

			return nil, config.ERROR_not_enough
		}
		if addr == nil {
			addr = item.Addr
		}
		pukBs, ok := keystore.GetPukByAddr(*item.Addr)
		if !ok {

			return nil, config.ERROR_public_key_not_exist
		}
		nonce := chain.GetBalance().FindNonce(item.Addr)
		vin := mining.Vin{
			Puk:   pukBs,
			Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
		}
		vins = append(vins, &vin)

		vout := mining.Vout{
			Value:   amount,
			Address: *addr,
		}
		vouts = append(vouts, &vout)
	} else if !isReg && nameinfo != nil {

	} else {

		if isHave && !isOvertime {
			return nil, config.ERROR_name_exist
		}

		return nil, config.ERROR_params_fail
	}

	var class uint64
	if isReg {

		class = config.Wallet_tx_type_account
	} else {

		class = config.Wallet_tx_type_account_cancel
	}

	currentHeight := chain.GetCurrentBlock()
	var txin *Tx_account
	for i := uint64(0); i < 10000; i++ {

		base := mining.TxBase{
			Type:       class,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
			Payload:    commentBs,
		}
		txin = &Tx_account{
			TxBase:              base,
			Account:             []byte(nameStr),
			NetIds:              netidsMHash,
			NetIdsMerkleHash:    utils.BuildMerkleRoot(netids),
			AddrCoins:           addrCoins,
			AddrCoinsMerkleHash: utils.BuildMerkleRoot(addrCoinBs),
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
