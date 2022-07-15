package mining

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/config"
	"sync"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

func (this *BalanceManager) FindNonce(addr *crypto.AddressCoin) (nonce big.Int) {

	txListItr, ok := this.cacheTxlock.Load(utils.Bytes2string(*addr))
	if ok {
		txList := txListItr.(*TxList)
		txList.Lock.Lock()
		if len(txList.Txs) != 0 {
			txOne := txList.Txs[len(txList.Txs)-1]
			nonce = (*txOne.GetVin())[0].Nonce
		}
		txList.Lock.Unlock()
	}
	if len(nonce.Bytes()) != 0 {
		return
	}
	var err error
	nonce, err = GetAddrNonce(addr)
	if err != nil {
		engine.Log.Info("GetAddrNonce error:%s", err.Error())
	}
	return

}

func (this *BalanceManager) AddLockTx(tx TxItr) {
	addr := (*tx.GetVin())[0].GetPukToAddr()
	txListItr, ok := this.cacheTxlock.Load(utils.Bytes2string(*addr))
	if ok {

		txList := txListItr.(*TxList)
		txList.AddTx(tx)
	} else {
		txList := NewTxList()
		txList.AddTx(tx)
		this.cacheTxlock.Store(utils.Bytes2string(*addr), txList)
	}
}

func (this *BalanceManager) DelLockTx(tx TxItr) {
	addr := (*tx.GetVin())[0].GetPukToAddr()
	txListItr, ok := this.cacheTxlock.Load(utils.Bytes2string(*addr))
	if ok {
		txList := txListItr.(*TxList)
		txList.DelTx(tx)
		if txList.GetLen() == 0 {
			this.cacheTxlock.Delete(utils.Bytes2string(*addr))
		}
	}
}

func (this *BalanceManager) FindLockTotalByAddr(addr *crypto.AddressCoin) (uint64, uint64) {
	txListItr, ok := this.cacheTxlock.Load(utils.Bytes2string(*addr))
	if ok {
		txList := txListItr.(*TxList)
		return txList.GetLockTotal()
	} else {
		return 0, 0
	}
}

func (this *BalanceManager) UnlockByHeight(blockHeight uint64, blockTime int64) {
	delTxs := make([]TxItr, 0)
	this.cacheTxlock.Range(func(k, v interface{}) bool {
		txList := v.(*TxList)
		for index, _ := range txList.Txs {
			one := txList.Txs[index]
			engine.Log.Info(":%d %d %d", one.GetLockHeight(), blockHeight, blockTime)
			if CheckFrozenHeightFree(one.GetLockHeight(), blockHeight, blockTime) {
				engine.Log.Info(":%d %d %d", one.GetLockHeight(), blockHeight, blockTime)
				delTxs = append(delTxs, one)
			}
		}
		return true
	})
	for _, one := range delTxs {
		engine.Log.Info("")
		this.DelLockTx(one)
	}

}

type TxList struct {
	Lock                *sync.RWMutex
	Txs                 []TxItr
	lockNotspendTotal   uint64
	lockVoteRewardTotal uint64
}

func (this *TxList) AddTx(tx TxItr) {
	this.Lock.Lock()
	if this.Txs == nil {
		this.Txs = make([]TxItr, 0)
	}
	this.Txs = append(this.Txs, tx)
	if tx.Class() == config.Wallet_tx_type_voting_reward {
		this.lockVoteRewardTotal = tx.GetSpend()
	} else {
		this.lockNotspendTotal = tx.GetSpend()
	}

	this.Lock.Unlock()
}

func (this *TxList) DelTx(tx TxItr) {

	this.Lock.Lock()
	for index, one := range this.Txs {
		if bytes.Equal(*one.GetHash(), *tx.GetHash()) {
			engine.Log.Info(":%s", hex.EncodeToString(*one.GetHash()))
			temp := this.Txs[:index]
			this.Txs = append(temp, this.Txs[index+1:]...)
			spend := one.GetSpend()
			if tx.Class() == config.Wallet_tx_type_voting_reward {
				this.lockVoteRewardTotal -= spend
			} else {
				this.lockNotspendTotal -= spend
			}

			break
		}
	}
	if len(this.Txs) == 0 {
		this.Txs = nil
	}
	this.Lock.Unlock()
}

func (this *TxList) GetLen() (lenght int) {
	this.Lock.RLock()
	if this.Txs == nil {
		lenght = 0
	} else {
		lenght = len(this.Txs)
	}
	this.Lock.RUnlock()
	return
}

func (this *TxList) GetLockTotal() (lockNotspendTotal, lockVoteRewardTotal uint64) {
	this.Lock.RLock()
	lockNotspendTotal = this.lockNotspendTotal
	lockVoteRewardTotal = this.lockVoteRewardTotal
	this.Lock.RUnlock()
	return
}

func NewTxList() *TxList {
	return &TxList{
		Lock: new(sync.RWMutex),
	}
}
