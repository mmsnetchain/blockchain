package mining

import (
	"github.com/prestonTao/keystore/crypto"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"strconv"

	"github.com/gogo/protobuf/proto"
)

var balanceHistoryManager = NewBalanceHistory()

type BalanceHistory struct {
	GenerateMaxId *big.Int
	ForkNo        uint64
}

type HistoryItem struct {
	GenerateId *big.Int
	IsIn       bool
	Type       uint64
	InAddr     []*crypto.AddressCoin
	OutAddr    []*crypto.AddressCoin
	Value      uint64
	Txid       []byte
	Height     uint64

	Payload []byte
}

func (this *HistoryItem) Proto() ([]byte, error) {

	inaddrs := make([][]byte, 0)
	for _, one := range this.InAddr {
		inaddrs = append(inaddrs, *one)
	}

	outaddrs := make([][]byte, 0)
	for _, one := range this.OutAddr {
		outaddrs = append(outaddrs, *one)
	}

	hip := go_protos.HistoryItem{
		GenerateId: this.GenerateId.Bytes(),
		IsIn:       this.IsIn,
		Type:       this.Type,
		InAddr:     inaddrs,
		OutAddr:    outaddrs,
		Value:      this.Value,
		Txid:       this.Txid,
		Height:     this.Height,
		Payload:    this.Payload,
	}
	return hip.Marshal()
}

func ParseHistoryItem(bs *[]byte) (*HistoryItem, error) {
	if bs == nil {
		return nil, nil
	}
	hip := new(go_protos.HistoryItem)
	err := proto.Unmarshal(*bs, hip)
	if err != nil {
		return nil, err
	}

	inaddrs := make([]*crypto.AddressCoin, 0)
	for _, one := range hip.InAddr {
		addrOne := crypto.AddressCoin(one)
		inaddrs = append(inaddrs, &addrOne)
	}

	outaddrs := make([]*crypto.AddressCoin, 0)
	for _, one := range hip.OutAddr {
		addrOne := crypto.AddressCoin(one)
		outaddrs = append(outaddrs, &addrOne)
	}
	hi := HistoryItem{
		GenerateId: new(big.Int).SetBytes(hip.GenerateId),
		IsIn:       hip.IsIn,
		Type:       hip.Type,
		InAddr:     inaddrs,
		OutAddr:    outaddrs,
		Value:      hip.Value,
		Txid:       hip.Txid,
		Height:     hip.Height,
		Payload:    hip.Payload,
	}
	return &hi, nil
}

func (this *BalanceHistory) Add(hi HistoryItem) error {

	if hi.GenerateId == nil {
		hi.GenerateId = this.GenerateMaxId
		this.GenerateMaxId = new(big.Int).Add(this.GenerateMaxId, big.NewInt(1))
	} else {
		if hi.GenerateId.Cmp(this.GenerateMaxId) == 0 {
			this.GenerateMaxId = new(big.Int).Add(this.GenerateMaxId, big.NewInt(1))
		}
	}
	bs, err := hi.Proto()

	if err != nil {
		return err
	}

	key := []byte(config.LEVELDB_Head_history_balance + strconv.Itoa(int(this.ForkNo)) + "_" + hi.GenerateId.String())

	return db.LevelTempDB.Save(key, &bs)
}

func (this *BalanceHistory) Get(start *big.Int, total int) []HistoryItem {
	if total == 0 {
		total = config.Wallet_balance_history
	}
	if start == nil {
		start = new(big.Int).Sub(this.GenerateMaxId, big.NewInt(1))
	}
	his := make([]HistoryItem, 0)

	key := config.LEVELDB_Head_history_balance + strconv.Itoa(int(this.ForkNo)) + "_"
	for i := 0; i < total; i++ {
		keyOne := key + new(big.Int).Sub(start, big.NewInt(int64(i))).String()
		bs, err := db.LevelTempDB.Find([]byte(keyOne))
		if err != nil {
			continue
		}

		hi, err := ParseHistoryItem(bs)

		if err != nil {
			continue
		}
		his = append(his, *hi)
	}
	return his
}

func NewBalanceHistory() *BalanceHistory {
	return &BalanceHistory{

		GenerateMaxId: big.NewInt(0),
	}
}
