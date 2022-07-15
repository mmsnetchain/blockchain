package main

import (
	"encoding/hex"
	"fmt"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {

	BuildFirstBlock()
}

func BuildMerkleRoot() {
	ids := [][]byte{
		[]byte("W1aLWC4unTJZhSFc4VNLFsazAJ1PyTocV7agmteQDL3J3N"),
		[]byte("W1gfVGa52yUJ4Gws4TiA9YbwGP8qCGgaYeeT8APjSiNk6U"),
		[]byte("W1j9RJ1xYHaoAuRk2HGBrVA82njoxFAoctYKQMH43k8hXu"),
		[]byte("W1atFt7bJ5Ubk4MXuV5GfsEYE7srWXR51exDgUEJcVr5fZ"),
		[]byte("W1n9XtbLAjRsh9sr2kbwfkfy3VGenyhazbHJwrEYsnDZ8M"),
	}

	root := mining.BuildMerkleRoot(ids)

	fmt.Println(root)
}

func BuildFirstBlock() {

	witness := make([]*keystore.KeyStore, 0)

	seed1 := keystore.NewKeyStore("wallet1", config.Wallet_seed)
	if n, _ := seed1.Load(); n <= 0 {
		seed1.NewLoad("wallet1_seed", "123456")
	}
	witness = append(witness, seed1)

	seed2 := keystore.NewKeyStore("wallet2", config.Wallet_seed)
	if n, _ := seed2.Load(); n <= 0 {
		seed2.NewLoad("wallet2_seed", "123456")
	}
	witness = append(witness, seed2)

	seed3 := keystore.NewKeyStore("wallet3", config.Wallet_seed)
	if n, _ := seed3.Load(); n <= 0 {
		seed3.NewLoad("wallet3_seed", "123456")
	}
	witness = append(witness, seed3)
	fmt.Println("-------------")

	db.InitDB("data")

	balanceTotal := uint64(3 * mining.Unit)

	now := time.Now().Unix()

	txs := make([][]byte, 0)

	txReward := make([]mining.Tx_reward, 0)
	for _, one := range witness {
		vout := mining.Vout{
			Value:   balanceTotal,
			Address: *one.GetAddr()[0].Hash,
		}
		base := mining.TxBase{
			Type: config.Wallet_tx_type_mining,
		}
		tx := mining.Tx_reward{
			TxBase:     base,
			CreateTime: now,
		}
		tx.TxBase.Vout = append(tx.TxBase.Vout, vout)
		tx.BuildHash()
		txs = append(txs, *tx.GetHash())
		txReward = append(txReward, tx)
	}

	deposits := make([]mining.Tx_deposit_in, 0)

	vouts := make([]mining.Vout, 0)
	vout := mining.Vout{
		Value:   config.Mining_deposit,
		Address: *witness[0].GetAddr()[0].Hash,
	}
	vouts = append(vouts, vout)

	vout = mining.Vout{
		Value:   balanceTotal - config.Mining_deposit,
		Address: *witness[0].GetAddr()[0].Hash,
	}
	vouts = append(vouts, vout)
	bs, err := json.Marshal(vouts)
	if err != nil {
		fmt.Println("json", err)
		return
	}
	voutSign, err := witness[0].GetAddr()[0].Sign(bs, "123456")
	if err != nil {
		fmt.Println("111", err)
		return
	}
	sign, err := txReward[0].Sign(&witness[0].GetAddr()[0], "123456")
	if err != nil {
		fmt.Println("222", err)
		return
	}
	vins := make([]mining.Vin, 0)
	vin := mining.Vin{
		Txid:     *txReward[0].GetHash(),
		Vout:     0,
		Puk:      witness[0].GetAddr()[0].GetPubKey(),
		Sign:     *sign,
		VoutSign: *voutSign,
	}
	vins = append(vins, vin)
	base := mining.TxBase{
		Type:       config.Wallet_tx_type_deposit_in,
		Vin_total:  uint64(len(vins)),
		Vin:        vins,
		Vout_total: uint64(len(vouts)),
		Vout:       vouts,
	}
	tx := mining.Tx_deposit_in{
		TxBase: base,
	}

	tx.BuildHash()
	txs = append(txs, *tx.GetHash())
	deposits = append(deposits, tx)

	for i, one := range witness {
		key := one.GetAddr()[0]
		newvouts := make([]mining.Vout, 0)
		vout := mining.Vout{
			Value:   config.Mining_deposit,
			Address: *one.GetAddr()[0].Hash,
		}
		newvouts = append(newvouts, vout)
		if i == 0 {

			vout = mining.Vout{
				Value:   1 * mining.Unit,
				Address: *one.GetAddr()[0].Hash,
			}
			newvouts = append(newvouts, vout)
		} else {

			vout = mining.Vout{
				Value:   2 * mining.Unit,
				Address: *one.GetAddr()[0].Hash,
			}
			newvouts = append(newvouts, vout)
		}
		bs, err := json.Marshal(newvouts)
		if err != nil {
			fmt.Println("json", err)
			return
		}
		voutSign, err := one.GetAddr()[0].Sign(bs, "123456")
		if err != nil {
			fmt.Println("111", err)
			return
		}

		vins := make([]mining.Vin, 0)
		var sign *[]byte
		if i == 0 {
			sign, err = deposits[0].Sign(&one.GetAddr()[0], "123456")
			if err != nil {
				fmt.Println("222", err)
				return
			}
			vin := mining.Vin{
				Txid:     *deposits[0].GetHash(),
				Vout:     1,
				Puk:      one.GetAddr()[0].GetPubKey(),
				Sign:     *sign,
				VoutSign: *voutSign,
			}
			vins = append(vins, vin)
		} else {
			sign, err = txReward[i].Sign(&key, "123456")
			if err != nil {
				fmt.Println("222", err)
				return
			}
			vin := mining.Vin{
				Txid:     *txReward[i].GetHash(),
				Vout:     0,
				Puk:      key.GetPubKey(),
				Sign:     *sign,
				VoutSign: *voutSign,
			}
			vins = append(vins, vin)
		}

		base := mining.TxBase{
			Type:       config.Wallet_tx_type_deposit_in,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(newvouts)),
			Vout:       newvouts,
		}
		tx := mining.Tx_deposit_in{
			TxBase: base,
		}

		tx.BuildHash()
		txs = append(txs, *tx.GetHash())
		deposits = append(deposits, tx)
	}

	backupMiner1 := mining.BackupMiners{
		Time:   time.Now().Unix(),
		Miners: make([]mining.BackupMiner, 0),
	}
	miner := mining.BackupMiner{
		Miner: *witness[0].GetAddr()[0].Hash,
		Count: 1,
	}
	backupMiner1.Miners = append(backupMiner1.Miners, miner)

	blockHead1 := mining.BlockHead{

		Height:      1,
		GroupHeight: 1,

		NTx:         uint64(len(txs)),
		Tx:          txs,
		Time:        time.Now().Unix(),
		Witness:     *witness[0].GetAddr()[0].Hash,
		BackupMiner: utils.Hash_SHA3_256(*backupMiner1.JSON()),
	}

	blockHead1.BuildMerkleRoot()
	blockHead1.BuildHash()
	db.Save(blockHead1.BackupMiner, backupMiner1.JSON())

	fmt.Println(hex.EncodeToString(blockHead1.Hash))
	bhbs, _ := blockHead1.Json()

	fmt.Println(string(*bhbs), "\n")
	db.Save(blockHead1.Hash, bhbs)

	db.Save(config.Key_block_start, &blockHead1.Hash)

	for _, one := range txReward {
		one.TxBase.BlockHash = blockHead1.Hash
		bs, err := one.Json()
		if err != nil {
			fmt.Println("2 json", err)
			return
		}
		db.Save(*one.GetHash(), bs)
		fmt.Println(string(*bs), "\n")
	}

	for _, one := range deposits {
		one.TxBase.BlockHash = blockHead1.Hash
		bs, err := one.Json()
		if err != nil {
			fmt.Println("3 json", err)
			return
		}
		db.Save(*one.GetHash(), bs)
		fmt.Println(string(*bs), "\n")
	}
	db.SaveBlockHeight(blockHead1.Height, &blockHead1.Hash)

}
