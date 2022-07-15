package main

import (
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"strconv"

	"bytes"
	"encoding/hex"
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

var CountAddrCoin = crypto.AddressFromB58String("MMSHj94iJaS245y5WhCzXPrEtRqhkTVBj8Dd4")

var heightStart = uint64(0)
var heightEnd = uint64(516831)

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

	mNotSpend, tim, timOld := CountAddrBalance(filepath.Join("wallet", "data"))

	PrintBalanceMap(mNotSpend, tim, timOld)

	time.Sleep(time.Second * 3)
}

func CountAddrBalance(dir string) (map[string]*map[string]*mining.TxItem, *mining.TxItemManager, *mining.TxItemManagerOld) {

	db.InitDB(config.DB_path, config.DB_path_temp)
	sqlite3_db.Init()

	tim := mining.NewTxItemManager()

	NotSpentBalance := make(map[string]*map[string]*mining.TxItem)

	beforBlockHash, err := db.LevelDB.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return nil, nil, nil
	}

	for beforBlockHash != nil {

		bhvo, err := mining.LoadBlockHeadVOByHash(beforBlockHash)
		if err != nil {
			engine.Log.Info(" :%s", err.Error())
			return nil, nil, nil
		}

		engine.Log.Info("%d -----------------------------------\n%s\nnext", bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Hash))

		for _, txBase := range bhvo.Txs {

			if txBase.Class() == config.Wallet_tx_type_mining {

			} else {

				for _, one := range *txBase.GetVin() {

					addrCoin := crypto.BuildAddr(config.AddrPre, one.Puk)
					ok := bytes.Equal(addrCoin, CountAddrCoin)
					if !ok {
						continue
					}

					addrStr := utils.Bytes2string(addrCoin)

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

				}
			}

			for voutIndex, _ := range *txBase.GetVout() {

				if txBase.Class() == config.Wallet_tx_type_deposit_in || txBase.Class() == config.Wallet_tx_type_vote_in || txBase.Class() == config.Wallet_tx_type_account {
					if voutIndex == 0 {
						continue
					}
				}
				one := (*txBase.GetVout())[voutIndex]

				ok := bytes.Equal(one.Address, CountAddrCoin)
				if !ok {
					continue
				}

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

			}

		}

		for _, items := range NotSpentBalance {
			for _, item := range *items {
				if mining.CheckFrozenHeightFree(item.LockupHeight, bhvo.BH.Height, bhvo.BH.Time) {

					item.Status = txItem_status_notSpent
				} else {
					item.Status = txItem_status_frozen

				}
			}
		}

		if bhvo.BH.Nextblockhash != nil && len(bhvo.BH.Nextblockhash) > 0 {
			beforBlockHash = &bhvo.BH.Nextblockhash
		} else {
			beforBlockHash = nil
		}

		tic := CountBalances(bhvo)
		tim.CountTxItem(tic, bhvo.BH.Height, bhvo.BH.Time)

		tim.Unfrozen(bhvo.BH.Height, bhvo.BH.Time)

		timNotSpendTotal := uint64(0)
		items := tim.FindBalanceAll()
		for _, one := range items {
			engine.Log.Info("txItem_status_notSpent:%+v", one)
			timNotSpendTotal += one.Value
		}

		thatNotSpendTotal := uint64(0)
		for _, items := range NotSpentBalance {
			for _, item := range *items {
				if item.Status == txItem_status_notSpent {
					engine.Log.Info("txItem_status_notSpent:%+v", item)
					thatNotSpendTotal += item.Value
				}
			}
		}
		if timNotSpendTotal != thatNotSpendTotal {
			engine.Log.Info("tim:%d that:%d", timNotSpendTotal, thatNotSpendTotal)

			engine.Log.Info("------------------------------")

			timNotSpendTotal := uint64(0)
			items := tim.FindBalanceAll()
			for _, one := range items {
				engine.Log.Info("txItem_status_notSpent:%+v", one)
				timNotSpendTotal += one.Value
			}
			engine.Log.Info("------------------------------")
			thatNotSpendTotal := uint64(0)
			for _, items := range NotSpentBalance {
				for _, item := range *items {
					if item.Status == txItem_status_notSpent {
						engine.Log.Info("txItem_status_notSpent:%+v", item)
						thatNotSpendTotal += item.Value
					}
				}
			}

			panic("")
		}

	}

	return NotSpentBalance, tim, nil
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

	itemCount := mining.TxItemCount{
		Additems: make([]*mining.TxItem, 0),
		SubItems: make([]*mining.TxSubItems, 0),
	}

	txItr.BuildHash()

	for _, vin := range *txItr.GetVin() {

		if txItr.Class() == config.Wallet_tx_type_mining {
			continue
		}

		addrOne := crypto.BuildAddr("MMS", vin.Puk)
		ok := bytes.Equal(addrOne, CountAddrCoin)

		if !ok {
			continue
		}

		itemCount.SubItems = append(itemCount.SubItems, &mining.TxSubItems{
			Txid:      vin.Txid,
			VoutIndex: vin.Vout,
			Addr:      *vin.GetPukToAddr(),
		})

	}

	for voutIndex, vout := range *txItr.GetVout() {

		ok := bytes.Equal(vout.Address, CountAddrCoin)

		if !ok {
			continue
		}

		if txItr.Class() == config.Wallet_tx_type_deposit_in || txItr.Class() == config.Wallet_tx_type_vote_in || txItr.Class() == config.Wallet_tx_type_account {
			if voutIndex == 0 {
				continue
			}
		}

		txItem := mining.TxItem{
			Addr: &(*txItr.GetVout())[voutIndex].Address,

			Value: vout.Value,
			Txid:  *txItr.GetHash(),

			VoutIndex:    uint64(voutIndex),
			Height:       height,
			LockupHeight: vout.FrozenHeight,
		}

		itemCount.Additems = append(itemCount.Additems, &txItem)

	}

	itemChan <- &itemCount

	<-tokenCPU
}
