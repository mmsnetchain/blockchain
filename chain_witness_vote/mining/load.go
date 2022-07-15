package mining

import (
	"bytes"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"runtime"
	"time"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/shirou/gopsutil/v3/mem"
)

func (this *Chain) LoadBlockChain() error {

	engine.Log.Info("Start loading blocks in database")
	forks.SetHighestBlock(db.GetHighstBlock())

	_, lastBlock := this.GetLastBlock()
	headid := lastBlock.Id
	bh, _, err := loadBlockForDB(&headid)
	if err != nil {
		return err
	}

	if bh.Nextblockhash == nil {
		return nil
	}
	blockhash := bh.Nextblockhash
	var bhvo *BlockHeadVO
	for blockhash != nil && len(blockhash) > 0 {

		memInfo, _ := mem.VirtualMemory()
		if memInfo.UsedPercent > config.Wallet_Memory_percentage_max {
			runtime.GC()
			time.Sleep(time.Second)
		}

		if bhvo != nil && bhvo.BH.Height+1 == config.CutBlockHeight && bhvo.BH.Nextblockhash != nil && bytes.Equal(bhvo.BH.Nextblockhash, config.CutBlockHash) {

			break
		}

		bhvo = this.deepCycleLoadBlock(&blockhash)
		if bhvo == nil {
			break
		}

		if bhvo.BH.Nextblockhash == nil || len(bhvo.BH.Nextblockhash) <= 0 {
			break
		}
		blockhash = bhvo.BH.Nextblockhash
	}
	engine.Log.Info("end loading blocks in database")
	return nil
}

func (this *Chain) deepCycleLoadBlock(bhash *[]byte) *BlockHeadVO {
	if len(config.BlockHashs) > 0 && bytes.Equal(*bhash, config.BlockHashs[0]) {
		bhash = this.testLoadBlock()

	}

	bh, txItrs, err := loadBlockForDB(bhash)
	if err != nil {
		engine.Log.Info("load block for db error:%s", err.Error())
		return nil
	}

	bhvo := &BlockHeadVO{FromBroadcast: false, BH: bh, Txs: txItrs}
	err = this.AddBlock(bhvo)
	if err != nil {
		if err.Error() == ERROR_repeat_import_block.Error() {

		} else {
			engine.Log.Info("add block error:%s", err.Error())
			return nil
		}
	}

	if bh.Nextblockhash == nil {
		engine.Log.Info("load block next blockhash nil")
		return nil
	}

	return bhvo
}

func (this *Chain) testLoadBlock() *[]byte {
	engine.Log.Info("testLoadBlock---------------------------")
	var nextBlockhash *[]byte
	for _, one := range config.BlockHashs {
		bh, txItrs, err := loadBlockForDB(&one)
		if err != nil {
			continue
		}

		bhvo := &BlockHeadVO{FromBroadcast: false, BH: bh, Txs: txItrs}
		err = this.AddBlock(bhvo)
		if err != nil {
			continue
		}

		if bh.Nextblockhash == nil {
			continue
		}
		nextBlockhash = &bh.Nextblockhash
	}
	return nextBlockhash
}

func LoadStartBlock() *BlockHeadVO {
	exist, err := db.LevelDB.CheckHashExist(config.Key_block_start)
	if err != nil {
		engine.Log.Info("load start block hash error:%s", err.Error())
		return nil
	}
	if !exist {
		return nil
	}

	headid, err := db.LevelDB.Find(config.Key_block_start)
	if err != nil {

		engine.Log.Info("This is an empty database")
		return nil
	}
	bh, txItrs, err := loadBlockForDB(headid)
	if err != nil {
		return nil
	}
	bhvo := BlockHeadVO{
		BH:  bh,
		Txs: txItrs,
	}
	return &bhvo
}
