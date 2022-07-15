package db

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

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

func Save(id, bs []byte) error {
	err := db.Put(id, bs, nil)
	if err != nil {
		fmt.Println("Leveldb save error", err)
	}
	return err
}

func Find(id []byte) ([]byte, error) {
	value, err := db.Get(id, nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func Remove(id []byte) error {
	return db.Delete(id, nil)
}

func CheckHashExist(hash []byte) bool {
	_, err := Find(hash)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false
		}
		return true
	}
	return true
}
