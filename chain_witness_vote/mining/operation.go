package mining

import (
	"bytes"
	"errors"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"mmschainnewaccount/core"
	"mmschainnewaccount/protos/go_protos"
	"mmschainnewaccount/sqlite3_db"
	"strconv"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/message_center/flood"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"golang.org/x/crypto/ed25519"
)

func FindBalanceValue() (uint64, uint64, uint64) {
	chain := forks.GetLongChain()
	if chain == nil {
		return 0, 0, 0
	}
	var notspend, frozen, lock uint64
	for _, one := range keystore.GetAddrAll() {
		n, f, l := GetBalanceForAddrSelf(one.Addr)
		notspend += n
		frozen += f
		lock += l
	}
	return notspend, frozen, lock
}

func GetBalanceForAddrSelf(addr crypto.AddressCoin) (uint64, uint64, uint64) {
	chain := forks.GetLongChain()
	if chain == nil {
		return 0, 0, 0
	}

	_, notspend := GetNotSpendBalance(&addr)

	frozenValue := GetAddrFrozenValue(&addr)

	communityFrozen := GetCommunityVoteRewardFrozen(&addr)
	lockValue, lockVoteReward := chain.Balance.FindLockTotalByAddr(&addr)

	frozenValue += communityFrozen
	notspend -= communityFrozen
	notspend -= lockValue

	frozenValue -= lockVoteReward
	lockValue += lockVoteReward

	return notspend, frozenValue, lockValue
}

func GetNotspendByAddrOther(addr crypto.AddressCoin) (uint64, uint64, uint64) {
	chain := forks.GetLongChain()
	if chain == nil {
		return 0, 0, 0
	}
	_, notspend := GetNotSpendBalance(&addr)

	frozenValue := GetAddrFrozenValue(&addr)

	communityFrozen := GetCommunityVoteRewardFrozen(&addr)

	frozenValue += communityFrozen
	notspend -= communityFrozen

	return notspend, frozenValue, 0
}

func SendToAddress(srcAddress, address *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string) (*Tx_Pay, error) {

	txpay, err := CreateTxPay(srcAddress, address, amount, gas, frozenHeight, pwd, comment)
	if err != nil {

		return nil, err
	}
	txpay.BuildHash()

	forks.GetLongChain().transactionManager.AddTx(txpay)

	MulticastTx(txpay)

	return txpay, nil
}

func SendToMoreAddress(addr *crypto.AddressCoin, address []PayNumber, gas uint64, pwd, comment string) (*Tx_Pay, error) {
	txpay, err := CreateTxsPay(addr, address, gas, pwd, comment)
	if err != nil {

		return nil, err
	}
	txpay.BuildHash()

	forks.GetLongChain().transactionManager.AddTx(txpay)
	MulticastTx(txpay)

	return txpay, nil
}

func SendToMoreAddressByPayload(addr *crypto.AddressCoin, address []PayNumber, gas uint64, pwd string, cs *CommunitySign, startHeight, endHeight uint64) (*Tx_Vote_Reward, error) {
	txpay, err := CreateTxVoteReward(addr, address, gas, pwd, startHeight, endHeight)
	if err != nil {

		return nil, err
	}
	txpay.BuildHash()

	forks.GetLongChain().transactionManager.AddTx(txpay)
	MulticastTx(txpay)

	return txpay, nil
}

func FindStartBlockForNeighbor() *ChainInfo {
	for _, key := range nodeStore.GetLogicNodes() {

		message, _ := message_center.SendNeighborMsg(config.MSGID_getStartBlockHead, &key, nil)

		bs, _ := flood.WaitRequest(message_center.CLASS_getBlockHead, utils.Bytes2string(message.Body.Hash), 0)

		if bs == nil {

			continue
		}
		chainInfo, err := ParseChainInfo(bs)

		if err != nil {
			return nil
		}

		return chainInfo

	}
	return nil
}

