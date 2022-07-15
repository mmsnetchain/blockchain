package mining

import (
	"errors"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"sync/atomic"
)

var forks = new(Forks)

func init() {

}

type Forks struct {
	Init bool

	LongChain *Chain
}

func (this *Forks) AddBlockHead(bhvo *BlockHeadVO) error {

	chain := forks.GetLongChain()

	return chain.AddBlock(bhvo)

}

func (this *Forks) GetLongChain() *Chain {
	return this.LongChain
}

func GetLongChain() *Chain {
	return forks.LongChain
}

func GetFirstBlock() error {

	chainInfo := FindStartBlockForNeighbor()
	if chainInfo == nil {
		return errors.New("Synchronization start block hash failed")
	}

	db.LevelDB.Save(config.Key_block_start, &chainInfo.StartBlockHash)
	config.StartBlockHash = chainInfo.StartBlockHash

	peerBlockinfo, _ := FindRemoteCurrentHeight()
	bhvo, _ := SyncBlockFlashDB(&chainInfo.StartBlockHash, peerBlockinfo)
	if bhvo == nil {
		return nil
	}

	bhvo.BH.BuildBlockHash()

	BuildFirstChain(bhvo)

	if forks.GetHighestBlock() < bhvo.BH.Height {
		forks.SetHighestBlock(bhvo.BH.Height)
	}

	return nil
}

func BuildFirstChain(bhvo *BlockHeadVO) {
	forks.buildFirstChain(bhvo)

}

func (this *Forks) buildFirstChain(bhvo *BlockHeadVO) {

	newChain := NewChain()
	newChain.SetStartingBlock(bhvo.BH.Height, uint64(bhvo.BH.Time))

	this.LongChain = newChain

	newChain.Balance.CountBalanceForBlock(bhvo)

	newChain.WitnessBackup.CountWitness(&bhvo.Txs)

	newChain.WitnessChain.BuildWitnessGroup(true, true)

	newChain.WitnessChain.SetWitnessBlock(bhvo)

	newChain.WitnessChain.BuildBlockGroup(bhvo, nil)

	return
}

func (this *Forks) GetHighestBlock() uint64 {
	if this.LongChain == nil {
		return 0
	}
	return atomic.LoadUint64(&this.LongChain.HighestBlock)
}

func (this *Forks) SetHighestBlock(n uint64) {

	if this.LongChain == nil {
		return
	}
	atomic.StoreUint64(&this.LongChain.HighestBlock, n)
	db.SaveHighstBlock(n)
}

func GetHighestBlock() uint64 {
	return forks.GetHighestBlock()
}

func SetHighestBlock(n uint64) {
	forks.SetHighestBlock(n)
}
