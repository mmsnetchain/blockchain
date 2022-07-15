package main

import (
	"encoding/hex"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"path/filepath"
	"runtime"
	"time"
)

func main() {

	go func() {
		for {
			engine.Log.Info("NumGoroutine:%d", runtime.NumGoroutine())
			time.Sleep(time.Minute)

		}
	}()
	Start()

}

func Start() {
	CountAddrBalance(filepath.Join("wallet", "data"))

	time.Sleep(time.Second * 3)
}

func CountAddrBalance(dir string) (map[string]uint64, map[string]uint64, *mining.TxItemManager, *mining.TxItemManagerOld) {
	sqlite3_db.Init()
	db.InitDB(dir)

	beforBlockHash, err := db.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return nil, nil, nil, nil
	}

	for beforBlockHash != nil {

		bhvo, err := mining.LoadBlockHeadVOByHash(beforBlockHash)
		if err != nil {
			engine.Log.Info(" :%s", err.Error())
			return nil, nil, nil, nil
		}

		engine.Log.Info("%d -----------------------------------\n%s\nnext", bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Hash))

		for _, txBase := range bhvo.Txs {

			if txBase.Class() != config.Wallet_tx_type_mining {
				for _, one := range *txBase.GetVin() {
					key := mining.BuildKeyForUnspentTransaction(one.Txid, one.Vout)
					if !db.CheckHashExist(key) {
						panic("")
						return nil, nil, nil, nil
					}

				}
			}

			for voutIndex, one := range *txBase.GetVout() {

				bs := utils.Uint64ToBytes(one.FrozenHeight)
				db.Save(mining.BuildKeyForUnspentTransaction(*txBase.GetHash(), uint64(voutIndex)), &bs)

			}

		}

		if bhvo.BH.Nextblockhash != nil && len(bhvo.BH.Nextblockhash) > 0 {
			beforBlockHash = &bhvo.BH.Nextblockhash
		} else {
			beforBlockHash = nil
		}

	}

	return nil, nil, nil, nil
}