func FindBlockForNeighbor(bhash *[]byte, peerBlockInfo *PeerBlockInfoDESC) (*BlockHeadVO, error) {
	var bhvo *BlockHeadVO
	var bs *[]byte
	var err error
	var newBhvo *BlockHeadVO

	peers := peerBlockInfo.Sort()

	addrs := make([]nodeStore.AddressNet, 0)
	for _, one := range peers {
		addrs = append(addrs, *one.Addr)
	}
	logicNodesInfo := core.SortNetAddrForSpeed(addrs)

	for i, _ := range logicNodesInfo {

		key := &logicNodesInfo[i].AddrNet

		engine.Log.Info("Send query message to node %s", key.B58String())
		bs, err = getBlockHeadVO(*key, bhash)
		if err != nil {
			engine.Log.Info("Send query message to node from:%s error:%s", key.B58String(), err.Error())
			continue
		}
		if bs == nil {
			engine.Log.Info("Send query message to node from:%s bs is nil", key.B58String())
			continue
		}
		newBhvo, err = ParseBlockHeadVOProto(bs)

		if err != nil {
			engine.Log.Info("Send query message to node from:%s error:%s", key.B58String(), err.Error())
			continue
		}
		bhvo = newBhvo

		if newBhvo.BH.Nextblockhash != nil && len(newBhvo.BH.Nextblockhash) > 0 {
			engine.Log.Info("this block next block hash not nil")
			return newBhvo, err
		}
		engine.Log.Info("this block next block hash is nil")

	}

	if bhvo != nil {
		return bhvo, nil
	}

	return bhvo, err
}

func FindBlockHeadNeighbor(bhash *[]byte, peerBlockInfo *PeerBlockInfoDESC) (*BlockHead, error) {
	var bh *BlockHead
	var bs *[]byte
	var err error
	var newBh *BlockHead

	peers := peerBlockInfo.Sort()
	addrs := make([]nodeStore.AddressNet, 0)
	for _, one := range peers {
		addrs = append(addrs, *one.Addr)
	}
	logicNodesInfo := core.SortNetAddrForSpeed(addrs)

	for i, _ := range logicNodesInfo {
		addrOne := &logicNodesInfo[i].AddrNet

		engine.Log.Info("Send query message to node %s", addrOne.B58String())
		bs, err = getKeyValue(*addrOne, bhash)
		if err != nil {
			engine.Log.Info("Send query message to node from:%s error:%s", addrOne.B58String(), err.Error())
			continue
		}
		if bs == nil {
			engine.Log.Info("Send query message to node from:%s bs is nil", addrOne.B58String())
			continue
		}
		newBh, err = ParseBlockHeadProto(bs)

		if err != nil {
			engine.Log.Info("Send query message to node from:%s error:%s", addrOne.B58String(), err.Error())
			continue
		}
		bh = newBh

		if newBh.Nextblockhash != nil && len(newBh.Nextblockhash) > 0 {
			engine.Log.Info("this block next block hash not nil")
			return newBh, err
		}
		engine.Log.Info("this block next block hash is nil")

	}

	if bh != nil {
		return bh, nil
	}

	return bh, err
}

func FindLastBlockForNeighbor(peerBlockInfo *PeerBlockInfoDESC) (*BlockHead, error) {
	var err error
	var bh *BlockHead

	peers := peerBlockInfo.Sort()
	addrs := make([]nodeStore.AddressNet, 0)
	for _, one := range peers {
		addrs = append(addrs, *one.Addr)
	}
	logicNodesInfo := core.SortNetAddrForSpeed(addrs)

	for i, _ := range logicNodesInfo {
		key := &logicNodesInfo[i].AddrNet

		message, _ := message_center.SendNeighborMsg(config.MSGID_getBlockLastCurrent, key, nil)
		bs, _ := flood.WaitRequest(message_center.CLASS_getBlockLastCurrent, utils.Bytes2string(message.Body.Hash), config.Wallet_sync_block_timeout)
		if bs == nil {
			continue
		}
		bh, err = ParseBlockHeadProto(bs)
		if err != nil {
			continue
		}
		return bh, nil
	}
	return nil, err
}

func getBlockHeadVO(key nodeStore.AddressNet, bhash *[]byte) (*[]byte, error) {

	start := time.Now()

	message, _ := message_center.SendNeighborMsg(config.MSGID_getBlockHeadVO, &key, bhash)

	bs, _ := flood.WaitRequest(message_center.CLASS_getTransaction, utils.Bytes2string(message.Body.Hash), config.Wallet_sync_block_timeout)
	if bs == nil {
		endTime := time.Now()

		engine.Log.Error("Receive %s message timeout %s", key.B58String(), time.Now().Sub(start))

		if (endTime.Unix() - start.Unix()) < config.Wallet_sync_block_timeout {
			core.AddNodeAddrSpeed(key, time.Second*(config.Wallet_sync_block_timeout+1))
		} else {

			core.AddNodeAddrSpeed(key, time.Now().Sub(start))
		}

		return nil, config.ERROR_chain_sync_block_timeout
	}
	core.AddNodeAddrSpeed(key, time.Now().Sub(start))

	return bs, nil
}

