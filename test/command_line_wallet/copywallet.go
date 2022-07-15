package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {
	Start()

}

func Start() {

	path := filepath.Join("copywallet", "data")
	CopyWallet(path)

	time.Sleep(time.Second * 3)
}

func CopyWallet(dir string) {

	db.InitDB(dir)
	beforBlockHash, err := db.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return nil, nil
	}

	balanceTotal := uint64(0)

	for beforBlockHash != nil {
		bs, err := db.Find(*beforBlockHash)
		if err != nil {
			engine.Log.Info(" " + err.Error())
			return nil, nil
		}
		bh, err := mining.ParseBlockHead(bs)
		if err != nil {

			engine.Log.Info(" " + err.Error())

			engine.Log.Info(string(*bs))
			return nil, nil
		}
		if bh.Nextblockhash == nil {
			engine.Log.Info("%d -----------------------------------\n%s\nnext", bh.Height,
				hex.EncodeToString(bh.Hash))
		} else {
			engine.Log.Info("%d -----------------------------------\n%s\nnext%d", bh.Height,
				hex.EncodeToString(bh.Hash), len(bh.Nextblockhash))
		}

		bhvo := mining.BlockHeadVO{

			BH:  bh,
			Txs: make([]mining.TxItr, 0),
		}

		for _, one := range bh.Tx {
			tx, err := db.Find(one)
			if err != nil {
				engine.Log.Info(" %d "+err.Error(), bh.Height)
				panic("error:")
				return nil, nil
			}
			txBase, err := mining.ParseTxBase(mining.ParseTxClass(one), tx)
			if err != nil {
				engine.Log.Info(" %d :%s %s", bh.Height, hex.EncodeToString(one), err.Error())
				panic("error:")
				return nil, nil
			}

			bhvo.Txs = append(bhvo.Txs, txBase)

			if txBase.Class() == config.Wallet_tx_type_mining {
				rewardTotal := uint64(0)
				for _, one := range *txBase.GetVout() {
					rewardTotal += one.Value
				}
				engine.Log.Info(" %d", rewardTotal)
			}

			if txBase.Class() != config.Wallet_tx_type_mining {
				for _, one := range *txBase.GetVin() {
					addrCoin := crypto.BuildAddr(config.AddrPre, one.Puk)
					addrStr := utils.Bytes2string(addrCoin)
					value, ok := balanceMap[addrStr]
					if !ok {
						panic("" + addrCoin.B58String())
						return nil, nil
					}

					tx, err := db.Find(one.Txid)
					if err != nil {

						panic("error:2")
						return nil, nil
					}
					txBase, err := mining.ParseTxBase(mining.ParseTxClass(one.Txid), tx)
					if err != nil {

						panic("error:2")
						return nil, nil
					}
					vout := (*txBase.GetVout())[one.Vout]

					balanceTotal = balanceTotal - vout.Value
					balanceMap[addrStr] = value - vout.Value

				}
			}

			for _, one := range *txBase.GetVout() {
				balanceTotal = balanceTotal + one.Value
				addrStr := utils.Bytes2string(one.Address)
				value, ok := balanceMap[addrStr]
				if !ok {
					balanceMap[addrStr] = one.Value
				} else {
					balanceMap[addrStr] = value + one.Value
				}

				now := time.Now().Unix()
				if uint64(now) > one.FrozenHeight {
					value, ok := balanceNotSpend[addrStr]
					if !ok {
						balanceNotSpend[addrStr] = one.Value
					} else {
						balanceNotSpend[addrStr] = value + one.Value
					}
				}

			}

		}

		if bh.Nextblockhash != nil {
			beforBlockHash = &bh.Nextblockhash
		} else {
			beforBlockHash = nil
		}

		tic := CountBalances(&bhvo)
		tim.CountTxItem(tic, bhvo.BH.Height, bhvo.BH.Time)
		timOld.CountTxItem(tic, bhvo.BH.Height, bhvo.BH.Time)
		tim.Unfrozen(bhvo.BH.Height, bhvo.BH.Time)
		timOld.Unfrozen(bhvo.BH.Height, bhvo.BH.Time)

		txItems := tim.FindBalance(CountAddrCoin)

		balanceAll := uint64(0)
		for _, one := range txItems {
			balanceAll = balanceAll + one.Value
		}
		if balanceAll == 0 {
			continue
		}

		txItemsOld := timOld.FindBalance(CountAddrCoin)

		balanceAllOld := uint64(0)
		for _, one := range txItemsOld {
			balanceAllOld = balanceAllOld + one.Value
		}

		if balanceAll != balanceAllOld {
			engine.Log.Info(" %d %d", balanceAll, balanceAllOld)
			itemsMap := make(map[string]*mining.TxItem)
			for _, one := range txItems {
				item := one

				key := utils.Bytes2string(one.Txid) + "_" + strconv.Itoa(int(one.VoutIndex))
				itemsMap[key] = one
			}
			itemsOldMap := make(map[string]*mining.TxItem)
			for _, one := range txItemsOld {
				key := utils.Bytes2string(one.Txid) + "_" + strconv.Itoa(int(one.VoutIndex))
				itemsOldMap[key] = one
			}

			for _, one := range txItemsOld {
				key := utils.Bytes2string(one.Txid) + "_" + strconv.Itoa(int(one.VoutIndex))
				_, ok := itemsMap[key]
				if !ok {
					engine.Log.Info("%+v", one)

					bs, _ := db.Find(one.Txid)
					txBase, _ := mining.ParseTxBase(mining.ParseTxClass(one.Txid), bs)
					tx := txBase.GetVOJSON()
					bsTx, _ := json.Marshal(tx)
					engine.Log.Info("%s", string(bsTx))

				}
			}

			panic("")
		}

	}
	engine.Log.Info("：%d", balanceTotal)
	return balanceMap, balanceNotSpend
}

func PrintBalanceMap(balances, balanceNotSpend map[string]uint64) {
	if balances == nil {
		engine.Log.Info("")
		return
	}
	for k, v := range balances {
		addrCoin := crypto.AddressCoin([]byte(k))
		engine.Log.Info("%s %d", addrCoin.B58String(), v)
	}

	engine.Log.Info("")

	if balanceNotSpend == nil {
		engine.Log.Info("")
		return
	}
	for k, v := range balanceNotSpend {
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
		for one := range itemsChan {
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

		itemCount.SubItems = append(itemCount.SubItems, &mining.TxSubItems{
			Txid:      vin.Txid,
			VoutIndex: vin.Vout,
			Addr:      *vin.GetPukToAddr(),
		})

	}

	txCtrl := mining.GetTransactionCtrl(txItr.Class())
	if txCtrl != nil {

		txCtrl.SyncCount()
		itemChan <- &itemCount
		<-tokenCPU
		return
	}

	for voutIndex, vout := range *txItr.GetVout() {

		if txItr.Class() == config.Wallet_tx_type_deposit_in || txItr.Class() == config.Wallet_tx_type_vote_in {
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
