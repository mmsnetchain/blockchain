package main

import (
	"bytes"
	"encoding/hex"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/token/payment"
	_ "mmschainnewaccount/chain_witness_vote/mining/token/publish"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

const (
	txItem_status_notSpent = int32(0)
	txItem_status_frozen   = int32(1)
	txItem_status_lock     = int32(2)

	addrOne = "MMS3E1LDivNFVBrqoPMcb3CXM291ikAuPoYJ4"
)

var maxHeight = uint64(0)

var addrSpecial = crypto.AddressFromB58String(addrOne)

func main() {
	Start()
}

func Start() {

	CountAddrBalance(filepath.Join("wallet", "data"))

	time.Sleep(time.Second * 3)
}

func CountAddrBalance(dir string) map[string]*map[string]*mining.TxItem {

	db.InitDB(config.DB_path, config.DB_path_temp)
	sqlite3_db.Init()

	balanceNotSpend := make(map[string]uint64)
	balanceFrozen := make(map[string]uint64)

	beforBlockHash, err := db.LevelDB.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return nil
	}

	maxHeight = db.GetHighstBlock()
	maxHeight = 93888
	engine.Log.Info(":%d", maxHeight)

	rewardTotal := uint64(0)
	balanceCirculationTotal := uint64(0)
	balanceFrozenTotal := uint64(0)

	for beforBlockHash != nil {
		runtime.GC()
		bhvo, err := mining.LoadBlockHeadVOByHash(beforBlockHash)
		if err != nil {
			engine.Log.Info(" :%s", err.Error())
			return nil
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

			if txBase.Class() == config.Wallet_tx_type_mining {
				for _, one := range *txBase.GetVout() {
					rewardTotal += one.Value
				}
				engine.Log.Info(" %d", rewardTotal)
			} else {

				for vinIndex, one := range *txBase.GetVin() {

					one.GetPukToAddr()

					preTxBase, err := mining.LoadTxBase(one.Txid)
					if err != nil {
						panic("error:2")
						return nil
					}
					vout := (*preTxBase.GetVout())[one.Vout]

					balanceCirculationTotal -= vout.Value

					oldFrozenValue, fok := balanceFrozen[utils.Bytes2string(vout.Address)]
					oldNotSpendValue, ok := balanceNotSpend[utils.Bytes2string(vout.Address)]

					if ok {
						if oldNotSpendValue < vout.Value {

							if !fok || oldFrozenValue < vout.Value {
								engine.Log.Info("hash:%s vinIndex:%d addr:%s oldValue:%d value:%d", hex.EncodeToString(*txBase.GetHash()),
									vinIndex, vout.Address.B58String(), oldFrozenValue, vout.Value)

							} else {
								balanceFrozen[utils.Bytes2string(vout.Address)] = oldFrozenValue - vout.Value
							}
						} else {
							balanceNotSpend[utils.Bytes2string(vout.Address)] = oldNotSpendValue - vout.Value
						}
					}
				}
			}

			for voutIndex, _ := range *txBase.GetVout() {
				one := (*txBase.GetVout())[voutIndex]

				if one.Value > 100000000*100000000 {
					panic("hash:" + hex.EncodeToString(*txBase.GetHash()) + " voutIndex:" + strconv.Itoa(voutIndex) + " value:" + strconv.Itoa(int(one.Value)))
				}

				if one.FrozenHeight > maxHeight {
					balanceFrozenTotal += one.Value
					oldValue, ok := balanceFrozen[utils.Bytes2string(one.Address)]
					if ok {
						balanceFrozen[utils.Bytes2string(one.Address)] = oldValue + one.Value
					} else {
						balanceFrozen[utils.Bytes2string(one.Address)] = one.Value
					}
				} else {
					balanceCirculationTotal += one.Value

					oldValue, ok := balanceNotSpend[utils.Bytes2string(one.Address)]
					if ok {
						balanceNotSpend[utils.Bytes2string(one.Address)] = oldValue + one.Value
					} else {
						balanceNotSpend[utils.Bytes2string(one.Address)] = one.Value
					}
				}

			}

		}

		if bhvo.BH.Height == maxHeight {
			break
		}

		if bhvo.BH.Nextblockhash != nil && len(bhvo.BH.Nextblockhash) > 0 {
			beforBlockHash = &bhvo.BH.Nextblockhash
		} else {
			beforBlockHash = nil
		}

	}
	engine.Log.Info(":%d :%d", rewardTotal, balanceCirculationTotal)

	for key, value := range balanceNotSpend {

		frozenValue := balanceFrozen[key]

		addr := crypto.AddressCoin([]byte(key))
		engine.Log.Info(":%s  :%d  :%d", addr.B58String(), value, frozenValue)
	}

	return nil
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
