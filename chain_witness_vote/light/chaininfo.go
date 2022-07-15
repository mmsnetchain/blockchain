package light

import (
	"bytes"
	"encoding/hex"
	"errors"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
	"runtime"
	"strconv"
	"time"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/shirou/gopsutil/v3/mem"
)

var CountBalanceHeight = uint64(0)

func SyncBlock() error {
	CountBalanceHeight = 0
	bhvo, err := LoadBlockChain()
	if err != nil {
		return err
	}
	chain := mining.GetLongChain()
	bhvo, err = SycnBlockChain(bhvo)
	if err != nil {
		for chain.Temp == nil {
			time.Sleep(time.Second * 6)
		}
		chain.ForkCheck(nil)
		return nil
	}
	chain.SyncBlockFinish = true

	LoopSycnBlockChain(bhvo)
	for chain.Temp == nil {
		time.Sleep(time.Second * 6)
	}
	chain.ForkCheck(nil)

	return nil
}

func LoadBlockChain() (*mining.BlockHeadVO, error) {
	engine.Log.Info("Start loading blocks in database")
	mining.SetHighestBlock(db.GetHighstBlock())

	chain := mining.GetLongChain()

	_, lastBlock := chain.GetLastBlock()
	headid := lastBlock.Id
	bhvo, err := mining.LoadBlockHeadVOByHash(&headid)
	if err != nil {
		return bhvo, err
	}

	if bhvo.BH.Nextblockhash == nil {
		return bhvo, nil
	}
	blockhash := bhvo.BH.Nextblockhash
	var bhvoTemp *mining.BlockHeadVO
	for blockhash != nil && len(blockhash) > 0 {

		memInfo, _ := mem.VirtualMemory()
		if memInfo.UsedPercent > config.Wallet_Memory_percentage_max {
			runtime.GC()
			time.Sleep(time.Second)
		}

		bhvoTemp, err = mining.LoadBlockHeadVOByHash(&blockhash)
		if err != nil || bhvoTemp == nil {
			break
		}

		bhvo = bhvoTemp

		CountBalance(chain, bhvoTemp)
		chain.SetPulledStates(bhvoTemp.BH.Height)
		chain.SetCurrentBlock(bhvoTemp.BH.Height)
		CountBalanceHeight = bhvoTemp.BH.Height

		if bhvoTemp.BH.Nextblockhash == nil || len(bhvoTemp.BH.Nextblockhash) <= 0 {
			break
		}
		blockhash = bhvoTemp.BH.Nextblockhash
	}
	engine.Log.Info("end loading blocks in database")
	return bhvo, nil
}

func SycnBlockChain(bhvo *mining.BlockHeadVO) (*mining.BlockHeadVO, error) {
	engine.Log.Info("Start synchronizing blocks from neighbor nodes")
	bhvoLocal := bhvo
	engine.Log.Info("Start Sync hash:%s", hex.EncodeToString(bhvoLocal.BH.Hash))

	chain := mining.GetLongChain()

	peerBlockinfo, rch := mining.FindRemoteCurrentHeight()

	bhvoRemote, err := mining.SyncBlockFlashDB(&bhvoLocal.BH.Hash, peerBlockinfo)
	if err != nil {
		engine.Log.Info("SyncBlockFlashDB error:%s", err.Error())

		return nil, err
	}
	bhvoLocal = bhvoRemote

	SyncFailedTotal := 0
	for SyncFailedTotal < 60 {

		if rch <= CountBalanceHeight {
			peerBlockinfo, rch = mining.FindRemoteCurrentHeight()
			if rch <= CountBalanceHeight {

				if !chain.SyncBlockFinish {
					chain.SyncBlockFinish = true
				}

				return bhvoLocal, nil
			}
		}

		mining.SetHighestBlock(rch)

		bhvoRemote, err = mining.SyncBlockFlashDB(&bhvoLocal.BH.Nextblockhash, peerBlockinfo)
		if err != nil {
			engine.Log.Info("SyncBlockFlashDB error:%s", err.Error())
			SyncFailedTotal++
			continue
		}
		engine.Log.Info("height:%d CountBalanceHeight:%d", bhvoRemote.BH.Height, CountBalanceHeight)

		if bhvoRemote.BH.Height > CountBalanceHeight {
			SyncFailedTotal = 0
			CountBalance(chain, bhvoRemote)
			chain.SetPulledStates(bhvoRemote.BH.Height)
			chain.SetCurrentBlock(bhvoRemote.BH.Height)
			CountBalanceHeight = bhvoRemote.BH.Height

		}

		if bhvoRemote.BH.Nextblockhash == nil || len(bhvoRemote.BH.Nextblockhash) <= 0 {
			engine.Log.Info("block not nextblockhash: %s", hex.EncodeToString(bhvoRemote.BH.Hash))
			SyncFailedTotal++
			continue
		}
		engine.Log.Info("sync next block hash")
		bhvoLocal = bhvoRemote
	}

	return bhvoLocal, errors.New("")
}

