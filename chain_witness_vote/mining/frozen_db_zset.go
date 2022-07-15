package mining

import (
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/utils"
)

func AddFrozenHeight(addr *crypto.AddressCoin, height uint64, value uint64) error {
	if height < config.Wallet_frozen_time_min {
		return AddZSetFrozenHeightForHeight(addr, height, value)
	} else {
		return AddZSetFrozenHeightForTime(addr, int64(height), value)
	}
}

func GetFrozenHeight(height uint64, time int64) (*[]*TxItem, error) {
	txItems, err := GetZSetFrozenHeightForHeight(height)
	if err != nil {
		return nil, err
	}
	txItemsTime, err := GetZSetFrozenHeightForTime(time)
	if err != nil {
		return nil, err
	}

	*txItems = append(*txItems, *txItemsTime...)
	return txItems, nil
}

func RemoveFrozenHeight(height uint64, time int64) error {
	sps, err := db.LevelTempDB.GetZSetPage(&config.DBKEY_zset_frozen_time, 0, time, int(time))
	if err != nil {
		return err
	}
	for i := 0; i < len(*sps); i++ {
		one := (*sps)[i]
		childKey := append(config.DBKEY_zset_frozen_time_children, utils.Uint64ToBytes(uint64(one.Score))...)
		err = db.LevelTempDB.DelZSetAll(&childKey)

		if err != nil {
			return err
		}
	}
	err = db.LevelTempDB.DelZSet(&config.DBKEY_zset_frozen_time, 0, time)
	if err != nil {
		return err
	}

	sps, err = db.LevelTempDB.GetZSetPage(&config.DBKEY_zset_frozen_height, 0, int64(height), int(height))
	if err != nil {
		return err
	}
	for i := 0; i < len(*sps); i++ {
		one := (*sps)[i]
		childKey := append(config.DBKEY_zset_frozen_height_children, utils.Uint64ToBytes(uint64(one.Score))...)
		err = db.LevelTempDB.DelZSetAll(&childKey)

		if err != nil {
			return err
		}
	}
	err = db.LevelTempDB.DelZSet(&config.DBKEY_zset_frozen_height, 0, int64(height))
	if err != nil {
		return err
	}
	return nil
}

func AddZSetFrozenHeightForHeight(addr *crypto.AddressCoin, height uint64, value uint64) error {
	txItemBs := UtilsSerializeTxItem(addr, value)
	childKey := append(config.DBKEY_zset_frozen_height_children, utils.Uint64ToBytes(height)...)

	err := db.LevelTempDB.AddZSetAutoincrId(&childKey, txItemBs)
	if err != nil {
		return err
	}

	err = db.LevelTempDB.AddZSet(&config.DBKEY_zset_frozen_height, &childKey, int64(height))
	return err
}

func GetZSetFrozenHeightForHeight(height uint64) (*[]*TxItem, error) {
	sps, err := db.LevelTempDB.GetZSetPage(&config.DBKEY_zset_frozen_height, 0, int64(height), int(height))
	if err != nil {
		return nil, err
	}

	addrMap := make(map[string]uint64)
	for i := 0; i < len(*sps); i++ {
		one := (*sps)[i]

		childKey := append(config.DBKEY_zset_frozen_height_children, utils.Uint64ToBytes(uint64(one.Score))...)
		spsChild, err := db.LevelTempDB.GetZSetAll(&childKey)
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(*spsChild); j++ {
			childOne := (*spsChild)[j]
			addr, value := UtilsParseTxItem(&childOne.Member)
			oldvalue, _ := addrMap[utils.Bytes2string(*addr)]
			oldvalue += value
			addrMap[utils.Bytes2string(*addr)] = oldvalue
		}
	}
	txItems := make([]*TxItem, 0)
	for addrStr, value := range addrMap {
		addr := crypto.AddressCoin([]byte(addrStr))

		item := TxItem{
			Addr:  &addr,
			Value: value,
		}
		txItems = append(txItems, &item)
	}

	return &txItems, nil
}

func AddZSetFrozenHeightForTime(addr *crypto.AddressCoin, height int64, value uint64) error {
	txItemBs := UtilsSerializeTxItem(addr, value)
	childKey := append(config.DBKEY_zset_frozen_time_children, utils.Uint64ToBytes(uint64(height))...)
	err := db.LevelTempDB.AddZSetAutoincrId(&childKey, txItemBs)
	if err != nil {
		return err
	}

	err = db.LevelTempDB.AddZSet(&config.DBKEY_zset_frozen_time, &childKey, int64(height))
	return err
}

func GetZSetFrozenHeightForTime(height int64) (*[]*TxItem, error) {
	sps, err := db.LevelTempDB.GetZSetPage(&config.DBKEY_zset_frozen_time, 0, int64(height), int(height))
	if err != nil {
		return nil, err
	}
	addrMap := make(map[string]uint64)
	for i := 0; i < len(*sps); i++ {
		one := (*sps)[i]
		childKey := append(config.DBKEY_zset_frozen_time_children, utils.Uint64ToBytes(uint64(one.Score))...)
		spsChild, err := db.LevelTempDB.GetZSetAll(&childKey)
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(*spsChild); j++ {
			childOne := (*spsChild)[j]
			addr, value := UtilsParseTxItem(&childOne.Member)
			oldvalue, _ := addrMap[utils.Bytes2string(*addr)]
			oldvalue += value
			addrMap[utils.Bytes2string(*addr)] = oldvalue
		}
	}
	txItems := make([]*TxItem, 0)
	for addrStr, value := range addrMap {
		addr := crypto.AddressCoin([]byte(addrStr))
		item := TxItem{
			Addr:  &addr,
			Value: value,
		}
		txItems = append(txItems, &item)
	}
	return &txItems, nil
}

func UtilsSerializeTxItem(addr *crypto.AddressCoin, value uint64) *[]byte {
	bs := make([]byte, 0, len(*addr)+8)
	bs = append(bs, utils.Uint64ToBytes(value)...)
	bs = append(bs, *addr...)
	return &bs
}

func UtilsParseTxItem(bs *[]byte) (*crypto.AddressCoin, uint64) {
	value := utils.BytesToUint64((*bs)[:8])
	addr := crypto.AddressCoin((*bs)[8:])
	return &addr, value
}
