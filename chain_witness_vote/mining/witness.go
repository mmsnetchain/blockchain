package mining

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"mmschainnewaccount/config"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/message_center/flood"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

type WitnessChain struct {
	chain           *Chain
	witnessBackup   *WitnessBackup
	witnessGroup    *WitnessGroup
	witnessNotGroup []*Witness
}

func NewWitnessChain(wb *WitnessBackup, chain *Chain) *WitnessChain {
	return &WitnessChain{
		chain:         chain,
		witnessBackup: wb,
	}
}

type WitnessGroup struct {
	Task         bool
	PreGroup     *WitnessGroup
	NextGroup    *WitnessGroup
	Height       uint64
	Witness      []*Witness
	BlockGroup   *Group
	IsBuildGroup bool
	tag          bool
	IsCount      bool
}

type Witness struct {
	Group           *WitnessGroup
	PreWitness      *Witness
	NextWitness     *Witness
	Addr            *crypto.AddressCoin
	Puk             []byte
	Block           *Block
	Score           uint64
	CommunityVotes  []*VoteScore
	Votes           []*VoteScore
	VoteNum         uint64
	StopMining      chan bool `json:"-"`
	BlockHeight     uint64
	CreateBlockTime int64

	WitnessBackupGroup *WitnessBackupGroup

	CheckIsMining bool
	syncBlockOnce *sync.Once
}

func (this *WitnessChain) BuildBlockGroupForGroupHeight(groupHeight uint64, blockhash *[]byte) {

	currentGroup := this.witnessGroup
	for currentGroup != nil {
		if currentGroup.Height < groupHeight {

			if currentGroup.NextGroup == nil {

				break
			}
			currentGroup = currentGroup.NextGroup
			continue
		}
		if currentGroup.Height > groupHeight {

			if currentGroup.PreGroup == nil {

				return
			}
			currentGroup = currentGroup.PreGroup
			continue
		}
		if currentGroup.Height == groupHeight {
			break
		}
	}

	for {

		if (this.witnessGroup.Height > currentGroup.Height) || (this.witnessGroup.NextGroup == nil) {
			break
		}
		this.witnessGroup.BuildGroup(blockhash)
		this.chain.CountBlock(this.witnessGroup)
		engine.Log.Info("3333333333333333 %d", this.witnessGroup.Height)
		this.witnessGroup = this.witnessGroup.NextGroup

	}

}

func (this *WitnessChain) BuildBlockGroup(bhvo *BlockHeadVO, preWitness *Witness) {

	if bhvo.BH.GroupHeight == config.Mining_group_start_height {

		this.witnessGroup.BuildGroup(nil)
		this.BuildWitnessGroup(false, false)
		this.witnessGroup = this.witnessGroup.NextGroup
		this.BuildWitnessGroup(false, true)
		return
	}

	if preWitness == nil {

		return
	}

	witness, _ := this.FindWitnessForBlock(bhvo)
	if witness == nil {
		return
	}

	isNewGroup := false

	if bhvo.BH.GroupHeight != preWitness.Group.Height {
		isNewGroup = true
		if bhvo.BH.Height < config.FixBuildGroupBUGHeightMax {
			preWitness.Group.BuildGroup(&bhvo.BH.Previousblockhash)
			wg := preWitness.Group.BuildGroup(&bhvo.BH.Previousblockhash)

			if wg != nil {
				for _, one := range wg.Witness {
					if !one.CheckIsMining {
						if one.Block == nil {
							this.chain.WitnessBackup.AddBlackList(*one.Addr)
						} else {
							this.chain.WitnessBackup.SubBlackList(*one.Addr)
						}
						one.CheckIsMining = true
					}
				}
			}
			this.chain.CountBlock(preWitness.Group)
		}

		if bhvo.BH.GroupHeight != config.Mining_group_start_height+1 {

			this.witnessGroup = witness.Group
		}
	}
	this.BuildWitnessGroup(false, isNewGroup)

}

func (this *WitnessGroup) BuildGroup(blockhash *[]byte) *WitnessGroup {

	if blockhash == nil || len(*blockhash) <= 0 {
		engine.Log.Info("count group:%d", this.Height)
	} else {
		engine.Log.Info("count group:%d preblockhash:%s", this.Height, hex.EncodeToString(*blockhash))
	}

	if this.IsBuildGroup {
		engine.Log.Info("Already built.")
		return nil
	}

	ok, group := this.CheckBlockGroup(blockhash)
	if !ok {
		engine.Log.Info("The number of people in this group is too small and unqualified. group:%d", this.Height)

		return nil
	}

	this.BlockGroup = group
	this.IsBuildGroup = true

	beforeGroup := this.PreGroup
	for beforeGroup = this.PreGroup; beforeGroup != nil; beforeGroup = beforeGroup.PreGroup {
		ok, _ = beforeGroup.CheckBlockGroup(blockhash)
		if ok {
			break
		}
	}

	if beforeGroup == nil {
		engine.Log.Info("The last group found is empty")
		return nil
	}

	beforeGroup.BlockGroup.NextGroup = this.BlockGroup
	this.BlockGroup.PreGroup = beforeGroup.BlockGroup

	beforeBlock := beforeGroup.BlockGroup.Blocks[len(beforeGroup.BlockGroup.Blocks)-1]

	beforeBlock.NextBlock = this.BlockGroup.Blocks[0]

	return beforeGroup
}

