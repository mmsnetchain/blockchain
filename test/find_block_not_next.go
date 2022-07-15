package main

import (
	"github.com/prestonTao/libp2parea/engine"
	"mmschainnewaccount/chain_witness_vote/mining"
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	path := filepath.Join("wallet", "data")
	InitDB(path)
	FindNotNext()
}

var Once_ConnLevelDB sync.Once
var db *leveldb.DB

func InitDB(name string) (err error) {
	Once_ConnLevelDB.Do(func() {

		db, err = leveldb.OpenFile(name, nil)

		if err != nil {
			return
		}
		return
	})
	return
}

func FindNotNext() {

	iter := db.NewIterator(nil, nil)
	for iter.Next() {

		valueBs := iter.Value()
		bh, err := mining.ParseBlockHead(&valueBs)
		if err != nil {

			continue
		}

		if bh.Nextblockhash == nil {
			bs, err := bh.Json()
			if err != nil {
				engine.Log.Info("2 " + err.Error())
				continue
			}
			engine.Log.Info(string(*bs))
		}

	}
	iter.Release()
	err := iter.Error()
	engine.Log.Info("3 " + err.Error())
}
