package mining

import (
	"bytes"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"sort"
	"sync"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type WitnessBackup struct {
	chain             *Chain
	lock              sync.RWMutex
	witnesses         []*BackupWitness
	witnessesMap      sync.Map
	VoteCommunity     sync.Map
	VoteCommunityList sync.Map
	Vote              sync.Map
	VoteList          sync.Map
	LightNode         sync.Map
	Blacklist         sync.Map
}

func (this WitnessBackup) Len() int {
	return len(this.witnesses)
}

func (this WitnessBackup) Less(i, j int) bool {
	return this.witnesses[i].VoteNum > this.witnesses[j].VoteNum
}

func (this WitnessBackup) Swap(i, j int) {
	this.witnesses[i], this.witnesses[j] = this.witnesses[j], this.witnesses[i]
}

func (this *WitnessBackup) GetBackupWitnessTotal() uint64 {
	this.lock.RLock()
	total := len(this.witnesses)
	this.lock.RUnlock()
	return uint64(total)
}

func (this *WitnessBackup) CountWitness(txs *[]TxItr) {

	for _, one := range *txs {
		switch one.Class() {

		case config.Wallet_tx_type_deposit_in:

			vout := (*one.GetVout())[0]
			score := vout.Value
			depositIn := one.(*Tx_deposit_in)
			this.addWitness(depositIn.Puk, score)

			if depositIn.Payload != nil && len(depositIn.Payload) > 0 {

				witnessAddr := crypto.BuildAddr(config.AddrPre, depositIn.Puk)

				SaveWitnessName(witnessAddr, string(one.GetPayload()))
			}
		case config.Wallet_tx_type_deposit_out:

			vinOne := (*one.GetVin())[0]
			addr := vinOne.GetPukToAddr()
			this.DelWitness(addr)
			this.DelBlackList(*addr)

			witnessName := FindWitnessName(*addr)
			if witnessName == "" {
				continue
			}

			DelWitnessName(witnessName)

		case config.Wallet_tx_type_vote_in:

			vout := (*one.GetVout())[0]
			score := vout.Value
			votein := one.(*Tx_vote_in)

			addr := votein.Vote
			this.addVote(votein.VoteType, &addr, &vout.Address, score)
			if addr != nil {
				votein.SetVoteAddr(addr)
			}

			if votein.Payload != nil && len(votein.Payload) > 0 {
				if votein.VoteType == VOTE_TYPE_community {

					SaveWitnessName(vout.Address, string(votein.Payload))
				}
			}

		case config.Wallet_tx_type_vote_out:
			vinOne := (*one.GetVin())[0]
			addr := vinOne.GetPukToAddr()
			voteOut := one.(*Tx_vote_out)
			voutOne := voteOut.Vout[0]
			this.DelVote(voteOut.VoteType, &voteOut.Vote, addr, voutOne.Value)

			if voteOut.VoteType == VOTE_TYPE_community {

				witnessName := FindWitnessName(*addr)
				if witnessName == "" {
					continue
				}

				DelWitnessName(witnessName)
			}

		}
	}

}

func (this *WitnessBackup) FindWitness(witnessAddr crypto.AddressCoin) bool {
	_, ok := this.witnessesMap.Load(utils.Bytes2string(witnessAddr))
	return ok
}

func (this *WitnessBackup) addWitness(puk []byte, score uint64) {
	witnessAddr := crypto.BuildAddr(config.AddrPre, puk)

	_, ok := this.witnessesMap.Load(utils.Bytes2string(witnessAddr))
	if ok {

		return
	}

	witness := &BackupWitness{
		Addr:  &witnessAddr,
		Puk:   puk,
		Score: score,
	}

	_, lvss := this.FindScore(&witnessAddr)
	for _, one := range lvss {

		witness.VoteNum = witness.VoteNum + one.Vote
	}

	this.lock.Lock()
	this.witnesses = append(this.witnesses, witness)
	this.lock.Unlock()
	this.witnessesMap.Store(utils.Bytes2string(witnessAddr), witness)
}

func (this *WitnessBackup) DelWitness(witnessAddr *crypto.AddressCoin) {
	this.lock.Lock()

	for i, one := range this.witnesses {
		if !bytes.Equal(*witnessAddr, *one.Addr) {
			continue
		}
		temp := this.witnesses[:i]
		this.witnesses = append(temp, this.witnesses[i+1:]...)
		break
	}

	this.lock.Unlock()
	this.witnessesMap.Delete(utils.Bytes2string(*witnessAddr))
}

func (this *WitnessBackup) addVote(voteType uint16, witnessAddr, voteAddr *crypto.AddressCoin, score uint64) {

	if bytes.Equal(*witnessAddr, *voteAddr) {
		return
	}

	isWitness := this.haveWitness(voteAddr)
	_, isCommunity := this.haveCommunityList(voteAddr)
	_, isLight := this.haveLight(voteAddr)

	newVote := new(VoteScore)
	newVote.Witness = witnessAddr
	newVote.Addr = voteAddr
	newVote.Scores = score

	switch voteType {
	case 1:
		if isLight || isWitness {

			return
		}
		if score != config.Mining_vote {

			return
		}

		vs, ok := this.haveCommunityList(voteAddr)
		if ok {
			if bytes.Equal(*vs.Witness, *witnessAddr) {
				vs.Scores = vs.Scores + score

				return
			}

			return
		}

		this.VoteCommunityList.Store(utils.Bytes2string(*voteAddr), newVote)

		v, ok := this.VoteCommunity.Load(utils.Bytes2string(*witnessAddr))
		if ok {
			vss := v.(*[]*VoteScore)
			*vss = append(*vss, newVote)
		} else {
			vss := make([]*VoteScore, 0)
			vss = append(vss, newVote)
			this.VoteCommunity.Store(utils.Bytes2string(*witnessAddr), &vss)
		}

		v, ok = this.Vote.Load(utils.Bytes2string(*voteAddr))
		if ok {
			vss := v.(*[]*VoteScore)
			voteNum := uint64(0)
			for _, one := range *vss {
				voteNum = voteNum + one.Scores
			}
			newVote.Vote = voteNum

			v, ok := this.witnessesMap.Load(utils.Bytes2string(*witnessAddr))
			if ok {
				bw := v.(*BackupWitness)
				bw.VoteNum = bw.VoteNum + voteNum
			}
		}

	case 2:

		if isCommunity || isWitness {

			return
		}
		vs, ok := this.haveVoteList(voteAddr)
		if ok {

			if !bytes.Equal(*vs.Witness, *witnessAddr) {

				return
			}

			vs.Scores = vs.Scores + score

		} else {

			this.VoteList.Store(utils.Bytes2string(*voteAddr), newVote)

			v, ok := this.Vote.Load(utils.Bytes2string(*witnessAddr))
			if ok {
				vss := v.(*[]*VoteScore)
				*vss = append(*vss, newVote)

			} else {
				vss := make([]*VoteScore, 0)
				vss = append(vss, newVote)
				this.Vote.Store(utils.Bytes2string(*witnessAddr), &vss)

			}

		}

		v, ok := this.VoteCommunityList.Load(utils.Bytes2string(*witnessAddr))
		if ok {
			vs := v.(*VoteScore)
			vs.Vote = vs.Vote + score

			if bytes.Equal(*voteAddr, config.SpecialAddrs) {

			}

			v, ok := this.witnessesMap.Load(utils.Bytes2string(*vs.Witness))
			if ok {
				bw := v.(*BackupWitness)
				bw.VoteNum = bw.VoteNum + score

			}
		}

	case 3:
		if isCommunity || isWitness {

			return
		}

		if score != config.Mining_light_min {

			return
		}

		v, ok := this.LightNode.Load(utils.Bytes2string(*voteAddr))
		if ok {
			vs := v.(*VoteScore)
			vs.Scores = vs.Scores + score
			return
		}
		this.LightNode.Store(utils.Bytes2string(*voteAddr), newVote)

	default:
		return
	}

}

func (this *WitnessBackup) DelVote(voteType uint16, witnessAddr, voteAddr *crypto.AddressCoin, score uint64) {

	switch voteType {
	case 1:

		v, ok := this.VoteCommunityList.Load(utils.Bytes2string(*voteAddr))
		if !ok {

			return
		}
		vs := v.(*VoteScore)
		vs.Scores = vs.Scores - score

		if vs.Scores > 0 {

			return
		}

		v, ok = this.Vote.Load(utils.Bytes2string(*voteAddr))
		if ok {
			vss := v.(*[]*VoteScore)
			voteNum := uint64(0)
			for _, one := range *vss {
				voteNum = voteNum + one.Scores
			}

			v, ok := this.witnessesMap.Load(utils.Bytes2string(*witnessAddr))
			if ok {
				bw := v.(*BackupWitness)
				bw.VoteNum = bw.VoteNum - voteNum
			}
		}

		db.LevelTempDB.Remove(config.BuildCommunityAddrStartHeight(*voteAddr))

		this.VoteCommunityList.Delete(utils.Bytes2string(*voteAddr))

		v, ok = this.VoteCommunity.Load(utils.Bytes2string(*witnessAddr))
		if !ok {
			return
		}
		vss := v.(*[]*VoteScore)
		for i, one := range *vss {
			if bytes.Equal(*one.Addr, *voteAddr) {
				temp := (*vss)[:i]
				temp = append(temp, (*vss)[i+1:]...)
				this.VoteCommunity.Store(utils.Bytes2string(*witnessAddr), &temp)
				break
			}
		}

	case 2:

		v, ok := this.VoteList.Load(utils.Bytes2string(*voteAddr))
		if !ok {

			return
		}
		vs := v.(*VoteScore)

		vs.Scores = vs.Scores - score

		v, ok = this.VoteCommunityList.Load(utils.Bytes2string(*witnessAddr))
		if ok {
			vs := v.(*VoteScore)
			vs.Vote = vs.Vote - score

			v, ok := this.witnessesMap.Load(utils.Bytes2string(*vs.Witness))
			if ok {
				bw := v.(*BackupWitness)
				bw.VoteNum = bw.VoteNum - score
			}
		}

		if vs.Scores > 0 {

			return
		}

		this.VoteList.Delete(utils.Bytes2string(*voteAddr))

		v, ok = this.VoteList.Load(utils.Bytes2string(*voteAddr))
		if !ok {

		} else {
			vs = v.(*VoteScore)

		}

		v, ok = this.Vote.Load(utils.Bytes2string(*witnessAddr))
		if !ok {
			return
		}
		vss := v.(*[]*VoteScore)
		for i, one := range *vss {
			if bytes.Equal(*one.Addr, *voteAddr) {
				temp := (*vss)[:i]
				temp = append(temp, (*vss)[i+1:]...)
				this.Vote.Store(utils.Bytes2string(*witnessAddr), &temp)
				break
			}
		}

	case 3:
		v, ok := this.LightNode.Load(utils.Bytes2string(*voteAddr))
		if !ok {
			return
		}
		vs := v.(*VoteScore)
		vs.Scores = vs.Scores - score

		if vs.Scores > 0 {
			return
		}

		this.LightNode.Delete(utils.Bytes2string(*voteAddr))
	}
}

func (this *WitnessBackup) haveWitness(witnessAddr *crypto.AddressCoin) (have bool) {
	this.lock.RLock()
	for _, one := range this.witnesses {
		have = bytes.Equal(*witnessAddr, *one.Addr)
		if have {
			break
		}
	}
	this.lock.RUnlock()
	return
}

func (this *WitnessBackup) haveCommunity(witnessAddr *crypto.AddressCoin) (*[]*VoteScore, bool) {
	v, ok := this.VoteCommunity.Load(utils.Bytes2string(*witnessAddr))
	if ok {
		value := v.(*[]*VoteScore)
		return value, ok
	}
	return nil, ok
}

func (this *WitnessBackup) haveCommunityList(addr *crypto.AddressCoin) (*VoteScore, bool) {
	v, ok := this.VoteCommunityList.Load(utils.Bytes2string(*addr))
	if ok {
		value := v.(*VoteScore)
		return value, ok
	}
	return nil, ok
}

func (this *WitnessBackup) haveVote(witnessAddr *crypto.AddressCoin) (*[]*VoteScore, bool) {
	v, ok := this.Vote.Load(utils.Bytes2string(*witnessAddr))
	if ok {
		value := v.(*[]*VoteScore)
		return value, ok
	}
	return nil, ok
}

func (this *WitnessBackup) haveVoteList(addr *crypto.AddressCoin) (*VoteScore, bool) {
	v, ok := this.VoteList.Load(utils.Bytes2string(*addr))
	if ok {
		value := v.(*VoteScore)
		return value, ok
	}
	return nil, ok
}

func (this *WitnessBackup) haveLight(addr *crypto.AddressCoin) (*VoteScore, bool) {
	v, ok := this.LightNode.Load(utils.Bytes2string(*addr))
	if ok {
		value := v.(*VoteScore)
		return value, ok
	}
	return nil, ok
}

type BackupWitness struct {
	Addr    *crypto.AddressCoin
	Puk     []byte
	Score   uint64
	VoteNum uint64
}

func (this *WitnessBackup) CreateWitnessGroup() []*Witness {
	if len(this.witnesses) < config.Witness_backup_min {

		return nil
	}
	wbg := this.GetWitnessListSort()

	IsBackup := this.chain.WitnessChain.FindWitness(keystore.GetCoinbase().Addr)
	if IsBackup {
		witnessAddrCoins := make([]*crypto.AddressCoin, 0)
		for _, one := range wbg.Witnesses {
			witnessAddrCoins = append(witnessAddrCoins, one.Addr)
		}

		go AddWitnessAddrNets(witnessAddrCoins)

	}

	for _, one := range wbg.Witnesses {
		total := uint64(0)
		for _, two := range one.CommunityVotes {

			total += two.Vote
		}

	}

	random := this.chain.HashRandom()

	wbg.Witnesses = OrderWitness(wbg.Witnesses, random)
	for i, _ := range wbg.Witnesses {

		wbg.Witnesses[i].WitnessBackupGroup = wbg
	}
	return wbg.Witnesses
}

func (this *WitnessBackup) PrintWitnessBackup() {

	this.lock.Lock()

	sort.Stable(this)
	this.lock.Unlock()

	for i, one := range this.witnesses {
		if i >= config.Witness_backup_max {

			break
		} else {
			engine.Log.Info("backup witness %d %s", i, one.Addr.B58String())
		}
	}
}

func (this *WitnessBackup) AddBlackList(addr crypto.AddressCoin) {
	addrStr := utils.Bytes2string(addr)

	v, ok := this.Blacklist.Load(addrStr)
	if ok {
		total := v.(uint64)

		total++

		this.Blacklist.Store(addrStr, total)

		return
	}
	this.Blacklist.Store(addrStr, uint64(1))

}

func (this *WitnessBackup) SubBlackList(addr crypto.AddressCoin) {
	addrStr := utils.Bytes2string(addr)

	_, ok := this.witnessesMap.Load(utils.Bytes2string(addr))
	if !ok {
		return
	}

	v, ok := this.Blacklist.Load(addrStr)
	if ok {
		total := v.(uint64)
		total--

		if total <= 0 {
			this.Blacklist.Delete(addrStr)
		} else {
			this.Blacklist.Store(addrStr, total)
		}
	}

}

func (this *WitnessBackup) DelBlackList(addr crypto.AddressCoin) {

	this.Blacklist.Delete(utils.Bytes2string(addr))
}

func NewWitnessBackup(chain *Chain) *WitnessBackup {
	wb := WitnessBackup{
		chain:             chain,
		lock:              *new(sync.RWMutex),
		witnesses:         make([]*BackupWitness, 0),
		witnessesMap:      *new(sync.Map),
		Vote:              *new(sync.Map),
		VoteList:          *new(sync.Map),
		VoteCommunity:     *new(sync.Map),
		VoteCommunityList: *new(sync.Map),
		LightNode:         *new(sync.Map),
	}
	return &wb
}

type VoteScore struct {
	Witness *crypto.AddressCoin
	Addr    *crypto.AddressCoin
	Scores  uint64
	Vote    uint64
}

func (this *WitnessBackup) FindWitnessInBlackList(addr crypto.AddressCoin) (have bool) {

	total, ok := this.Blacklist.Load(utils.Bytes2string(addr))

	if !ok {
		return false
	}
	t := total.(uint64)
	if config.CheckAddBlacklist(this.GetBackupWitnessTotal(), t) {
		return false
	}
	return true

}

func (this *WitnessBackup) GetWitnessListSort() *WitnessBackupGroup {
	currentHeight := this.chain.GetCurrentBlock()

	this.Blacklist.Range(func(k, v interface{}) bool {
		total := v.(uint64)
		ok := config.CheckAddBlacklist(this.GetBackupWitnessTotal(), total)
		if currentHeight < config.CheckAddBlacklistChangeHeight {
			ok = !ok
		}
		if ok {

			addrStr := k.(string)

			addr := crypto.AddressCoin([]byte(addrStr))

			this.DelWitness(&addr)

		}
		return true
	})

	this.lock.Lock()

	sort.Stable(this)
	this.lock.Unlock()

	wbg := WitnessBackupGroup{
		Witnesses:     make([]*Witness, 0),
		WitnessBackup: make([]*Witness, 0),
	}

	for i, _ := range this.witnesses {
		newWitness := new(Witness)
		newWitness.Addr = this.witnesses[i].Addr
		newWitness.Puk = this.witnesses[i].Puk
		newWitness.Score = this.witnesses[i].Score
		newWitness.CommunityVotes, newWitness.Votes = this.FindScore(newWitness.Addr)
		newWitness.VoteNum = this.witnesses[i].VoteNum
		newWitness.StopMining = make(chan bool, 1)

		if i < config.Witness_backup_max {
			wbg.Witnesses = append(wbg.Witnesses, newWitness)
		} else {
			wbg.WitnessBackup = append(wbg.WitnessBackup, newWitness)
		}
	}

	return &wbg
}

func (this *WitnessBackup) FindScore(addr *crypto.AddressCoin) ([]*VoteScore, []*VoteScore) {

	vssAll := make([]*VoteScore, 0)
	communityScore := make([]*VoteScore, 0)

	v, ok := this.VoteCommunity.Load(utils.Bytes2string(*addr))
	if !ok {

		return communityScore, vssAll
	}
	vss := v.(*[]*VoteScore)

	for _, one := range *vss {
		vs := VoteScore{
			Witness: one.Witness,
			Addr:    one.Addr,
			Scores:  one.Scores,
			Vote:    one.Vote,
		}
		communityScore = append(communityScore, &vs)

	}

	for _, one := range *vss {
		v, ok := this.Vote.Load(utils.Bytes2string(*one.Addr))
		if !ok {
			continue
		}
		vss := v.(*[]*VoteScore)
		voteNum := uint64(0)
		for _, one := range *vss {
			voteNum = voteNum + one.Scores
			vs := VoteScore{
				Witness: one.Witness,
				Addr:    one.Addr,
				Vote:    one.Scores,
			}
			vssAll = append(vssAll, &vs)

		}
		one.Vote = voteNum

	}
	return communityScore, vssAll
}

type VoteScoreVO struct {
	Witness string
	Addr    string
	Payload string
	Score   uint64
	Vote    uint64
}

func (this *WitnessBackup) GetCommunityListSort() []*VoteScoreVO {
	vssVO := make([]*VoteScoreVO, 0)
	this.VoteCommunityList.Range(func(k, v interface{}) bool {
		vsOne := v.(*VoteScore)

		voteNum := uint64(0)
		v, ok := this.Vote.Load(utils.Bytes2string(*vsOne.Addr))
		if ok {
			vss := v.(*[]*VoteScore)
			for _, one := range *vss {
				voteNum = voteNum + one.Scores
			}
		}

		name := FindWitnessName(*vsOne.Addr)

		vsVOOne := VoteScoreVO{
			Witness: vsOne.Witness.B58String(),
			Addr:    vsOne.Addr.B58String(),
			Payload: name,
			Score:   vsOne.Scores,
			Vote:    voteNum,
		}

		vssVO = append(vssVO, &vsVOOne)
		return true
	})
	return vssVO
}
