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

	m, tim, timOld := CountAddrBalance(filepath.Join("wallet", "data"))

	PrintBalanceMap(m, tim, timOld)

	time.Sleep(time.Second * 3)
}

func CountAddrBalance(dir string) (map[string]*map[string]*mining.TxItem, *mining.TxItemManager, *mining.TxItemManagerOld) {

	db.InitDB(config.DB_path, config.DB_path_temp)
	sqlite3_db.Init()

	tim := mining.NewTxItemManager()

	beforBlockHash, err := db.LevelDB.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return nil, nil, nil
	}

	balanceTotal := uint64(0)
	rewardTotal := uint64(0)

	for beforBlockHash != nil {

		bhvo, err := mining.LoadBlockHeadVOByHash(beforBlockHash)
		if err != nil {
			engine.Log.Info(" :%s", err.Error())
			return nil, nil, nil
		}

		engine.Log.Info("%d -----------------------------------\n%s\nnext", bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Hash))

		if bhvo.BH.Nextblockhash != nil && len(bhvo.BH.Nextblockhash) > 0 {
			beforBlockHash = &bhvo.BH.Nextblockhash
		} else {
			beforBlockHash = nil
		}

		for _, one := range config.Exclude_Tx {
			if bhvo.BH.Height != one.Height {
				continue
			}
			for j, two := range bhvo.Txs {
				if !bytes.Equal(one.TxByte, *two.GetHash()) {

					continue
				}

				notExcludeTx := bhvo.Txs[:j]
				bhvo.Txs = append(notExcludeTx, bhvo.Txs[j+1:]...)

				break
			}
		}

		tic := CountBalances(bhvo)
		tim.CountTxItem(tic, bhvo.BH.Height, bhvo.BH.Time)

		tim.Unfrozen(bhvo.BH.Height, bhvo.BH.Time)

	}
	engine.Log.Info(":%d :%d", rewardTotal, balanceTotal)

	return nil, tim, nil
}

func PrintBalanceMap(notSpentBalance map[string]*map[string]*mining.TxItem, tim *mining.TxItemManager, timOld *mining.TxItemManagerOld) {

	nb, fb, lb := tim.FindBalanceAllAddrs()

	engine.Log.Info("")
	for k, v := range nb {
		addrCoin := crypto.AddressCoin([]byte(k))
		engine.Log.Info("%s %d", addrCoin.B58String(), v)
	}

	engine.Log.Info("")
	for k, v := range fb {
		addrCoin := crypto.AddressCoin([]byte(k))
		engine.Log.Info("%s %d", addrCoin.B58String(), v)
	}
	engine.Log.Info("")
	for k, v := range lb {
		addrCoin := crypto.AddressCoin([]byte(k))
		engine.Log.Info("%s %d", addrCoin.B58String(), v)
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
