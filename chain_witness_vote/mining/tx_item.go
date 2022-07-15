package mining

import (
	"math/big"

	"github.com/prestonTao/keystore/crypto"
)

const (
	txItem_status_notSpent = int32(0)
	txItem_status_frozen   = int32(1)
	txItem_status_lock     = int32(2)
)

type Balance struct {
	WitnessAddr *crypto.AddressCoin

	item *TxItem
}

type DepositInfo struct {
	WitnessAddr crypto.AddressCoin
	SelfAddr    crypto.AddressCoin
	Value       uint64
}

type TxItem struct {
	Addr  *crypto.AddressCoin
	Value uint64

	VoteType     uint16
	LockupHeight uint64
}

func (this *TxItem) Clean() {
	this.Addr = nil

	this.Value = 0

	this.VoteType = 0
}

type TxItemSort []*TxItem

func (this *TxItemSort) Len() int {
	return len(*this)
}

func (this *TxItemSort) Less(i, j int) bool {
	if (*this)[i].Value < (*this)[j].Value {
		return false
	} else {
		return true
	}
}

func (this *TxItemSort) Swap(i, j int) {
	(*this)[i], (*this)[j] = (*this)[j], (*this)[i]
}

type TxItemCount struct {
	Additems []*TxItem
	SubItems []*TxSubItems
}

type TxSubItems struct {
	Txid      []byte
	VoutIndex uint64
	Addr      crypto.AddressCoin
}

type TxItemCountMap struct {
	AddItems map[string]*map[uint64]int64
	Nonce    map[string]big.Int
}
