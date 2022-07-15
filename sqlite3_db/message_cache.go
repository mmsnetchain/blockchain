package sqlite3_db

import (
	"errors"

	_ "github.com/go-xorm/xorm"
)

type MessageCache struct {
	Id         int64  `xorm:"pk autoincr unique 'id'"`
	Hash       []byte `xorm:"Blob index unique 'hash'"`
	Head       []byte `xorm:"Blob 'head'"`
	Body       []byte `xorm:"Blob 'body'"`
	CreateTime int64  `xorm:"created 'createtime'"`
}

func (this *MessageCache) Add(hash []byte, head, body []byte) error {
	dblock.Lock()
	mc := MessageCache{
		Hash: hash,
		Head: head,
		Body: body,
	}
	_, err := engineDB.Insert(&mc)
	dblock.Unlock()
	return err
}

func (this *MessageCache) FindByHash(hash []byte) (*MessageCache, error) {
	dblock.Lock()
	defer dblock.Unlock()
	mls := make([]MessageCache, 0)

	err := engineDB.Where("hash=?", hash).Find(&mls)
	if err != nil {
		return nil, err
	}
	if len(mls) <= 0 {
		return nil, errors.New("not find")
	}
	return &mls[0], nil
}

func (this *MessageCache) Remove(createtime int64) error {
	if engineDB == nil {
		return nil
	}
	dblock.Lock()
	_, err := engineDB.Where("createtime < ?", createtime).Unscoped().Delete(this)
	dblock.Unlock()
	return err
}
