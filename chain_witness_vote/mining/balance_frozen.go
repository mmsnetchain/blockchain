package mining

import (
	"github.com/prestonTao/libp2parea/engine"
)

func (this *BalanceManager) Frozen(items []*TxItem, tx TxItr) {

}

func (this *BalanceManager) Unfrozen(blockHeight uint64, blockTime int64) {

	txItems, err := GetFrozenHeight(blockHeight, blockTime)
	if err != nil {
		engine.Log.Error("GetFrozenHeight error:%s", err.Error())
		return
	}
	for _, one := range *txItems {

		_, oldValue := GetNotSpendBalance(one.Addr)

		oldValue += one.Value
		SetNotSpendBalance(one.Addr, oldValue)

		oldValue = GetAddrFrozenValue(one.Addr)

		oldValue -= one.Value
		SetAddrFrozenValue(one.Addr, oldValue)
	}
	RemoveFrozenHeight(blockHeight, blockTime)

	this.UnlockByHeight(blockHeight, blockTime)
}
