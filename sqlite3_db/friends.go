package sqlite3_db

import (
	_ "github.com/go-xorm/xorm"
)

type Friends struct {
	Id           uint64 `xorm:"pk autoincr unique 'id'"`
	NodeId       string `xorm:"varchar(25) 'nodeid'"`
	Nickname     string `xorm:"varchar(25) 'nick_name'"`
	Notename     string `xorm:"varchar(25) 'note_name'"`
	Basecoinaddr string `xorm:"varchar(100) 'basecoinaddr'"`
	Note         string `xorm:"varchar(25) 'note'"`
	Status       int    `xorm:"int 'status'"`
	IsAdd        int    `xorm:"int 'isadd'"`
	Hello        string `xorm:"varchar(25) 'hello'"`
	Read         int    `xorm:"int 'read'"`
}

func (this *Friends) Add(f *Friends) error {
	_, err := table_friends.Insert(f)
	return err
}

func (this *Friends) Del(id string) error {
	_, err := engineDB.Where("nodeid = ?", id).Unscoped().Delete(this)
	return err
}

func (this *Friends) Update() error {
	_, err := engineDB.Where("nodeid = ?", this.NodeId).Update(this)
	return err
}

func (this *Friends) UpdateNoteName() error {
	_, err := engineDB.Nullable("note_name").Where("nodeid = ?", this.NodeId).Update(this)
	return err
}
func (this *Friends) Getall() ([]Friends, error) {

	fs := make([]Friends, 0)
	err := engineDB.Find(&fs)
	return fs, err
}

func (this *Friends) FindById(id string) (*Friends, error) {
	fs := make([]Friends, 0)
	err := engineDB.Where("nodeid = ?", id).Find(&fs)
	if err != nil {
		return nil, err
	}
	if len(fs) <= 0 {
		return nil, nil
	}
	return &fs[0], nil
}
