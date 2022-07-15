package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/store/fs"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"

	_ "github.com/logoove/sqlite"
)

const dbpath = "D:/test/hzzfiles/sqlkey10table.db"

var engineDB *xorm.Engine

func main() {
	fmt.Println("start")
	s := fs.NewStorage(dbpath)

	time.Sleep(time.Second * 10)

	start := time.Now()

	s.Find("ff226c33120b3ea63659303ec89055fc332700b09d33d72d70977a194e61d773")
	fmt.Println(time.Now().Sub(start))
	fmt.Println("end")
	return

	fmt.Println("end")
}

func getRandBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

type Key1 struct {
	Key
}
type Key2 struct {
	Key
}
type Key3 struct {
	Key
}
type Key4 struct {
	Key
}
type Key5 struct {
	Key
}
type Key6 struct {
	Key
}
type Key7 struct {
	Key
}
type Key8 struct {
	Key
}
type Key9 struct {
	Key
}
type Key10 struct {
	Key
}

type KeyInterface interface {
	Set(key string, value []byte)
}

type Key struct {
	Id     uint64 `xorm:"pk autoincr unique 'id'"`
	Key    string `xorm:"varchar(25) 'key'"`
	Value  []byte `xorm:"Blob 'value'"`
	Status int    `xorm:"int 'status'"`
}

func (this *Key) Set(key string, value []byte) {
	this.Key = key
	this.Value = value
	this.Status = 1
}

type Storage struct {
	sqldb    *xorm.Engine
	DbPath   string
	SpaceNum int
	PerSpace int
	TableNum int
}

func (this *Storage) Find(key string) *[]byte {
	for i := 0; i < this.TableNum; i++ {
		sqlStr := "select value from key" + strconv.Itoa(i) + " where key=?"
		result, err := this.sqldb.Query(sqlStr, key)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if len(result) <= 0 {
			continue
		}

		bs := result[0]["value"]
		return &bs
	}
	return nil
}

func (this *Storage) FullTable() {

	for i := 0; i < this.TableNum; i++ {
		this.fullTableOne(i)
	}

}

func (this *Storage) fullTableOne(n int) {

	sqlStr := "INSERT INTO `key" + strconv.Itoa(n) + "` (`key`,`value`) VALUES (?,?)"

	num := this.SpaceNum / this.TableNum / this.PerSpace
	fmt.Println("init data space ...", num)
	for i := 0; i < num; i++ {
		b := make([]byte, this.PerSpace)
		if _, err := rand.Read(b); err != nil {
			fmt.Println(err)
			return
		}

		key := utils.Hash_SHA3_256(b)
		keyStr := hex.EncodeToString(key)

		this.sqldb.Exec(sqlStr, keyStr, b)

	}
}

func NewStorage(abspath string, index uint64) *Storage {
	var err error
	engineDB, err = xorm.NewEngine("sqlite3", "file:"+abspath+"?cache=shared")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	engineDB.ShowSQL(config.ShowSQL)

	s := &Storage{
		sqldb:    engineDB,
		DbPath:   abspath,
		Index:    index,
		SpaceNum: 1024 * 1024 * 1024 * 10,
		PerSpace: 512 * 1024,
		TableNum: 10,
	}

	for i := 0; i < s.TableNum; i++ {

		sqlStr := "CREATE TABLE IF NOT EXISTS `key" + strconv.Itoa(i) + "` (`key` TEXT NULL, `value` BOLB NULL)"

		engineDB.Exec(sqlStr)
	}
	return s

}
