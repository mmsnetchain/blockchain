package mining

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	mc "github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/message_center/flood"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"github.com/shirou/gopsutil/v3/mem"
)

var syncSaveBlockHead = make(chan *BlockHeadVO, 1)

var syncForNeighborChain = make(chan bool, 1)

func (this *Chain) FirstDownloadBlock() error {

	engine.Log.Info("After the node is started, the block is synchronized for the first time %v", config.LoadNode)

	greenChannel := false
	var err error
	count := 0
	this.SyncBlockLock.Lock()
	this.LoadBlockChain()
	oldCurrentBlockHeight := uint64(0)
	total := 0
	for {
		count++

		engine.Log.Info("Which round of synchronization block %d", count)
		currentBlock := GetLongChain().GetCurrentBlock()
		if greenChannel {
			greenChannel = false
		} else {
			if currentBlock > oldCurrentBlockHeight {
				total = 0
				oldCurrentBlockHeight = currentBlock
			} else {

				total++
				if total > 3 {

					engine.Log.Info("FirstDownloadBlock fail")

					peerBlockinfo, _ := FindRemoteCurrentHeight()
					bh, err := FindLastBlockForNeighbor(peerBlockinfo)
					if err != nil || bh == nil {
						break
					}

					if bh.Height <= currentBlock {
						break
					}
					engine.Log.Error("ForkCheck")
					this.ForkCheck(&bh.Hash)
					break
				}
			}
		}
		peerBlockinfo, remoteCuurentHeightMax := FindRemoteCurrentHeight()
		if currentBlock >= remoteCuurentHeightMax {
			engine.Log.Info("FirstDownloadBlock finish")
			break
		}
		err = this.SyncBlockHead(peerBlockinfo)

		if err == nil {
			continue
		} else {
			engine.Log.Info("FirstDownloadBlock error:%s", err.Error())
			if err.Error() == config.ERROR_chain_sync_block_timeout.Error() {
				time.Sleep(time.Second * 60)
				greenChannel = true
				continue
			}
		}
		time.Sleep(time.Second * 11)
	}
	this.SyncBlockLock.Unlock()

	FinishFirstLoadBlockChain()

	return err
}
func (this *Chain) GoSyncBlock() {

	utils.Go(func() {

		beforBlockHeight := uint64(0)
		stopBlockHeightTotal := 0
		greenChannel := false

		for range syncForNeighborChain {

			if atomic.LoadUint32(&this.StopSyncBlock) == 1 {

				engine.Log.Info("stop sync queue !!!")
				continue
			}
			currentHeight := GetLongChain().GetCurrentBlock()
			if greenChannel {
				greenChannel = false
			} else {
				if currentHeight <= beforBlockHeight {
					time.Sleep(time.Second * 6)
					stopBlockHeightTotal++
				} else {
					stopBlockHeightTotal = 0
					beforBlockHeight = currentHeight
				}
				if stopBlockHeightTotal >= 5 {
					engine.Log.Error("ForkCheck")
					this.ForkCheck(nil)
					continue
				}
			}

			_, _, isKickOut, _, _ := GetWitnessStatus()

			if isKickOut {
				engine.Log.Info("isKickOut")
				continue
			}

			this.SyncBlockLock.Lock()

			this.LoadBlockChain()
			peerBlockinfo, _ := FindRemoteCurrentHeight()

			peers := peerBlockinfo.Sort()
			if len(peers) <= 0 || peers[0].CurrentHeight <= currentHeight {
				time.Sleep(time.Second * 60)
				greenChannel = true
				this.SyncBlockLock.Unlock()
				continue
			}

			err := this.SyncBlockHead(peerBlockinfo)
			if err != nil && err.Error() == config.ERROR_chain_sync_block_timeout.Error() {
				time.Sleep(time.Second * 60)
				greenChannel = true
				this.SyncBlockLock.Unlock()
				continue
			}

			this.WitnessChain.StopAllMining()
			this.WitnessChain.BuildMiningTime()

			this.SyncBlockLock.Unlock()

		}
	})
}

