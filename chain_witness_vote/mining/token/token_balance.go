package token

import (
	"mmschainnewaccount/chain_witness_vote/mining"
	"sync"

	"github.com/prestonTao/keystore/crypto"
)

type TokenTxItemManager struct {
	lock *sync.RWMutex
}

type TokenTxItemCount struct {
	PublishTxidStr string
	Name           string
	Symbol         string
	Supply         uint64
	Additems       *mining.TxItem
}

func (this *TokenTxItemManager) CountTxItem(txItemCounts *map[string]*map[string]int64) {

	for txidStr, itemValue := range *txItemCounts {
		txid := []byte(txidStr)
		for addrStr, value := range *itemValue {
			addr := crypto.AddressCoin([]byte(addrStr))
			_, oldvalue := mining.GetNotSpendBalanceToken(&txid, &addr)
			if value > 0 {
				mining.SetNotSpendBalanceToken(&txid, &addr, oldvalue+uint64(value))
			} else if value < 0 {
				mining.SetNotSpendBalanceToken(&txid, &addr, oldvalue-uint64(value))
			}
		}
	}

}

func (this *TokenTxItemManager) FindBalanceByAddr(txid, addrStr string) *TokenBalance {

	tb := new(TokenBalance)
	this.lock.RLock()
	defer func() {
		this.lock.RUnlock()

	}()

	return tb
}

func (this *TokenTxItemManager) FindTokenBalanceForTxid(txidbs []byte) (map[string]uint64, map[string]uint64, map[string]uint64) {

	balance := make(map[string]uint64)
	balanceFrozen := make(map[string]uint64)
	balanceLockup := make(map[string]uint64)

	return balance, balanceFrozen, balanceLockup
}

func (this *TokenTxItemManager) FindTokenBalanceForAll() []TokenBalance {

	tbs := make([]TokenBalance, 0)

	return tbs
}

func (this *TokenTxItemManager) GetReadyPayToken(txidbs []byte, srcAddr crypto.AddressCoin) []*mining.TxItem {

	itemsPay := make([]*mining.TxItem, 0)

	return itemsPay
}

func (this *TokenTxItemManager) FrozenToken(txidBs []byte, items []*mining.TxItem, tx mining.TxItr) {

}

func (this *TokenTxItemManager) Unfrozen(blockHeight uint64) {

}

func NewTokenTxItemManager() *TokenTxItemManager {
	return &TokenTxItemManager{
		lock: new(sync.RWMutex),
	}
}

var tokenManager = NewTokenTxItemManager()

func CountToken(txItemCounts *map[string]*map[string]int64) {
	tokenManager.CountTxItem(txItemCounts)
}

func FindTokenBalanceForTxid(txidbs []byte) (map[string]uint64, map[string]uint64, map[string]uint64) {
	return tokenManager.FindTokenBalanceForTxid(txidbs)
}

func FindBalanceByAddr(txid string, addr string) *TokenBalance {
	return tokenManager.FindBalanceByAddr(txid, addr)
}

func FindTokenBalanceForAll() []TokenBalance {
	tbs := tokenManager.FindTokenBalanceForAll()
	for i, one := range tbs {
		tokeninfo, err := FindTokenInfo(one.TokenId)
		if err != nil {
			continue
		}
		tbs[i].Name = tokeninfo.Name
		tbs[i].Symbol = tokeninfo.Symbol
		tbs[i].Supply = tokeninfo.Supply
	}
	return tbs
}

func FrozenToken(txid []byte, items []*mining.TxItem, tx mining.TxItr) {
	tokenManager.FrozenToken(txid, items, tx)
}

func UnfrozenToken(blockHeight uint64) {
	tokenManager.Unfrozen(blockHeight)
}

func GetReadyPayToken(txid *[]byte, srcAddress *crypto.AddressCoin, amount uint64) (total uint64, txItems *mining.TxItem) {

	if srcAddress != nil && len(*srcAddress) > 0 {

		tis, value := mining.GetNotSpendBalanceToken(txid, srcAddress)
		return value, tis
	} else {
		tis, value := mining.FindNotSpendBalanceToken(txid, amount)
		return value, tis
	}

}

type TokenBalance struct {
	TokenId       []byte
	Name          string
	Symbol        string
	Supply        uint64
	Balance       uint64
	BalanceFrozen uint64
	BalanceLockup uint64
}

type TokenBalanceVO struct {
	TokenId       string
	Name          string
	Symbol        string
	Supply        uint64
	Balance       uint64
	BalanceFrozen uint64
	BalanceLockup uint64
}
