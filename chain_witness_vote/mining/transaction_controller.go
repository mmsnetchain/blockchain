package mining

import (
	"sync"

	"github.com/prestonTao/keystore/crypto"
)

type TxController interface {
	Factory() interface{}
	CountBalance(deposit *sync.Map, bhvo *BlockHeadVO)
	SyncCount()
	RollbackBalance()

	BuildTx(deposit *sync.Map, srcAddr, addr *crypto.AddressCoin, amount,
		gas, frozenHeight uint64, pwd, comment string, params ...interface{}) (TxItr, error)
	CheckMultiplePayments(txItr TxItr) error

	ParseProto(bs *[]byte) (interface{}, error)
}

var txCtrlMap = new(sync.Map)

func RegisterTransaction(class uint64, txCtrl TxController) {
	txCtrlMap.Store(class, txCtrl)
}

func GetTransactionCtrl(class uint64) TxController {
	itr, ok := txCtrlMap.Load(class)
	if ok {
		txCtrl := itr.(TxController)
		return txCtrl
	} else {
		return nil
	}
}

func GetNewTransaction(class uint64, bs *[]byte) interface{} {
	itr, ok := txCtrlMap.Load(class)
	if ok {
		txCtrl := itr.(TxController)
		if bs == nil {
			return txCtrl.Factory()
		}
		if bs == nil || len(*bs) <= 0 {
			return txCtrl
		}
		tx, err := txCtrl.ParseProto(bs)
		if err != nil {
			return nil
		}
		return tx
	} else {
		return nil
	}
}

func CountBalanceOther(deposit *sync.Map, bhvo *BlockHeadVO) {
	txCtrlMap.Range(func(k, v interface{}) bool {
		txCtrl := v.(TxController)
		txCtrl.CountBalance(deposit, bhvo)
		return true
	})
}