func (this *Chain) NoticeLoadBlockForDB() {

	select {
	case syncForNeighborChain <- false:
		engine.Log.Info("Put in sync queue")
	default:
	}
}

func (this *Chain) ForkCheck(bhash *[]byte) {
	if atomic.LoadUint32(&this.StopSyncBlock) == 1 {
		return
	}
	atomic.StoreUint32(&this.StopSyncBlock, 1)

	engine.Log.Error("ForkCheck")

	if bhash == nil && this.Temp == nil {
		return
	}

	var preBlockHash []byte
	if bhash != nil {
		preBlockHash = *bhash
	} else {
		preBlockHash = this.Temp.BH.Previousblockhash
	}
	var nextBlockHash []byte
	for {
		engine.Log.Info("load block:%s", hex.EncodeToString(preBlockHash))

		bh, err := LoadBlockHeadByHash(&preBlockHash)
		if err != nil || bh == nil {

			peerBlockinfo, _ := FindRemoteCurrentHeight()
			bhvo, err := SyncBlockFlashDB(&preBlockHash, peerBlockinfo)
			if err != nil {
				engine.Log.Error("SyncBlockHead error:%s", err.Error())
				return
			}
			bh = bhvo.BH
		}

		engine.Log.Info("check fork height:%d", bh.Height)

		bhash, err := db.LevelDB.Find([]byte(config.BlockHeight + strconv.Itoa(int(bh.Height))))
		if err == nil && bytes.Equal(*bhash, bh.Hash) {
			engine.Log.Info("fork block:%s", hex.EncodeToString(*bhash))

			if nextBlockHash != nil && len(nextBlockHash) > 0 {
				bh.Nextblockhash = nextBlockHash
				bs, err := bh.Proto()
				err = db.LevelDB.Save(bh.Hash, bs)
				if err != nil {
					engine.Log.Error("save block error:%s", err.Error())
					return
				}
			}

			break
		}

		preBlockHash = bh.Previousblockhash
		nextBlockHash = bh.Hash

	}

	engine.Log.Error("ForkCheck finish!")
}

func SaveBlockHead(bhvo *BlockHeadVO) (bool, error) {
	bhvo.BH.BuildBlockHash()
	ok := true

	for i, _ := range bhvo.Txs {
		one := bhvo.Txs[i]
		one.BuildHash()

		TxCache.FlashTxInCache(*one.GetHash(), one)

		err := SaveTempTx(one)
		if err != nil {
			return false, err
		}
	}

	_, err := db.LevelDB.Find(bhvo.BH.Previousblockhash)

	if err != nil {

		ok = false
	}

	bs, err := bhvo.BH.Proto()
	if err != nil {

		return false, err
	}

	TxCache.AddBlockHeadCache(bhvo.BH.Hash, bhvo.BH)

	if ok, _ := db.LevelDB.CheckHashExist(bhvo.BH.Hash); !ok {

		err = db.LevelDB.Save(bhvo.BH.Hash, bs)
		if err != nil {

			return false, err
		}
	}
	return ok, nil
}

