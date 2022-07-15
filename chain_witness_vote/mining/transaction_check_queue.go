package mining

import (
	"encoding/hex"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"runtime"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

var checkTxQueue = make(chan TxItr, 100000)

func init() {
	go startCheckQueue()
}

func startCheckQueue() {
	utils.Go(func() {
		NumCPUTokenChan := make(chan bool, config.CPUNUM)
		for txItr := range checkTxQueue {
			NumCPUTokenChan <- false
			go checkMulticastTransaction(txItr, NumCPUTokenChan)
		}
	})
}

func checkMulticastTransaction(txbase TxItr, tokenCPU chan bool) {
	defer func() {
		<-tokenCPU
	}()

	if len(*txbase.GetVin()) > config.Mining_pay_vin_max {

		engine.Log.Warn(config.ERROR_pay_vin_too_much.Error())
		return
	}

	if err := txbase.CheckLockHeight(GetLongChain().GetCurrentBlock()); err != nil {

		engine.Log.Warn("Failed to verify transaction lock height")
		return
	}

	if GetHighestBlock() > config.Mining_block_start_height+config.Mining_block_start_height_jump {
		if err := txbase.CheckSign(); err != nil {

			runtime.GC()
			engine.Log.Warn("Failed to verify transaction signature %s %s", hex.EncodeToString(*txbase.GetHash()), err.Error())
			return
		}
		runtime.GC()
	}

	exist, _ := db.LevelDB.CheckHashExist(*txbase.GetHash())
	notImport, _ := db.LevelDB.CheckHashExist(config.BuildTxNotImport(*txbase.GetHash()))
	if exist && !notImport {

		engine.Log.Warn("Transaction hash collision is the same %s", hex.EncodeToString(*txbase.GetHash()))
		return
	}

	forks.GetLongChain().transactionManager.AddTx(txbase)
}
