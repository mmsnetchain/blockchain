package mining

import (
	"encoding/hex"
	"errors"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"sync"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type TransactionManager struct {
	witnessBackup *WitnessBackup

	unpackedNonceDiscontinuity *sync.Map
	unpackedTransaction        *UnpackedTransaction
	tempTxLock                 *sync.RWMutex
	tempTx                     []TransactionRatio
	tempTxsignal               chan bool
}

func (this *TransactionManager) AddTx(txItr TxItr) bool {
	txItr.BuildHash()
	if len((*txItr.GetVin())[0].Nonce.Bytes()) == 0 {
		engine.Log.Info("txid:%s nonce is nil", hex.EncodeToString(*txItr.GetHash()))
		return false
	}

	if ok, _ := db.LevelDB.CheckHashExist(*txItr.GetHash()); !ok {

		err := SaveTempTx(txItr)
		if err != nil {
			engine.Log.Info("txid:%s save error:%s", txItr.GetHash(), err.Error())
			return false
		}
	}

	spend := txItr.GetSpend()

	txbs := txItr.Serialize()

	div1 := new(big.Int).Mul(big.NewInt(int64(txItr.GetGas())), big.NewInt(100000000))
	div2 := big.NewInt(int64(len(*txbs)))

	ratioValue := new(big.Int).Div(div1, div2)

	ratio := TransactionRatio{
		tx:    txItr,
		size:  uint64(len(*txbs)),
		gas:   txItr.GetGas(),
		Ratio: ratioValue,
	}
	if txItr.Class() == config.Wallet_tx_type_voting_reward {
		ratio.spendLock = spend
	} else {
		ratio.spend = spend
	}

	this.tempTxLock.Lock()
	this.tempTx = append(this.tempTx, ratio)
	this.tempTxLock.Unlock()

	select {
	case this.tempTxsignal <- false:

	default:

	}

	return true
}

func (this *TransactionManager) loopCheckTxs() {
	utils.Go(func() {

		for range this.tempTxsignal {

			this.tempTxLock.Lock()
			if len(this.tempTx) <= 0 {
				this.tempTxLock.Unlock()
				continue
			}

			temp := this.tempTx

			this.tempTx = make([]TransactionRatio, 0)
			this.tempTxLock.Unlock()

			for i, _ := range temp {
				one := temp[i]

				err := this.checkTx(&one)
				if err != nil {

					engine.Log.Debug("Transaction validation failed %s %s", hex.EncodeToString(*one.tx.GetHash()), err.Error())

					continue
				}

			}
		}
	})
}

func (this *TransactionManager) checkTx(ratio *TransactionRatio) error {

	var err error

	txItr := ratio.tx
	txItr.BuildHash()

	if this.CheckAddrValueBig(ratio) {
		return config.ERROR_addr_value_big
	}

	txCtrl := GetTransactionCtrl(txItr.Class())
	if txCtrl != nil {
		err = txCtrl.CheckMultiplePayments(txItr)
		if err != nil {

			return err
		}
	}

	err = this.CheckNotspend(ratio)
	if err != nil {

		return err
	}

	ok, err := this.CheckNonce(txItr)
	if err != nil {

		return err
	}

	if !ok {

		this.unpackedNonceDiscontinuity.Store(utils.Bytes2string(*txItr.GetHash()), ratio)

		return nil
	}
	this.unpackedTransaction.AddTx(ratio)

	this.unpackedNonceDiscontinuity.Range(func(k, v interface{}) bool {
		ratio := v.(*TransactionRatio)

		err = this.CheckNotspend(ratio)
		if err != nil {
			return false
		}
		return true

		ok, err := this.CheckNonce(ratio.tx)
		if err != nil {
			return false
		}
		if !ok {

			return false
		}
		this.unpackedTransaction.AddTx(ratio)
		return true
	})

	return nil
}

func (this *TransactionManager) CheckNonce(txItr TxItr) (bool, error) {

	var err error
	fromAddr := (*txItr.GetVin())[0].GetPukToAddr()
	nonce := this.unpackedTransaction.FindAddrNonce(fromAddr)

	if len(nonce.Bytes()) == 0 {

		nonce, err = GetAddrNonce(fromAddr)
		if err != nil {

			return false, err
		}
	}

	nonceOne := (*txItr.GetVin())[0].Nonce
	cmp := new(big.Int).Add(&nonce, big.NewInt(1)).Cmp(&nonceOne)
	if cmp == 0 {

	} else if cmp < 0 {

		if new(big.Int).Sub(&nonceOne, &nonce).Cmp(big.NewInt(int64(config.Wallet_addr_tx_count_max))) > 0 {

			return false, errors.New("nonce too big")
		} else {

			return false, nil
		}
	} else {

		return false, errors.New("nonce too small")
	}

	return true, nil
}

func (this *TransactionManager) CheckNotspend(ratio *TransactionRatio) error {
	fromAddr := (*ratio.tx.GetVin())[0].GetPukToAddr()
	spend, spendLock := this.unpackedTransaction.FindAddrSpend(fromAddr)

	if ratio.tx.Class() == config.Wallet_tx_type_voting_reward {
		rfValue := GetCommunityVoteRewardFrozen(fromAddr)
		if spendLock+ratio.spendLock > rfValue {

			return config.ERROR_not_enough
		}
		return nil
	}

	notSpend, _, _ := GetNotspendByAddrOther(*fromAddr)
	if notSpend < spend+ratio.spend {

		return config.ERROR_not_enough
	}
	return nil
}

func (this *TransactionManager) CheckAddrValueBig(ratio *TransactionRatio) bool {
	for _, one := range *ratio.tx.GetVin() {
		addr := one.GetPukToAddr()
		if GetAddrValueBig(addr) {
			return true
		}
	}
	return false
}

func (this *TransactionManager) AddTxs(txs ...TxItr) {
	for _, one := range txs {
		this.AddTx(one)
	}
}

func (this *TransactionManager) DelTx(txs []TxItr) {
	for _, one := range txs {

		this.unpackedNonceDiscontinuity.Delete(utils.Bytes2string(*one.GetHash()))
		this.unpackedTransaction.DelTx(one)

	}
}

func (this *TransactionManager) Package(reward *Tx_reward, height uint64, blocks []Block, createBlockTime int64) ([]TxItr, [][]byte) {

	unacknowledgedTxs := make([]TxItr, 0)

	exclude := make(map[string]string)
	for _, one := range blocks {

		_, txs, err := one.LoadTxs()
		if err != nil {
			return nil, nil
		}
		for _, txOne := range *txs {

			exclude[utils.Bytes2string(*txOne.GetHash())] = ""
			unacknowledgedTxs = append(unacknowledgedTxs, txOne)
		}
	}

	tsrs := this.unpackedTransaction.FindTx()
	txs := make([]TxItr, 0)
	txids := make([][]byte, 0)
	repeatTransaction := true
	checkRepeated := false
	sizeTotal := uint64(0)
	for _, one := range *tsrs {

		repeatTransaction = true
		checkRepeated = false
		addrTxs := make([]TxItr, 0, len(one.trs))
		for _, tsr := range one.trs {

			if sizeTotal+tsr.size > config.Block_size_max {
				repeatTransaction = false
			}
			_, ok := exclude[utils.Bytes2string(*tsr.tx.GetHash())]
			if ok {

				repeatTransaction = false
				break
			}
			if err := tsr.tx.CheckLockHeight(height); err != nil {

				repeatTransaction = false
				break
			}

			if !tsr.tx.CheckRepeatedTx(append(unacknowledgedTxs, addrTxs...)...) {

				repeatTransaction = false
				checkRepeated = true
				break
			}
			addrTxs = append(addrTxs, tsr.tx)
		}
		if checkRepeated {
			txs := make([]TxItr, 0, len(one.trs))
			for _, tsr := range one.trs {
				txs = append(txs, tsr.tx)
			}
			this.DelTx(txs)
			continue
		}
		if repeatTransaction {
			txs = append(txs, addrTxs...)
			sizeTotal = sizeTotal + one.size
			unacknowledgedTxs = append(unacknowledgedTxs, addrTxs...)
			for _, txOne := range addrTxs {
				txids = append(txids, *txOne.GetHash())
			}
		}
	}

	return txs, txids
}

func (this *TransactionManager) CleanIxOvertime(height uint64) {

	if GetHighestBlock() < config.Mining_block_start_height+config.Mining_block_start_height_jump {
		return
	}
	this.unpackedTransaction.CleanOverTimeTx(height)

}

func NewTransactionManager(wb *WitnessBackup) *TransactionManager {
	tm := TransactionManager{
		witnessBackup: wb,

		unpackedNonceDiscontinuity: new(sync.Map),

		unpackedTransaction: NewUnpackedTransaction(),

		tempTxLock:   new(sync.RWMutex),
		tempTx:       make([]TransactionRatio, 0),
		tempTxsignal: make(chan bool, 1),
	}
	tm.loopCheckTxs()

	return &tm
}
