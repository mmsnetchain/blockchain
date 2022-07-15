package sqlite3_db

import (
	"bytes"
	"fmt"

	_ "github.com/go-xorm/xorm"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Property struct {
	Hash       string `xorm:"varchar(25) pk notnull unique 'hash'"`
	Nickname   string `xorm:"varchar(25) 'nick_name'"`
	CreateTime uint64 `xorm:"created 'createtime'"`
	UpdateTime uint64 `xorm:"updated 'updated'"`
}

func (this *Property) Json() []byte {
	rs, err := json.Marshal(this)
	if err != nil {
		fmt.Println(err)
	}
	return rs
}

func ParseProperty(bs []byte) Property {
	p := Property{}

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	return p
}
func (this *Property) Update() error {
	ok, err := engineDB.Where("hash=?", this.Hash).Exist(&Property{})
	if err != nil {
		return err
	}
	if ok {
		_, err = engineDB.Where("hash=?", this.Hash).Update(this)
	} else {
		_, err = engineDB.Insert(this)
	}

	return err
}
func (this *Property) Get(hash string) (dp Property) {
	_, err := engineDB.Where("hash = ?", hash).Get(&dp)
	if err != nil {
		fmt.Println(err)
	}
	return
}