func (this *WitnessGroup) CheckBlockGroup(blockhash *[]byte) (bool, *Group) {

	if this.BlockGroup != nil {

		return true, this.BlockGroup
	}

	group := this.SelectionChain(blockhash)

	if group == nil {

		return false, nil
	}
	totalWitness := len(this.Witness)
	totalHave := len(group.Blocks)

	if (totalHave * 2) <= totalWitness {

		return false, group
	}

	return true, group
}

func (this *WitnessGroup) SelectionChain(blockhash *[]byte) *Group {

	this.CollationRelationship(blockhash)

	groupMap := make(map[string]*Group)
	for _, one := range this.Witness {

		if one.Block == nil {

			continue
		}
		if one.Block.Group == nil {

			continue
		}
		blocks := one.Block.Group.Blocks

		groupMap[utils.Bytes2string(blocks[len(blocks)-1].Id)] = one.Block.Group

	}

	var group *Group
	var groupByBlockhash *Group
	for _, v := range groupMap {
		if blockhash != nil && groupByBlockhash == nil {
			for _, blockOne := range v.Blocks {

				if bytes.Equal(blockOne.Id, *blockhash) {

					groupByBlockhash = v
					break
				}
			}
		}

		if group == nil {

			group = v
			continue
		}

		if len(v.Blocks) > len(group.Blocks) {

			group = v
		}

	}

	if blockhash != nil {
		group = groupByBlockhash
	}
	if group == nil || group.Blocks == nil || len(group.Blocks) <= 0 {

		return group
	}

	for _, lastBlock := range group.Blocks {
		for _, preBlock := range group.Blocks {
			if bytes.Equal(lastBlock.PreBlockID, preBlock.Id) {
				preBlock.NextBlock = lastBlock
				lastBlock.PreBlock = preBlock

				break
			}
		}
	}

	firstBlock := group.Blocks[0]

	for preWitnessGroup := this.PreGroup; preWitnessGroup != nil; preWitnessGroup = preWitnessGroup.PreGroup {

		have := false
		for _, witnessOne := range preWitnessGroup.Witness {

			if witnessOne.Block == nil {

				continue
			}

			if bytes.Equal(witnessOne.Block.Id, firstBlock.PreBlockID) {

				witnessOne.Block.NextBlock = firstBlock
				firstBlock.PreBlock = witnessOne.Block
				have = true
				break
			}
		}
		if have {
			break
		}
	}

	return group
}

func (this *WitnessGroup) CleanRelationship() {
	for _, one := range this.Witness {
		if one.Block == nil {
			continue
		}
		one.Block.Group = nil
	}
}

func (this *WitnessGroup) CollationRelationship(blockhash *[]byte) {

	if blockhash == nil {

		for i := len(this.Witness); i > 0; i-- {
			witness := this.Witness[i-1]

			if witness.Block == nil {

				continue
			}
			newBlock := witness.Block

			for witnessTemp := witness.NextWitness; witnessTemp != nil; witnessTemp = witnessTemp.NextWitness {

				if witnessTemp.Group.Height != witness.Group.Height {

					newGroup := new(Group)
					newGroup.Height = witness.Group.Height
					newGroup.Blocks = []*Block{newBlock}
					newBlock.Group = newGroup
					break
				}
				if witnessTemp.Block == nil {

					continue
				}
				if bytes.Equal(witnessTemp.Block.PreBlockID, newBlock.Id) {

					witnessTemp.Block.Group.Blocks = append([]*Block{newBlock}, witnessTemp.Block.Group.Blocks...)
					witness.Block.Group = witnessTemp.Block.Group
					newBlock.NextBlock = witnessTemp.Block

					witnessTemp.Block.PreBlock = newBlock
					break
				}
			}

			if witness.Block.Height == config.Mining_block_start_height {

				newGroup := new(Group)
				newGroup.Height = witness.Group.Height
				newGroup.Blocks = []*Block{newBlock}

				newBlock.Group = newGroup
				continue
			}
		}
	} else {

		var findBlock *Block
		for i := len(this.Witness); i > 0; i-- {
			witness := this.Witness[i-1]
			if witness.Block == nil {

				continue
			}

			newBlock := witness.Block
			if findBlock == nil {

				if bytes.Equal(newBlock.Id, *blockhash) {

					newGroup := new(Group)
					newGroup.Height = witness.Group.Height
					newGroup.Blocks = []*Block{newBlock}
					newBlock.Group = newGroup
					findBlock = newBlock
				}
				continue
			}

			if bytes.Equal(witness.Block.Id, findBlock.PreBlockID) {

				findBlock.Group.Blocks = append([]*Block{newBlock}, findBlock.Group.Blocks...)
				witness.Block.Group = findBlock.Group
				findBlock = witness.Block

				continue
			}

			if witness.Block.Height == config.Mining_block_start_height {

				newGroup := new(Group)
				newGroup.Height = witness.Group.Height
				newGroup.Blocks = []*Block{newBlock}

				newBlock.Group = newGroup
				continue
			}
		}
	}

}

