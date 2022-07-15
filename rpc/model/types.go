package model

import (
	"mmschainnewaccount/chain_witness_vote/mining/token"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Result struct {
	Result interface{} `json:"result"`
}

type Getinfo struct {
	Netid          []byte                 `json:"netid"`
	TotalAmount    uint64                 `json:"TotalAmount"`
	Balance        uint64                 `json:"balance"`
	BalanceFrozen  uint64                 `json:"BalanceFrozen"`
	BalanceLockup  uint64                 `json:"BalanceLockup"`
	Testnet        bool                   `json:"testnet"`
	Blocks         uint64                 `json:"blocks"`
	Group          uint64                 `json:"group"`
	StartingBlock  uint64                 `json:"StartingBlock"`
	StartBlockTime uint64                 `json:"StartBlockTime"`
	HighestBlock   uint64                 `json:"HighestBlock"`
	CurrentBlock   uint64                 `json:"CurrentBlock"`
	PulledStates   uint64                 `json:"PulledStates"`
	BlockTime      uint64                 `json:"BlockTime"`
	LightNode      uint64                 `json:"LightNode"`
	CommunityNode  uint64                 `json:"CommunityNode"`
	WitnessNode    uint64                 `json:"WitnessNode"`
	NameDepositMin uint64                 `json:"NameDepositMin"`
	AddrPre        string                 `json:"AddrPre"`
	TokenBalance   []token.TokenBalanceVO `json:"TokenBalance"`
}

type GetNewAddress struct {
	Address string `json:"address"`
}

type GetAccount struct {
	Balance       uint64 `json:"Balance"`
	BalanceFrozen uint64 `json:"BalanceFrozen"`
}

func Tojson(data interface{}) ([]byte, error) {
	res, err := json.Marshal(Result{Result: data})
	return res, err
}