func getKeyValue(addrOne nodeStore.AddressNet, key *[]byte) (*[]byte, error) {

	start := time.Now()

	message, _ := message_center.SendNeighborMsg(config.MSGID_getDBKey_one, &addrOne, key)

	bs, err := flood.WaitRequest(config.CLASS_getKeyValue, utils.Bytes2string(message.Body.Hash), config.Wallet_sync_block_timeout)
	if err != nil || bs == nil {
		endTime := time.Now()

		engine.Log.Error("Receive %s message timeout %s", addrOne.B58String(), time.Now().Sub(start))

		if (endTime.Unix() - start.Unix()) < config.Wallet_sync_block_timeout {
			core.AddNodeAddrSpeed(addrOne, time.Second*(config.Wallet_sync_block_timeout+1))
		} else {

			core.AddNodeAddrSpeed(addrOne, time.Now().Sub(start))
		}

		return nil, config.ERROR_chain_sync_block_timeout
	}
	core.AddNodeAddrSpeed(addrOne, time.Now().Sub(start))

	return bs, nil
}

func GetUnconfirmedBlockForNeighbor(height uint64, peerBlockInfo *PeerBlockInfoDESC) ([]*BlockHeadVO, error) {
	engine.Log.Info("Synchronize unacknowledged chunks from neighbor nodes")

	heightBs := utils.Uint64ToBytes(height)

	var bs *[]byte
	var err error

	logicNodes := peerBlockInfo.Sort()
	for j, _ := range logicNodes {
		engine.Log.Info("Synchronize unacknowledged from:%s height:%d", logicNodes[j].Addr.B58String(), height)

		message, err := message_center.SendNeighborMsg(config.MSGID_getUnconfirmedBlock, logicNodes[j].Addr, &heightBs)
		if err != nil {

			continue
		}

		bs, _ = flood.WaitRequest(message_center.CLASS_getUnconfirmedBlock, utils.Bytes2string(message.Body.Hash), 0)
		if bs == nil {
			engine.Log.Info("Failed to get unconfirmed block from neighbor node, sending shared file message, maybe timeout")
			err = errors.New("Failed to get unconfirmed block from neighbor node, sending shared file message, maybe timeout")
			continue
		} else {
			err = nil
		}
		engine.Log.Info("Synchronize unacknowledged ok")
		break
	}

	if bs == nil {
		engine.Log.Warn("Get unacknowledged block BS error")
		return nil, err
	}

	rbsp := new(go_protos.RepeatedBytes)
	err = proto.Unmarshal(*bs, rbsp)
	if err != nil {
		engine.Log.Warn("Get unacknowledged block BS error:", err.Error())
		return nil, err
	}

	blockHeadVOs := make([]*BlockHeadVO, 0)
	for _, one := range rbsp.Bss {
		bhvo, err := ParseBlockHeadVOProto(&one)
		if err != nil {
			engine.Log.Warn("Get unacknowledged block BS error:", err.Error())
			return nil, err
		}
		blockHeadVOs = append(blockHeadVOs, bhvo)
	}

	engine.Log.Info("Get unacknowledged block Success")
	return blockHeadVOs, nil

}

func DepositIn(amount, gas uint64, pwd, payload string) error {

	err := forks.GetLongChain().Balance.DepositIn(amount, gas, pwd, payload)
	if err != nil {

	}

	return err
}

func DepositOut(addr string, amount, gas uint64, pwd string) error {

	err := forks.GetLongChain().Balance.DepositOut(addr, amount, gas, pwd)

	return err
}

func VoteIn(t uint16, witnessAddr crypto.AddressCoin, addr crypto.AddressCoin, amount, gas uint64, pwd, payload string) error {

	err := forks.GetLongChain().Balance.VoteIn(t, witnessAddr, addr, amount, gas, pwd, payload)
	if err != nil {

	}

	return err
}

func VoteOut(voteType uint16, addr crypto.AddressCoin, amount, gas uint64, pwd, payload string) error {

	return forks.GetLongChain().Balance.VoteOut(voteType, addr, amount, gas, pwd, payload)

}

func GetWitnessStatus() (IsCandidate bool, IsBackup bool, IsKickOut bool, Addr crypto.AddressCoin, value uint64) {
	addrInfo := keystore.GetCoinbase()
	Addr = addrInfo.Addr
	IsCandidate = forks.GetLongChain().WitnessBackup.FindWitness(Addr)
	IsBackup = forks.GetLongChain().WitnessChain.FindWitness(addrInfo.Addr)
	IsKickOut = forks.GetLongChain().WitnessBackup.FindWitnessInBlackList(Addr)
	txItem := forks.GetLongChain().Balance.GetDepositIn()
	if txItem == nil {
		value = 0
	} else {
		value = uint64(txItem.Value)
	}
	return
}

func GetWitnessListSort() *WitnessBackupGroup {
	return forks.GetLongChain().WitnessBackup.GetWitnessListSort()
}

