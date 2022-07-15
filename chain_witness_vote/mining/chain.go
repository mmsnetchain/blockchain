package mining

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

func FindSurplus(addr utils.Multihash) uint64 {
	return 0
}

func FindLastGroupMiner() []utils.Multihash {
	return []utils.Multihash{}
}

type Chain struct {
	No                 uint64
	StartingBlock      uint64
	StartBlockTime     uint64
	CurrentBlock       uint64
	PulledStates       uint64
	HighestBlock       uint64
	SyncBlockLock      *sync.RWMutex
	SyncBlockFinish    bool
	WitnessBackup      *WitnessBackup
	WitnessChain       *WitnessChain
	Balance            *BalanceManager
	transactionManager *TransactionManager

	StopSyncBlock uint32
	Temp          *BlockHeadVO
}

func (this *Chain) GetStartingBlock() uint64 {
	return atomic.LoadUint64(&this.StartingBlock)
}

func (this *Chain) GetStartBlockTime() uint64 {
	return atomic.LoadUint64(&this.StartBlockTime)
}

func (this *Chain) SetStartingBlock(n, startBlockTime uint64) {
	atomic.StoreUint64(&this.StartingBlock, n)
	atomic.StoreUint64(&this.StartBlockTime, startBlockTime)
}

func (this *Chain) GetCurrentBlock() uint64 {
	return atomic.LoadUint64(&this.CurrentBlock)
}

func (this *Chain) SetCurrentBlock(n uint64) {

	atomic.StoreUint64(&this.CurrentBlock, n)
}

func (this *Chain) GetPulledStates() uint64 {
	return atomic.LoadUint64(&this.PulledStates)
}

func (this *Chain) SetPulledStates(n uint64) {

	atomic.StoreUint64(&this.PulledStates, n)
}

func NewChain() *Chain {
	chain := &Chain{}
	chain.StopSyncBlock = 0
	chain.SyncBlockLock = new(sync.RWMutex)
	chain.SyncBlockFinish = false

	wb := NewWitnessBackup(chain)
	wc := NewWitnessChain(wb, chain)
	tm := NewTransactionManager(wb)
	b := NewBalanceManager(wb, tm, chain)

	chain.WitnessBackup = wb
	chain.WitnessChain = wc
	chain.Balance = b
	chain.transactionManager = tm

	utils.Go(chain.GoSyncBlock)
	return chain
}

type Group struct {
	PreGroup  *Group
	NextGroup *Group
	Height    uint64
	Blocks    []*Block
}

type Block struct {
	Id         []byte
	PreBlockID []byte
	PreBlock   *Block
	NextBlock  *Block
	Group      *Group
	Height     uint64
	witness    *Witness
}

func (this *Block) Load() (*BlockHead, error) {

	blockHead, ok := TxCache.FindBlockHeadCache(this.Id)
	if !ok {

		var err error
		blockHead, err = LoadBlockHeadByHash(&this.Id)
		if err != nil {
			return nil, err
		}
		TxCache.AddBlockHeadCache(this.Id, blockHead)
	}

	return blockHead, nil
}

func (this *Block) LoadTxs() (*BlockHead, *[]TxItr, error) {
	bh, err := this.Load()
	if err != nil {

		return nil, nil, err
	}
	txs := make([]TxItr, 0)
	for i, one := range bh.Tx {

		var txItr TxItr
		ok := false

		if config.EnableCache {

			txItr, ok = TxCache.FindTxInCache(bh.Tx[i])
		}
		if !ok {

			var err error
			txItr, err = LoadTxBase(one)
			if err != nil {
				return nil, nil, err
			}

			if config.EnableCache {
				TxCache.AddTxInCache(bh.Tx[i], txItr)
			}
		}

		txs = append(txs, txItr)
	}
	return bh, &txs, nil
}

var AddBlockLock = new(sync.RWMutex)

