package mining

import (
	"bytes"
	"errors"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"mmschainnewaccount/sqlite3_db"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type BalanceManager struct {
	countTotal    uint64
	chain         *Chain
	syncBlockHead chan *BlockHeadVO

	depositWitness   *TxItem
	depositCommunity *sync.Map
	depositLight     *sync.Map
	depositVote      *sync.Map
	witnessBackup    *WitnessBackup
	txManager        *TransactionManager
	otherDeposit     *sync.Map

	cacheTxlock *sync.Map
}

func NewBalanceManager(wb *WitnessBackup, tm *TransactionManager, chain *Chain) *BalanceManager {
	bm := &BalanceManager{
		chain:         chain,
		syncBlockHead: make(chan *BlockHeadVO, 1),

		witnessBackup:    wb,
		txManager:        tm,
		depositCommunity: new(sync.Map),
		depositLight:     new(sync.Map),
		depositVote:      new(sync.Map),
		otherDeposit:     new(sync.Map),

		cacheTxlock: new(sync.Map),
	}
	utils.Go(bm.run)

	return bm
}

func (this *BalanceManager) GetDepositIn() *TxItem {
	return this.depositWitness
}

func (this *BalanceManager) GetDepositCommunity(addr *crypto.AddressCoin) (item *DepositInfo) {
	v, ok := this.depositCommunity.Load(utils.Bytes2string(*addr))
	if !ok {
		return
	}
	b := v.(*DepositInfo)

	return b
}

func (this *BalanceManager) GetDepositLight(addr *crypto.AddressCoin) (item *DepositInfo) {
	v, ok := this.depositLight.Load(utils.Bytes2string(*addr))
	if !ok {
		return
	}
	item = v.(*DepositInfo)
	return
}

func (this *BalanceManager) GetDepositVote(addr *crypto.AddressCoin) *DepositInfo {
	v, ok := this.depositVote.Load(utils.Bytes2string(*addr))
	if !ok {
		return nil
	}
	b := v.(*DepositInfo)
	return b
}

func (this *BalanceManager) BuildPayVinNew(srcAddress *crypto.AddressCoin, amount uint64) (uint64, *TxItem) {
	if srcAddress != nil && len(*srcAddress) > 0 {

		notspend, _, _ := GetBalanceForAddrSelf(*srcAddress)
		item := TxItem{
			Addr:  srcAddress,
			Value: notspend,
		}

		return notspend, &item
	} else {

		var addr *crypto.AddressCoin
		var notspend uint64
		for _, one := range keystore.GetAddrAll() {
			addr = &one.Addr
			notspend, _, _ = GetBalanceForAddrSelf(one.Addr)
			if notspend < amount {
				continue
			}
			lockValue, _ := this.FindLockTotalByAddr(&one.Addr)
			notspend = notspend - lockValue
			if notspend < amount {
				continue
			}
			break
		}
		item := TxItem{
			Addr:  addr,
			Value: notspend,
		}

		return notspend, &item
	}
}

func (this *BalanceManager) CountBalanceForBlock(bhvo *BlockHeadVO) {
	this.countBlock(bhvo)

	db.LevelDB.Save([]byte(config.BlockHeight+strconv.Itoa(int(bhvo.BH.Height))), &bhvo.BH.Hash)
}

func (this *BalanceManager) run() {
	for bhvo := range this.syncBlockHead {
		this.countBlock(bhvo)
	}
}

func (this *BalanceManager) countBlock(bhvo *BlockHeadVO) {
	this.countTotal++
	if this.countTotal == bhvo.BH.Height {
		engine.Log.Info("count block group:%d height:%d total:%d", bhvo.BH.GroupHeight, bhvo.BH.Height, this.countTotal)
	} else {
		engine.Log.Info("count block group:%d height:%d total:%d Unequal", bhvo.BH.GroupHeight, bhvo.BH.Height, this.countTotal)
	}

	SaveTxToBlockHead(bhvo)

	CountBalanceOther(this.otherDeposit, bhvo)

	this.countDepositAndVote(bhvo)

	countCommunityVoteReward(bhvo)

	this.CleanCacheTx(bhvo)

	start := time.Now()

	this.countBalancesNew(bhvo)
	engine.Log.Info("countBalancesNew time:%s", time.Now().Sub(start))
	if time.Now().Sub(start).Seconds() > 1 {
		engine.Log.Info("countBalancesNew time too long")
	}

	this.countTxHistory(bhvo)

	this.Unfrozen(bhvo.BH.Height-1, bhvo.BH.Time)

	config.UpdateSoreName(bhvo.BH.Height)
}

func (this *BalanceManager) CleanCacheTx(bhvo *BlockHeadVO) {
	if bhvo.BH.Height == config.Mining_block_start_height {
		return
	}
	for _, one := range bhvo.Txs {
		if one.Class() == config.Wallet_tx_type_mining {
			continue
		}
		this.DelLockTx(one)
	}
}

func (this *BalanceManager) countDepositAndVote(bhvo *BlockHeadVO) {

	for _, txItr := range bhvo.Txs {

		txItr.BuildHash()

		txCtrl := GetTransactionCtrl(txItr.Class())
		if txCtrl != nil {

			return
		}

		switch txItr.Class() {
		case config.Wallet_tx_type_mining:
		case config.Wallet_tx_type_deposit_in:
			voutOne := (*txItr.GetVout())[0]
			SetDepositWitnessAddr(&voutOne.Address, voutOne.Value)

			if !voutOne.CheckIsSelf() {
				continue
			}
			this.depositWitness = &TxItem{
				Addr:  &voutOne.Address,
				Value: voutOne.Value,
			}

			if config.InitNode {
				config.AlreadyMining = true
			}

			if config.SubmitDepositin {
				config.AlreadyMining = true
			}
		case config.Wallet_tx_type_deposit_out:
			vin := (*txItr.GetVin())[0]
			RemoveDepositWitnessAddr(vin.GetPukToAddr())

			if !vin.CheckIsSelf() {
				continue
			}
			this.depositWitness = nil
		case config.Wallet_tx_type_pay:
		case config.Wallet_tx_type_vote_in:
			voutOne := (*txItr.GetVout())[0]
			voteIn := txItr.(*Tx_vote_in)
			switch voteIn.VoteType {
			case VOTE_TYPE_community:
				SetDepositCommunityAddr(&voutOne.Address, voutOne.Value)

				if !voutOne.CheckIsSelf() {
					continue
				}
				txItem := DepositInfo{
					WitnessAddr: voteIn.Vote,
					SelfAddr:    voutOne.Address,
					Value:       voutOne.Value,
				}

				this.depositCommunity.Store(utils.Bytes2string(voutOne.Address), &txItem)

				blockhash, _ := db.GetTxToBlockHash(txItr.GetHash())
				db.LevelTempDB.Save(config.BuildCommunityAddrStartHeight(voutOne.Address), blockhash)

			case VOTE_TYPE_vote:
				value := GetDepositLightVoteValue(&voutOne.Address, &voteIn.Vote)
				value += voutOne.Value
				SetDepositLightVoteValue(&voutOne.Address, &voteIn.Vote, value)
				SetVoteAddr(&voutOne.Address, &voteIn.Vote)

				if !voutOne.CheckIsSelf() {
					continue
				}

				lightAddr := utils.Bytes2string(voutOne.Address)
				itemItr, ok := this.depositVote.Load(lightAddr)
				if ok {
					item := itemItr.(*DepositInfo)
					item.Value += voutOne.Value
				} else {
					txItem := DepositInfo{
						WitnessAddr: voteIn.Vote,
						SelfAddr:    voutOne.Address,
						Value:       voutOne.Value,
					}
					this.depositVote.Store(lightAddr, &txItem)
				}
			case VOTE_TYPE_light:
				SetDepositLightAddr(&voutOne.Address, voutOne.Value)

				if !voutOne.CheckIsSelf() {
					continue
				}
				txItem := DepositInfo{
					WitnessAddr: voteIn.Vote,
					SelfAddr:    voutOne.Address,
					Value:       voutOne.Value,
				}
				this.depositLight.Store(utils.Bytes2string(voutOne.Address), &txItem)
			}
		case config.Wallet_tx_type_vote_out:
			vin := (*txItr.GetVin())[0]
			voutOne := (*txItr.GetVout())[0]
			voteOut := txItr.(*Tx_vote_out)
			switch voteOut.VoteType {
			case VOTE_TYPE_community:
				RemoveDepositCommunityAddr(vin.GetPukToAddr())

				if !vin.CheckIsSelf() {
					continue
				}
				this.depositCommunity.Delete(utils.Bytes2string(*vin.GetPukToAddr()))
			case VOTE_TYPE_vote:
				value := GetDepositLightVoteValue(vin.GetPukToAddr(), &voteOut.Vote)
				value -= voutOne.Value
				if value == 0 {
					RemoveDepositLightVoteValue(vin.GetPukToAddr(), &voteOut.Vote)
					RemoveVoteAddr(vin.GetPukToAddr())
				} else {
					SetDepositLightVoteValue(vin.GetPukToAddr(), &voteOut.Vote, value)
				}

				if !voutOne.CheckIsSelf() {
					continue
				}
				lightAddr := utils.Bytes2string(*vin.GetPukToAddr())
				itemItr, ok := this.depositVote.Load(lightAddr)
				if ok {
					item := itemItr.(*DepositInfo)
					item.Value -= voutOne.Value
					if item.Value == 0 {
						this.depositVote.Delete(lightAddr)
					}
				}
			case VOTE_TYPE_light:
				RemoveDepositLightAddr(vin.GetPukToAddr())

				if !vin.CheckIsSelf() {
					continue
				}
				this.depositLight.Delete(utils.Bytes2string(*vin.GetPukToAddr()))
			}

		}

	}
}

func countCommunityVoteReward(bhvo *BlockHeadVO) {

	for _, txItr := range bhvo.Txs {

		switch txItr.Class() {
		case config.Wallet_tx_type_mining:
			for _, vout := range *txItr.GetVout() {

				if !ExistCommunityAddr(&vout.Address) {

					continue
				}
				value := GetCommunityVoteRewardFrozen(&vout.Address)
				SetCommunityVoteRewardFrozen(&vout.Address, value+vout.Value)

			}
		case config.Wallet_tx_type_vote_in:
			voteIn := txItr.(*Tx_vote_in)
			if voteIn.VoteType == VOTE_TYPE_community {

				addr := voteIn.Vout[0].Address

				db.LevelTempDB.Save(config.BuildCommunityAddrStartHeight(addr), &bhvo.BH.Hash)

				SetCommunityAddr(&addr)

				if !ExistCommunityVoteRewardFrozen(&addr) {
					SetCommunityVoteRewardFrozen(&addr, 0)
				}
			}

		case config.Wallet_tx_type_vote_out:
			vinOne := (*txItr.GetVin())[0]
			voteOut := txItr.(*Tx_vote_out)

			if voteOut.VoteType == VOTE_TYPE_community {
				addr := vinOne.GetPukToAddr()
				RemoveCommunityAddr(addr)
			}

		case config.Wallet_tx_type_voting_reward:
			vinOne := (*txItr.GetVin())[0]
			communitAddr := vinOne.GetPukToAddr()
			totalValue := txItr.GetGas()
			for _, voutOne := range *txItr.GetVout() {

				totalValue += voutOne.Value
			}
			value := GetCommunityVoteRewardFrozen(communitAddr)

			SetCommunityVoteRewardFrozen(communitAddr, value-totalValue)

			if value-totalValue == 0 {

			}

		}

	}
}

func (this *BalanceManager) countBalancesNew(bhvo *BlockHeadVO) {

	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, 0),
		Nonce:    make(map[string]big.Int),
	}
	itemsChan := make(chan *TxItemCountMap, len(bhvo.Txs))

	start := time.Now()

	wg := new(sync.WaitGroup)
	wg.Add(len(bhvo.Txs))
	utils.Go(
		func() {

			for i := 0; i < len(bhvo.Txs); i++ {
				one := <-itemsChan
				if one != nil {

					for addrStr, itemMap := range one.AddItems {
						oldItemMap, ok := itemCount.AddItems[addrStr]
						if ok {
							for frozenHeight, value := range *itemMap {
								oldValue, ok := (*oldItemMap)[frozenHeight]
								if ok {
									oldValue += value
									(*oldItemMap)[frozenHeight] = oldValue
								} else {
									(*oldItemMap)[frozenHeight] = value
								}
							}
						} else {
							itemCount.AddItems[addrStr] = itemMap
						}
					}

					for addrStr, nonce := range one.Nonce {
						oldNonce, ok := itemCount.Nonce[addrStr]
						if ok {
							if oldNonce.Cmp(&nonce) < 0 {
								itemCount.Nonce[addrStr] = nonce
							}
						} else {
							itemCount.Nonce[addrStr] = nonce
						}
					}
				}
				wg.Done()
			}
		})

	NumCPUTokenChan := make(chan bool, runtime.NumCPU()*6)
	for _, txItr := range bhvo.Txs {
		go this.countBalancesNewOne(txItr, bhvo.BH.Height, NumCPUTokenChan, itemsChan)
	}

	wg.Wait()
	engine.Log.Info("count block spend time:%s", time.Now().Sub(start))
	start = time.Now()

	for addrStr, itemMap := range itemCount.AddItems {
		addr := crypto.AddressCoin([]byte(addrStr))
		_, oldvalue := GetNotSpendBalance(&addr)
		for frozenHeight, value := range *itemMap {
			if frozenHeight <= bhvo.BH.Height {
				if value > 0 {

					oldvalue += uint64(value)

					SetNotSpendBalance(&addr, oldvalue)
				} else if value < 0 {

					oldvalue -= uint64(1 - value - 1)

					SetNotSpendBalance(&addr, oldvalue)
				}
			} else {
				oldValue := GetAddrFrozenValue(&addr)

				oldValue += uint64(value)
				SetAddrFrozenValue(&addr, oldValue)
				AddFrozenHeight(&addr, frozenHeight, uint64(value))
			}
		}
	}
	engine.Log.Info("count block spend time:%s", time.Now().Sub(start))

	var err error
	for addrStr, nonce := range itemCount.Nonce {
		addr := crypto.AddressCoin([]byte(addrStr))
		err = SetAddrNonce(&addr, &nonce)
		if err != nil {
			engine.Log.Error("SetAddrNonce error:%s", err.Error())
		}
	}

}