func GetCommunityListSort() []*VoteScoreVO {
	return forks.GetLongChain().WitnessBackup.GetCommunityListSort()
}

func GetDepositCommunityList() []*DepositInfo {
	items := make([]*DepositInfo, 0)
	for _, one := range keystore.GetAddrAll() {
		txItem := GetLongChain().Balance.GetDepositCommunity(&one.Addr)
		if txItem != nil {
			items = append(items, txItem)
		}
	}
	return items
}

func GetDepositLightList() []*DepositInfo {
	items := make([]*DepositInfo, 0)
	for _, one := range keystore.GetAddrAll() {
		txItem := GetLongChain().Balance.GetDepositLight(&one.Addr)
		if txItem != nil {
			items = append(items, txItem)
		}
	}
	return items
}

func GetDepositVoteList() []*DepositInfo {
	items := make([]*DepositInfo, 0)
	for _, one := range keystore.GetAddrAll() {
		txItem := GetLongChain().Balance.GetDepositVote(&one.Addr)
		if txItem != nil {
			items = append(items, txItem)
		}
	}
	return items
}

func FindTx(txid []byte) (TxItr, uint64) {

	txItr, err := LoadTxBase(txid)
	if err != nil {
		return nil, 1
	}

	_, err = db.GetTxToBlockHash(&txid)
	if err != nil {
		lockheight := txItr.GetLockHeight()

		height := GetLongChain().GetCurrentBlock()
		if height > lockheight {

			return txItr, 3
		}
		return txItr, 1
	}
	return txItr, 2
}

func FindTxJsonVo(txid []byte) (interface{}, uint64) {
	txItr, code := FindTx(txid)
	return txItr.GetVOJSON(), code
}

func GetAddrState(addr crypto.AddressCoin) int {
	witnessBackup := forks.GetLongChain().WitnessBackup

	_, isLight := witnessBackup.haveLight(&addr)
	if isLight {
		return 3
	}

	_, isCommunity := witnessBackup.haveCommunityList(&addr)
	if isCommunity {
		return 2
	}

	isWitness := witnessBackup.haveWitness(&addr)
	if isWitness {
		return 1
	}
	return 4
}

func AddTx(txItr TxItr) error {
	if txItr == nil {

		return errors.New("Failure to pay deposit")
	}
	txItr.BuildHash()

	ok := forks.GetLongChain().transactionManager.AddTx(txItr)
	if !ok {

		return errors.New("Waiting for the chain, please try again later")
	}
	MulticastTx(txItr)
	return nil
}

func CreateTxPayM(height uint64, items []*TxItem, pubs map[string]ed25519.PublicKey, address *crypto.AddressCoin, amount, gas uint64, comment string, returnaddr crypto.AddressCoin) (*Tx_Pay, error) {
	if len(items) == 0 {

		return nil, config.ERROR_not_enough
	}

	vins := make([]*Vin, 0)
	total := uint64(0)

	for _, item := range items {

		addrstr := *item.Addr
		puk, ok := pubs[addrstr.B58String()]
		if !ok {
			continue
		}

		vin := Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)

		total = total + uint64(item.Value)
		if total >= amount+gas {

			break
		}

	}

	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	vouts := make([]*Vout, 0)
	vout := Vout{
		Value:   amount,
		Address: *address,
	}
	vouts = append(vouts, &vout)

	if total > amount+gas {
		vout := Vout{
			Value: total - amount - gas,

			Address: returnaddr,
		}
		vouts = append(vouts, &vout)
	}
	var pay *Tx_Pay

	base := TxBase{
		Type:       config.Wallet_tx_type_pay,
		Vin_total:  uint64(len(vins)),
		Vin:        vins,
		Vout_total: uint64(len(vouts)),
		Vout:       vouts,
		Gas:        gas,

		LockHeight: height + 100,
		Payload:    []byte(comment),
	}
	pay = &Tx_Pay{
		TxBase: base,
	}

	for i, _ := range pay.Vin {
		sign := pay.GetWaitSign(uint64(i))
		if sign == nil {

			return nil, errors.New("Data error while pre signing")
		}

		pay.Vin[i].Sign = *sign
	}

	return pay, nil
}

