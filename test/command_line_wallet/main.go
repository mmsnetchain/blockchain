package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/prestonTao/libp2parea/engine"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc"
	"os"
	"path/filepath"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	scan()
}

func scan() {
	engine.Log.Info("start init leveldb")
	db.InitDB(filepath.Join("wallet", "data"), filepath.Join("wallet", "temp"))
	input := bufio.NewScanner(os.Stdin)

	engine.Log.Info("Please type in something:")

	for input.Scan() {
		line := input.Text()
		engine.Log.Info("command: %s", line)
		ok := parseHeight(line)
		if ok {
			continue
		}
		ok = parsehash(line)
		if ok {
			continue
		}

	}
}

func parseHeight(line string) bool {

	height, err := strconv.Atoi(line)
	if err == nil {

		bhash, err := db.LevelDB.Find([]byte(config.BlockHeight + strconv.Itoa(int(height))))
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		bs, err := db.LevelDB.Find(*bhash)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		bh, err := mining.ParseBlockHeadProto(bs)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
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
		*bs, _ = json.Marshal(bhvo)
		engine.Log.Info(string(*bs))
		engine.Log.Info("finish!")
		return true
	}
	return false
}

func parsehash(line string) bool {
	hash, err := hex.DecodeString(line)
	if err != nil {
		engine.Log.Info(err.Error())
		return false
	}

	bs, err := db.LevelDB.Find(hash)
	if err != nil {
		engine.Log.Info(err.Error())
		return false
	}

	bh, err := mining.ParseBlockHeadProto(bs)
	if err != nil {

		txBase, err := mining.ParseTxBaseProto(mining.ParseTxClass(hash), bs)
		if err != nil {
			engine.Log.Info(err.Error())
			return false
		}

		itr := txBase.GetVOJSON()
		bs, _ := json.Marshal(itr)
		engine.Log.Info(string(bs))

		blockHash, err := db.LevelTempDB.Find(config.BuildTxToBlockHash(hash))
		if err == nil {
			engine.Log.Info("blockhash:%s", hex.EncodeToString(*blockHash))
		}

		engine.Log.Info("finish!")
		return true
	}

	if bh.Height <= 0 {

		txBase, err := mining.ParseTxBaseProto(mining.ParseTxClass(hash), bs)
		if err != nil {
			engine.Log.Info(err.Error())
			return false
		}

		itr := txBase.GetVOJSON()
		bs, _ := json.Marshal(itr)
		engine.Log.Info(string(bs))

		blockHash, err := db.LevelTempDB.Find(config.BuildTxToBlockHash(hash))
		if err == nil {
			engine.Log.Info("blockhash:%s", hex.EncodeToString(*blockHash))
		}

		engine.Log.Info("")
		err = txBase.Check()
		if err != nil {
			engine.Log.Info(err.Error())
		}

		engine.Log.Info("finish!")
		return true
	}

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
	*bs, _ = json.Marshal(bhvo)
	engine.Log.Info(string(*bs))

	engine.Log.Info("finish!")
	return true
}
