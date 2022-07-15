package sqlite3_db

import (
	"errors"
	"sync/atomic"
	"time"

	_ "github.com/go-xorm/xorm"
)

const Self = "self"

var generateMsgLogId int64 = 0

func LoadMsgLogGenerateID() error {
	mls := make([]MsgLog, 0)
	err := table_msglog.Desc("id").Limit(1, 0).Find(&mls)
	if err != nil {
		return err
	}
	if len(mls) <= 0 {
		return nil
	}
	generateMsgLogId = mls[0].Id
	return nil
}

func GenerateMsgLogId() int64 {
	return atomic.AddInt64(&generateMsgLogId, 1)
}

type MsgLog struct {
	Id         int64     `xorm:"int64 pk notnull unique 'id'"`
	Sender     string    `xorm:"varchar(25) index(recipient) 'sender'"`
	Recipient  string    `xorm:"varchar(25) 'recipient'"`
	Content    string    `xorm:"varchar(25) 'content'"`
	Read       int       `xorm:"int 'read'"`
	Successful int       `xorm:"int 'successful'"`
	Path       string    `xorm:"varchar(25) 'path'"`
	Class      int       `xorm:"varchar(25) 'class'"`
	PayStatus  int       `xorm:"int 'paystatus'"`
	Size       int64     `xorm:"int64 'size'"`
	Rate       int64     `xorm:"int64 'rate'"`
	Speed      int64     `xorm:"int64 'speed'"`
	Index      int64     `xorm:"int64 'index'"`
	Toid       int64     `xorm:"int64 'toid'"`
	CreateTime time.Time `xorm:"created 'createtime'"`
	UpdateTime uint64    `xorm:"updated 'updated'"`
}

func (this *MsgLog) Add(sender, recipient, content, path string, class int) (int64, error) {

	ml := MsgLog{
		Id:         GenerateMsgLogId(),
		Sender:     sender,
		Recipient:  recipient,
		Content:    content,
		Path:       path,
		Class:      class,
		Read:       1,
		Successful: 1,
	}
	_, err := engineDB.Insert(&ml)
	return ml.Id, err
}

func (this *MsgLog) IsRead(id int64) error {
	ml := MsgLog{
		Read: 2,
	}
	_, err := engineDB.Id(id).Update(&ml)
	return err
}

func (this *MsgLog) IsSuccessful(id int64) error {
	ml := MsgLog{
		Successful: 2,
	}
	_, err := engineDB.Id(id).Update(&ml)
	return err
}

func (this *MsgLog) IsPaySuccess(id int64, status int) error {
	ml := MsgLog{
		PayStatus: status,
	}
	_, err := engineDB.Id(id).Update(&ml)
	return err
}

func (this *MsgLog) IsDefault(id int64) error {
	ml := MsgLog{
		Successful: 1,
	}
	_, err := engineDB.Id(id).Update(&ml)
	return err
}

func (this *MsgLog) UpRate(id, toid, index, rate, speed, size int64) error {
	if index > size {
		index = size
	}
	if rate > 100 {
		rate = 100
		speed = 0
	}
	ml := MsgLog{
		Toid:  toid,
		Index: index,
		Rate:  rate,
		Size:  size,
		Speed: speed,
	}
	_, err := engineDB.Id(id).Update(&ml)
	return err
}

func (this *MsgLog) IsFalse(id int64) error {
	ml := MsgLog{
		Successful: 3,
	}
	_, err := engineDB.Id(id).Update(&ml)
	return err
}
func (this *MsgLog) GetPage(recipient string, startId int64) ([]MsgLog, error) {
	mls := make([]MsgLog, 0)
	var err error
	if startId == 0 {

		err = engineDB.Alias("t").Where("t.sender = ? or t.recipient = ?",
			recipient, recipient).Limit(10, 0).Desc("createtime").Find(&mls)
	} else {
		err = engineDB.Alias("t").Where("t.id < ? and (t.sender = ? or t.recipient = ?)",
			startId, recipient, recipient).Limit(10, 0).Desc("createtime").Find(&mls)
	}

	return mls, err
}

func (this *MsgLog) FindById(id int64) (*MsgLog, error) {
	mls := make([]MsgLog, 0)

	err := engineDB.Where("id=?", id).Find(&mls)
	if err != nil {
		return nil, err
	}
	if len(mls) <= 0 {
		return nil, errors.New("")
	}
	return &mls[0], nil
}

func (this *MsgLog) Remove(ids ...int64) error {
	_, err := engineDB.In("id", ids).Unscoped().Delete(this)
	return err
}

func (this *MsgLog) RemoveAllForFriend(recipient string) error {
	_, err := engineDB.Where("sender = ? or recipient = ?",
		recipient, recipient).Unscoped().Delete(this)
	return err
}

func (this *MsgLog) FindState(ids []int64) ([]MsgLog, error) {
	res := []MsgLog{}
	err := engineDB.In("id", ids).Unscoped().Find(&res)
	return res, err
}

func (this *MsgLog) FindSize(toid string) (int64, error) {
	size, err := engineDB.Where("Recipient=? and successful=?", toid, 2).SumInt(this, "size")
	return size, err
}
