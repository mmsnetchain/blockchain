package mining

import (
	"bytes"
	"encoding/hex"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"runtime"
	"sync"
	"time"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

const (
	Mining_Status_Start           = 1
	Mining_Status_WaitMulticas    = 2
	Mining_Status_WaitImportBlock = 3
	Mining_Status_ImportBlock     = 4
)

var MiningStatusLock = new(sync.Mutex)
var MiningStatus = Mining_Status_Start
var BhvoMulticasCache = make(map[string]*BlockHeadVO)

func init() {

}

func SetMiningStatus_() {

}

func AddBlockToCache(bhvo *BlockHeadVO) {

	MiningStatusLock.Lock()
	now := time.Now()
	for k, bhvo := range BhvoMulticasCache {
		if bhvo.BH.Time < now.Unix()-(config.Mining_block_time*10) {
			delete(BhvoMulticasCache, k)
		}
	}

	BhvoMulticasCache[utils.Bytes2string(bhvo.BH.Hash)] = bhvo
	MiningStatusLock.Unlock()
}

func ImportBlockByCache(hash *[]byte) {
	var bhvo, prebhvo *BlockHeadVO

	MiningStatusLock.Lock()
	bhvo, _ = BhvoMulticasCache[utils.Bytes2string(*hash)]
	MiningStatusLock.Unlock()
	if bhvo == nil {
		return
	}

	MiningStatusLock.Lock()
	prebhvo, _ = BhvoMulticasCache[utils.Bytes2string(bhvo.BH.Previousblockhash)]
	MiningStatusLock.Unlock()
	if prebhvo != nil {

		err := forks.AddBlockHead(prebhvo)
		if err == nil {

			go MulticastBlock(*prebhvo)

			MiningStatusLock.Lock()
			delete(BhvoMulticasCache, utils.Bytes2string(prebhvo.BH.Hash))
			MiningStatusLock.Unlock()
		}
	}

	err := forks.AddBlockHead(bhvo)
	if err == nil {

		go MulticastBlock(*bhvo)

		MiningStatusLock.Lock()
		delete(BhvoMulticasCache, utils.Bytes2string(bhvo.BH.Hash))
		MiningStatusLock.Unlock()
	}
}

func (this *Witness) FindUnconfirmedBlock() (*Block, []Block) {

	var preBlock *Block

	isFirst := false

	group := this.Group.SelectionChain(nil)
	if group == nil {
		isFirst = true
	} else {
		isFirst = false

		preBlock = group.Blocks[len(group.Blocks)-1]

	}

	preGroup := this.Group
	var preGroupBlock *Group
	var ok bool
	for {

		ok = false
		preGroup = preGroup.PreGroup
		ok, preGroupBlock = preGroup.CheckBlockGroup(nil)
		if ok {

			if isFirst {

				preBlock = preGroupBlock.Blocks[len(preGroupBlock.Blocks)-1]

			}
			break
		}

	}

	blocks := make([]Block, 0)
	if preGroup.Height != this.Group.Height {
		for _, one := range preGroupBlock.Blocks {
			blocks = append(blocks, *one)
		}
	}
	if group != nil {
		for _, one := range group.Blocks {
			blocks = append(blocks, *one)
		}
	}
	return preBlock, blocks
}

func (this *Witness) CheckUnconfirmedBlock(blockhash *[]byte) (*Block, []Block) {

	var preBlock *Block

	isFirst := false

	group := this.Group.SelectionChain(blockhash)

	if group == nil {
		isFirst = true
	} else {
		isFirst = false

		preBlock = group.Blocks[len(group.Blocks)-1]

	}

	preGroup := this.Group
	var preGroupBlock *Group
	var ok bool
	for {

		ok = false
		preGroup = preGroup.PreGroup

		ok, preGroupBlock = preGroup.CheckBlockGroup(blockhash)

		if ok {

			if isFirst {

				preBlock = preGroupBlock.Blocks[len(preGroupBlock.Blocks)-1]

			}
			break
		}

	}

	preWitness := this.PreWitness
	for {
		if preWitness == nil {
			break
		}
		if preWitness.Block == nil {
			preWitness = preWitness.PreWitness
			continue
		}

		if bytes.Equal(preWitness.Block.Id, *blockhash) {

			preBlock = preWitness.Block
			break
		}
		preWitness = preWitness.PreWitness

	}

	blocks := make([]Block, 0)
	if preGroup.Height != this.Group.Height {
		for _, one := range preGroupBlock.Blocks {
			blocks = append(blocks, *one)
		}
	}
	if group != nil {
		for _, one := range group.Blocks {
			blocks = append(blocks, *one)
		}
	}
	return preBlock, blocks
}

func (this *Witness) BuildBlock() {

	addrInfo := keystore.GetCoinbase()

	if !bytes.Equal(*this.Addr, addrInfo.Addr) {
		return
	}

	if !GetLongChain().SyncBlockFinish {
		return
	}

	engine.Log.Info("=== start building blocks === group height:%d", this.Group.Height)

	preBlock, blocks := this.FindUnconfirmedBlock()

	tx := make([]TxItr, 0)
	txids := make([][]byte, 0)

	var reward *Tx_reward

	if this.WitnessBackupGroup != preBlock.witness.WitnessBackupGroup {

		reward = preBlock.witness.WitnessBackupGroup.CountRewardToWitnessGroup(preBlock.Height+1, blocks, preBlock)
		tx = append(tx, reward)

		txids = append(txids, reward.Hash)
	}

	chain := forks.GetLongChain()

	txs, ids := chain.transactionManager.Package(reward, preBlock.Height+1, blocks, this.CreateBlockTime)
	tx = append(tx, txs...)
	txids = append(txids, ids...)

	coinbase := keystore.GetCoinbase()

	var bh *BlockHead
	now := utils.GetNow()
	for i := int64(0); i < (config.Mining_block_time*2)-1; i++ {

		bh = &BlockHead{
			Height:            preBlock.Height + 1,
			GroupHeight:       this.Group.Height,
			Previousblockhash: preBlock.Id,
			NTx:               uint64(len(tx)),
			Tx:                txids,
			Time:              now + i,
			Witness:           coinbase.Addr,
		}

		bh.BuildMerkleRoot()
		bh.BuildSign(coinbase.Addr)
		bh.BuildBlockHash()
		if ok, _ := bh.CheckHashExist(); ok {
			bh = nil
			continue
		} else {
			break
		}
	}
	if bh == nil {
		engine.Log.Info("Block out failed, all hash have collisions")

		return
	}

	bhvo := CreateBlockHeadVO(config.StartBlockHash, bh, tx)

	engine.Log.Info("=== build block Success === group height:%d block height:%d", bhvo.BH.GroupHeight, bhvo.BH.Height)
	engine.Log.Info("=== build block Success === Block hash %s", hex.EncodeToString(bhvo.BH.Hash))
	engine.Log.Info("=== build block Success === pre Block hash %s", hex.EncodeToString(bhvo.BH.Previousblockhash))

	bhvo.FromBroadcast = true

	UniformityMulticastBlock(bhvo)

}

func MulticastBlock(bhVO BlockHeadVO) {
	goroutineId := utils.GetRandomDomain() + utils.TimeFormatToNanosecondStr()
	_, file, line, _ := runtime.Caller(0)
	engine.AddRuntime(file, line, goroutineId)
	defer engine.DelRuntime(file, line, goroutineId)
	bs, err := bhVO.Proto()
	if err != nil {
		return
	}
	message_center.SendMulticastMsg(config.MSGID_multicast_blockhead, bs)
}

func MulticastBlockAndImport(bhVO *BlockHeadVO) error {

	bs, err := bhVO.Proto()
	if err != nil {
		return err
	}

	head := message_center.NewMessageHead(nil, nil, false)
	body := message_center.NewMessageBody(config.MSGID_multicast_witness_blockhead, bs, 0, nil, 0)
	message := message_center.NewMessage(head, body)
	message.BuildHash()

	err = new(sqlite3_db.MessageCache).Add(message.Body.Hash, head.Proto(), body.Proto())
	if err != nil {
		engine.Log.Error(err.Error())
		return err
	}
	engine.Log.Info("multicast message hash:%s", hex.EncodeToString(message.Body.Hash))
	return MulticastBlockSync(message)
}

func MulticastBlockSync(message *message_center.Message) error {
	whiltlistNodes := nodeStore.GetWhiltListNodes()
	return message_center.BroadcastsAll(1, config.MSGID_multicast_witness_blockhead, whiltlistNodes, nil, nil, &message.Body.Hash)
}

func UniformityMulticastBlock(bhVO *BlockHeadVO) {

	bs, err := bhVO.Proto()
	if err != nil {
		return
	}

	head := message_center.NewMessageHead(nil, nil, false)
	body := message_center.NewMessageBody(config.MSGID_multicast_witness_blockhead, bs, 0, nil, 0)
	message := message_center.NewMessage(head, body)
	message.BuildHash()

	err = new(sqlite3_db.MessageCache).Add(message.Body.Hash, head.Proto(), body.Proto())
	if err != nil {
		engine.Log.Error("save message cache error:%s", err.Error())
		return
	}

	AddBlockToCache(bhVO)

	err = UniformityBroadcasts(&message.Body.Hash, config.MSGID_uniformity_multicast_witness_blockhead, config.CLASS_uniformity_witness_multicas_blockhead, config.Wallet_sync_block_timeout)
	if err != nil {
		engine.Log.Info("multcas blockhash error:%s", err.Error())
		return
	}

	err = UniformityBroadcasts(&bhVO.BH.Hash, config.MSGID_uniformity_multicast_witness_block_import, config.CLASS_uniformity_witness_multicas_block_import, config.Wallet_sync_block_timeout)
	if err != nil {
		engine.Log.Info("multcas import block error:%s", err.Error())
		return
	}
	engine.Log.Info("start import block")
	ImportBlockByCache(&bhVO.BH.Hash)

}

func UniformityBroadcasts(hash *[]byte, msgid uint64, waitRequestClass string, timeout int64) error {

	whiltlistNodes := nodeStore.GetWhiltListNodes()

	allNodes := make(map[string]bool)

	var timeouterrorlock = new(sync.Mutex)
	var timeouterror error

	cs := make(chan bool, config.CPUNUM)
	group := new(sync.WaitGroup)
	for i, _ := range whiltlistNodes {
		sessionid := whiltlistNodes[i]

		if bytes.Equal(nodeStore.NodeSelf.IdInfo.Id, sessionid) {
			continue
		}
		_, ok := allNodes[utils.Bytes2string(sessionid)]
		if ok {

			continue
		}
		allNodes[utils.Bytes2string(sessionid)] = false
		cs <- false
		group.Add(1)

		utils.Go(func() {
			success := false
			_, err := message_center.SendNeighborWithReplyMsg(msgid, &sessionid, hash, waitRequestClass, timeout)
			if err == nil {
				success = true
			} else {
				if err.Error() == config.ERROR_wait_msg_timeout.Error() {
				} else {

					success = true
				}
			}
			if !success {
				timeouterrorlock.Lock()
				timeouterror = config.ERROR_wait_msg_timeout
				timeouterrorlock.Lock()
			}
			<-cs
			group.Done()
		})
	}
	group.Wait()

	return timeouterror
}
