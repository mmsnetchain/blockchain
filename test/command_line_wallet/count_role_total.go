package main

import (
	"encoding/hex"
	"github.com/prestonTao/keystore/crypto"
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
	"strconv"
	"sync"
	"time"
)

const (
	txItem_status_notSpent = int32(0)
	txItem_status_frozen   = int32(1)
	txItem_status_lock     = int32(2)
)

func main() {

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

	NotSpentBalance := make(map[string]*map[string]*mining.TxItem)
	NotSpentBalanceHex := make(map[string]*map[string]*mining.TxItem)

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

		for _, txBase := range bhvo.Txs {

			if txBase.Class() == config.Wallet_tx_type_mining {

				for _, one := range *txBase.GetVout() {
					rewardTotal += one.Value
				}
				engine.Log.Info(" %d", rewardTotal)
			} else {

				for _, one := range *txBase.GetVin() {

					addrCoin := crypto.BuildAddr(config.AddrPre, one.Puk)
					addrStr := utils.Bytes2string(addrCoin)

					txBase, err := mining.LoadTxBase(one.Txid)
					if err != nil {
						panic("error:2")
						return nil, nil, nil
					}

					vout := (*txBase.GetVout())[one.Vout]

					free := mining.CheckFrozenHeightFree(vout.FrozenHeight-1, bhvo.BH.Height, bhvo.BH.Time)
					if free {

					} else {
						engine.Log.Error("，:%s %d %d %d %d", hex.EncodeToString(one.Txid), one.Vout, vout.FrozenHeight, bhvo.BH.Height, bhvo.BH.Time)

					}

					key := utils.Bytes2string(one.Txid) + "_" + strconv.Itoa(int(one.Vout))

					items, ok := NotSpentBalance[addrStr]
					if !ok {
						panic("")
						return nil, nil, nil
					}
					_, ok = (*items)[key]
					if !ok {
						engine.Log.Info("%s %d", hex.EncodeToString(one.Txid), one.Vout)

					}

					delete(*items, key)

					balanceTotal = balanceTotal - vout.Value

					key = hex.EncodeToString(one.Txid) + "_" + strconv.Itoa(int(one.Vout))

					items, ok = NotSpentBalanceHex[addrCoin.B58String()]
					if !ok {
						panic("")
						return nil, nil, nil
					}
					_, ok = (*items)[key]
					if !ok {
						engine.Log.Info("%s %d", hex.EncodeToString(one.Txid), one.Vout)

					}
					delete(*items, key)

				}
			}

			for voutIndex, _ := range *txBase.GetVout() {

				if txBase.Class() == config.Wallet_tx_type_deposit_in || txBase.Class() == config.Wallet_tx_type_vote_in || txBase.Class() == config.Wallet_tx_type_account {
					if voutIndex == 0 {
						continue
					}
				}

				one := (*txBase.GetVout())[voutIndex]
				addrStr := utils.Bytes2string(one.Address)
				key := utils.Bytes2string(*txBase.GetHash()) + "_" + strconv.Itoa(int(voutIndex))

				txItem := mining.TxItem{
					Addr: &one.Address,

					Value: one.Value,
					Txid:  *txBase.GetHash(),

					VoutIndex:    uint64(voutIndex),
					Height:       bhvo.BH.Height,
					LockupHeight: one.FrozenHeight,
				}

				items, ok := NotSpentBalance[addrStr]
				if !ok {
					itemsTemp := make(map[string]*mining.TxItem)
					items = &itemsTemp
					NotSpentBalance[addrStr] = items
				}
				(*items)[key] = &txItem

				balanceTotal = balanceTotal + one.Value

				key = hex.EncodeToString(*txBase.GetHash()) + "_" + strconv.Itoa(int(voutIndex))
				items, ok = NotSpentBalanceHex[one.Address.B58String()]
				if !ok {
					itemsTemp := make(map[string]*mining.TxItem)
					items = &itemsTemp
					NotSpentBalanceHex[one.Address.B58String()] = items
				}
				(*items)[key] = &txItem
			}

		}

		if bhvo.BH.Nextblockhash != nil && len(bhvo.BH.Nextblockhash) > 0 {
			beforBlockHash = &bhvo.BH.Nextblockhash
		} else {
			beforBlockHash = nil
		}

		thatNotSpendTotalHex := uint64(0)
		for _, items := range NotSpentBalanceHex {
			for _, item := range *items {

				thatNotSpendTotalHex += item.Value

			}
		}

		thatNotSpendTotal := uint64(0)
		for _, items := range NotSpentBalance {
			for _, item := range *items {

				thatNotSpendTotal += item.Value

			}
		}

		if thatNotSpendTotalHex != thatNotSpendTotal {
			engine.Log.Info("that:%d thatHex:%d", thatNotSpendTotal, thatNotSpendTotalHex)
			panic("")
		}

	}
	engine.Log.Info(":%d :%d", rewardTotal, balanceTotal)

	return NotSpentBalance, nil, nil
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
		items := tim.FindBalance(addrCoin)
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
