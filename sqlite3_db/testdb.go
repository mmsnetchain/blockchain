package sqlite3_db

import (
	"github.com/prestonTao/libp2parea/engine"
	"time"

	_ "github.com/go-xorm/xorm"
)

type TestDB struct {
	Id         int64     `xorm:"pk autoincr unique 'id'"`
	Txid       []byte    `xorm:"Blob 'txid'"`
	CreateTime time.Time `xorm:"created 'createtime'"`
}

func (this *TestDB) AddTxItems(txitems *[]TestDB) error {
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

func (this *TestDB) RemoveMoreKey(keys [][]byte) error {
	if keys == nil || len(keys) <= 0 {
		return nil
	}

	tis := make([]TestDB, 0)

	err := engineDB.In("txid = ?", keys).Find(&tis)
	engine.Log.Info(":%d", len(tis))

	for i, _ := range keys {
		n, err := engineDB.Where("txid = ?", keys[i]).Unscoped().Delete(this)
		engine.Log.Info(":%d error:%v", n, err)

	}

	n, err := engineDB.In("txid = ?", keys).Unscoped().Delete(this)
	engine.Log.Info(":%d", n)

	return err
}