func CreateTxsPayM(height uint64, items []*TxItem, pubs map[string]ed25519.PublicKey, address []PayNumber, gas uint64, comment string, returnaddr crypto.AddressCoin) (*Tx_Pay, error) {
	if len(items) == 0 {

		return nil, config.ERROR_not_enough
	}

	vins := make([]*Vin, 0)
	total := uint64(0)
	amount := uint64(0)
	for _, one := range address {
		amount += one.Amount
	}

	for _, item := range items {

		addrstr := *item.Addr
		puk, ok := pubs[addrstr.B58String()]
		if !ok {
			continue
		}

		vin := Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)

		total = total + uint64(item.Value)
		if total >= amount+gas {

			break
		}

	}

	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	vouts := make([]*Vout, 0)
	for _, one := range address {
		vout := Vout{
			Value:   one.Amount,
			Address: one.Address,
		}
		vouts = append(vouts, &vout)
	}

	if total > amount+gas {
		vout := Vout{
			Value: total - amount - gas,

			Address: returnaddr,
		}
		vouts = append(vouts, &vout)
	}
	var pay *Tx_Pay

	base := TxBase{
		Type:       config.Wallet_tx_type_pay,
		Vin_total:  uint64(len(vins)),
		Vin:        vins,
		Vout_total: uint64(len(vouts)),
		Vout:       vouts,
		Gas:        gas,

		LockHeight: height + 100,
		Payload:    []byte(comment),
	}
	pay = &Tx_Pay{
		TxBase: base,
	}

	for i, _ := range pay.Vin {
		sign := pay.GetWaitSign(uint64(i))
		if sign == nil {

			return nil, errors.New("Data error while pre signing")
		}

		pay.Vin[i].Sign = *sign
	}

	return pay, nil
}

func CreateTxVoteInM(height uint64, items []*TxItem, pubs map[string]ed25519.PublicKey, voteType uint16, witnessAddr crypto.AddressCoin, addr string, amount, gas uint64, comment string, returnaddr crypto.AddressCoin) (*Tx_vote_in, error) {
	if len(items) == 0 {

		return nil, config.ERROR_not_enough
	}
	if voteType == 1 && amount != config.Mining_vote {

		return nil, errors.New("Minimum deposit required" + strconv.FormatUint(config.Mining_vote, 10))
	}
	if voteType == 3 && amount != config.Mining_light_min {

		return nil, errors.New("Minimum deposit required" + strconv.FormatUint(config.Mining_light_min, 10))
	}

	vins := make([]*Vin, 0)
	total := uint64(0)

	for _, item := range items {

		addrstr := *item.Addr
		puk, ok := pubs[addrstr.B58String()]
		if !ok {
			continue
		}

		vin := Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)

		total = total + uint64(item.Value)
		if total >= amount+gas {

			break
		}

	}

	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	var dstAddr crypto.AddressCoin
	if addr == "" {

		dstAddr = returnaddr
	} else {

		dstAddr = crypto.AddressFromB58String(addr)
	}

	vouts := make([]*Vout, 0)
	vout := Vout{
		Value:   amount,
		Address: dstAddr,
	}
	vouts = append(vouts, &vout)

	if total > amount+gas {
		vout := Vout{
			Value: total - amount - gas,

			Address: returnaddr,
		}
		vouts = append(vouts, &vout)
	}
	var txin *Tx_vote_in

	base := TxBase{
		Type:       config.Wallet_tx_type_vote_in,
		Vin_total:  uint64(len(vins)),
		Vin:        vins,
		Vout_total: uint64(len(vouts)),
		Vout:       vouts,
		Gas:        gas,

		LockHeight: height + 100,
		Payload:    []byte(comment),
	}

	txin = &Tx_vote_in{
		TxBase:   base,
		Vote:     witnessAddr,
		VoteType: voteType,
	}

	for i, _ := range txin.Vin {
		sign := txin.GetWaitSign(uint64(i))
		if sign == nil {
			return nil, config.ERROR_get_sign_data_fail
		}

		txin.Vin[i].Sign = *sign
	}

	return txin, nil
}

