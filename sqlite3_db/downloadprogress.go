package sqlite3_db

import (
	"fmt"

	_ "github.com/go-xorm/xorm"
)

type Downprogress struct {
	Id         uint64  `xorm:"pk autoincr unique 'id'"`
	Hash       string  `xorm:"varchar(25) unique 'hash'"`
	State      int     `xorm:"'state'"`
	FileInfo   []byte  `xorm:"varchar(255) 'fileinfo'"`
	Rate       float32 `xorm:"'rate'"`
	Speed      uint64  `xorm:"'speed'"`
	UpdateTime uint64  `xorm:"updated 'updated'"`
	CreateTime uint64  `xorm:"created 'created'"`
}

func (this *Downprogress) Add() error {
	ok, err := table_downloadprogress.Where("hash=?", this.Hash).Exist(&Downprogress{})
	if err != nil {
		return err
	}
	if ok {
		_, err = table_downloadprogress.Where("hash=?", this.Hash).Update(this)
	} else {
		_, err = table_downloadprogress.Insert(this)
	}

	return err
}
func (this *Downprogress) Update() error {
	_, err := table_downloadprogress.Where("hash=?", this.Hash).Update(this)
	return err
}
func (this *Downprogress) Delete() error {
	_, err := table_downloadprogress.Where("hash=?", this.Hash).Delete(this)
	return err
}
func (this *Downprogress) Get(hash string) (dp Downprogress) {
	_, err := table_downloadprogress.Where("hash = ?", hash).Get(&dp)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (this *Downprogress) List() (dps []Downprogress) {
	err := table_downloadprogress.Find(&dps)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (this *Downprogress) Listcomplete() (dps []Downprogress) {
	err := table_downloadprogress.Where("Rate=?", 100).Find(&dps)
	if err != nil {
		fmt.Println(err)
	}
	return
}
