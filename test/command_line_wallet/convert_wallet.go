package main

import (
	"encoding/hex"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/mining"
	_ "mmschainnewaccount/chain_witness_vote/mining/token/payment"
	_ "mmschainnewaccount/chain_witness_vote/mining/token/publish"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	_ "mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc"
)

func main() {

	findNextBlock()
	engine.Log.Info("finish!")
}

func findNextBlock() {
	oldWallet := utils.CreateLevelDB("wallet/data")
	newWallet := utils.CreateLevelDB("newwallet")

	beforBlockHash, err := oldWallet.Find(config.Key_block_start)
	if err != nil {
		engine.Log.Info("111 id " + err.Error())
		return
	}
	newWallet.Save(config.Key_block_start, beforBlockHash)

	beforGroupHeight := uint64(0)

	for beforBlockHash != nil {
		bs, err := oldWallet.Find(beforBlockHash)
		if err != nil {
			engine.Log.Info(" " + err.Error())
			return
		}
		bh, err := mining.ParseBlockHead(&bs)
		if err != nil {

			engine.Log.Info(" " + err.Error())

			engine.Log.Info(string(bs))
			return
		}
		tempBs, _ := bh.Proto()
		newWallet.Save(beforBlockHash, tempBs)
		if bh.Nextblockhash == nil {
			engine.Log.Info("%d -----------------------------------\n%s\n", bh.Height)
		} else {
			engine.Log.Info("%d -----------------------------------\n%s\n", bh.Height,
				hex.EncodeToString(bh.Hash))

		}
		engine.Log.Info(" %d", len(bh.Tx))

		txs := make([]string, 0)
		for _, one := range bh.Tx {
			txs = append(txs, hex.EncodeToString(one))
		}
		bhvo := rpc.BlockHeadVO{
			Hash:              hex.EncodeToString(bh.Hash),
			Height:            bh.Height,
			GroupHeight:       bh.GroupHeight,
			GroupHeightGrowth: bh.GroupHeightGrowth,
			Previousblockhash: hex.EncodeToString(bh.Previousblockhash),
			Nextblockhash:     hex.EncodeToString(bh.Nextblockhash),
			NTx:               bh.NTx,
			MerkleRoot:        hex.EncodeToString(bh.MerkleRoot),
			Tx:                txs,
			Time:              bh.Time,
			Witness:           bh.Witness.B58String(),
			Sign:              hex.EncodeToString(bh.Sign),
		}
		bs, _ = json.Marshal(bhvo)
		engine.Log.Info(string(bs))

		intervalGroup := bhvo.GroupHeight - beforGroupHeight
		if intervalGroup > 1 {
			engine.Log.Warn(" %d", intervalGroup)
		}
		beforGroupHeight = bhvo.GroupHeight

		for _, one := range bh.Tx {
			tx, err := oldWallet.Find(one)
			if err != nil {
				engine.Log.Info(" %d "+err.Error(), bh.Height)
				panic("error:")
				return
			}
			txBase, err := mining.ParseTxBase(mining.ParseTxClass(one), tx)
			if err != nil {
				engine.Log.Info(" %d "+err.Error(), bh.Height)

				panic("error:")
				return
			}

			tempBs, _ := txBase.Proto()
			newWallet.Save(one, tempBs)

			if txBase.Class() == config.Wallet_tx_type_mining {
				rewardTotal := uint64(0)
				for _, one := range *txBase.GetVout() {
					rewardTotal += one.Value
				}
				engine.Log.Info(" %d", rewardTotal)
			}

		}

		engine.Log.Info("hash %s \n", hex.EncodeToString(bh.Nextblockhash))

		if bh.Nextblockhash != nil {
			beforBlockHash = bh.Nextblockhash
		} else {
			beforBlockHash = nil
		}
	}
}