func CreateTxVoteOutM(height uint64, voteitems, items []*TxItem, pubs map[string]ed25519.PublicKey, witness *crypto.AddressCoin, addr string, amount, gas uint64, returnaddr crypto.AddressCoin) (*Tx_vote_out, error) {

	vins := make([]*Vin, 0)
	total := uint64(0)

	for _, item := range voteitems {

		voutaddr := *item.Addr
		puk, ok := pubs[voutaddr.B58String()]
		if !ok {
			continue
		}

		vin := Vin{

			Puk: puk,
		}
		vins = append(vins, &vin)

		total = total + uint64(item.Value)
		if total >= amount+gas {
			break
		}
	}

	if total < amount+gas {
		for _, item := range items {

			addrstr := *item.Addr
			puk, ok := pubs[addrstr.B58String()]
			if !ok {
				continue
			}

			vin := Vin{

				Puk: puk,
			}
			vins = append(vins, &vin)

			total = total + uint64(item.Value)
			if total >= amount+gas {
				break
			}
		}
	}

	if total < (amount + gas) {

		return nil, config.ERROR_not_enough
	}

	var dstAddr crypto.AddressCoin
	if addr == "" {

		dstAddr = returnaddr
	} else {

		dstAddr = crypto.AddressFromB58String(addr)
	}

	vouts := make([]*Vout, 0)

	vout := Vout{
		Value:   total - gas,
		Address: dstAddr,
	}
	vouts = append(vouts, &vout)

	var txout *Tx_vote_out

	base := TxBase{
		Type:       config.Wallet_tx_type_vote_out,
		Vin_total:  uint64(len(vins)),
		Vin:        vins,
		Vout_total: uint64(len(vouts)),
		Vout:       vouts,
		Gas:        gas,
		LockHeight: height + 100,
	}
	txout = &Tx_vote_out{
		TxBase: base,
	}

	for i, _ := range txout.Vin {
		sign := txout.GetWaitSign(uint64(i))
		if sign == nil {
			return nil, config.ERROR_get_sign_data_fail
		}

		txout.Vin[i].Sign = *sign
	}

	return txout, nil
}

type BlockVotesVO struct {
	EndHeight uint64
	Group     []GroupVO
}

type GroupVO struct {
	StartHeight    uint64
	EndHeight      uint64
	CommunityVotes []VoteScoreRewadrVO
}

type VoteScoreRewadrVO struct {
	VoteScore
	LightVotes []VoteScore
	Reward     uint64
}

func FindLightVote(startHeight, endHeight uint64) (*BlockVotesVO, error) {

	bvVO := &BlockVotesVO{
		EndHeight: endHeight,
		Group:     make([]GroupVO, 0),
	}

	var preBlock *Block
	preGroup := forks.LongChain.WitnessChain.witnessGroup
	for {

		preGroup = preGroup.PreGroup
		ok, preGroupBlock := preGroup.CheckBlockGroup(nil)
		if ok {
			preBlock = preGroupBlock.Blocks[len(preGroupBlock.Blocks)-1]
			break
		}
	}

	for {

		if preBlock.Height > startHeight {
			if preBlock.PreBlock == nil {
				break
			}
			preBlock = preBlock.PreBlock
		}
		if preBlock.Height < startHeight {
			if preBlock.NextBlock == nil {
				break
			}
			preBlock = preBlock.NextBlock
		}
		if preBlock.Height == startHeight {
			break
		}
	}

	for {

		if preBlock.PreBlock == nil {

			break
		}
		temp := preBlock.witness.WitnessBackupGroup
		if preBlock.PreBlock.witness.WitnessBackupGroup == temp {
			preBlock = preBlock.PreBlock
		} else {
			break
		}
	}

	isFind := false
	for ; preBlock != nil && preBlock.NextBlock != nil && preBlock.Height <= endHeight; preBlock = preBlock.NextBlock {

		if isFind && preBlock.NextBlock.witness.WitnessBackupGroup == preBlock.witness.WitnessBackupGroup {

			continue
		} else {
			isFind = false
		}
		_, txs, err := preBlock.NextBlock.LoadTxs()
		if err != nil {

			return nil, err
		}
		if txs == nil || len(*txs) <= 0 {

			continue
		}

		reward, ok := (*txs)[0].(*Tx_reward)
		if !ok {
			continue
		}
		isFind = true

		groupVO := new(GroupVO)
		groupVO.CommunityVotes = make([]VoteScoreRewadrVO, 0)

		groupVO.EndHeight = preBlock.Height

		for _, one := range preBlock.witness.WitnessBackupGroup.Witnesses {

			m := make(map[string]*[]VoteScore)
			for i, two := range one.Votes {
				vo := VoteScore{
					Witness: one.Votes[i].Witness,
					Addr:    one.Votes[i].Addr,
					Scores:  one.Votes[i].Scores,
					Vote:    one.Votes[i].Vote,
				}
				v, ok := m[utils.Bytes2string(*two.Witness)]
				if ok {

				} else {
					temp := make([]VoteScore, 0)
					v = &temp
				}
				*v = append(*v, vo)
				m[utils.Bytes2string(*two.Witness)] = v
			}

			for _, two := range one.CommunityVotes {

				vo := VoteScore{
					Witness: two.Witness,
					Addr:    two.Addr,
					Scores:  two.Scores,
					Vote:    two.Vote,
				}
				vsVOone := VoteScoreRewadrVO{
					VoteScore:  vo,
					LightVotes: make([]VoteScore, 0),
					Reward:     0,
				}

				for _, one := range *reward.GetVout() {
					if bytes.Equal(one.Address, *vsVOone.VoteScore.Addr) {
						vsVOone.Reward = one.Value
						break
					}
				}

				v, ok := m[utils.Bytes2string(*two.Addr)]
				if ok {
					vsVOone.LightVotes = *v
				}

				groupVO.CommunityVotes = append(groupVO.CommunityVotes, vsVOone)
			}
		}
		bvVO.Group = append(bvVO.Group, *groupVO)

	}

	return bvVO, nil

}