func (this *WitnessGroup) SelectionChainOld() *Group {
	this.CollationRelationshipOld()

	groupMap := make(map[string]*Group)
	for _, one := range this.Witness {

		if one.Block == nil {

			continue
		}
		if one.Block.Group == nil {

			continue
		}

		groupMap[utils.Bytes2string(one.Block.Group.Blocks[0].Id)] = one.Block.Group
	}

	var group *Group
	for _, v := range groupMap {

		if group == nil {

			group = v
			continue
		}

		if len(v.Blocks) > len(group.Blocks) {

			group = v
		}

	}

	return group
}

func (this *WitnessGroup) CollationRelationshipOld() {

	for _, witness := range this.Witness {

		if witness.Block == nil {

			continue
		}

		newBlock := witness.Block

		var beforeWitness *Witness
		for witnessTemp := witness.PreWitness; witnessTemp != nil; witnessTemp = witnessTemp.PreWitness {

			if witnessTemp.Block == nil {

				continue
			}

			if witnessTemp.Group.Height+uint64(config.Witness_backup_group) <= witness.Group.Height {

				for syncWitness := witness.PreWitness; syncWitness != nil &&
					witness.Group.Height == syncWitness.Group.Height+1; syncWitness = syncWitness.PreWitness {
					if syncWitness.Block != nil {
						continue
					}

				}

			}

			if bytes.Equal(witnessTemp.Block.Id, newBlock.PreBlockID) {

				beforeWitness = witnessTemp
				break
			}
		}

		if witness.Block.Height == config.Mining_block_start_height {

			newGroup := new(Group)
			newGroup.Height = witness.Group.Height
			newGroup.Blocks = []*Block{newBlock}

			newBlock.Group = newGroup
			continue
		}

		if beforeWitness == nil {

			continue
		}

		if beforeWitness.Group.Height == witness.Group.Height {

			have := false
			for _, one := range beforeWitness.Block.Group.Blocks {
				if newBlock.witness == one.witness {
					have = true
					break
				}
			}
			if !have {

				beforeWitness.Block.Group.Blocks = append(beforeWitness.Block.Group.Blocks, newBlock)
				newBlock.Group = beforeWitness.Block.Group
				newBlock.PreBlock = beforeWitness.Block
				beforeWitness.Block.NextBlock = newBlock
			}
		} else {

			if newBlock.Group == nil {

				newBlock.Group = new(Group)
				newBlock.Group.Height = witness.Group.Height
				newBlock.Group.Blocks = []*Block{newBlock}

			}
			newBlock.PreBlock = beforeWitness.Block
		}

	}
}

func (this *WitnessChain) BuildWitnessGroup(first, isPrint bool) {
	backupGroup := uint64(config.Witness_backup_group)

	totalBackupGroup := 0
	tag := false
	lastGroupHeight := uint64(config.Mining_group_start_height)
	var lastGroup *WitnessGroup
	for lastGroup = this.witnessGroup; lastGroup != nil; lastGroup = lastGroup.NextGroup {

		totalBackupGroup++
		tag = lastGroup.tag
		lastGroupHeight = lastGroup.Height
		if lastGroup.NextGroup == nil {
			break
		}
	}

	if first {
		backupGroup = 0
	} else {
		total := this.witnessBackup.GetBackupWitnessTotal()
		groupNum := total / config.Mining_group_max
		if groupNum > backupGroup {
			backupGroup = groupNum
		}
	}

	for i := uint64(totalBackupGroup); i <= backupGroup; i++ {
		newGroup := this.AdditionalWitnessBackup(lastGroup, lastGroupHeight, tag)
		if newGroup == nil {
			break
		}
		lastGroupHeight++
		tag = !tag
		lastGroup = newGroup

		if this.witnessGroup == nil {

			this.witnessGroup = newGroup
		}
	}
	if isPrint {
		this.PrintWitnessList()
	}
}