func (this *BalanceManager) countBalancesNewOne(txItr TxItr, height uint64, tokenCPU chan bool, itemChan chan *TxItemCountMap) {
	tokenCPU <- false
	txItr.BuildHash()

	itemCount := txItr.CountTxItemsNew(height)

	itemChan <- itemCount
	<-tokenCPU
}

func (this *BalanceManager) countCommunityReward(bhvo *BlockHeadVO) {

	var err error
	var addr crypto.AddressCoin
	var ok bool
	var cs *CommunitySign
	var sn *sqlite3_db.SnapshotReward
	var rt *RewardTotal
	var r *[]sqlite3_db.RewardLight
	for _, txItr := range bhvo.Txs {

		if txItr.Class() != config.Wallet_tx_type_pay {
			continue
		}

		addr, ok, cs = CheckPayload(txItr)
		if !ok {

			continue
		}

		_, ok = keystore.FindAddress(addr)
		if !ok {

			continue
		}

		sn, _, err = FindNotSendReward(&addr)
		if err != nil && err.Error() != xorm.ErrNotExist.Error() {
			engine.Log.Error("querying database Error %s", err.Error())
			return
		}

		if sn == nil || sn.EndHeight < cs.EndHeight {

			rt, r, err = GetRewardCount(&addr, cs.StartHeight, cs.EndHeight)
			if err != nil {
				return
			}
			err = CreateRewardCount(addr, rt, *r)
			if err != nil {
				return
			}
		}
	}

}