type RewardTotal struct {
	CommunityReward uint64
	LightReward     uint64
	StartHeight     uint64
	Height          uint64
	IsGrant         bool
	AllLight        uint64
	RewardLight     uint64
	IsNew           bool
}

type RewardTotalDetail struct {
	startHeight uint64
	endHeight   uint64
	RT          *RewardTotal
	RL          *[]sqlite3_db.RewardLight
}

var rewardCountProcessMapLock = new(sync.Mutex)
var rewardCountProcessMap = make(map[string]*RewardTotalDetail)

func GetRewardCount(addr *crypto.AddressCoin, startHeight, endHeight uint64) (*RewardTotal, *[]sqlite3_db.RewardLight, error) {

	currentHeight := forks.LongChain.GetCurrentBlock()
	if endHeight <= 0 || endHeight > currentHeight {
		endHeight = currentHeight
	}

	var rt *RewardTotal
	var rl *[]sqlite3_db.RewardLight
	var err error

	have := false
	rewardCountProcessMapLock.Lock()
	rtd, ok := rewardCountProcessMap[utils.Bytes2string(*addr)]
	if ok {
		have = true
		if rtd != nil {

			if endHeight > rtd.endHeight && endHeight-rtd.endHeight > config.Mining_community_reward_time {
				have = false
			} else {
				rt = rtd.RT
				rl = rtd.RL
			}
		}
	}

	if !have {
		err = config.ERROR_get_reward_count_sync

		rewardCountProcessMap[utils.Bytes2string(*addr)] = nil
		utils.Go(func() {
			rt, rl, err := RewardCountProcess(addr, startHeight, endHeight)
			if err != nil {
				engine.Log.Info("RewardCountProcess error:%s", err.Error())
				return
			}
			rtd := RewardTotalDetail{
				startHeight: startHeight,
				endHeight:   endHeight,
				RT:          rt,
				RL:          rl,
			}
			rewardCountProcessMapLock.Lock()
			rewardCountProcessMap[utils.Bytes2string(*addr)] = &rtd
			rewardCountProcessMapLock.Unlock()
		})
	} else {

	}
	rewardCountProcessMapLock.Unlock()
	return rt, rl, err
}

func CleanRewardCountProcessMap(addr *crypto.AddressCoin) {
	rewardCountProcessMapLock.Lock()
	delete(rewardCountProcessMap, utils.Bytes2string(*addr))
	rewardCountProcessMapLock.Unlock()
}

func RewardCountProcess(addr *crypto.AddressCoin, startHeight, endHeight uint64) (*RewardTotal, *[]sqlite3_db.RewardLight, error) {

	rewardDBs := make([]sqlite3_db.RewardLight, 0)

	bvvo, err := FindLightVote(startHeight, endHeight)
	if err != nil {

		return nil, nil, err
	}

	allReward := uint64(0)
	light := uint64(0)
	for _, one := range bvvo.Group {
		for _, voteOne := range one.CommunityVotes {

			if bytes.Equal(*voteOne.Addr, *addr) {
				allReward += voteOne.Reward

				for _, two := range voteOne.LightVotes {
					temp := new(big.Int).Mul(big.NewInt(int64(voteOne.Reward)), big.NewInt(int64(two.Vote)))
					value := new(big.Int).Div(temp, big.NewInt(int64(voteOne.Vote)))

					reward := value.Uint64()
					reward = reward - (reward / 10)
					light += reward
					r := sqlite3_db.RewardLight{
						Addr:         *two.Addr,
						Reward:       reward,
						Distribution: 0,
					}

					rewardDBs = append(rewardDBs, r)
				}
				break
			}
		}
	}
	community := allReward - light

	reward := RewardTotal{
		CommunityReward: community,
		LightReward:     light,
		StartHeight:     startHeight,
		Height:          bvvo.EndHeight,
		IsGrant:         (bvvo.EndHeight - startHeight) > (config.Mining_community_reward_time / config.Mining_block_time),
		AllLight:        0,
		RewardLight:     0,
	}

	voutMap := make(map[string]*sqlite3_db.RewardLight)
	for i, _ := range rewardDBs {
		one := rewardDBs[i]

		if one.Reward == 0 {
			continue
		}
		v, ok := voutMap[utils.Bytes2string(one.Addr)]
		if ok {
			v.Reward = v.Reward + one.Reward
			continue
		}
		voutMap[utils.Bytes2string(one.Addr)] = &(rewardDBs)[i]
	}
	vouts := make([]sqlite3_db.RewardLight, 0)
	for _, v := range voutMap {

		vouts = append(vouts, *v)
	}

	communityReward := sqlite3_db.RewardLight{
		Addr:         *addr,
		Reward:       community,
		Distribution: 0,
	}
	vouts = append(vouts, communityReward)

	reward.AllLight = uint64(len(vouts))

	return &reward, &vouts, nil
}