func (this *WitnessChain) AdditionalWitnessBackup(lastGroup *WitnessGroup, lastGroupHeight uint64, tag bool) *WitnessGroup {

	witnessGroup := this.GetOneGroupWitness()

	if witnessGroup == nil || len(witnessGroup) < config.Mining_group_min {
		engine.Log.Info("Too few witnesses current:%d min:%d", len(witnessGroup), config.Mining_group_min)
		return nil
	}

	startTime := int64(0)
	if lastGroup != nil {
		w := lastGroup.Witness[len(lastGroup.Witness)-1]
		startTime = w.CreateBlockTime
	}

	for _, one := range witnessGroup {
		startTime = startTime + config.Mining_block_time
		one.CreateBlockTime = startTime
	}
	if lastGroup != nil {
		lastGroupHeight++
	}
	newGroup := &WitnessGroup{
		PreGroup: lastGroup,
		Height:   lastGroupHeight,
		Witness:  witnessGroup,
		tag:      !tag,
	}

	for i, _ := range newGroup.Witness {
		newGroup.Witness[i].Group = newGroup
	}

	if lastGroup != nil {
		lastGroup.Witness[len(lastGroup.Witness)-1].NextWitness = newGroup.Witness[0]
		newGroup.Witness[0].PreWitness = lastGroup.Witness[len(lastGroup.Witness)-1]
		lastGroup.NextGroup = newGroup
		newGroup.PreGroup = lastGroup
	}
	return newGroup
}

func (this *WitnessChain) CompensateWitnessGroup() {
	backupGroup := uint64(config.Witness_backup_group)

	totalBackupGroup := 0
	tag := false
	lastGroupHeight := uint64(config.Mining_group_start_height)
	var lastGroup *WitnessGroup
	for lastGroup = this.witnessGroup; lastGroup != nil; lastGroup = lastGroup.NextGroup {

		totalBackupGroup++
		tag = lastGroup.tag
		lastGroupHeight = lastGroup.Height
		if lastGroup.NextGroup == nil {
			break
		}
	}

	witness := lastGroup.Witness[len(lastGroup.Witness)-1]
	for {
		witness = lastGroup.Witness[len(lastGroup.Witness)-1]
		if witness.CreateBlockTime > utils.GetNow() {
			break
		}

		this.BuildBlockGroupForGroupHeight(lastGroup.Height-1, nil)

		totalBackupGroup = 0
		newGroup := this.AdditionalWitnessBackup(lastGroup, lastGroupHeight, tag)
		if newGroup == nil {
			break
		}
		lastGroupHeight++
		tag = !tag
		lastGroup = newGroup

		engine.Log.Info("New group height %d %d", lastGroupHeight, lastGroup.Height)

		engine.Log.Info("Estimated block out time %d %s", witness.Group.Height, time.Unix(witness.CreateBlockTime, 0).Format("2006-01-02 15:04:05"))
		continue
	}

	total := this.witnessBackup.GetBackupWitnessTotal()
	groupNum := total / config.Mining_group_max
	if groupNum > backupGroup {
		backupGroup = groupNum
	}

	for i := uint64(totalBackupGroup); i <= backupGroup; i++ {

		newGroup := this.AdditionalWitnessBackup(lastGroup, lastGroupHeight, tag)
		if newGroup == nil {
			break
		}
		lastGroupHeight++
		tag = !tag
		lastGroup = newGroup

		if this.witnessGroup == nil {

			this.witnessGroup = newGroup
		}
	}

	this.PrintWitnessList()

	_, isBackup, _, _, _ := GetWitnessStatus()
	if isBackup {
		config.AlreadyMining = true
	}

	this.StopAllMining()

	this.BuildMiningTime()

	this.chain.SyncBlockFinish = true
}

func (this *WitnessChain) CompensateWitnessGroupByGroupHeight(groupHeight uint64) {
	engine.Log.Info("build group height:%d", groupHeight)

	totalBackupGroup := 0
	tag := false
	lastGroupHeight := uint64(config.Mining_group_start_height)
	var lastGroup *WitnessGroup
	for lastGroup = this.witnessGroup; lastGroup != nil; lastGroup = lastGroup.NextGroup {

		totalBackupGroup++
		tag = lastGroup.tag
		lastGroupHeight = lastGroup.Height
		if lastGroup.NextGroup == nil {
			break
		}
	}

	for i := lastGroupHeight; i <= groupHeight; i++ {
		engine.Log.Info("start build group height:%d", lastGroupHeight)

		totalBackupGroup = 0
		newGroup := this.AdditionalWitnessBackup(lastGroup, lastGroupHeight, tag)
		if newGroup == nil {
			break
		}
		lastGroupHeight++
		tag = !tag
		lastGroup = newGroup

		continue
	}

	this.PrintWitnessList()
}

func (this *WitnessChain) StopAllMining() {

	addr := keystore.GetCoinbase()

	witnessTemp := this.witnessGroup.Witness[0]
	for {
		witnessTemp = witnessTemp.NextWitness
		if witnessTemp == nil || witnessTemp.Group == nil {

			break
		}
		if !witnessTemp.Group.Task {

			continue
		}

		if !bytes.Equal(*witnessTemp.Addr, addr.Addr) {

			continue
		}

		select {
		case witnessTemp.StopMining <- false:

		default:

		}
		witnessTemp.Group.Task = false
	}
}

