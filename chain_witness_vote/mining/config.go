package mining

import (
	"bytes"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/config"
	"strconv"

	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

const (
	path_blocks     = "blocks"
	path_chainstate = "chainstate"

	Unit uint64 = 1e8
)

func BuildKeyForUnspentTransaction(txid []byte, voutIndex uint64) []byte {
	txidStr := utils.Bytes2string(txid)
	return []byte(config.AlreadyUsed_tx + txidStr + "_" + strconv.Itoa(int(voutIndex)))
}

func CheckFrozenHeightFree(frozenHeight uint64, freeHeight uint64, freeTime int64) bool {

	if frozenHeight > config.Wallet_frozen_time_min {

		if int64(frozenHeight) > freeTime {
			return false
		} else {
			return true
		}
	} else {

		if frozenHeight > freeHeight {
			return false
		} else {
			return true
		}
	}
}

func CheckNameStore() bool {

	nameinfo := name.FindName(config.Name_store)
	if nameinfo == nil {

		return false
	}

	for _, one := range nameinfo.NetIds {
		if bytes.Equal(nodeStore.NodeSelf.IdInfo.Id, one) {
			return true
		}
	}
	return false
}
