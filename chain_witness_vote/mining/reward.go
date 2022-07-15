package mining

import (
	"bytes"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type WitnessBackupGroup struct {
	Witnesses     []*Witness
	WitnessBackup []*Witness
}

func (this *WitnessBackupGroup) CountRewardWitness(blockhash *[]byte, height uint64) *[]*Witness {

	rewardWitness := make([]*Witness, 0)
	preBlockHash := blockhash
	if height < config.Reward_witness_height_new {
		for i := len(this.Witnesses); i > 0; i-- {
			witnessOne := this.Witnesses[i-1]

			if witnessOne.Block == nil {

				continue
			}
			if bytes.Equal(*preBlockHash, witnessOne.Block.Id) {

				rewardWitness = append(rewardWitness, witnessOne)
				preBlockHash = &witnessOne.Block.PreBlockID
			} else {

			}
		}
	} else {

		var lastWitness *Witness
		for _, one := range this.Witnesses {
			if one.Block == nil {
				continue
			}
			if bytes.Equal(*preBlockHash, one.Block.Id) {
				lastWitness = one
				rewardWitness = append(rewardWitness, one)
				preBlockHash = &one.Block.PreBlockID

				break
			}
		}
		if lastWitness != nil {
			for lastWitness = lastWitness.PreWitness; lastWitness != nil && lastWitness.WitnessBackupGroup == this; lastWitness = lastWitness.PreWitness {
				if lastWitness.Block == nil {

					continue
				}
				if bytes.Equal(*preBlockHash, lastWitness.Block.Id) {

					rewardWitness = append(rewardWitness, lastWitness)
					preBlockHash = &lastWitness.Block.PreBlockID
				} else {

				}
			}
		}
	}
	return &rewardWitness
}

func (this *WitnessBackupGroup) BuildRewardVouts(blocks []Block, height uint64, blockhash *[]byte, preBlock *Block) []*Vout {

	vouts := make([]*Vout, 0)

	witneses := make([]*crypto.AddressCoin, 0)

	allCommiunty := make([]*VoteScore, 0)

	allVotePos := uint64(0)
	allGas := uint64(0)
	allReward := uint64(0)

	if height >= config.Reward_witness_height {
		witnessAll := this.CountRewardWitness(&preBlock.Id, height)
		for _, one := range *witnessAll {

			oneReward := config.ClacRewardForBlockHeight(one.Block.Height)

			allReward = allReward + oneReward

			_, txs, _ := one.Block.LoadTxs()
			for _, one := range *txs {
				allGas = allGas + one.GetGas()

			}

			witneses = append(witneses, one.Addr)
		}

		for _, one := range this.Witnesses {

			for _, v := range one.CommunityVotes {
				v.Scores = 0
				allCommiunty = append(allCommiunty, v)
				allVotePos = allVotePos + v.Vote

			}
		}

	} else {

		for _, one := range this.Witnesses {

			for _, v := range one.CommunityVotes {
				v.Scores = 0
				allCommiunty = append(allCommiunty, v)
				allVotePos = allVotePos + v.Vote

			}

			isUnconfirmed := false

			nowWitnessGroup := GetLongChain().WitnessChain.witnessGroup
			if nowWitnessGroup.Height == one.Group.Height {
				for _, oneBlock := range blocks {

					if oneBlock.Group.Height == one.Group.Height && bytes.Equal(*one.Addr, *oneBlock.witness.Addr) {
						isUnconfirmed = true
						break
					}
				}
				if !isUnconfirmed {

					continue
				}
			}

			if !isUnconfirmed {

				if one.Block == nil {
					continue
				}

				if one.Block.Group == nil {
					continue
				}

				ok, group := one.Group.CheckBlockGroup(blockhash)
				if !ok {
					continue
				}

				if one.Block.Group != group {
					continue
				}

			}

			witneses = append(witneses, one.Addr)

			_, txs, _ := one.Block.LoadTxs()
			for _, one := range *txs {
				allGas = allGas + one.GetGas()

			}

			oneReward := config.ClacRewardForBlockHeight(one.Block.Height)
			allReward = allReward + oneReward

		}
	}

	allReward = allReward + allGas

	foundationReward := uint64(0)

	nameinfo := name.FindNameToNet(config.Name_Foundation)
	if nameinfo != nil && len(nameinfo.AddrCoins) != 0 {

		temp := new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(6)))
		value := new(big.Int).Div(temp, big.NewInt(int64(100)))
		foundationReward = value.Uint64()

		addrCoin := nameinfo.AddrCoins[utils.GetRandNum(int64(len(nameinfo.AddrCoins)))]

		voutsOne := LinearRelease180Day(addrCoin, foundationReward, height)
		vouts = append(vouts, voutsOne...)

	}

	investorReward := uint64(0)

	nameinfo = name.FindNameToNet(config.Name_investor)
	if nameinfo != nil && len(nameinfo.AddrCoins) != 0 {

		temp := new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(8)))
		value := new(big.Int).Div(temp, big.NewInt(int64(100)))
		investorReward = value.Uint64()

		addrCoin := nameinfo.AddrCoins[utils.GetRandNum(int64(len(nameinfo.AddrCoins)))]

		voutsOne := LinearRelease180Day(addrCoin, investorReward, height)
		vouts = append(vouts, voutsOne...)

	}

	teamReward := uint64(0)

	nameinfo = name.FindNameToNet(config.Name_team)
	if nameinfo != nil && len(nameinfo.AddrCoins) != 0 {

		temp := new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(14)))
		value := new(big.Int).Div(temp, big.NewInt(int64(100)))
		teamReward = value.Uint64()

		addrCoin := nameinfo.AddrCoins[utils.GetRandNum(int64(len(nameinfo.AddrCoins)))]

		voutsOne := LinearRelease180Day(addrCoin, teamReward, height)
		vouts = append(vouts, voutsOne...)

	}

	storeReward := uint64(0)

	nameinfo = name.FindNameToNet(config.Name_store)
	if nameinfo != nil && len(nameinfo.AddrCoins) != 0 {

		temp := new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(39)))
		value := new(big.Int).Div(temp, big.NewInt(int64(100)))
		storeReward = value.Uint64()

		addrCoin := nameinfo.AddrCoins[utils.GetRandNum(int64(len(nameinfo.AddrCoins)))]

		vout := &Vout{
			Value:   storeReward,
			Address: addrCoin,
		}

		vouts = append(vouts, vout)

	}

	temp := new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(66)))
	value := new(big.Int).Div(temp, big.NewInt(int64(10000)))
	witnessReward99 := value.Uint64()

	temp = new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(264)))
	value = new(big.Int).Div(temp, big.NewInt(int64(1000)))
	witnessReward31 := value.Uint64()

	temp = new(big.Int).Mul(big.NewInt(int64(allReward)), big.NewInt(int64(297)))
	value = new(big.Int).Div(temp, big.NewInt(int64(1000)))
	communityReward := value.Uint64()

	surplus := allReward - (foundationReward + investorReward + teamReward + storeReward + witnessReward99 + witnessReward31 + communityReward)
	if surplus > 0 {
		witnessReward31 = witnessReward31 + surplus
	}

	use := uint64(0)
	oneReward := uint64(0)

	{
		if height < config.Mining_witness_average_height {
			use = uint64(0)
			temp = new(big.Int).Mul(big.NewInt(int64(witnessReward99)), big.NewInt(int64(1)))
			value = new(big.Int).Div(temp, big.NewInt(int64(config.Witness_backup_reward_max)))
			oneReward = value.Uint64()
			for i, one := range append(this.Witnesses, this.WitnessBackup...) {
				if i >= config.Witness_backup_reward_max {
					break
				}

				use = use + oneReward
				voutsOne := LinearRelease180Day(*one.Addr, oneReward, height)
				vouts = append(vouts, voutsOne...)

			}

			if len(vouts) > 0 {
				vouts[len(vouts)-1].Value = vouts[len(vouts)-1].Value + (witnessReward99 - use)
			}
		} else {

			averageBackupTotalMax := config.Witness_backup_reward_max - config.Witness_backup_max
			witnessBackupTotal := len(this.WitnessBackup)
			averageWitness := this.Witnesses
			if witnessBackupTotal >= averageBackupTotalMax {
				averageWitness = append(averageWitness, this.WitnessBackup[:averageBackupTotalMax]...)
			} else if witnessBackupTotal > 0 && witnessBackupTotal < averageBackupTotalMax {
				averageWitness = append(averageWitness, this.WitnessBackup...)
			}

			use = uint64(0)
			temp = new(big.Int).Mul(big.NewInt(int64(witnessReward99)), big.NewInt(int64(1)))
			value = new(big.Int).Div(temp, big.NewInt(int64(len(averageWitness))))
			oneReward = value.Uint64()
			for _, one := range averageWitness {

				use = use + oneReward
				voutsOne := LinearRelease180Day(*one.Addr, oneReward, height)
				vouts = append(vouts, voutsOne...)

			}

			if len(vouts) > 0 {
				vouts[len(vouts)-1].Value = vouts[len(vouts)-1].Value + (witnessReward99 - use)
			}
		}
	}

	{
		use = uint64(0)
		if allVotePos <= 0 {

			temp = new(big.Int).Mul(big.NewInt(int64(witnessReward31)), big.NewInt(1))
			value = new(big.Int).Div(temp, big.NewInt(int64(len(this.Witnesses))))
			oneReward = value.Uint64()
			for _, one := range this.Witnesses {

				use = use + oneReward
				voutsOne := LinearRelease180Day(*one.Addr, oneReward, height)
				vouts = append(vouts, voutsOne...)
			}
		} else {
			for _, one := range this.Witnesses {
				temp = new(big.Int).Mul(big.NewInt(int64(witnessReward31)), big.NewInt(int64(one.VoteNum)))
				value = new(big.Int).Div(temp, big.NewInt(int64(allVotePos)))
				oneReward = value.Uint64()

				use = use + oneReward
				voutsOne := LinearRelease180Day(*one.Addr, oneReward, height)
				vouts = append(vouts, voutsOne...)
			}
		}

		if len(vouts) > 0 {

			vouts[len(vouts)-1].Value = vouts[len(vouts)-1].Value + (witnessReward31 - use)
		}
	}

	{

		use = uint64(0)

		if allVotePos <= 0 {

			temp = new(big.Int).Mul(big.NewInt(int64(communityReward)), big.NewInt(int64(1)))
			value = new(big.Int).Div(temp, big.NewInt(int64(len(this.Witnesses))))
			oneReward = value.Uint64()
			for i, _ := range this.Witnesses {
				use = use + oneReward
				vout := Vout{
					Value:   oneReward,
					Address: *this.Witnesses[i].Addr,
				}
				vouts = append(vouts, &vout)

			}
		} else {

			for i, one := range allCommiunty {

				if one.Vote == 0 {
					continue
				}
				temp = new(big.Int).Mul(big.NewInt(int64(communityReward)), big.NewInt(int64(one.Vote)))
				value = new(big.Int).Div(temp, big.NewInt(int64(allVotePos)))
				oneReward = value.Uint64()

				use = use + oneReward
				vout := Vout{
					Value:   oneReward,
					Address: *allCommiunty[i].Addr,
				}
				vouts = append(vouts, &vout)

			}
		}

		if len(vouts) > 0 {

			vouts[len(vouts)-1].Value = vouts[len(vouts)-1].Value + (communityReward - use)
		}

	}

	vouts = CleanZeroVouts(&vouts)

	vouts = MergeVouts(&vouts)

	return vouts
}

