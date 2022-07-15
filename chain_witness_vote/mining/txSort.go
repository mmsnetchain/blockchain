package mining

import (
	"bytes"
	"math/big"
	"mmschainnewaccount/config"
	"sync"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type TransactionRatio struct {
	tx        TxItr
	size      uint64
	gas       uint64
	Ratio     *big.Int
	spend     uint64
	spendLock uint64
}

type RatioSort struct {
	txRatio []TransactionRatio
}

func (this *RatioSort) Len() int {
	return len(this.txRatio)
}

func (this *RatioSort) Less(i, j int) bool {
	if this.txRatio[i].Ratio.Cmp(this.txRatio[j].Ratio) < 0 {
		return false
	} else {
		return true
	}
}

func (this *RatioSort) Swap(i, j int) {
	this.txRatio[i], this.txRatio[j] = this.txRatio[j], this.txRatio[i]
}

type TransactionsRatio struct {
	lock      *sync.RWMutex
	trs       []*TransactionRatio
	size      uint64
	gas       uint64
	Ratio     *big.Int
	spend     uint64
	spendLock uint64
}

type AddrRatioSort struct {
	txRatio []TransactionsRatio
}

func (this *AddrRatioSort) Len() int {
	return len(this.txRatio)
}

func (this *AddrRatioSort) Less(i, j int) bool {
	if this.txRatio[i].Ratio.Cmp(this.txRatio[j].Ratio) < 0 {
		return false
	} else {
		return true
	}
}

func (this *AddrRatioSort) Swap(i, j int) {
	this.txRatio[i], this.txRatio[j] = this.txRatio[j], this.txRatio[i]
}

type UnpackedTransaction struct {
	addrs *sync.Map
}

func (this *UnpackedTransaction) AddTx(tr *TransactionRatio) {
	addr := (*tr.tx.GetVin())[0].GetPukToAddr()
	addrStr := utils.Bytes2string(*addr)
	tsrItr, ok := this.addrs.Load(addrStr)
	if ok {
		tsr := tsrItr.(*TransactionsRatio)
		if len(tsr.trs) >= config.Wallet_addr_tx_count_max || tsr.size+tr.size > config.Block_size_max {
			return
		}
		tsr.lock.Lock()
		tsr.gas += tr.gas
		tsr.size += tr.size
		tsr.trs = append(tsr.trs, tr)
		tsr.spend += tr.spend
		tsr.spendLock += tr.spendLock
		div1 := new(big.Int).Mul(big.NewInt(int64(tsr.gas)), big.NewInt(100000000))
		div2 := big.NewInt(int64(tsr.size))
		tsr.Ratio = new(big.Int).Div(div1, div2)
		tsr.lock.Unlock()
	} else {
		tsr := TransactionsRatio{
			lock:      new(sync.RWMutex),
			trs:       make([]*TransactionRatio, 0),
			size:      tr.size,
			gas:       tr.gas,
			Ratio:     tr.Ratio,
			spend:     tr.spend,
			spendLock: tr.spendLock,
		}
		tsr.trs = append(tsr.trs, tr)
		this.addrs.Store(addrStr, &tsr)
	}
}

func (this *UnpackedTransaction) FindAddrNonce(addr *crypto.AddressCoin) big.Int {
	engine.Log.Info("FindAddrNonce")
	tsrItr, ok := this.addrs.Load(utils.Bytes2string(*addr))
	if ok {
		engine.Log.Info("FindAddrNonce")
		tsr := tsrItr.(*TransactionsRatio)
		engine.Log.Info("FindAddrNonce:%d", len(tsr.trs))
		trOne := tsr.trs[len(tsr.trs)-1]
		engine.Log.Info("FindAddrNonce")

		vins := trOne.tx.GetVin()
		engine.Log.Info("FindAddrNonce:%d", len(tsr.trs))
		nonce := (*vins)[0].Nonce
		engine.Log.Info("FindAddrNonce:%d", len(tsr.trs))
		return nonce
	} else {
		engine.Log.Info("FindAddrNonce")
		return big.Int{}
	}
}

func (this *UnpackedTransaction) FindAddrSpend(addr *crypto.AddressCoin) (uint64, uint64) {
	tsrItr, ok := this.addrs.Load(utils.Bytes2string(*addr))
	if ok {
		tsr := tsrItr.(*TransactionsRatio)
		return tsr.spend, tsr.spendLock
	} else {
		return 0, 0
	}
}

func (this *UnpackedTransaction) DelTx(tx TxItr) {
	addr := (*tx.GetVin())[0].GetPukToAddr()
	addrStr := utils.Bytes2string(*addr)
	tsrItr, ok := this.addrs.Load(addrStr)
	if ok {
		tsr := tsrItr.(*TransactionsRatio)
		tsr.lock.Lock()
		for i, one := range tsr.trs {
			if bytes.Equal(*one.tx.GetHash(), *tx.GetHash()) {
				tsr.gas -= one.gas
				tsr.size -= one.size
				tsr.spend -= one.spend
				tsr.spendLock -= one.spendLock
				temp := tsr.trs[:i]
				tsr.trs = append(temp, tsr.trs[i+1:]...)
			}
		}
		tsr.lock.Unlock()

		if len(tsr.trs) <= 0 {
			this.addrs.Delete(addrStr)
		}
		return
	} else {
		return
	}
}

func (this *UnpackedTransaction) FindTx() *[]*TransactionsRatio {
	engine.Log.Info("")
	tsrs := make([]*TransactionsRatio, 0)
	this.addrs.Range(func(k, v interface{}) bool {
		engine.Log.Info(" 111")
		tsr := v.(*TransactionsRatio)

		tsrs = append(tsrs, tsr)
		return true
	})
	return &tsrs
}

func (this *UnpackedTransaction) ExistTxByAddrTxid(tx TxItr) bool {
	fromAddr := (*tx.GetVin())[0].GetPukToAddr()
	addrStr := utils.Bytes2string(*fromAddr)
	tsrItr, ok := this.addrs.Load(addrStr)
	if ok {
		tsr := tsrItr.(*TransactionsRatio)
		for _, one := range tsr.trs {
			if bytes.Equal(*one.tx.GetHash(), *tx.GetHash()) {
				return true
			}
		}
	}
	return false
}

func (this *UnpackedTransaction) CleanOverTimeTx(height uint64) {
	this.addrs.Range(func(k, v interface{}) bool {
		tsr := v.(*TransactionsRatio)
		for index, one := range tsr.trs {
			if one.tx.GetLockHeight() > height {
				continue
			}
			temp := tsr.trs[:index]
			tsr.trs = append(temp, tsr.trs[index+1:]...)
		}

		if len(tsr.trs) <= 0 {
			this.addrs.Delete(k)
		}
		return true
	})
}

func NewUnpackedTransaction() *UnpackedTransaction {
	ut := UnpackedTransaction{
		addrs: new(sync.Map),
	}
	return &ut
}
