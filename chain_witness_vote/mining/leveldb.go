package mining

import (
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"strconv"
	"sync"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

func LoadBlockHeadByHash(hash *[]byte) (*BlockHead, error) {
	bh, err := db.LevelDB.Find(*hash)
	if err != nil {
		return nil, err
	}
	return ParseBlockHeadProto(bh)
}

func LoadBlockHeadByHeight(height uint64) *BlockHead {
	bhash, err := db.LevelDB.Find([]byte(config.BlockHeight + strconv.Itoa(int(height))))
	if err != nil {
		return nil
	}

	bh, err := LoadBlockHeadByHash(bhash)

	if err != nil {
		return nil
	}
	return bh
}

func LoadBlockHashByHeight(height uint64) *[]byte {
	bhash, err := db.LevelDB.Find([]byte(config.BlockHeight + strconv.Itoa(int(height))))
	if err != nil {
		return nil
	}
	return bhash
}

func LoadTxBase(txid []byte) (TxItr, error) {
	var err error
	var txItr TxItr
	ok := false

	if config.EnableCache {

		txItr, ok = TxCache.FindTxInCache(txid)
	}
	if !ok {

		var bs *[]byte
		bs, err = db.LevelDB.Find(txid)
		if err != nil {
			return nil, err
		}

		txItr, err = ParseTxBaseProto(ParseTxClass(txid), bs)
		if err != nil {
			return nil, err
		}
	}
	return txItr, err
}

func loadBlockForDB(bhash *[]byte) (*BlockHead, []TxItr, error) {

	hB, err := LoadBlockHeadByHash(bhash)
	if err != nil {
		return nil, nil, err
	}
	txItrs := make([]TxItr, 0)
	for _, one := range hB.Tx {

		txItr, err := LoadTxBase(one)

		if err != nil {

			return nil, nil, err
		}

		txItrs = append(txItrs, txItr)
	}

	return hB, txItrs, nil
}

func LoadBlockHeadVOByHash(hash *[]byte) (*BlockHeadVO, error) {
	bh, txs, err := loadBlockForDB(hash)
	if err != nil {
		return nil, err
	}

	bhvo := new(BlockHeadVO)
	bhvo.Txs = make([]TxItr, 0)

	bhvo.BH = bh
	bhvo.Txs = txs
	return bhvo, nil
}

func GetCommunityAddrStartHeight() {

}

func ExistCommunityVoteRewardFrozen(addr *crypto.AddressCoin) bool {
	dbkey := config.BuildDBKeyCommunityAddrFrozen(*addr)
	ok, err := db.LevelTempDB.CheckHashExist(*dbkey)
	if err != nil {
		return false
	}
	return ok
}

func GetCommunityVoteRewardFrozen(addr *crypto.AddressCoin) uint64 {
	dbkey := config.BuildDBKeyCommunityAddrFrozen(*addr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return 0
	}
	return utils.BytesToUint64(*valueBs)
}

func SetCommunityVoteRewardFrozen(addr *crypto.AddressCoin, value uint64) {
	dbkey := config.BuildDBKeyCommunityAddrFrozen(*addr)
	valueNewBs := utils.Uint64ToBytes(value)
	db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func ExistCommunityAddr(addr *crypto.AddressCoin) bool {
	dbkey := config.BuildDBKeyCommunityAddr(*addr)
	ok, err := db.LevelTempDB.CheckHashExist(*dbkey)

	if err != nil {
		return false
	}
	return ok
}

func SetCommunityAddr(addr *crypto.AddressCoin) {
	dbkey := config.BuildDBKeyCommunityAddr(*addr)
	err := db.LevelTempDB.Save(*dbkey, nil)
	if err != nil {

	}
}

func RemoveCommunityAddr(addr *crypto.AddressCoin) {
	dbkey := config.BuildDBKeyCommunityAddr(*addr)
	err := db.LevelTempDB.Remove(*dbkey)
	if err != nil {

	}
}

var notSpendBalanceLock = new(sync.RWMutex)
var notSpendBalance = make(map[string]uint64)

func GetNotSpendBalance(addr *crypto.AddressCoin) (*TxItem, uint64) {
	dbkey := config.BuildAddrValue(*addr)
	item := TxItem{
		Addr:  addr,
		Value: 0,
	}
	notSpendBalanceLock.Lock()
	value, _ := notSpendBalance[utils.Bytes2string(dbkey)]
	item.Value = value
	notSpendBalanceLock.Unlock()
	return &item, value

	dbkey = config.BuildAddrValue(*addr)
	valueBs, err := db.LevelTempDB.Find(dbkey)
	if err != nil {
		return nil, 0
	}
	value = utils.BytesToUint64(*valueBs)

	item = TxItem{
		Addr:  addr,
		Value: value,
	}

	return &item, value
}

func SetNotSpendBalance(addr *crypto.AddressCoin, value uint64) error {
	if value > config.Mining_coin_total {
		engine.Log.Info(":%d", value)
		SetAddrValueBig(addr)
	}
	dbkey := config.BuildAddrValue(*addr)
	notSpendBalanceLock.Lock()
	notSpendBalance[utils.Bytes2string(dbkey)] = value
	notSpendBalanceLock.Unlock()
	return nil

	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(dbkey, &valueNewBs)
}

func GetVoteAddr(addr *crypto.AddressCoin) *crypto.AddressCoin {
	dbkey := config.BuildDBKeyVoteAddr(*addr)
	bs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return nil
	}
	if bs == nil || len(*bs) <= 0 {
		return nil
	}
	voteAddr := crypto.AddressCoin(*bs)
	return &voteAddr
}

func SetVoteAddr(addr, voteAddr *crypto.AddressCoin) error {
	dbkey := config.BuildDBKeyVoteAddr(*addr)
	bs := []byte(*voteAddr)
	return db.LevelTempDB.Save(*dbkey, &bs)
}

func RemoveVoteAddr(addr *crypto.AddressCoin) error {
	dbkey := config.BuildDBKeyVoteAddr(*addr)
	return db.LevelTempDB.Remove(*dbkey)
}

func GetDepositWitnessAddr(addr *crypto.AddressCoin) uint64 {
	dbkey := config.BuildDBKeyDepositWitnessAddr(*addr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return 0
	}
	value := utils.BytesToUint64(*valueBs)
	return value
}

func SetDepositWitnessAddr(addr *crypto.AddressCoin, value uint64) error {
	dbkey := config.BuildDBKeyDepositWitnessAddr(*addr)
	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func RemoveDepositWitnessAddr(addr *crypto.AddressCoin) error {
	dbkey := config.BuildDBKeyDepositWitnessAddr(*addr)
	return db.LevelTempDB.Remove(*dbkey)
}

func GetDepositLightAddr(addr *crypto.AddressCoin) uint64 {
	dbkey := config.BuildDBKeyDepositLightAddr(*addr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return 0
	}
	value := utils.BytesToUint64(*valueBs)
	return value
}

func SetDepositLightAddr(addr *crypto.AddressCoin, value uint64) error {
	dbkey := config.BuildDBKeyDepositLightAddr(*addr)
	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func RemoveDepositLightAddr(addr *crypto.AddressCoin) error {
	dbkey := config.BuildDBKeyDepositLightAddr(*addr)
	return db.LevelTempDB.Remove(*dbkey)
}

func GetDepositCommunityAddr(addr *crypto.AddressCoin) uint64 {
	dbkey := config.BuildDBKeyDepositCommunityAddr(*addr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return 0
	}
	value := utils.BytesToUint64(*valueBs)
	return value
}

func SetDepositCommunityAddr(addr *crypto.AddressCoin, value uint64) error {
	dbkey := config.BuildDBKeyDepositCommunityAddr(*addr)
	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func RemoveDepositCommunityAddr(addr *crypto.AddressCoin) error {
	dbkey := config.BuildDBKeyDepositCommunityAddr(*addr)
	return db.LevelTempDB.Remove(*dbkey)
}

func GetDepositLightVoteValue(lightAddr, communityAddr *crypto.AddressCoin) uint64 {
	dbkey := config.BuildDBKeyDepositLightVoteValue(*lightAddr, *communityAddr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return 0
	}
	value := utils.BytesToUint64(*valueBs)
	return value
}

func SetDepositLightVoteValue(lightAddr, communityAddr *crypto.AddressCoin, value uint64) error {
	dbkey := config.BuildDBKeyDepositLightVoteValue(*lightAddr, *communityAddr)
	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func RemoveDepositLightVoteValue(lightAddr, communityAddr *crypto.AddressCoin) error {
	dbkey := config.BuildDBKeyDepositLightVoteValue(*lightAddr, *communityAddr)
	return db.LevelTempDB.Remove(*dbkey)
}

func GetAddrNonce(addr *crypto.AddressCoin) (big.Int, error) {
	dbkey := config.BuildDBKeyAddrNonce(*addr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return big.Int{}, err
	}
	nonce := new(big.Int).SetBytes(*valueBs)
	return *nonce, nil
}

func SetAddrNonce(addr *crypto.AddressCoin, value *big.Int) error {
	dbkey := config.BuildDBKeyAddrNonce(*addr)
	valueNewBs := value.Bytes()
	return db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func SetNotSpendBalanceToken(txid *[]byte, addr *crypto.AddressCoin, value uint64) error {
	dbkey := config.BuildDBKeyTokenAddrValue(*txid, *addr)
	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(*dbkey, &valueNewBs)
}

func GetNotSpendBalanceToken(txid *[]byte, addr *crypto.AddressCoin) (*TxItem, uint64) {
	dbkey := config.BuildDBKeyTokenAddrValue(*txid, *addr)
	valueBs, err := db.LevelTempDB.Find(*dbkey)
	if err != nil {
		return nil, 0
	}
	value := utils.BytesToUint64(*valueBs)

	item := TxItem{
		Addr:  addr,
		Value: value,
	}

	return &item, value
}

func FindNotSpendBalanceToken(txid *[]byte, amount uint64) (*TxItem, uint64) {
	addrs := keystore.GetAddrAll()
	for _, one := range addrs {
		item, value := GetNotSpendBalanceToken(txid, &one.Addr)
		if amount <= value {
			return item, value
		}
	}
	return nil, 0
}

func GetAddrFrozenValue(addr *crypto.AddressCoin) uint64 {
	dbkey := config.BuildAddrFrozen(*addr)
	valueBs, err := db.LevelTempDB.Find(dbkey)
	if err != nil {
		return 0
	}
	value := utils.BytesToUint64(*valueBs)
	return value
}

func SetAddrFrozenValue(addr *crypto.AddressCoin, value uint64) error {
	dbkey := config.BuildAddrFrozen(*addr)
	valueNewBs := utils.Uint64ToBytes(value)
	return db.LevelTempDB.Save(dbkey, &valueNewBs)
}

func RemoveAddrFrozenValue(addr *crypto.AddressCoin) error {
	dbkey := config.BuildAddrFrozen(*addr)
	return db.LevelTempDB.Remove(dbkey)
}

func GetAddrValueBig(addr *crypto.AddressCoin) bool {
	dbkey := config.BuildAddrValueBig(*addr)
	ok, err := db.LevelTempDB.CheckHashExist(dbkey)
	if err != nil {
		return true
	}
	return ok
}

func SetAddrValueBig(addr *crypto.AddressCoin) error {
	dbkey := config.BuildAddrValueBig(*addr)
	return db.LevelTempDB.Save(dbkey, nil)
}