func LinearRelease180Day(addr crypto.AddressCoin, total uint64, height uint64) []*Vout {

	vouts := make([]*Vout, 0)

	first25 := new(big.Int).Div(big.NewInt(int64(total)), big.NewInt(int64(4)))

	surplus := new(big.Int).Sub(big.NewInt(int64(total)), first25)

	vout := &Vout{
		Value:   first25.Uint64(),
		Address: addr,
	}
	vouts = append(vouts, vout)

	dayOne := new(big.Int).Div(surplus, big.NewInt(int64(18))).Uint64()
	intervalHeight := 60 * 60 * 24 * 10 / 10

	totalUse := uint64(0)
	for i := 0; i < 18; i++ {
		vout := &Vout{
			Value:        dayOne,
			Address:      addr,
			FrozenHeight: height + uint64((i+1)*intervalHeight),
		}
		vouts = append(vouts, vout)
		totalUse = totalUse + dayOne
	}

	if totalUse < surplus.Uint64() {

		vouts[len(vouts)-1].Value = vouts[len(vouts)-1].Value + (surplus.Uint64() - totalUse)
	}
	return vouts
}

func (this *WitnessBackupGroup) CountRewardToWitnessGroup(blockHeight uint64, blocks []Block, preBlock *Block) *Tx_reward {

	vouts := this.BuildRewardVouts(blocks, blockHeight, nil, preBlock)

	baseCoinAddr := keystore.GetCoinbase()

	vins := make([]*Vin, 0)
	vin := Vin{
		Puk:  baseCoinAddr.Puk,
		Sign: nil,
	}
	vins = append(vins, &vin)

	var txReward *Tx_reward
	for i := uint64(0); i < 10000; i++ {
		base := TxBase{
			Type:       config.Wallet_tx_type_mining,
			Vin_total:  1,
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			LockHeight: blockHeight + i,
		}
		txReward = &Tx_reward{
			TxBase: base,
		}

		txReward.MergeVout()

		for i, one := range txReward.Vin {

			_, prk, err := keystore.GetKeyByPuk(one.Puk, config.Wallet_keystore_default_pwd)
			if err != nil {
				engine.Log.Error("build reward error:%s", err.Error())
				return nil
			}

			sign := txReward.GetSign(&prk, uint64(i))
			txReward.Vin[i].Sign = *sign

		}

		txReward.BuildHash()
		if txReward.CheckHashExist() {
			txReward = nil

			continue
		} else {
			break
		}
	}

	return txReward
}