func (this *Chain) AddBlock(bhvo *BlockHeadVO) error {
	goroutineId := utils.GetRandomDomain() + utils.TimeFormatToNanosecondStr()
	_, file, line, _ := runtime.Caller(0)
	engine.AddRuntime(file, line, goroutineId)
	defer engine.DelRuntime(file, line, goroutineId)

	start := time.Now()

	AddBlockLock.Lock()
	defer AddBlockLock.Unlock()

	if config.DBUG_import_height_max != 0 && bhvo.BH.Height > config.DBUG_import_height_max {
		return nil
	}

	engine.Log.Info("====== Import block group:%d block:%d prehash:%s hash:%s witness:%s", bhvo.BH.GroupHeight,
		bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Previousblockhash), hex.EncodeToString(bhvo.BH.Hash), bhvo.BH.Witness.B58String())

	if bhvo.BH.Height < this.GetCurrentBlock() {
		engine.Log.Info("Block height too low")
		return nil
	}

	bhvo.BH.BuildBlockHash()

	ok, err := SaveBlockHead(bhvo)
	if err != nil {
		engine.Log.Warn("save block error %s", err.Error())
		return err
	}

	if !ok {

		this.NoticeLoadBlockForDB()
		return nil
	}

	now := utils.GetNow()
	if bhvo.BH.Time > now+config.Mining_block_time {
		engine.Log.Warn("Build block It's too late %d %d %s", bhvo.BH.Time, now, time.Unix(bhvo.BH.Time, 0).String())

		return errors.New("Build block It's too late")
	}

	if bhvo.BH.Height > forks.GetHighestBlock() {

		forks.SetHighestBlock(bhvo.BH.Height)
		this.Temp = bhvo
	}

	for _, one := range config.Exclude_Tx {
		if bhvo.BH.Height != one.Height {
			continue
		}
		for j, two := range bhvo.Txs {
			if !bytes.Equal(one.TxByte, *two.GetHash()) {

				continue
			}

			notExcludeTx := bhvo.Txs[:j]
			bhvo.Txs = append(notExcludeTx, bhvo.Txs[j+1:]...)

			break
		}
	}

	if this.WitnessChain.CheckRepeatImportBlock(bhvo) {

		engine.Log.Info("Repeat import block")
		return ERROR_repeat_import_block
	}

	preWitness := this.WitnessChain.FindPreWitnessForBlock(bhvo.BH.Previousblockhash)
	if preWitness == nil {

		if this.GetCurrentBlock() > bhvo.BH.Height {
			return nil
		}

		engine.Log.Warn("The front block cannot be found, and the new block is discontinuous with a new height:%d preblockhash:%s v:%+v", bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Previousblockhash), bhvo.BH)

		this.NoticeLoadBlockForDB()
		return ERROR_fork_import_block
	} else {

		engine.Log.Info("pre witness height:%d bhvo height:%d", preWitness.Block.Height, bhvo.BH.Height)
		if preWitness.Block.Height+1 != bhvo.BH.Height {
			engine.Log.Error("new block height fail,pre height:%d new height:%d", preWitness.Block.Height, bhvo.BH.Height)
			return ERROR_import_block_height_not_continuity
		}
	}

	var currentWitness *Witness
	var isOverWitnessGroupChain bool

	if bhvo.BH.GroupHeight != config.Mining_group_start_height {

		currentWitness, isOverWitnessGroupChain = this.WitnessChain.FindWitnessForBlockOnly(bhvo)

		if bhvo.BH.Height < config.FixBuildGroupBUGHeightMax {
			if currentWitness == nil {

				engine.Log.Info("The new height of the witness for this block cannot be found:%d", bhvo.BH.Height)

				this.WitnessChain.BuildBlockGroupForGroupHeight(bhvo.BH.GroupHeight-1, &bhvo.BH.Previousblockhash)

				this.WitnessChain.CompensateWitnessGroupByGroupHeight(bhvo.BH.GroupHeight)

				this.WitnessChain.BuildBlockGroupForGroupHeight(bhvo.BH.GroupHeight-1, &bhvo.BH.Previousblockhash)

			}
			this.WitnessChain.BuildBlockGroup(bhvo, preWitness)
		} else {
			if isOverWitnessGroupChain {
				engine.Log.Info("The new height of the witness for this block cannot be found:%d", bhvo.BH.Height)

				this.WitnessChain.BuildBlockGroupForGroupHeight(bhvo.BH.GroupHeight-1, &bhvo.BH.Previousblockhash)

				this.WitnessChain.CompensateWitnessGroupByGroupHeight(bhvo.BH.GroupHeight)

				this.WitnessChain.BuildBlockGroupForGroupHeight(bhvo.BH.GroupHeight-1, &bhvo.BH.Previousblockhash)

			}
		}

	}

	currentWitness, _ = this.WitnessChain.FindWitnessForBlockOnly(bhvo)
	if currentWitness == nil {

		return errors.New("not font witness")
	}

	ok = this.WitnessChain.SetWitnessBlock(bhvo)
	if !ok {

		engine.Log.Debug("Setting witness block failed")

		return errors.New("Setting witness block failed")
	}

	if bhvo.BH.Height >= config.FixBuildGroupBUGHeightMax {

		if bhvo.BH.GroupHeight == config.Mining_group_start_height {
			this.WitnessChain.witnessGroup.BuildGroup(nil)
			this.WitnessChain.BuildWitnessGroup(false, false)
			this.WitnessChain.witnessGroup = this.WitnessChain.witnessGroup.NextGroup
			this.WitnessChain.BuildWitnessGroup(false, true)
		} else {

			ok, group := currentWitness.Group.CheckBlockGroup(nil)
			if ok {

				witness := this.WitnessChain.FindPreWitnessForBlock(group.Blocks[0].PreBlockID)

				wg := witness.Group.BuildGroup(&witness.Block.Id)

				if wg != nil {

					for _, one := range wg.Witness {
						if !one.CheckIsMining {
							if one.Block == nil {
								this.WitnessBackup.AddBlackList(*one.Addr)
							} else {
								this.WitnessBackup.SubBlackList(*one.Addr)
							}
							one.CheckIsMining = true
						}
					}
				}

				this.CountBlock(witness.Group)

				this.WitnessChain.witnessGroup = currentWitness.Group

				this.WitnessChain.BuildWitnessGroup(false, true)
			}
		}
	}

	engine.Log.Info("Save block Time spent %s", time.Now().Sub(start))

	this.WitnessChain.BuildMiningTime()

	this.WitnessChain.GCWitnessOld()

	if bhvo.BH.Height == config.Witness_backup_group_overheight {
		config.Witness_backup_group = config.Witness_backup_group_new
	}

	return nil

}