func (this *BalanceManager) countTxHistory(bhvo *BlockHeadVO) {
	for _, txItr := range bhvo.Txs {

		txItr.BuildHash()

		txItr.CountTxHistory(bhvo.BH.Height)

	}
}

func (this *BalanceManager) DepositIn(amount, gas uint64, pwd, payload string) error {

	if this.depositWitness != nil {
		return errors.New("Deposit cannot be paid repeatedly")
	}

	if amount != config.Mining_deposit {
		return errors.New("Deposit not less than" + strconv.Itoa(int(uint64(config.Mining_deposit)/Unit)))
	}

	deposiIn, err := CreateTxDepositIn(amount, gas, pwd, payload)
	if err != nil {
		return err
	}
	if deposiIn == nil {
		return errors.New("Failure to pay deposit")
	}
	deposiIn.BuildHash()
	MulticastTx(deposiIn)

	this.txManager.AddTx(deposiIn)
	return nil
}

func (this *BalanceManager) DepositOut(addr string, amount, gas uint64, pwd string) error {

	if this.depositWitness == nil {
		return errors.New("I didn't pay the deposit")
	}

	deposiOut, err := CreateTxDepositOut(addr, amount, gas, pwd)
	if err != nil {
		return err
	}
	if deposiOut == nil {
		return errors.New("Failure to pay deposit")
	}
	deposiOut.BuildHash()

	MulticastTx(deposiOut)

	this.txManager.AddTx(deposiOut)
	return nil
}