func LoopSycnBlockChain(bhvo *mining.BlockHeadVO) *[]byte {
	engine.Log.Info("Start loop synchronizing blocks from neighbor nodes")
	bhvoLocal := bhvo
	engine.Log.Info("Start Sync hash:%s", hex.EncodeToString(bhvoLocal.BH.Hash))

	chain := mining.GetLongChain()

	var bhRemote *mining.BlockHead
	var bhvoRemote *mining.BlockHeadVO
	var err error

	ticker := time.NewTicker(time.Second * config.Mining_block_time).C
	nowChan := make(chan bool, 1)

	SyncFailedTotal := 0

	var peerBlockinfo *mining.PeerBlockInfoDESC
	var rch uint64

	for SyncFailedTotal < 60 {

		select {
		case <-ticker:
			peerBlockinfo, rch = mining.FindRemoteCurrentHeight()
			if rch <= CountBalanceHeight {
				SyncFailedTotal = 0
				continue
			}
		case <-nowChan:
		}
		engine.Log.Info("remote height:%d local height:%d", rch, CountBalanceHeight)
		if rch <= CountBalanceHeight {
			continue
		}

		bhRemote, err = mining.FindBlockHeadNeighbor(&bhvoLocal.BH.Hash, peerBlockinfo)
		if err != nil {
			engine.Log.Info("FindBlockHeadNeighbor error:%s", err.Error())
			SyncFailedTotal++
			continue
		}

		if !bytes.Equal(bhRemote.Nextblockhash, bhvoLocal.BH.Nextblockhash) {
			bhvoLocal.BH.Nextblockhash = bhRemote.Nextblockhash
			bs, err := bhvoLocal.BH.Proto()
			if err != nil {
				return nil
			}
			db.LevelDB.Save(bhvoLocal.BH.Hash, bs)
		}

		mining.SetHighestBlock(rch)

		bhvoRemote, err = mining.LoadBlockHeadVOByHash(&bhvoLocal.BH.Nextblockhash)
		if err != nil {

			bhvoRemote, err = mining.SyncBlockFlashDB(&bhvoLocal.BH.Nextblockhash, peerBlockinfo)
			if err != nil {
				engine.Log.Info("SyncBlockFlashDB error:%s", err.Error())
				SyncFailedTotal++
				continue
			}
		}

		if bhvoRemote.BH.Height > CountBalanceHeight {
			SyncFailedTotal = 0
			CountBalance(chain, bhvoRemote)
			chain.SetPulledStates(bhvoRemote.BH.Height)
			chain.SetCurrentBlock(bhvoRemote.BH.Height)
			CountBalanceHeight = bhvoRemote.BH.Height
			select {
			case nowChan <- false:
			default:
			}
		}

		if bhvoRemote.BH.Nextblockhash == nil || len(bhvoRemote.BH.Nextblockhash) <= 0 {
			engine.Log.Info("block not nextblockhash: %s", hex.EncodeToString(bhvoRemote.BH.Hash))
			SyncFailedTotal++
			continue
		}
		engine.Log.Info("sync next block hash")
		bhvoLocal = bhvoRemote
	}
	chain.SyncBlockFinish = true
	return &bhvoLocal.BH.Hash
}

func CountBalance(chain *mining.Chain, bhvo *mining.BlockHeadVO) {
	engine.Log.Info("====== Count block group:%d block:%d prehash:%s hash:%s witness:%s", bhvo.BH.GroupHeight,
		bhvo.BH.Height, hex.EncodeToString(bhvo.BH.Previousblockhash), hex.EncodeToString(bhvo.BH.Hash), bhvo.BH.Witness.B58String())
	mining.GetLongChain().Balance.CountBalanceForBlock(bhvo)
	chain.WitnessBackup.CountWitness(&bhvo.Txs)
}

func ForkCheck(bhash *[]byte) {
	engine.Log.Error("ForkCheck hash:%s", hex.EncodeToString(*bhash))

	if bhash == nil {
		return
	}
	preBlockHash := bhash

	for {
		engine.Log.Info("LoadBlockHeadByHash:%s", hex.EncodeToString(*preBlockHash))

		bh, err := mining.LoadBlockHeadByHash(preBlockHash)
		if err != nil || bh == nil {
			engine.Log.Error("ForkCheck faile!")
			return
		}

		peerBlockinfo, _ := mining.FindRemoteCurrentHeight()
		bhvo, err := mining.SyncBlockFlashDB(preBlockHash, peerBlockinfo)
		if err == nil && bhvo != nil {

			bhash, err := db.LevelDB.Find([]byte(config.BlockHeight + strconv.Itoa(int(bhvo.BH.Height))))
			if err == nil && bytes.Equal(*bhash, bhvo.BH.Hash) {
				engine.Log.Info("fork block:%s", hex.EncodeToString(*bhash))
				break
			}
		}

		preBlockHash = &bh.Previousblockhash

	}
	engine.Log.Error("ForkCheck finish!")
}