func SaveTxToBlockHead(bhvo *BlockHeadVO) error {
	bhvo.BH.BuildBlockHash()
	var err error

	for i, _ := range bhvo.Txs {
		one := bhvo.Txs[i]
		one.BuildHash()
		err = db.SaveTxToBlockHash(one.GetHash(), &bhvo.BH.Hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveTempTx(txItr TxItr) error {

	bs, err := txItr.Proto()
	if err != nil {

		return err
	}

	err = db.LevelDB.Save(*txItr.GetHash(), bs)
	if err != nil {

		return err
	}

	db.LevelDB.Save(config.BuildTxNotImport(*txItr.GetHash()), nil)
	return nil
}

func FindBlockHeight() {
	goroutineId := utils.GetRandomDomain() + utils.TimeFormatToNanosecondStr()
	_, file, line, _ := runtime.Caller(0)
	engine.AddRuntime(file, line, goroutineId)
	defer engine.DelRuntime(file, line, goroutineId)
	syncHeightBlock := new(sync.Map)

	for _, key := range nodeStore.GetLogicNodes() {
		sessionName := ""
		session, ok := engine.GetSession(utils.Bytes2string(key))
		if ok {
			sessionName = session.GetName()
		}
		message, err := mc.SendNeighborMsg(config.MSGID_heightBlock, &key, nil)
		if err == nil {

			bs, _ := flood.WaitRequest(mc.CLASS_findHeightBlock, utils.Bytes2string(message.Body.Hash), 0)

			if bs == nil {

				continue
			}
			chain := forks.GetLongChain()

			heightBlock := binary.LittleEndian.Uint64((*bs)[8:])

			if chain.GetCurrentBlock() > heightBlock {
				continue
			}

			if GetHighestBlock() < heightBlock {
				SetHighestBlock(heightBlock)
			}

			syncHeightBlock.Store(sessionName, heightBlock)
		}

	}

	count := 0
	syncHeightBlock.Range(func(key, value interface{}) bool {
		count++
		return false
	})
	if count <= 0 {
		return
	}

	heightBlockVote := new(sync.Map)
	syncHeightBlock.Range(func(key, value interface{}) bool {

		height := value.(uint64)
		v, ok := heightBlockVote.Load(height)
		if ok {
			total := v.(uint64)
			heightBlockVote.Store(height, uint64(total+1))
		} else {
			heightBlockVote.Store(height, uint64(1))
		}
		return true
	})

	heightBlockMax := uint64(0)
	heightBlock := uint64(0)
	heightTotal := uint64(0)

	heightBlockVote.Range(func(k, v interface{}) bool {

		height := k.(uint64)
		if height == 0 {
			return true
		}
		if height > heightBlockMax {
			heightBlockMax = height
		}
		total := v.(uint64)
		if total > heightTotal {
			heightTotal = total
			heightBlock = height
		} else if total == heightTotal {

		}
		return true
	})

	SetHighestBlock(heightBlock)

}

func FindRemoteCurrentHeight() (*PeerBlockInfoDESC, uint64) {
	remoteCuurentHeightMax := uint64(0)

	peers := make([]*PeerBlockInfo, 0)

	logicNodes := nodeStore.GetLogicNodes()
	logicNodes = append(logicNodes, nodeStore.GetNodesClient()...)

	for i, _ := range logicNodes {
		key := logicNodes[i]

		message, err := mc.SendNeighborMsg(config.MSGID_heightBlock, &key, nil)
		if err != nil {

			continue
		}

		bs, _ := flood.WaitRequest(mc.CLASS_findHeightBlock, utils.Bytes2string(message.Body.Hash), 5)

		if bs == nil {

			continue
		}

		heightBlock := binary.LittleEndian.Uint64((*bs)[8:])

		if remoteCuurentHeightMax > heightBlock {

			continue
		}
		remoteCuurentHeightMax = heightBlock

		peers = append(peers, &PeerBlockInfo{
			Addr:          &key,
			CurrentHeight: heightBlock,
		})

	}

	peersDESC := NewPeerBlockInfoDESC(peers)
	return peersDESC, remoteCuurentHeightMax

}

func (this *Chain) SyncBlockHead(peerBlockInfo *PeerBlockInfoDESC) error {
	engine.Log.Info("Start synchronizing blocks from neighbor nodes")

	chain := forks.GetLongChain()
	_, block := chain.GetLastBlock()
	bhash := block.Id

	bhvo, err := SyncBlockFlashDB(&bhash, peerBlockInfo)
	if err != nil {
		engine.Log.Error("SyncBlockHead error:%s", err.Error())
		return err
	}

	if bhvo == nil {
		engine.Log.Error("SyncBlockHead finish")
		return nil
	}

	bhvo.BH.BuildBlockHash()
	engine.Log.Info("Print blocks synchronized to %d %s", bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Hash))

	tiker := time.NewTicker(time.Minute)

	bhvo, err = this.deepCycleSyncBlock(bhvo, tiker.C, bhvo.BH.Height+1, peerBlockInfo)
	tiker.Stop()
	if err != nil {
		engine.Log.Error("SyncBlockHead error:%s", err.Error())
		return err
	}

	this.Balance.Unfrozen(bhvo.BH.Height-1, bhvo.BH.Time)

	_, block = this.GetLastBlock()
	bhvos, err := GetUnconfirmedBlockForNeighbor(block.Height, peerBlockInfo)
	if err != nil {
		engine.Log.Error("GetUnconfirmedBlockForNeighbor error:%s", err.Error())
		return err
	}
	for _, one := range bhvos {
		engine.Log.Info("Import GetUnconfirmedBlockForNeighbor height:%d", one.BH.Height)
		one.FromBroadcast = false
		this.AddBlock(one)
	}

	this.SyncBlockFinish = true
	engine.Log.Info("Sync block complete")
	return nil
}

func (this *Chain) deepCycleSyncBlock(bhvo *BlockHeadVO, c <-chan time.Time, height uint64, peerBlockInfo *PeerBlockInfoDESC) (*BlockHeadVO, error) {

	memInfo, _ := mem.VirtualMemory()
	if memInfo.UsedPercent > config.Wallet_Memory_percentage_max {
		runtime.GC()
		time.Sleep(time.Second)
	}

	bhash := &bhvo.BH.Nextblockhash

	if bhash == nil || len(*bhash) <= 0 {

		engine.Log.Warn("The next block hash of the query is empty")
		return bhvo, nil
	}

	bh, txItrs, err := this.syncBlockForDBAndNeighbor(bhash, peerBlockInfo)
	if err != nil {
		engine.Log.Info("Error synchronizing block: %s", err.Error())
		return bhvo, err
	}

	bhvo = &BlockHeadVO{FromBroadcast: false, BH: bh, Txs: txItrs}
	if err = this.AddBlock(bhvo); err != nil {
		if err.Error() == ERROR_repeat_import_block.Error() {

		} else {
			engine.Log.Info("deepCycleSyncBlock error: %s", err.Error())
			return bhvo, err
		}
	}

	if GetHighestBlock() < bh.Height {
		SetHighestBlock(bh.Height)
	}

	select {
	case <-c:

		utils.Go(FindBlockHeight)
	default:
	}

	engine.Log.Info("Next block %s", hex.EncodeToString(bh.Nextblockhash))

	return this.deepCycleSyncBlock(bhvo, c, bh.Height+1, peerBlockInfo)
}

func (this *Chain) syncBlockForDBAndNeighbor(bhash *[]byte, peerBlockInfo *PeerBlockInfoDESC) (*BlockHead, []TxItr, error) {

	bhvo, err := FindBlockForNeighbor(bhash, peerBlockInfo)
	if err != nil {
		engine.Log.Error("find next block error:%s", err.Error())
		return nil, nil, err
	}
	if bhvo == nil {

		engine.Log.Error("find next block fail")
		return nil, nil, config.ERROR_chain_sysn_block_fail
	}

	for i, _ := range bhvo.Txs {
		bhvo.Txs[i].BuildHash()

		bs, err := bhvo.Txs[i].Proto()
		if err != nil {

			engine.Log.Error("load tx error:%s", err.Error())
			return nil, nil, err
		}

		db.LevelDB.Save(*bhvo.Txs[i].GetHash(), bs)

	}

	if this.GetStartingBlock() > config.Mining_block_start_height {

		bh, err := LoadBlockHeadByHash(&bhvo.BH.Previousblockhash)
		if err != nil {
			engine.Log.Error("load blockhead error:%s", err.Error())
			return nil, nil, err
		}

		bh.Nextblockhash = bhvo.BH.Hash

		bs, err := bh.Proto()
		if err != nil {

			engine.Log.Error("parse blockhead error:%s", err.Error())
			return nil, nil, err
		}
		db.LevelDB.Save(bh.Hash, bs)
	}

	bs, err := bhvo.BH.Proto()
	if err != nil {

		engine.Log.Error("parse blockhead error:%s", err.Error())
		return nil, nil, err
	}

	db.LevelDB.Save(bhvo.BH.Hash, bs)

	return bhvo.BH, bhvo.Txs, nil
}

func SyncBlockFlashDB(bhash *[]byte, peerBlockInfo *PeerBlockInfoDESC) (*BlockHeadVO, error) {

	bhvo, err := FindBlockForNeighbor(bhash, peerBlockInfo)
	if err != nil {
		return nil, err
	}
	if bhvo == nil {
		return nil, config.ERROR_chain_sysn_block_fail
	}
	bhvo.BH.BuildBlockHash()

	bs, err := bhvo.BH.Proto()
	if err != nil {
		return nil, err
	}

	db.LevelDB.Save(*bhash, bs)
	for _, one := range bhvo.Txs {

		bs, err := one.Proto()
		if err != nil {
			return nil, err
		}
		db.LevelDB.Save(*one.GetHash(), bs)
	}
	return bhvo, nil
}

func GetRemoteTxAndSave(txid []byte) TxItr {
	bs := GetRemoteKeyFlashLocalKey(txid)

	txItr, err := ParseTxBaseProto(ParseTxClass(txid), bs)
	if err != nil {
		return nil
	}
	return txItr
}

func GetRemoteKeyFlashLocalKey(txid []byte) *[]byte {

	var bs *[]byte

	logicNodes := nodeStore.GetLogicNodes()
	logicNodes = OrderNodeAddr(logicNodes)
	for _, key := range logicNodes {

		message, _ := message_center.SendNeighborMsg(config.MSGID_getDBKey_one, &key, &txid)

		bs, _ = flood.WaitRequest(mc.CLASS_getTransaction, utils.Bytes2string(message.Body.Hash), config.Mining_block_time)

		if bs == nil {

			engine.Log.Error("Receive message timeout %s", key.B58String())

			continue
		}

		db.LevelDB.Save(txid, bs)
		break
	}
	return bs
}

type PeerBlockInfo struct {
	Addr          *nodeStore.AddressNet
	CurrentHeight uint64
}

type PeerBlockInfoDESC struct {
	Peers []*PeerBlockInfo
}

func (this PeerBlockInfoDESC) Len() int {
	return len(this.Peers)
}

func (this PeerBlockInfoDESC) Less(i, j int) bool {
	if this.Peers[i].CurrentHeight < this.Peers[j].CurrentHeight {
		return false
	} else {
		return true
	}
}

func (this PeerBlockInfoDESC) Swap(i, j int) {
	this.Peers[i], this.Peers[j] = this.Peers[j], this.Peers[i]
}

func (this PeerBlockInfoDESC) Sort() []*PeerBlockInfo {
	if len(this.Peers) <= 0 {
		return nil
	}
	sort.Sort(this)

	height := uint64(0)
	for i, one := range this.Peers {
		if i == 0 {
			height = one.CurrentHeight
			continue
		}
		engine.Log.Info(":%d %d", one.CurrentHeight, height)
		if one.CurrentHeight < height {
			engine.Log.Info(":%d", i)
			return this.Peers[:i]
		}
	}
	return this.Peers
}

func NewPeerBlockInfoDESC(peers []*PeerBlockInfo) *PeerBlockInfoDESC {
	return &PeerBlockInfoDESC{
		Peers: peers,
	}
}
