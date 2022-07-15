package wallet

import (
	"mmschainnewaccount/chain_witness_vote"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc"
	"mmschainnewaccount/rpc/model"

	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Index struct {
	beego.Controller
}

func (this *Index) Index() {

	this.Data["CheckKey"] = chain_witness_vote.CheckKey()

	this.TplName = "wallet/index.tpl"
}

func (this *Index) Getinfo() {
	info := model.Getinfo{
		Netid:          nil,
		TotalAmount:    config.Mining_coin_total,
		Balance:        0,
		BalanceFrozen:  0,
		BalanceLockup:  0,
		Testnet:        false,
		Blocks:         mining.GetLongChain().GetCurrentBlock(),
		Group:          0,
		StartingBlock:  mining.GetLongChain().GetStartingBlock(),
		StartBlockTime: mining.GetLongChain().GetStartBlockTime(),
		HighestBlock:   mining.GetHighestBlock(),
		CurrentBlock:   mining.GetLongChain().GetCurrentBlock(),
		PulledStates:   mining.GetLongChain().GetPulledStates(),
		BlockTime:      config.Mining_block_time,
		LightNode:      config.Mining_light_min,
		CommunityNode:  config.Mining_vote,
		WitnessNode:    config.Mining_deposit,
		NameDepositMin: config.Mining_name_deposit_min,
		AddrPre:        config.AddrPre,
		TokenBalance:   nil,
	}
	this.Data["json"] = info
	this.ServeJSON()
	return
}

func (this *Index) Block() {

	out := make(map[string]interface{})
	paramsMap := make(map[string]interface{})

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &paramsMap)
	if err != nil {
		out["Msg"] = "not find param"
		out["Code"] = 1
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	value, ok := paramsMap["height"]
	if !ok {
		out["Msg"] = "not find height param"
		out["Code"] = 1
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	height, ok := value.(float64)
	if !ok {
		out["Msg"] = "height param fail"
		out["Code"] = 1
		this.Data["json"] = out
		this.ServeJSON()
		return
	}

	bhvo := mining.BlockHeadVO{
		Txs: make([]mining.TxItr, 0),
	}

	bh := mining.LoadBlockHeadByHeight(uint64(height))

	if bh == nil {
		out["Code"] = 1
		out["Msg"] = "not find blockhead"
		this.Data["json"] = out
		this.ServeJSON()
		return
	}
	bhvo.BH = bh

	for _, one := range bh.Tx {
		txItr, e := mining.LoadTxBase(one)

		if e != nil {
			out["Msg"] = "not find tx"
			out["Code"] = 1
			this.Data["json"] = out
			this.ServeJSON()
			return
		}
		bhvo.Txs = append(bhvo.Txs, txItr)
	}

	out["Code"] = 0
	out["Data"] = bhvo
	this.Data["json"] = out
	this.ServeJSON()
	return
}

func (this *Index) GetWitnessList() {

	out := make(map[string]interface{})

	wbg := mining.GetWitnessListSort()

	wvos := make([]rpc.WitnessVO, 0)
	for _, one := range append(wbg.Witnesses, wbg.WitnessBackup...) {

		name := mining.FindWitnessName(*one.Addr)

		wvo := rpc.WitnessVO{
			Addr:            one.Addr.B58String(),
			Payload:         name,
			Score:           one.Score,
			Vote:            one.VoteNum,
			CreateBlockTime: one.CreateBlockTime,
		}
		wvos = append(wvos, wvo)
	}

	out["Code"] = 0
	out["Data"] = wvos
	this.Data["json"] = out
	this.ServeJSON()
	return
}
