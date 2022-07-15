package sqlite3_db

import (
	"database/sql"
	"fmt"
	"mmschainnewaccount/config"
	"sync"

	"github.com/go-xorm/xorm"

	_ "github.com/logoove/sqlite"
)

var once sync.Once

var db *sql.DB

var engineDB *xorm.Engine
var dblock = new(sync.Mutex)

var (
	table_friends          *xorm.Session
	table_sharefolder      *xorm.Session
	table_msglog           *xorm.Session
	table_downloadprogress *xorm.Session
	table_storefolder      *xorm.Session
	table_storefile        *xorm.Session
	table_property         *xorm.Session
	table_peerinfo         *xorm.Session
	table_snapshot         *xorm.Session
	table_reward           *xorm.Session
	table_msgcache         *xorm.Session

	table_wallet_txitem *xorm.Session
	table_wallet_testdb *xorm.Session
)

func Init() {
	once.Do(connect)

}

func connect() {
	var err error
	engineDB, err = xorm.NewEngine("sqlite3", "file:"+config.SQLITE3DB_path+"?cache=shared")
	if err != nil {
		fmt.Println(err)
	}
	engineDB.ShowSQL(config.SQL_SHOW)

	ok, err := engineDB.IsTableExist(Friends{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(Friends{}).CreateTable(Friends{})
		table_friends = engineDB.Table(Friends{})
		table_friends.CreateIndexes(Friends{})
		table_friends.CreateUniques(Friends{})

		engineDB.Table(ShareFolder{}).CreateTable(ShareFolder{})
		table_sharefolder = engineDB.Table(ShareFolder{})
		table_sharefolder.CreateIndexes(ShareFolder{})
		table_sharefolder.CreateUniques(ShareFolder{})
	} else {
		table_friends = engineDB.Table(Friends{})
		table_sharefolder = engineDB.Table(ShareFolder{})
	}

	ok, err = engineDB.IsTableExist(MsgLog{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(MsgLog{}).CreateTable(MsgLog{})
		table_msglog = engineDB.Table(MsgLog{})
		table_msglog.CreateIndexes(MsgLog{})
		table_msglog.CreateUniques(MsgLog{})
	} else {
		table_msglog = engineDB.Table(MsgLog{})
	}

	LoadMsgLogGenerateID()

	ok, err = engineDB.IsTableExist(Downprogress{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(Downprogress{}).CreateTable(Downprogress{})
	}
	table_downloadprogress = engineDB.Table(Downprogress{})

	ok, err = engineDB.IsTableExist(StoreFolder{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(StoreFolder{}).CreateTable(StoreFolder{})
	}
	table_storefolder = engineDB.Table(StoreFolder{})

	ok, err = engineDB.IsTableExist(StoreFolderFile{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(StoreFolderFile{}).CreateTable(StoreFolderFile{})
	}
	table_storefile = engineDB.Table(StoreFolderFile{})

	ok, err = engineDB.IsTableExist(Property{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(Property{}).CreateTable(Property{})
	}
	table_property = engineDB.Table(Property{})

	ok, err = engineDB.IsTableExist(PeerInfo{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(PeerInfo{}).CreateTable(PeerInfo{})
		table_peerinfo = engineDB.Table(PeerInfo{})
		table_peerinfo.CreateIndexes(PeerInfo{})
		table_peerinfo.CreateUniques(PeerInfo{})
	} else {

		engineDB.Table(PeerInfo{}).CreateTable(PeerInfo{})
		table_peerinfo = engineDB.Table(PeerInfo{})
		table_peerinfo.CreateIndexes(PeerInfo{})
		table_peerinfo.CreateUniques(PeerInfo{})
	}

	ok, err = engineDB.IsTableExist(SnapshotReward{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(SnapshotReward{}).CreateTable(SnapshotReward{})
		table_snapshot = engineDB.Table(SnapshotReward{})
		table_snapshot.CreateIndexes(SnapshotReward{})
		table_snapshot.CreateUniques(SnapshotReward{})
	} else {

		engineDB.Table(SnapshotReward{}).CreateTable(SnapshotReward{})
		table_snapshot = engineDB.Table(SnapshotReward{})
		table_snapshot.CreateIndexes(SnapshotReward{})
		table_snapshot.CreateUniques(SnapshotReward{})
	}

	ok, err = engineDB.IsTableExist(RewardLight{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(RewardLight{}).CreateTable(RewardLight{})
		table_reward = engineDB.Table(RewardLight{})
		table_reward.CreateIndexes(RewardLight{})
		table_reward.CreateUniques(RewardLight{})
	} else {

		engineDB.Table(RewardLight{}).CreateTable(RewardLight{})
		table_reward = engineDB.Table(RewardLight{})
		table_reward.CreateIndexes(RewardLight{})
		table_reward.CreateUniques(RewardLight{})
	}

	ok, err = engineDB.IsTableExist(MessageCache{})
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		engineDB.Table(MessageCache{}).CreateTable(MessageCache{})
		table_msgcache = engineDB.Table(MessageCache{})
		table_msgcache.CreateIndexes(MessageCache{})
		table_msgcache.CreateUniques(MessageCache{})
	} else {

		table_msgcache = engineDB.Table(MessageCache{})
		table_msgcache.CreateIndexes(MessageCache{})
		table_msgcache.CreateUniques(MessageCache{})
	}
	initWalletTable()
}

func initWalletTable() {

	engineDB.DropTables(TxItem{})
	engineDB.Table(TxItem{}).CreateTable(TxItem{})
	table_wallet_txitem = engineDB.Table(TxItem{})
	table_wallet_txitem.CreateIndexes(TxItem{})
	table_wallet_txitem.CreateUniques(TxItem{})

	engineDB.DropTables(TestDB{})
	engineDB.Table(TestDB{}).CreateTable(TestDB{})
	table_wallet_testdb = engineDB.Table(TestDB{})
	table_wallet_testdb.CreateIndexes(TestDB{})
	table_wallet_testdb.CreateUniques(TestDB{})

}