func (this *WitnessChain) BuildMiningTime() error {

	addr := keystore.GetCoinbase()

	for witnessGroup := this.witnessGroup; witnessGroup != nil; witnessGroup = witnessGroup.NextGroup {

		if witnessGroup.Task {
			continue
		}
		for _, witnessTemp := range witnessGroup.Witness {
			if witnessTemp.Block != nil {
				continue
			}
			if bytes.Equal(*witnessTemp.Addr, addr.Addr) {

				atomic.StoreUint32(&this.chain.StopSyncBlock, 0)
				future := int64(0)
				now := utils.GetNow()

				if witnessTemp.CreateBlockTime > now {
					future = witnessTemp.CreateBlockTime - now
				} else if witnessTemp.CreateBlockTime == now {
					future = 0
				} else {
					difference := now - witnessTemp.CreateBlockTime
					if difference < config.Mining_block_time {
						future = 0
					} else {

						continue
					}

				}

				if !config.InitNode && future <= 20 {

					continue
				}

				engine.Log.Info("Build blocks in %d seconds", future)

				witnessTemp.Group.Task = true
				go witnessTemp.SyncBuildBlock(int64(future))

			} else {

			}

		}

	}

	return nil
}

func (this *WitnessChain) GetOneGroupWitness() []*Witness {

	groupNum := config.Mining_group_min
	total := this.witnessBackup.GetBackupWitnessTotal()
	if total > config.Mining_group_max {
		groupNum = config.Mining_group_max
	} else if total < config.Mining_group_min {
		groupNum = config.Mining_group_min
	} else {
		groupNum = int(total)
	}

	if len(this.witnessNotGroup) < groupNum {

		witness := this.witnessBackup.CreateWitnessGroup()
		if witness == nil {
			return nil
		}
		this.witnessNotGroup = append(this.witnessNotGroup, witness...)

	}

	index := 0
	moveWitness := make([]*Witness, 0)
	witnessGroup := make([]*Witness, 0)
	for i, tempWitness := range this.witnessNotGroup {
		index = i

		isHave := false
		for _, one := range witnessGroup {
			if bytes.Equal(*tempWitness.Addr, *one.Addr) {

				moveWitness = append(moveWitness, tempWitness)

				isHave = true
				break
			}
		}

		if isHave {
			continue
		}
		witnessGroup = append(witnessGroup, tempWitness)
		if len(witnessGroup) >= groupNum {
			break
		}
	}
	newWitnessNotGroup := this.witnessNotGroup[index+1:]
	this.witnessNotGroup = make([]*Witness, 0)

	for _, one := range moveWitness {
		this.witnessNotGroup = append(this.witnessNotGroup, one)
	}
	for _, one := range newWitnessNotGroup {
		this.witnessNotGroup = append(this.witnessNotGroup, one)
	}

	tempWitness := witnessGroup[0]
	for i, _ := range witnessGroup {
		if i == 0 {
			continue
		}
		tempWitness.NextWitness = witnessGroup[i]
		witnessGroup[i].PreWitness = tempWitness
		tempWitness = witnessGroup[i]
	}

	return witnessGroup

}

func (this *WitnessChain) PrintWitnessList() {

	return

	count := 0
	group := this.witnessGroup
	for group != nil {
		engine.Log.Info("--------------")
		for _, one := range group.Witness {

			engine.Log.Info("witness tag %s %t %d %s %s %d", fmt.Sprintf("%p", one.WitnessBackupGroup), group.tag, group.Height, one.Addr.B58String(),
				time.Unix(one.CreateBlockTime, 0).Format("2006-01-02 15:04:05"), one.VoteNum)
		}

		count++
		if count >= 2 {
			return
		}

		group = group.NextGroup
	}
}

func (this *WitnessChain) FindBlockInCurrent(bh *BlockHead) bool {

	currentGroup := this.witnessGroup
	for currentGroup != nil {
		if currentGroup.Height < bh.GroupHeight {
			currentGroup = currentGroup.NextGroup
			continue
		}
		if currentGroup.Height > bh.GroupHeight {
			currentGroup = currentGroup.PreGroup
			continue
		}
		break
	}
	if currentGroup == nil {
		return false
	}

	if currentGroup.BlockGroup == nil {
		return false
	}

	for _, one := range currentGroup.Witness {
		if one.Block == nil {
			continue
		}
		if bytes.Equal(one.Block.Id, bh.Hash) {
			return true
		}
	}

	engine.Log.Info("find the group %+v", currentGroup)
	return false
}