func (this *BalanceManager) VoteIn(voteType uint16, witnessAddr crypto.AddressCoin, addr crypto.AddressCoin, amount, gas uint64, pwd, payload string) error {

	if bytes.Equal(witnessAddr, addr) {
		return errors.New("You can't vote for yourself")
	}
	dstAddr := addr

	isWitness := this.witnessBackup.haveWitness(&dstAddr)
	_, isCommunity := this.witnessBackup.haveCommunityList(&dstAddr)
	_, isLight := this.witnessBackup.haveLight(&dstAddr)

	switch voteType {
	case 1:
		if isLight || isWitness {
			return errors.New("The voting address is already another role")
		}
		vs, ok := this.witnessBackup.haveCommunityList(&dstAddr)
		if ok {
			if bytes.Equal(*vs.Witness, witnessAddr) {
				return errors.New("Can't vote again")
			}
			return errors.New("Cannot vote for multiple witnesses")
		}

		if amount != config.Mining_vote {
			return errors.New("Community node deposit is " + strconv.Itoa(int(config.Mining_vote/1e8)))
		}

	case 2:

		if isCommunity || isWitness {

			return errors.New("The voting address is already another role")
		}

		if !isLight {

			return errors.New("Become a light node first")
		}

		vs, ok := this.witnessBackup.haveVoteList(&dstAddr)
		if ok {
			if !bytes.Equal(*vs.Witness, witnessAddr) {

				return errors.New("Cannot vote for multiple community nodes")
			}
		}

	case 3:
		if isCommunity || isWitness {

			return errors.New("The voting address is already another role")
		}
		if isLight {

			return errors.New("It's already a light node")
		}

		if amount != config.Mining_light_min {

			return errors.New("Light node deposit is " + strconv.Itoa(int(config.Mining_light_min/1e8)))
		}
		witnessAddr = nil
	default:

		return errors.New("Unrecognized voting type")

	}

	voetIn, err := CreateTxVoteIn(voteType, witnessAddr, addr, amount, gas, pwd, payload)
	if err != nil {
		return err
	}
	if voetIn == nil {

		return errors.New("Failure to pay deposit")
	}
	voetIn.BuildHash()

	MulticastTx(voetIn)
	this.txManager.AddTx(voetIn)

	return nil
}

func (this *BalanceManager) VoteOut(voteType uint16, addr crypto.AddressCoin, amount, gas uint64, pwd, payload string) error {

	deposiOut, err := CreateTxVoteOut(voteType, addr, amount, gas, pwd, payload)
	if err != nil {
		return err
	}
	if deposiOut == nil {

		return errors.New("Failure to pay deposit")
	}
	deposiOut.BuildHash()
	MulticastTx(deposiOut)
	this.txManager.AddTx(deposiOut)
	return nil
}

func (this *BalanceManager) BuildOtherTx(class uint64, srcAddr, addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string, params ...interface{}) (TxItr, error) {

	ctrl := GetTransactionCtrl(class)
	txItr, err := ctrl.BuildTx(this.otherDeposit, srcAddr, addr, amount, gas, frozenHeight, pwd, comment, params...)
	if err != nil {
		return nil, err
	}
	txItr.BuildHash()
	MulticastTx(txItr)

	this.txManager.AddTx(txItr)
	return txItr, nil
}