func FindNotSendReward(addr *crypto.AddressCoin) (*sqlite3_db.SnapshotReward, *[]sqlite3_db.RewardLight, error) {

	s, err := new(sqlite3_db.SnapshotReward).Find(*addr)
	if err != nil {
		return nil, nil, err
	}
	if s == nil {
		return nil, nil, nil
	} else {
		rewardNotSend, err := new(sqlite3_db.RewardLight).FindNotSend(s.Id)
		if err != nil {
			return nil, nil, err
		}
		return s, rewardNotSend, nil

	}
}

func CreateRewardCount(addr crypto.AddressCoin, rt *RewardTotal, rs []sqlite3_db.RewardLight) error {
	ss := &sqlite3_db.SnapshotReward{
		Addr:        addr,
		StartHeight: rt.StartHeight,
		EndHeight:   rt.Height,
		Reward:      rt.LightReward + rt.CommunityReward,
		LightNum:    uint64(len(rs)),
	}

	err := new(sqlite3_db.SnapshotReward).Add(ss)
	if err != nil {
		return err
	}

	ss, err = new(sqlite3_db.SnapshotReward).Find(addr)
	if err != nil {
		return err
	}

	count := uint64(0)
	for _, one := range rs {
		count++
		one.Sort = count
		one.SnapshotId = ss.Id
		err := new(sqlite3_db.RewardLight).Add(&one)
		if err != nil {

			return err
		}
	}
	return nil
}

func DistributionReward(addr *crypto.AddressCoin, notSend *[]sqlite3_db.RewardLight, gas uint64, pwd string, cs *CommunitySign, startHeight, endHeight, currentHeight uint64) error {

	if notSend == nil || len(*notSend) <= 0 {
		return nil
	}

	if len(*notSend) > config.Wallet_community_reward_max {
		temp := (*notSend)[:config.Wallet_community_reward_max]
		notSend = &temp

	}

	value := new(big.Int).Div(big.NewInt(int64(gas)), big.NewInt(int64(len(*notSend)))).Uint64()

	payNum := make([]PayNumber, 0)
	for i := 0; i < len(*notSend); i++ {
		one := (*notSend)[i]
		addr := crypto.AddressCoin(one.Addr)

		pns := LinearRelease180DayForLight(addr, one.Reward-value, currentHeight)
		payNum = append(payNum, pns...)
	}
	tx, err := SendToMoreAddressByPayload(addr, payNum, gas, pwd, cs, startHeight, endHeight)
	if err != nil {

		engine.Log.Error(err.Error())
		return err
	} else {

	}

	for i := 0; i < len(*notSend); i++ {
		one := (*notSend)[i]
		one.Txid = *tx.GetHash()
		one.LockHeight = tx.GetLockHeight()
		err := one.UpdateTxid(one.Id)
		if err != nil {
			engine.Log.Error(err.Error())
		}
	}
	return err
}

func LinearRelease180DayForLight(addr crypto.AddressCoin, total uint64, height uint64) []PayNumber {

	pns := make([]PayNumber, 0)

	first25 := new(big.Int).Div(big.NewInt(int64(total)), big.NewInt(int64(4)))

	surplus := new(big.Int).Sub(big.NewInt(int64(total)), first25)

	pnOne := PayNumber{
		Address: addr,
		Amount:  first25.Uint64(),
	}

	pns = append(pns, pnOne)

	dayOne := new(big.Int).Div(surplus, big.NewInt(int64(18))).Uint64()

	intervalHeight := 60 * 60 * 24 * 10 / 10

	totalUse := uint64(0)
	for i := 0; i < 18; i++ {
		pnOne := PayNumber{
			Address:      addr,
			Amount:       dayOne,
			FrozenHeight: height + uint64((i+1)*intervalHeight),
		}
		pns = append(pns, pnOne)
		totalUse = totalUse + dayOne
	}

	if totalUse < surplus.Uint64() {

		pns[len(pns)-1].Amount = pns[len(pns)-1].Amount + (surplus.Uint64() - totalUse)
	}
	return pns
}
