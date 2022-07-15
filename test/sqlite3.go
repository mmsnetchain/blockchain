package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/go-xorm/xorm"

	_ "github.com/logoove/sqlite"
)

var once sync.Once

var db *sql.DB

var engineDB *xorm.Engine

func main() {
	Init()

	example()

}

func Init() {
	var err error
	engineDB, err = xorm.NewEngine("sqlite3", "file:sqlite3.db?cache=shared")
	if err != nil {
		fmt.Println(err)
	}
	engineDB.ShowSQL(true)

	ok, err := engineDB.IsTableExist(FileindexSelf{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(FileindexSelf{}).CreateTable(FileindexSelf{})
	}
}

func example() {
	vidBs := []byte("vid")
	valueBs := []byte("value")
	fileidBs := []byte("fileid")
	table := FileindexSelf{
		Name:   "123",
		Vid:    vidBs,
		FileId: fileidBs,
		Value:  valueBs,
	}
	table.Add(&table)

	one, _ := table.FindByVid(vidBs)
	fmt.Println("", one)
}

type FileindexSelf struct {
	Id     uint64 `xorm:"pk autoincr unique 'id'"`
	Name   string `xorm:"Text 'name'"`
	Vid    []byte `xorm:"Blob 'vid'"`
	FileId []byte `xorm:"Blob 'fileid'"`
	Value  []byte `xorm:"Blob 'value'"`
}

func (this *FileindexSelf) TableName() string {
	return "fileindex_self"
}

func (this *FileindexSelf) Add(f *FileindexSelf) error {
	_, err := engineDB.Insert(f)
	return err
}

func (this *FileindexSelf) Del(fid string) error {
	_, err := engineDB.Where("fileid = ?", fid).Unscoped().Delete(this)
	return err
}

func (this *FileindexSelf) Update() error {
	_, err := engineDB.Where("nodeid = ?", this.FileId).Update(this)
	return err
}

func (this *FileindexSelf) Getall() ([]FileindexSelf, error) {

	fs := make([]FileindexSelf, 0)
	err := engineDB.Find(&fs)
	return fs, err
}

func (this *FileindexSelf) FindByVid(vid []byte) (*FileindexSelf, error) {
	fs := make([]FileindexSelf, 0)
	err := engineDB.Where("vid = ?", vid).Find(&fs)
	if err != nil {
		return nil, err
	}
	if len(fs) <= 0 {
		return nil, nil
	}
	return &fs[0], nil
}
