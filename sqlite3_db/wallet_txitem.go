package sqlite3_db

import (
	"github.com/prestonTao/libp2parea/engine"
	"time"

	_ "github.com/go-xorm/xorm"
)

type TxItem struct {
	Id           int64     `xorm:"pk autoincr unique 'id'"`
	Key          string    `xorm:"index 'key'"`
	Addr         []byte    `xorm:"Blob 'addr'"`
	Value        uint64    `xorm:"uint64 'value'"`
	Txid         []byte    `xorm:"Blob 'txid'"`
	VoutIndex    uint64    `xorm:"uint64 'voutindex'"`
	Height       uint64    `xorm:"uint64 'height'"`
	VoteType     uint16    `xorm:"uint64 'votetype'"`
	FrozenHeight uint64    `xorm:"uint64 'frozenheight'"`
	LockupHeight uint64    `xorm:"uint64 'lockupheight'"`
	Status       int32     `xorm:"int32 'status'"`
	CreateTime   time.Time `xorm:"created 'createtime'"`
}

func (this *TxItem) AddTxItems(txitems *[]TxItem) error {
	if txitems == nil || len(*txitems) <= 0 {
		return nil
	}
	lenght := len(*txitems)
	pageOne := 100
	var err error

	for i := 0; i < lenght/pageOne; i++ {
		items := (*txitems)[i*pageOne : (i+1)*pageOne]
		_, err = engineDB.Insert(&items)
		if err != nil {

			return err
		}
	}
	if lenght%pageOne > 0 {
		i := lenght / pageOne
		items := (*txitems)[i*pageOne : lenght]
		_, err = engineDB.Insert(&items)
		if err != nil {

			return err
		}
	}

	return nil
}

func (this *TxItem) RemoveMoreKey(keys [][]byte) error {
	if keys == nil || len(keys) <= 0 {
		return nil
	}

	for _, one := range keys {
		_, err := engineDB.Where("txid = ?", one).Unscoped().Delete(this)
		if err != nil {
			engine.Log.Error("TxItem RemoveMoreKey error:%s", err.Error())
		}
	}

	return nil
}

func (this *TxItem) UpdateLockHeight(ids *[]int64, lockTime uint64) error {
	if ids == nil || len(*ids) <= 0 {
		return nil
	}
	txitem := TxItem{LockupHeight: lockTime}

	_, err := engineDB.In("id", *ids).Update(&txitem)

	return err
}

func (this *TxItem) FindNotSpentSum(height uint64) (uint64, error) {
	ss := new(TxItem)

	total, err := engineDB.Where("frozenheight < ? and lockupheight < ?", height, height).SumInt(ss, "value")

	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (this *TxItem) FindLockHeightSum(lockTime uint64) (uint64, error) {
	ss := new(TxItem)

	total, err := engineDB.Where("lockupheight >?", lockTime).SumInt(ss, "value")

	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (this *TxItem) FindFrozenHeightSum(frozenheight uint64) (uint64, error) {
	ss := new(TxItem)

	total, err := engineDB.Where("frozenheight >?", frozenheight).SumInt(ss, "value")

	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (this *TxItem) FindNotSpentTxitem(height uint64, amount uint64, limit int) (*[]TxItem, error) {
	tis := make([]TxItem, 0)

	err := engineDB.Where("frozenheight < ? and lockupheight < ? and value >= ?", height, height, amount).Limit(1, 0).Find(&tis)

	if err != nil {
		return &tis, err
	}
	if len(tis) > 0 {
		return &tis, nil
	}

	err = engineDB.Where("frozenheight < ? and lockupheight < ? and value < ?", height, height, amount).Limit(limit, 0).Find(&tis)

	if err != nil {
		return &tis, err
	}
	return &tis, nil
}

func (this *TxItem) FindNotSpentTxitemByAddr(addr []byte, height uint64, amount uint64, limit int) (*[]TxItem, error) {
	tis := make([]TxItem, 0)

	err := engineDB.Where("addr = ? and frozenheight < ? and lockupheight < ? and value >= ?", addr, height, height, amount).Limit(1, 0).Find(&tis)

	if err != nil {
		return &tis, err
	}
	if len(tis) > 0 {
		return &tis, nil
	}

	err = engineDB.Where("addr = ? and frozenheight < ? and lockupheight < ? and value < ?", addr, height, height, amount).Limit(limit, 0).Find(&tis)

	if err != nil {
		return &tis, err
	}
	return &tis, nil
}