func (this *Block) FlashNextblockhash() error {

	bh, err := this.Load()
	if err != nil {
		return err
	}

	bh.Nextblockhash = this.NextBlock.Id

	bs, err := bh.Proto()
	if err != nil {
		return err
	}

	TxCache.FlashBlockHeadCache(this.Id, bh)

	if bh.Nextblockhash == nil {
		engine.Log.Error("save block nextblockhash nil %s", string(*bs))
	}
	err = db.LevelDB.Save(this.Id, bs)
	if err != nil {
		return err
	}

	return nil

}

func (this *Chain) CountBlock(witnessGroup *WitnessGroup) {

	if witnessGroup.BlockGroup == nil {
		return
	}
	if witnessGroup.IsCount {
		return
	}

	for _, one := range witnessGroup.BlockGroup.Blocks {

		if one.Height == config.Mining_block_start_height {
			continue
		}

		bh, txs, err := one.LoadTxs()
		if err != nil {

			continue
		}
		bhvo := &BlockHeadVO{BH: bh, Txs: *txs}

		for _, one := range config.Exclude_Tx {
			if bhvo.BH.Height != one.Height {
				continue
			}
			for j, two := range bhvo.Txs {
				if !bytes.Equal(one.TxByte, *two.GetHash()) {

					continue
				}

				notExcludeTx := bhvo.Txs[:j]
				bhvo.Txs = append(notExcludeTx, bhvo.Txs[j+1:]...)

				break
			}
		}

		one.PreBlock.FlashNextblockhash()

		this.Balance.CountBalanceForBlock(bhvo)

		this.WitnessBackup.CountWitness(&bhvo.Txs)

		this.transactionManager.DelTx(bhvo.Txs)

		this.SetCurrentBlock(one.Height)

		this.transactionManager.CleanIxOvertime(one.Height)

	}
	witnessGroup.IsCount = true

}

func (this *Chain) GetLastBlock() (witness *Witness, block *Block) {

	witnessGroup := this.WitnessChain.witnessGroup
	if witnessGroup == nil {

		return
	}

	if witnessGroup.Height != config.Mining_group_start_height {

		for {
			witnessGroup = witnessGroup.PreGroup
			if witnessGroup == nil {

				break
			}

			if witnessGroup.BlockGroup != nil {

				break
			}
		}
	}

	block = witnessGroup.BlockGroup.Blocks[len(witnessGroup.BlockGroup.Blocks)-1]
	witness = block.witness

	return

}

func (this *Chain) GetBalance() *BalanceManager {
	return this.Balance
}

func (this *Chain) PrintBlockList() {

}

func (this *Chain) HashRandom() *[]byte {
	_, lastBlock := this.GetLastBlock()

	if lastBlock != nil {

		if random, ok := config.RandomMap.Load(utils.Bytes2string(lastBlock.Id)); ok {
			bs := random.(*[]byte)

			return bs
		}
	}

	var preHash *[]byte
	bs := make([]byte, 0)

	for i := 0; lastBlock != nil && i < config.Mining_block_hash_count; i++ {
		if lastBlock.Height < config.NextHashHeightMax {
			if preHash == nil {
				if random, ok := config.NextHash.Load(utils.Bytes2string(lastBlock.Id)); ok {
					bsOne := random.(*[]byte)
					bs = append(bs, *bsOne...)
					preHash = bsOne

					continue
				}
			} else {
				if random, ok := config.NextHash.Load(utils.Bytes2string(*preHash)); ok {
					bsOne := random.(*[]byte)
					bs = append(bs, *bsOne...)
					preHash = bsOne

					continue
				}
			}
		}

		bs = append(bs, lastBlock.Id...)
		if lastBlock.PreBlock == nil {
			break
		}
		lastBlock = lastBlock.PreBlock
	}

	bs = utils.Hash_SHA3_256(bs)
	if lastBlock != nil {

	}

	return &bs
}

func (this *Chain) GetHistoryBalance(start *big.Int, total int) []HistoryItem {

	return balanceHistoryManager.Get(start, total)
}
