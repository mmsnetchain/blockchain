package main

import (
	"bytes"
	"encoding/hex"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	txItem_status_notSpent = int32(0)
	txItem_status_frozen   = int32(1)
	txItem_status_lock     = int32(2)

	maxHeight = 1100591
)

func main() {
	Start()
}

func Start() {

	CountAddrBalance(filepath.Join("wallet", "data"))

	time.Sleep(time.Second * 3)
}

func CountAddrBalance(dir string) (map[string]*map[string]*mining.TxItem, *mining.TxItemManager, *mining.TxItemManagerOld) {

	db.InitDB(config.DB_path, config.DB_path_temp)
	sqlite3_db.Init()

	beforBlockHash, err := db.LevelDB.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return nil, nil, nil
	}

	rewardTotal := uint64(0)
	balanceCirculationTotal := uint64(0)
	balanceFrozenTotal := uint64(0)

	for beforBlockHash != nil {
		runtime.GC()
		bhvo, err := mining.LoadBlockHeadVOByHash(beforBlockHash)
		if err != nil {
			engine.Log.Info(" :%s", err.Error())
			return nil, nil, nil
		}

		engine.Log.Info("%d -----------------------------------\n%s\nnext", bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Hash))

		for _, txBase := range bhvo.Txs {

			isExclude := false

			for _, one := range config.Exclude_Tx {
				if bhvo.BH.Height != one.Height {
					continue
				}

				if bytes.Equal(one.TxByte, *txBase.GetHash()) {

					isExclude = true
					break
				}
			}
			if isExclude {
				continue
			}
			if bhvo.BH.Height > config.Mining_block_start_height_jump {
				if err := txBase.Check(); err != nil {
					panic(err.Error())
				}
			}

			if txBase.Class() == config.Wallet_tx_type_mining {
				for _, one := range *txBase.GetVout() {
					rewardTotal += one.Value
				}
				engine.Log.Info(" %d", rewardTotal)
			} else {

				for _, one := range *txBase.GetVin() {

					txBase, err := mining.LoadTxBase(one.Txid)
					if err != nil {
						panic("error:2")
						return nil, nil, nil
					}
					vout := (*txBase.GetVout())[one.Vout]

					balanceCirculationTotal -= vout.Value
				}
			}

			for _, one := range *txBase.GetVout() {
				if one.FrozenHeight > maxHeight {
					balanceFrozenTotal += one.Value
				} else {
					balanceCirculationTotal += one.Value
				}
			}
		}

		if bhvo.BH.Nextblockhash != nil && len(bhvo.BH.Nextblockhash) > 0 {
			beforBlockHash = &bhvo.BH.Nextblockhash
		} else {
			beforBlockHash = nil
		}

	}
	engine.Log.Info(":%d :%d :%d", rewardTotal, balanceCirculationTotal, balanceFrozenTotal)

	return nil, nil, nil
}

func PrintBalanceMap(notSpentBalance map[string]*map[string]*mining.TxItem, tim *mining.TxItemManager, timOld *mining.TxItemManagerOld) {

	for addrStr, items := range notSpentBalance {
		addrCoin := crypto.AddressCoin([]byte(addrStr))

		notSpend := uint64(0)
		frozen := uint64(0)
		for _, item := range *items {
			if item.Status == txItem_status_frozen {
				notSpend += item.Value
			} else {
				frozen += item.Value
			}
		}

		notSpend2 := uint64(0)
		frozen2 := uint64(0)
		items := tim.FindBalanceNotSpent(addrCoin)
		for _, one := range items {
			notSpend2 += one.Value
		}
		items = tim.FindBalanceFrozen(addrCoin)
		for _, one := range items {
			frozen2 += one.Value
		}

		engine.Log.Info(":%s :%d :%d  2:%d 2:%d", addrCoin.B58String(), notSpend, frozen, notSpend2, frozen2)
	}

}

func CountBalances(bhvo *mining.BlockHeadVO) mining.TxItemCount {

	itemCount := mining.TxItemCount{
		Additems: make([]*mining.TxItem, 0),
		SubItems: make([]*mining.TxSubItems, 0),
	}
	itemsChan := make(chan *mining.TxItemCount, len(bhvo.Txs))

	wg := new(sync.WaitGroup)
	wg.Add(len(bhvo.Txs))
	go func() {
		for i := 0; i < len(bhvo.Txs); i++ {

			one := <-itemsChan

			if one != nil {
				itemCount.Additems = append(itemCount.Additems, one.Additems...)
				itemCount.SubItems = append(itemCount.SubItems, one.SubItems...)

			}
			wg.Done()
		}
	}()

	NumCPUTokenChan := make(chan bool, runtime.NumCPU()*6)
	for _, txItr := range bhvo.Txs {

		go countBalancesTxOne(txItr, bhvo.BH.Height, NumCPUTokenChan, itemsChan)
	}

	wg.Wait()

	return itemCount
}

func countBalancesTxOne(txItr mining.TxItr, height uint64, tokenCPU chan bool, itemChan chan *mining.TxItemCount) {
	tokenCPU <- false

	txItr.BuildHash()

	itemCount := txItr.CountTxItems(height)

	itemChan <- itemCount

	<-tokenCPU
}
