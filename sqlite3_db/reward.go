package sqlite3_db

import (
	_ "github.com/go-xorm/xorm"
)

type SnapshotReward struct {
	Id          uint64 `xorm:"pk autoincr unique 'id'"`
	Addr        []byte `xorm:"Blob 'addr'"`
	StartHeight uint64 `xorm:"int64 'startheight'"`
	EndHeight   uint64 `xorm:"int64 'endheight'"`
	Reward      uint64 `xorm:"int64 'reward'"`
	LightNum    uint64 `xorm:"int64 'lightnum'"`
	CreateTime  uint64 `xorm:"created 'createtime'"`
}

func (this *SnapshotReward) Add(s *SnapshotReward) error {
	_, err := engineDB.Insert(s)
	return err
}

func (this *SnapshotReward) Find(addr []byte) (*SnapshotReward, error) {
	ss := make([]SnapshotReward, 0)
	err := engineDB.Where("addr = ?", addr).Limit(1, 0).Desc("endheight").Find(&ss)
	if err != nil {
		return nil, err
	}
	if len(ss) <= 0 {
		return nil, nil
	}
	return &ss[0], nil
}

type RewardLight struct {
	Id           uint64 `xorm:"pk autoincr unique 'id'"`
	Sort         uint64 `xorm:"int64 'sort'"`
	SnapshotId   uint64 `xorm:"int64 'snapshotid'"`
	Addr         []byte `xorm:"Blob 'addr'"`
	Reward       uint64 `xorm:"int64 'reward'"`
	Txid         []byte `xorm:"Blob 'txid'"`
	LockHeight   uint64 `xorm:"int64 'lock_height'"`
	Distribution uint64 `xorm:"int64 'distribution'"`
	CreateTime   uint64 `xorm:"created 'createtime'"`
}

func (this *RewardLight) Add(r *RewardLight) error {
	_, err := engineDB.Insert(r)
	return err
}

func (this *RewardLight) FindNotSend(id uint64) (*[]RewardLight, error) {
	ssNotSend := make([]RewardLight, 0)
	err := engineDB.Where("snapshotid = ? and reward <> distribution", id).Find(&ssNotSend)
	if err != nil {
		return nil, err
	}

	return &ssNotSend, nil
}

func (this *RewardLight) UpdateTxid(id uint64) error {
	_, err := engineDB.Where("id = ?", id).Update(this)
	return err
}

func (this *RewardLight) RemoveTxid(ids []uint64) error {
	one := new(RewardLight)
	_, err := engineDB.In("id", ids).Cols("txid", "lock_height").Update(one)
	return err
}

func (this *RewardLight) UpdateDistribution(id uint64, distribution uint64) error {
	reward := &RewardLight{Distribution: distribution}
	_, err := engineDB.Where("id = ?", id).Update(reward)
	return err
}
