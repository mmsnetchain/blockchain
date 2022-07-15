package sqlite3_db

import (
	_ "github.com/go-xorm/xorm"
)

type PeerInfo struct {
	Id         string `xorm:"varchar(25) pk notnull unique 'id'"`
	SuperId    string `xorm:"varchar(25) 'sid'"`
	Puk        string `xorm:"varchar(25) 'puk'"`
	SK         string `xorm:"varchar(25) 'sk'"`
	SharedHka  string `xorm:"varchar(25) 'ska'"`
	SharedNhkb string `xorm:"varchar(25) 'skb'"`
}

func (this *PeerInfo) Add() error {
	_, err := engineDB.Insert(this)
	return err
}

func (this *PeerInfo) FindByid(id string) (*PeerInfo, error) {
	fs := make([]PeerInfo, 0)
	err := engineDB.Where("id = ?", id).Find(&fs)
	if err != nil {
		return nil, err
	}
	if len(fs) <= 0 {
		return nil, nil
	}
	return &fs[0], nil
}

func (this *PeerInfo) Update() error {
	_, err := engineDB.Id(this.Id).Update(this)
	return err
}