func (this *WitnessChain) FindPreWitnessForBlock(preBlockHash []byte) *Witness {

	engine.Log.Info(" %d", this.witnessGroup.Height)
	currentGroup := this.witnessGroup
	for {
		if currentGroup.BlockGroup != nil {
			break
		}
		if currentGroup.PreGroup == nil {
			engine.Log.Info("find the group %+v", currentGroup)
			return nil
		}
		currentGroup = currentGroup.PreGroup
	}
	engine.Log.Info(" %+v", currentGroup)

	for currentGroup != nil {
		for _, one := range currentGroup.Witness {
			if one.Block == nil {
				continue
			}
			if bytes.Equal(one.Block.Id, preBlockHash) {
				return one
			}
		}
		currentGroup = currentGroup.NextGroup
	}

	engine.Log.Info("find the group %+v", currentGroup)
	return nil
}

func (this *WitnessChain) CheckWitnessBuildBlockTime(bhvo *BlockHeadVO) bool {
	return true
}

func (this *WitnessChain) FindWitnessForBlockOnly(bhvo *BlockHeadVO) (*Witness, bool) {
	engine.Log.Info("find group height:%d this.witnessGroup:%d", bhvo.BH.GroupHeight, this.witnessGroup.Height)

	currentGroup := this.witnessGroup
	for currentGroup != nil {
		if currentGroup.Height < bhvo.BH.GroupHeight {

			currentGroup = currentGroup.NextGroup
			continue
		}
		if currentGroup.Height > bhvo.BH.GroupHeight {

			return nil, false
		}
		break
	}

	if currentGroup == nil {

		return nil, true
	}

	if bhvo.BH.Height < config.WitnessOrderCorrectEnd && bhvo.BH.Height > config.WitnessOrderCorrectStart {

		for _, one := range currentGroup.Witness {

			if one.Block != nil {

				continue
			}

			return one, false
		}
	}

	for _, one := range currentGroup.Witness {

		if bytes.Equal(bhvo.BH.Witness, *one.Addr) {

			return one, false
		}
	}
	return nil, false
}

