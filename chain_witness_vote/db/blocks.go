package db

import (
	"github.com/prestonTao/utils"
	"mmschainnewaccount/config"
)

func SaveHighstBlock(number uint64) error {
	bs := utils.Uint64ToBytes(number)
	return LevelDB.Save([]byte(config.Block_Highest), &bs)
}

func GetHighstBlock() uint64 {
	bs, err := LevelDB.Find([]byte(config.Block_Highest))
	if err != nil {
		return 0
	}
	return utils.BytesToUint64(*bs)
}

func SaveTxToBlockHash(txid, blockhash *[]byte) error {
	key := config.BuildTxToBlockHash(*txid)

	return LevelTempDB.Save(key, blockhash)
}

func GetTxToBlockHash(txid *[]byte) (*[]byte, error) {

	key := config.BuildTxToBlockHash(*txid)
	value, err := LevelTempDB.Find(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}
