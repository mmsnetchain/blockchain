package main

import (
	"encoding/hex"
	"fmt"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"path/filepath"
)

func main() {

	path := filepath.Join("../example/peer_root/wallet", "data")

	findNextBlock(path)
	fmt.Println("finish!")
}

var tempBlockHeight = uint64(1000)

func findSomeBlock(dir string) {
	nums := []uint64{}
	for i := uint64(1); i < tempBlockHeight; i++ {
		nums = append(nums, i)
	}

	db.InitDB(dir)
	beforBlockHash, err := db.Find(config.Key_block_start)
	if err != nil {
		fmt.Println("111 id", err)
		return
	}
	maxBlock := uint64(0)
	for _, one := range nums {
		if one > maxBlock {
			maxBlock = one
		}
	}

	for i := uint64(1); i <= maxBlock; i++ {
		bs, err := db.Find(*beforBlockHash)
		if err != nil {
			fmt.Println("", i, "", err)
			return
		}
		bh, err := mining.ParseBlockHead(bs)
		if err != nil {
			fmt.Println("", i, "", err)
			return
		}
		beforBlockHash = &bh.Nextblockhash
		isPrint := false
		for _, one := range nums {
			if one == i {
				isPrint = true
				break
			}
		}
		if isPrint {
			fmt.Println("", i, " ----------------------------------\n",
				hex.EncodeToString(bh.Hash), "\n", string(*bs), "\n")

			fmt.Println("hash", hex.EncodeToString(bh.Nextblockhash))

			for _, one := range bh.Tx {
				tx, err := db.Find(one)
				if err != nil {
					fmt.Println("", i, "", err)
					return
				}
				txBase, err := mining.ParseTxBase(tx)
				if err != nil {
					fmt.Println("", i, "", err)
					return
				}

				txid := txBase.GetHash()

				fmt.Println(string(hex.EncodeToString(*txid)), "\n", string(*tx), "\n")
			}
		}
	}
}

func findNextBlock(dir string) {

	db.InitDB(dir)
	beforBlockHash, err := db.Find(config.Key_block_start)
	if err != nil {
		fmt.Println("111 id", err)
		return
	}

	for beforBlockHash != nil {
		bs, err := db.Find(*beforBlockHash)
		if err != nil {
			fmt.Println("", "", err)
			return
		}
		bh, err := mining.ParseBlockHead(bs)
		if err != nil {

			fmt.Println("", "", err)
			fmt.Println(string(*bs))
			return
		}
		if bh.Nextblockhash == nil {
			fmt.Println("", bh.Height, " -----------------------------------\n",
				hex.EncodeToString(bh.Hash), "\n", string(*bs), "\nnext")
		} else {
			fmt.Println("", bh.Height, " -----------------------------------\n",
				hex.EncodeToString(bh.Hash), "\n", string(*bs), "\nnext", len(bh.Nextblockhash))
		}

		fmt.Println("hash", hex.EncodeToString(bh.Nextblockhash))
		for _, one := range bh.Tx {
			tx, err := db.Find(one)
			if err != nil {
				fmt.Println("", bh.Height, "", err)
				panic("error:")
				return
			}
			txBase, err := mining.ParseTxBase(tx)
			if err != nil {
				fmt.Println("", bh.Height, "", err)
				panic("error:")
				return
			}
			fmt.Println("##########start#############")
			color.Green.Printf("txBase:%+v\n", txBase)
			vt := txBase.GetVout()
			for voutIndex, vout := range *vt {
				txItem := mining.TxItem{
					Addr:     &vout.Address,
					Value:    vout.Value,
					Txid:     *txBase.GetHash(),
					OutIndex: uint64(voutIndex),
					Height:   bh.Height,
				}
				color.Red.Printf("txItem:%+v\n", txItem)
			}
			fmt.Println("##########end#############")

			txid := txBase.GetHash()

			fmt.Println(string(hex.EncodeToString(*txid)), "\n", string(*tx), "\n")
		}

		if bh.Nextblockhash != nil {
			beforBlockHash = &bh.Nextblockhash
		} else {
			beforBlockHash = nil
		}
	}
}