func (this *WitnessChain) CheckRepeatImportBlock(bhvo *BlockHeadVO) bool {

	currentGroup := this.witnessGroup
	for currentGroup != nil {
		if currentGroup.Height < bhvo.BH.GroupHeight {
			currentGroup = currentGroup.NextGroup
			continue
		}
		if currentGroup.Height > bhvo.BH.GroupHeight {
			currentGroup = currentGroup.PreGroup
			continue
		}
		break
	}
	if currentGroup == nil {

		return false
	}

	if currentGroup.IsBuildGroup {
		engine.Log.Info("this group is build: %d", currentGroup.Height)
		return true
	}

	if bhvo.BH.Height < config.WitnessOrderCorrectEnd && bhvo.BH.Height > config.WitnessOrderCorrectStart {

		for _, one := range currentGroup.Witness {

			if one.Block != nil {

				continue
			}

			return false
		}
	}

	for _, one := range currentGroup.Witness {

		if bytes.Equal(*one.Addr, bhvo.BH.Witness) {
			if one.Block != nil {

				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (this *WitnessChain) CheckBifurcationBlock(groupHeight, blockHeight uint64, preBlockHash []byte) (*Witness, bool, bool) {
	return nil, false, false
}

func (this *WitnessChain) FindWitnessForBlock(bhvo *BlockHeadVO) (*Witness, bool) {

	var witness *Witness
	for group := this.witnessGroup; group != nil; group = group.NextGroup {

		if group.Height < bhvo.BH.GroupHeight {

			if group.NextGroup != nil {
				continue
			}
			if utils.GetNow() < bhvo.BH.Time {
				continue
			}

		}
		if group.Height > bhvo.BH.GroupHeight {

			return nil, false
		}

		if bhvo.BH.Height < config.WitnessOrderCorrectEnd && bhvo.BH.Height > config.WitnessOrderCorrectStart {

			for _, one := range group.Witness {

				if one.Block != nil {

					continue
				}

				witness = one

				break
			}
		} else {

			for _, one := range group.Witness {

				if !bytes.Equal(bhvo.BH.Witness, *one.Addr) {

					continue
				}
				now := utils.GetNow()

				if one.CreateBlockTime > now+config.Mining_block_time {

					engine.Log.Warn("， %s %s %s", time.Unix(one.CreateBlockTime, 0).Format("2006-01-02 15:04:05"),
						time.Unix(bhvo.BH.Time, 0).Format("2006-01-02 15:04:05"), time.Unix(now, 0).Format("2006-01-02 15:04:05"))

					break
				}

				witness = one

				break
			}
		}

		if witness != nil {

			break
		}
	}
	return witness, false
}

func (this *WitnessChain) SetWitnessBlock(bhvo *BlockHeadVO) bool {

	witness, needSync := this.FindWitnessForBlock(bhvo)
	if witness != nil && witness.Block != nil {

		engine.Log.Warn("You don't need to set it again if it's already set")
		return false
	}

	if witness == nil {

		engine.Log.Warn("No witness found")

		if needSync {

			this.chain.NoticeLoadBlockForDB()
		}
		return false
	}

	if !bhvo.BH.CheckBlockHead(witness.Puk) {

		engine.Log.Warn("Block verification failed, block is illegal group:%d block:%d", bhvo.BH.GroupHeight, bhvo.BH.Height)

		return false
	}

	if bhvo.BH.Height != config.Mining_block_start_height {

		preBlock, blocks := witness.CheckUnconfirmedBlock(&bhvo.BH.Previousblockhash)

		if bhvo.BH.Height-1 != preBlock.Height {
			engine.Log.Warn("block height not continuity height:%d inport height:%d", preBlock.Height, bhvo.BH.Height)
			engine.Log.Warn("%+v %+v", witness, bhvo.BH)
			return false
		}

		for _, one := range bhvo.Txs {

			err := one.CheckLockHeight(bhvo.BH.Height)
			if err != nil {
				engine.Log.Error("Illegal transaction 111 %s %s", hex.EncodeToString(*one.GetHash()), err.Error())

				return false
			}

			if this.chain.transactionManager.unpackedTransaction.ExistTxByAddrTxid(one) {
				continue
			}

			if bhvo.BH.Height > config.Mining_block_start_height+config.Mining_block_start_height_jump {

				err = one.CheckSign()

				if err != nil {
					engine.Log.Error("Illegal transaction 333 %s %s", hex.EncodeToString(*one.GetHash()), err.Error())

					return false
				}
			}

		}

		unacknowledgedTxs := make([]TxItr, 0)

		exclude := make(map[string]string)
		for _, one := range blocks {

			_, txs, err := one.LoadTxs()
			if err != nil {
				engine.Log.Warn("not find transaction %s", err.Error())

				return false
			}
			for _, txOne := range *txs {

				exclude[utils.Bytes2string(*txOne.GetHash())] = ""
				unacknowledgedTxs = append(unacknowledgedTxs, txOne)
			}
		}

		sizeTotal := uint64(0)
		for i, one := range bhvo.Txs {

			if !one.CheckRepeatedTx(unacknowledgedTxs...) {
				engine.Log.Warn("Transaction verification failed")

				return false
			}
			unacknowledgedTxs = append(unacknowledgedTxs, bhvo.Txs[i])
			sizeTotal = sizeTotal + uint64(len(*one.Serialize()))
		}

		if sizeTotal > config.Block_size_max {
			engine.Log.Error(" %d  %d", sizeTotal, len(bhvo.Txs))
			engine.Log.Warn("Transaction over size %d", sizeTotal)

			return false
		}

		if bhvo.BH.Height > config.Mining_block_start_height+config.Mining_block_start_height_jump {

			if witness.WitnessBackupGroup != preBlock.witness.WitnessBackupGroup {

				vouts := preBlock.witness.WitnessBackupGroup.BuildRewardVouts(blocks, bhvo.BH.Height, &bhvo.BH.Previousblockhash, preBlock)

				haveReward := false
				for _, one := range bhvo.Txs {
					if one.Class() != config.Wallet_tx_type_mining {
						continue
					}
					if haveReward {

						engine.Log.Warn("Illegal if there are multiple reward transactions in a block")
						return false
					}
					haveReward = true

					m := make(map[string]uint64)
					for _, one := range *one.GetVout() {
						m[utils.Bytes2string(one.Address)+strconv.Itoa(int(one.FrozenHeight))] = one.Value
					}
					for _, one := range vouts {

						value, ok := m[utils.Bytes2string(one.Address)+strconv.Itoa(int(one.FrozenHeight))]
						if !ok {

							engine.Log.Warn("Without this person's reward, the verification fails %s", one.Address.B58String())
							return false
						}
						if value != one.Value {

							engine.Log.Warn("If the reward amount is incorrect, the verification fails %s %d want:%d", one.Address.B58String(), value, one.Value)
							return false
						}
					}
				}
				if !haveReward {

					engine.Log.Warn("If there is no block reward, the block is illegal")
					return false
				}
			} else {

				for _, one := range bhvo.Txs {
					if one.Class() == config.Wallet_tx_type_mining {
						engine.Log.Warn("")
						return false
					}
				}
			}
		}
	}

	select {
	case witness.StopMining <- false:
	default:
	}

	newBlock := new(Block)
	newBlock.Id = bhvo.BH.Hash
	newBlock.Height = bhvo.BH.Height
	newBlock.PreBlockID = bhvo.BH.Previousblockhash

	witness.Block = newBlock
	witness.Block.witness = witness

	witness.Group.SelectionChainOld()

	if bhvo.BH.Height == config.Mining_block_start_height {
		witness.CreateBlockTime = bhvo.BH.Time
	}

	this.chain.SetPulledStates(bhvo.BH.Height)

	return true

}

func (this *WitnessGroup) FirstWitness() bool {
	for _, one := range this.Witness {
		if one.Block != nil {
			return false
		}
	}
	return true
}

func (this *WitnessGroup) DistributionRewards() {

}

func (this *WitnessChain) FindWitness(addr crypto.AddressCoin) bool {
	if this.witnessGroup == nil {
		return false
	}
	witnessTemp := this.witnessGroup.Witness[0]
	for {
		witnessTemp = witnessTemp.NextWitness
		if witnessTemp == nil || witnessTemp.Group == nil {
			break
		}
		if bytes.Equal(*witnessTemp.Addr, addr) {
			return true
		}
	}
	return false
}

func (this *Witness) syncBlockTiming() {

	if !forks.GetLongChain().SyncBlockFinish {
		return
	}
	if this.syncBlockOnce != nil {
		return
	}
	this.syncBlockOnce = new(sync.Once)
	var syncBlock = func() {
		utils.Go(func() {
			goroutineId := utils.GetRandomDomain() + utils.TimeFormatToNanosecondStr()
			_, file, line, _ := runtime.Caller(0)
			engine.AddRuntime(file, line, goroutineId)
			defer engine.DelRuntime(file, line, goroutineId)

			bfw := BlockForWitness{
				GroupHeight: this.Group.Height,
				Addr:        *this.Addr,
			}

			now := utils.GetNow()

			if this.CreateBlockTime < now {
				intervalTime := now - this.CreateBlockTime
				if intervalTime > config.Mining_block_time*2 {

					return
				}
			}

			waitTime := this.CreateBlockTime - utils.GetNow()

			if waitTime < 0 {
				waitTime = 0
			}

			delayTime_min := time.Duration(config.Mining_block_time * time.Second / 3)
			delayTime_max := time.Duration(config.Mining_block_time * time.Second / 2)
			delayTime := delayTime_min + time.Duration(utils.GetRandNum(int64(delayTime_max-delayTime_min)))
			delayTime = (time.Duration(waitTime) * time.Second) + delayTime

			engine.Log.Info("Groups %d Time to wait for synchronization %d %s", this.Group.Height, waitTime, delayTime)

			time.Sleep(delayTime)

			intervalTime := time.Second
			intervalTotal := ((config.Mining_block_time * time.Second) - delayTime) / intervalTime

			this.syncBlock(int(intervalTotal), intervalTime, &bfw)

		})
	}
	this.syncBlockOnce.Do(syncBlock)
}

func (this *Witness) syncBlock(total int, intervalTime time.Duration, bfw *BlockForWitness) {
	bs, err := bfw.Proto()

	if err != nil {
		return
	}
	for i := int64(0); i < int64(total); i++ {
		if this.Block != nil {

			return
		}

		broadcasts := append(nodeStore.GetLogicNodes(), nodeStore.GetProxyAll()...)

		for j, _ := range broadcasts {
			if this.Block != nil {

				return
			}
			engine.Log.Info("%d Synchronize blocks from neighbor nodes %s", this.Group.Height, broadcasts[j].B58String())

			message, err := message_center.SendNeighborMsg(config.MSGID_getblockforwitness, &broadcasts[j], bs)
			if err != nil {

				continue
			}

			bs, _ := flood.WaitRequest(config.CLASS_wallet_getblockforwitness, utils.Bytes2string(message.Body.Hash), 1)
			if bs == nil {

				continue
			}

			bhVO, err := ParseBlockHeadVOProto(bs)

			if err != nil {

				continue
			}
			bhVO.FromBroadcast = true
			forks.AddBlockHead(bhVO)

		}

		time.Sleep(intervalTime)
	}
}

func (this *WitnessChain) GCWitnessOld() {

	IsBackup := this.chain.WitnessChain.FindWitness(keystore.GetCoinbase().Addr)
	if !IsBackup {
		return
	}

	total := config.Mining_block_hash_count * 3
	currentWitnessGroup := this.witnessGroup.PreGroup
	for {

		if currentWitnessGroup == nil {
			break
		}

		for _, one := range currentWitnessGroup.Witness {
			if one.Block == nil {
				continue
			}
			total = total - 1

			if total <= 0 {
				break
			}
		}

		if total <= 0 {
			break
		}

		currentWitnessGroup = currentWitnessGroup.PreGroup
	}
	if currentWitnessGroup == nil {
		return
	}

	currentWitnessGroup.PreGroup = nil
	if currentWitnessGroup.BlockGroup != nil {
		currentWitnessGroup.BlockGroup.PreGroup = nil
	}

	firstWitness := currentWitnessGroup.Witness[0]
	if firstWitness.Block != nil {
		if firstWitness.Block.witness != nil {
			firstWitness.Block.witness.PreWitness = nil
		}
		if firstWitness.Block.Group != nil {
			firstWitness.Block.Group.PreGroup = nil
		}
		firstWitness.Block.PreBlock = nil
	}
	firstWitness.PreWitness = nil
	if firstWitness.Group != nil {
		firstWitness.Group.PreGroup = nil
	}

}
