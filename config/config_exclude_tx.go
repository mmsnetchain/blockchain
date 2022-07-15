package config

import (
	"encoding/hex"
	"sync"
)

const BlockRewardHeightOffset = 1600000

const Mining_block_start_height_jump = 0
const WitnessOrderCorrectStart = 0
const WitnessOrderCorrectEnd = 0
const RandomHashHeightMin = 0
const NextHashHeightMax = 0

const CheckAddBlacklistChangeHeight = 100000

const FixBuildGroupBUGHeightMax = 0

var RandomHashFixed = []byte{}

const RandomHashFixedStr = "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"

var SpecialAddrs = []byte{}

var SpecialBlockHash = []byte{}

var CutBlockHeight = uint64(0)
var CutBlockHash = []byte{}

var NextHash = new(sync.Map)

var RandomMap = new(sync.Map)

var Exclude_Tx = []ExcludeTx{}

var BlockHashs = make([][]byte, 0)

const Mining_witness_average_height = 0
const Reward_witness_height = 0

const Reward_witness_height_new = 0

var DBUG_import_height_max = uint64(0)

const Store_name_new = "storereward"
const Store_name_new_height = 160000

func init() {

	BuildRandom()
	CustomOrderBlockHash()

	for i, one := range Exclude_Tx {
		bs, err := hex.DecodeString(one.TxStr)
		if err != nil {
			panic("hash:" + err.Error())
		}
		Exclude_Tx[i].TxByte = bs
		Exclude_Tx[i].TxStr = ""

	}

}

type ExcludeTx struct {
	Height uint64
	TxStr  string
	TxByte []byte
}

func BuildRandom() {

}

func CustomOrderBlockHash() {
	if false {

	}
}

func UpdateSoreName(height uint64) {
	if height == Store_name_new_height {
		Name_store = Store_name_new
	}
}
