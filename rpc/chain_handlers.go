package rpc

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/chain_witness_vote/mining/token/payment"
	"mmschainnewaccount/chain_witness_vote/mining/token/publish"
	"mmschainnewaccount/cloud_reward/server"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"mmschainnewaccount/sqlite3_db"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

type HistoryItemVO struct {
	GenerateId string
	IsIn       bool
	Type       uint64
	InAddr     []string
	OutAddr    []string
	Value      uint64
	Txid       string
	Height     uint64
	Payload    string
}

func GetTransactionHistoty(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	if mining.GetLongChain() == nil {
		res, err = model.Tojson(model.Success)
		return
	}

	if ok, _, _, _, _ := mining.GetWitnessStatus(); ok {
		utils.SetTimeToken(config.TIMETOKEN_GetTransactionHistoty, time.Second*5)
		if allow := utils.GetTimeToken(config.TIMETOKEN_GetTransactionHistoty, false); !allow {
			res, err = model.Tojson(model.Success)
			return
		}
	}

	id := ""
	idItr, ok := rj.Get("id")
	if ok {
		if !rj.VerifyType("id", "string") {
			res, err = model.Errcode(model.TypeWrong, "id")
			return
		}
		id = idItr.(string)
	}
	var startId *big.Int
	if id != "" {
		var ok bool
		startId, ok = new(big.Int).SetString(id, 10)
		if !ok {
			res, err = model.Errcode(model.TypeWrong, "id")
			return
		}
	}

	total := 0
	totalItr, ok := rj.Get("total")
	if ok {
		total = int(totalItr.(float64))

	}
	hivos := make([]HistoryItemVO, 0)

	chain := mining.GetLongChain()
	if chain == nil {
		res, err = model.Tojson(hivos)
		return
	}
	his := chain.GetHistoryBalance(startId, total)

	for _, one := range his {
		hivo := HistoryItemVO{
			GenerateId: one.GenerateId.String(),
			IsIn:       one.IsIn,
			Type:       one.Type,
			InAddr:     make([]string, 0),
			OutAddr:    make([]string, 0),
			Value:      one.Value,
			Txid:       hex.EncodeToString(one.Txid),
			Height:     one.Height,
			Payload:    string(one.Payload),
		}

		for _, two := range one.InAddr {
			hivo.InAddr = append(hivo.InAddr, two.B58String())
		}

		for _, two := range one.OutAddr {
			hivo.OutAddr = append(hivo.OutAddr, two.B58String())
		}

		hivos = append(hivos, hivo)
	}

	res, err = model.Tojson(hivos)
	return
}

func GetNodeTotal(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	outMap := make(map[string]interface{})

	nameinfo := name.FindName(config.Name_store)
	if nameinfo == nil {

		engine.Log.Debug("Domain name does not exist")
		outMap["total_addr"] = 0
		outMap["total_space"] = 0
		res, err = model.Tojson(outMap)
		return
	}

	have := false
	for _, one := range nameinfo.NetIds {
		if bytes.Equal(nodeStore.NodeSelf.IdInfo.Id, one) {
			have = true
			break
		}
	}

	if !have {
		engine.Log.Debug("You are not in the super node address")
		outMap["total_addr"] = config.GetSpaceTotalAddr()
		outMap["total_space"] = config.GetSpaceTotal()
		res, err = model.Tojson(outMap)
		return
	}

	outMap["total_addr"], outMap["total_space"] = server.CountStoreTotal()

	res, err = model.Tojson(outMap)
	return
}

func GetNonce(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	var addr *crypto.AddressCoin
	addrItr, ok := rj.Get("address")
	if !ok {

		res, err = model.Errcode(model.NoField, "address")
		return
	}
	addrStr := addrItr.(string)
	if addrStr != "" {
		addrMul := crypto.AddressFromB58String(addrStr)
		addr = &addrMul
	}
	if addr == nil {
		res, err = model.Errcode(model.NoField, "address")
		return
	}

	nonceInt, e := mining.GetAddrNonce(addr)
	if e != nil {
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}
	out := make(map[string]interface{})
	out["nonce"] = nonceInt.Uint64()
	res, err = model.Tojson(out)
	return
}

type WitnessInfo struct {
	IsCandidate bool
	IsBackup    bool
	IsKickOut   bool
	Addr        string
	Payload     string
	Value       uint64
}

func GetWitnessInfo(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	winfo := WitnessInfo{}

	chain := mining.GetLongChain()
	if chain == nil {
		res, err = model.Tojson(winfo)
		return
	}
	var witnessAddr crypto.AddressCoin
	winfo.IsCandidate, winfo.IsBackup, winfo.IsKickOut, witnessAddr, winfo.Value = mining.GetWitnessStatus()
	winfo.Addr = witnessAddr.B58String()

	addr := keystore.GetCoinbase()
	winfo.Payload = mining.FindWitnessName(addr.Addr)

	res, err = model.Tojson(winfo)
	return
}

type WitnessVO struct {
	Addr            string
	Payload         string
	Score           uint64
	Vote            uint64
	CreateBlockTime int64
}

func GetCandidateList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	wbg := mining.GetWitnessListSort()

	wvos := make([]WitnessVO, 0)
	for _, one := range append(wbg.Witnesses, wbg.WitnessBackup...) {

		name := mining.FindWitnessName(*one.Addr)
		engine.Log.Info(":%s", name)

		wvo := WitnessVO{
			Addr:            one.Addr.B58String(),
			Payload:         name,
			Score:           one.Score,
			Vote:            one.VoteNum,
			CreateBlockTime: one.CreateBlockTime,
		}
		wvos = append(wvos, wvo)
	}

	res, err = model.Tojson(wvos)
	return
}

func GetCommunityList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	vss := mining.GetCommunityListSort()
	res, err = model.Tojson(vss)
	return
}

type VoteInfoVO struct {
	Txid        string
	WitnessAddr string
	Value       uint64
	Height      uint64
	AddrSelf    string
	Payload     string
}

type Vinfos struct {
	infos []VoteInfoVO
}

func (this *Vinfos) Len() int {
	return len(this.infos)
}

func (this *Vinfos) Less(i, j int) bool {
	if this.infos[i].Height < this.infos[j].Height {
		return false
	} else {
		return true
	}
}

func (this *Vinfos) Swap(i, j int) {
	this.infos[i], this.infos[j] = this.infos[j], this.infos[i]
}

func GetVoteList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	vtItr, ok := rj.Get("votetype")
	if !ok {
		res, err = model.Errcode(model.NoField, "votetype")
		return
	}
	voteType := uint16(vtItr.(float64))

	var items []*mining.DepositInfo
	switch voteType {
	case mining.VOTE_TYPE_community:

		items = mining.GetDepositCommunityList()
	case mining.VOTE_TYPE_light:

		items = mining.GetDepositLightList()
	case mining.VOTE_TYPE_vote:

		items = mining.GetDepositVoteList()
	}
	vinfos := Vinfos{
		infos: make([]VoteInfoVO, 0, len(items)),
	}

	for _, item := range items {
		name := mining.FindWitnessName(item.WitnessAddr)
		viVO := VoteInfoVO{

			WitnessAddr: item.WitnessAddr.B58String(),
			Value:       item.Value,

			AddrSelf: item.SelfAddr.B58String(),
			Payload:  name,
		}
		vinfos.infos = append(vinfos.infos, viVO)
	}

	sort.Stable(&vinfos)

	res, err = model.Tojson(vinfos.infos)
	return
}

func FindTx(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	txItr, ok := rj.Get("txid")
	if !ok {
		res, err = model.Errcode(model.NoField, "txid")
		return
	}
	txidStr := txItr.(string)

	txid, err := hex.DecodeString(txidStr)
	if err != nil {
		res, err = model.Errcode(model.TypeWrong, "txid")
		return
	}

	outMap := make(map[string]interface{})

	txItr, code := mining.FindTxJsonVo(txid)
	outMap["txinfo"] = txItr
	outMap["upchaincode"] = code
	res, err = model.Tojson(outMap)
	return
}

func GetValueForKey(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	keyItr, ok := rj.Get("key")
	if !ok {
		res, err = model.Errcode(model.NoField, "key")
		return
	}
	keyBs, e := hex.DecodeString(keyItr.(string))
	if e != nil {
		res, err = model.Errcode(model.TypeWrong, "key")
		return
	}
	value, e := db.LevelDB.Find(keyBs)
	if e != nil {
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}
	res, err = model.Tojson(value)
	return
}

type BlockHeadVO struct {
	Hash              string
	Height            uint64
	GroupHeight       uint64
	GroupHeightGrowth uint64
	Previousblockhash string
	Nextblockhash     string
	NTx               uint64
	MerkleRoot        string
	Tx                []string
	Time              int64
	Witness           string
	Sign              string
}

func FindBlock(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	heightItr, ok := rj.Get("height")
	if !ok {
		res, err = model.Errcode(model.NoField, "height")
		return
	}
	height := uint64(heightItr.(float64))
	bh := mining.LoadBlockHeadByHeight(height)

	if bh == nil {
		res, err = model.Errcode(model.NotExist)
		return
	}

	txs := make([]string, 0)
	for _, one := range bh.Tx {
		txs = append(txs, hex.EncodeToString(one))
	}

	bhvo := BlockHeadVO{
		Hash:              hex.EncodeToString(bh.Hash),
		Height:            bh.Height,
		GroupHeight:       bh.GroupHeight,
		Previousblockhash: hex.EncodeToString(bh.Previousblockhash),
		Nextblockhash:     hex.EncodeToString(bh.Nextblockhash),
		NTx:               bh.NTx,
		MerkleRoot:        hex.EncodeToString(bh.MerkleRoot),
		Tx:                txs,
		Time:              bh.Time,
		Witness:           bh.Witness.B58String(),
		Sign:              hex.EncodeToString(bh.Sign),
	}

	res, err = model.Tojson(bhvo)
	return

}

var findCommunityStartHeightByAddrOnceLock = new(sync.Mutex)
var findCommunityStartHeightByAddrOnce = make(map[string]bool)

func findCommunityStartHeightByAddr(addr crypto.AddressCoin) {

	ok := false
	findCommunityStartHeightByAddrOnceLock.Lock()
	_, ok = findCommunityStartHeightByAddrOnce[utils.Bytes2string(addr)]
	findCommunityStartHeightByAddrOnce[utils.Bytes2string(addr)] = false
	findCommunityStartHeightByAddrOnceLock.Unlock()
	if ok {

		return
	}

	utils.Go(func() {
		bhHash, err := db.LevelTempDB.Find(config.BuildCommunityAddrStartHeight(addr))
		if err != nil {
			engine.Log.Error("this addr not community:%s error:%s", addr.B58String(), err.Error())
			return
		}
		var sn *sqlite3_db.SnapshotReward

		sn, _, err = mining.FindNotSendReward(&addr)
		if err != nil && err.Error() != xorm.ErrNotExist.Error() {
			engine.Log.Error("querying database Error %s", err.Error())
			return
		}

		if sn != nil {

			return
		}

		snapshots := make([]sqlite3_db.SnapshotReward, 0)

		var bhvo *mining.BlockHeadVO
		var txItr mining.TxItr

		var ok bool

		var have bool
		for {
			if bhHash == nil || len(*bhHash) <= 0 {

				break
			}
			bhvo, err = mining.LoadBlockHeadVOByHash(bhHash)
			if err != nil {
				engine.Log.Error("findCommunityStartHeightByAddr load blockhead error:%s", err.Error())
				return
			}

			bhHash = &bhvo.BH.Nextblockhash
			if len(snapshots) <= 0 {

				snapshotsOne := sqlite3_db.SnapshotReward{
					Addr:        addr,
					StartHeight: bhvo.BH.Height,
					EndHeight:   bhvo.BH.Height,
					Reward:      0,
					LightNum:    0,
				}
				snapshots = append(snapshots, snapshotsOne)
			}
			for _, txItr = range bhvo.Txs {

				if txItr.Class() != config.Wallet_tx_type_voting_reward {
					continue
				}

				_, ok = keystore.FindPuk((*txItr.GetVin())[0].Puk)

				if !ok {

					continue
				}

				txVoteReward := txItr.(*mining.Tx_Vote_Reward)

				have = false
				for _, one := range snapshots {
					if bytes.Equal(addr, one.Addr) && one.StartHeight == txVoteReward.StartHeight && one.EndHeight == txVoteReward.EndHeight {
						have = true
						break
					}
				}
				if have {
					continue
				}

				snapshotsOne := sqlite3_db.SnapshotReward{
					Addr:        addr,
					StartHeight: txVoteReward.StartHeight,
					EndHeight:   txVoteReward.EndHeight,
					Reward:      0,
					LightNum:    0,
				}

				snapshots = append(snapshots, snapshotsOne)
			}
		}

		for i, _ := range snapshots {
			one := snapshots[i]
			err = one.Add(&one)
			if err != nil {
				engine.Log.Info(err.Error())
			}
		}

	})

}

func GetCommunityReward(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	engine.Log.Info("11111111111111111111")
	var addr *crypto.AddressCoin
	addrItr, ok := rj.Get("address")
	if !ok {

		res, err = model.Errcode(model.NoField, "address")
		return
	}
	addrStr := addrItr.(string)
	if addrStr != "" {
		addrMul := crypto.AddressFromB58String(addrStr)
		addr = &addrMul
	}

	if addrStr != "" {
		dst := crypto.AddressFromB58String(addrStr)
		if !crypto.ValidAddr(config.AddrPre, dst) {
			res, err = model.Errcode(model.ContentIncorrectFormat, "address")

			return
		}
	}

	if mining.GetAddrState(*addr) != 2 {

		res, err = model.Errcode(model.RuleField, "address")

		return
	}

	chain := mining.GetLongChain()
	if chain == nil {
		res, err = model.Errcode(model.Nomarl, "The chain end is not synchronized")

		return
	}
	currentHeight := chain.GetCurrentBlock()

	ns, notSend, err := mining.FindNotSendReward(addr)
	if err != nil {

		res, err = model.Errcode(model.Nomarl, err.Error())

		return
	}
	if ns != nil && notSend != nil && len(*notSend) > 0 {

		*notSend, _ = checkTxUpchain(*notSend, currentHeight)

		community := ns.Reward / 10
		light := ns.Reward - community
		rewardTotal := mining.RewardTotal{
			CommunityReward: community,
			LightReward:     light,
			StartHeight:     ns.StartHeight,
			Height:          ns.EndHeight,
			IsGrant:         false,
			AllLight:        ns.LightNum,
			RewardLight:     ns.LightNum - uint64(len(*notSend)),
			IsNew:           false,
		}

		res, err = model.Tojson(rewardTotal)
		return
	}
	if ns == nil {

		findCommunityStartHeightByAddr(*addr)
		res, err = model.Errcode(model.Nomarl, "load reward history")

		return
	}

	startHeight := ns.EndHeight + 1

	rt, _, err := mining.GetRewardCount(addr, startHeight, 0)
	if err != nil {
		if err.Error() == config.ERROR_get_reward_count_sync.Error() {
			res, err = model.Errcode(model.RewardCountSync, err.Error())
			return
		}

		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	rt.IsNew = true
	res, err = model.Tojson(rt)

	return
}

func SendCommunityReward(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	addrItr, ok := rj.Get("address")
	if !ok {
		res, err = model.Errcode(model.NoField, "address")

		return
	}
	addrStr := addrItr.(string)
	if addrStr == "" {
		res, err = model.Errcode(model.ContentIncorrectFormat, "address")

		return
	}
	addr := crypto.AddressFromB58String(addrStr)
	if !crypto.ValidAddr(config.AddrPre, addr) {
		res, err = model.Errcode(model.ContentIncorrectFormat, "address")

		return
	}

	if mining.GetAddrState(addr) != 2 {

		res, err = model.Errcode(model.RuleField, "address")

		return
	}

	puk, ok := keystore.GetPukByAddr(addr)
	if !ok {
		res, err = model.Errcode(model.Nomarl, config.ERROR_public_key_not_exist.Error())

		return
	}
	gasItr, ok := rj.Get("gas")
	if !ok {
		res, err = model.Errcode(model.NoField, "gas")

		return
	}
	gas := uint64(gasItr.(float64))

	pwdItr, ok := rj.Get("pwd")
	if !ok {
		res, err = model.Errcode(model.NoField, "pwd")

		return
	}
	pwd := pwdItr.(string)

	startHeightItr, ok := rj.Get("startheight")
	if !ok {
		res, err = model.Errcode(model.NoField, "startheight")

		return
	}
	startHeight := uint64(startHeightItr.(float64))

	endheightItr, ok := rj.Get("endheight")
	if !ok {
		res, err = model.Errcode(model.NoField, "endheight")

		return
	}
	endheight := uint64(endheightItr.(float64))
	chain := mining.GetLongChain()
	if chain == nil {
		res, err = model.Errcode(model.Nomarl, "The chain end is not synchronized")

		return
	}
	currentHeight := chain.GetCurrentBlock()

	ns, notSend, err := mining.FindNotSendReward(&addr)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, err.Error())

		return
	}

	if ns != nil && notSend != nil && len(*notSend) > 0 {

		*notSend, ok = checkTxUpchain(*notSend, currentHeight)
		if ok {

			res, err = model.Errcode(model.Nomarl, "There are rewards that are not linked")

			return
		}
		cs := mining.NewCommunitySign(puk, ns.StartHeight, ns.EndHeight)

		err = mining.DistributionReward(&addr, notSend, gas, pwd, cs, ns.StartHeight, ns.EndHeight, currentHeight)
		if err != nil {
			if err.Error() == config.ERROR_password_fail.Error() {

				res, err = model.Errcode(model.FailPwd)
				return
			}
			if err.Error() == config.ERROR_not_enough.Error() {
				res, err = model.Errcode(model.NotEnough)
				return
			}
			res, err = model.Errcode(model.Nomarl, err.Error())

			return
		}

		res, err = model.Tojson(model.Success)

		return
	}

	if ns != nil && startHeight <= ns.EndHeight {
		res, err = model.Errcode(model.Nomarl, "Repeat reward")

		return
	}

	rt, notSend, err := mining.GetRewardCount(&addr, startHeight, endheight)
	if err != nil {
		if err.Error() == config.ERROR_get_reward_count_sync.Error() {
			res, err = model.Errcode(model.RewardCountSync, err.Error())
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())

		return
	}
	if !rt.IsGrant {
		now := time.Now().Unix()
		blockNum := (config.Mining_community_reward_time / config.Mining_block_time) - (rt.Height - rt.StartHeight)
		wait := blockNum * config.Mining_block_time
		futuer := time.Unix(now+int64(wait), 0)

		res, err = model.Errcode(model.Nomarl, "Please distribute the reward after "+futuer.Format("2006-01-02 15:04:05"))
		return
	}

	err = mining.CreateRewardCount(addr, rt, *notSend)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, err.Error())

		return
	}

	mining.CleanRewardCountProcessMap(&addr)

	res, err = model.Tojson(model.Success)

	return
}

func checkTxUpchain(notSend []sqlite3_db.RewardLight, currentHeight uint64) ([]sqlite3_db.RewardLight, bool) {

	txidUpchain := make(map[string]int)
	txidNotUpchain := make(map[string]int)
	resultUpchain := make([]sqlite3_db.RewardLight, 0)
	resultUnLockHeight := make([]sqlite3_db.RewardLight, 0)
	resultReward := make([]sqlite3_db.RewardLight, 0)
	haveNotUpchain := false
	for i, _ := range notSend {
		one := notSend[i]
		if one.Txid == nil {
			resultReward = append(resultReward, one)
			continue
		}

		_, ok := txidUpchain[utils.Bytes2string(one.Txid)]
		if ok {

			resultUpchain = append(resultUpchain, one)
			continue
		}
		_, ok = txidNotUpchain[utils.Bytes2string(one.Txid)]
		if ok {

			if one.LockHeight < currentHeight {

				resultUnLockHeight = append(resultUnLockHeight, one)
			} else {

				resultReward = append(resultReward, one)
				haveNotUpchain = true
			}
			continue
		}

		txItr, err := mining.LoadTxBase(one.Txid)
		blockhash, berr := db.GetTxToBlockHash(&one.Txid)

		if err != nil || txItr == nil || berr != nil || blockhash == nil {

			txidNotUpchain[utils.Bytes2string(one.Txid)] = 0

			if one.LockHeight < currentHeight {

				resultUnLockHeight = append(resultUnLockHeight, one)
			} else {

				resultReward = append(resultReward, one)
				haveNotUpchain = true
			}
		} else {

			txidUpchain[utils.Bytes2string(one.Txid)] = 0
			resultUpchain = append(resultUpchain, one)
		}
	}

	if len(resultUnLockHeight) > 0 {
		ids := make([]uint64, 0)
		for _, one := range resultUnLockHeight {
			ids = append(ids, one.Id)
		}
		err := new(sqlite3_db.RewardLight).RemoveTxid(ids)
		if err != nil {
			engine.Log.Error(err.Error())
		}
	}

	if len(resultUpchain) > 0 {
		var err error
		for _, one := range resultUpchain {
			err = one.UpdateDistribution(one.Id, one.Reward)
			if err != nil {
				engine.Log.Error(err.Error())
			}
		}
	}
	return resultReward, haveNotUpchain
}

func TokenPublish(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	var src crypto.AddressCoin
	addrItr, ok := rj.Get("srcaddress")
	if ok {
		srcaddr := addrItr.(string)
		if srcaddr != "" {
			src = crypto.AddressFromB58String(srcaddr)
			_, ok := keystore.FindAddress(src)
			if !ok {
				res, err = model.Errcode(model.ContentIncorrectFormat, "srcaddress")
				return
			}
		}
	}

	var addr *crypto.AddressCoin
	addrItr, ok = rj.Get("address")
	if ok {
		addrStr := addrItr.(string)
		if addrStr != "" {
			addrMul := crypto.AddressFromB58String(addrStr)
			addr = &addrMul
		}

		if addrStr != "" {
			dst := crypto.AddressFromB58String(addrStr)
			if !crypto.ValidAddr(config.AddrPre, dst) {
				res, err = model.Errcode(model.ContentIncorrectFormat, "address")
				return
			}
		}
	}

	amount := uint64(0)

	gasItr, ok := rj.Get("gas")
	if !ok {
		res, err = model.Errcode(5002, "gas")
		return
	}
	gas := uint64(gasItr.(float64))

	pwdItr, ok := rj.Get("pwd")
	if !ok {
		res, err = model.Errcode(5002, "pwd")
		return
	}
	pwd := pwdItr.(string)

	comment := ""
	commentItr, ok := rj.Get("comment")
	if ok && rj.VerifyType("comment", "string") {
		comment = commentItr.(string)
	}

	name := ""
	nameItr, ok := rj.Get("name")
	if ok {
		name = nameItr.(string)
	}

	if strings.Contains(name, ".") || strings.Contains(name, " ") {
		res, err = model.Errcode(5002, "name")
		return
	}

	symbol := ""
	symbolItr, ok := rj.Get("symbol")
	if ok {
		symbol = symbolItr.(string)
	}

	if strings.Contains(symbol, ".") || strings.Contains(symbol, " ") {
		res, err = model.Errcode(5002, "symbol")
		return
	}

	supplyItr, ok := rj.Get("supply")
	if !ok {
		res, err = model.Errcode(5002, "supply")
		return
	}
	supply := uint64(supplyItr.(float64))
	if supply < config.Witness_token_supply_min {
		res, err = model.Errcode(model.Nomarl, config.ERROR_token_min_fail.Error())
		return
	}

	var owner crypto.AddressCoin
	ownerItr, ok := rj.Get("owner")
	if ok {
		ownerStr := ownerItr.(string)
		if ownerStr != "" {
			ownerMul := crypto.AddressFromB58String(ownerStr)
			owner = ownerMul
		}

		if ownerStr != "" {
			dst := crypto.AddressFromB58String(ownerStr)
			if !crypto.ValidAddr(config.AddrPre, dst) {
				res, err = model.Errcode(model.ContentIncorrectFormat, "owner")
				return
			}
		}
	}

	frozenHeight := uint64(0)
	frozenHeightItr, ok := rj.Get("frozen_height")
	if ok {
		frozenHeight = uint64(frozenHeightItr.(float64))
	}

	txItr, e := publish.PublishToken(&src, addr, amount, gas, frozenHeight, pwd, comment, name, symbol, supply, owner)
	if e == nil {

		result, e := utils.ChangeMap(txItr)
		if e != nil {
			res, err = model.Errcode(model.Nomarl, err.Error())
			return
		}
		result["hash"] = hex.EncodeToString(*txItr.GetHash())

		res, err = model.Tojson(result)
		return
	}
	if e.Error() == config.ERROR_password_fail.Error() {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	if e.Error() == config.ERROR_not_enough.Error() {
		res, err = model.Errcode(model.NotEnough)
		return
	}
	if e.Error() == config.ERROR_name_exist.Error() {
		res, err = model.Errcode(model.Exist)
		return
	}
	res, err = model.Errcode(model.Nomarl, e.Error())

	return
}

func TokenPay(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	var src crypto.AddressCoin
	addrItr, ok := rj.Get("srcaddress")
	if ok {
		srcaddr := addrItr.(string)
		if srcaddr != "" {
			src = crypto.AddressFromB58String(srcaddr)
			_, ok := keystore.FindAddress(src)
			if !ok {
				res, err = model.Errcode(model.ContentIncorrectFormat, "srcaddress")
				return
			}
		}
	}

	var addr *crypto.AddressCoin
	addrItr, ok = rj.Get("address")
	if ok {
		addrStr := addrItr.(string)
		if addrStr != "" {
			addrMul := crypto.AddressFromB58String(addrStr)
			addr = &addrMul
		}

		if addrStr != "" {
			dst := crypto.AddressFromB58String(addrStr)
			if !crypto.ValidAddr(config.AddrPre, dst) {
				res, err = model.Errcode(model.ContentIncorrectFormat, "address")
				return
			}
		}
	}

	amountItr, ok := rj.Get("amount")
	if !ok {
		res, err = model.Errcode(model.NoField, "amount")
		return
	}
	amount := uint64(amountItr.(float64))
	if amount <= 0 {
		res, err = model.Errcode(model.AmountIsZero, "amount")
		return
	}

	gasItr, ok := rj.Get("gas")
	if !ok {
		res, err = model.Errcode(5002, "gas")
		return
	}
	gas := uint64(gasItr.(float64))

	pwdItr, ok := rj.Get("pwd")
	if !ok {
		res, err = model.Errcode(5002, "pwd")
		return
	}
	pwd := pwdItr.(string)

	comment := ""
	commentItr, ok := rj.Get("comment")
	if ok && rj.VerifyType("comment", "string") {
		comment = commentItr.(string)
	}

	txidItr, ok := rj.Get("txid")
	if !ok {
		res, err = model.Errcode(5002, "txid")
		return
	}
	txid := txidItr.(string)

	frozenHeight := uint64(0)
	frozenHeightItr, ok := rj.Get("frozen_height")
	if ok {
		frozenHeight = uint64(frozenHeightItr.(float64))
	}

	txItr, e := payment.TokenPay(&src, addr, amount, gas, frozenHeight, pwd, comment, txid)
	if e == nil {
		result, e := utils.ChangeMap(txItr)
		if e != nil {
			res, err = model.Errcode(model.Nomarl, err.Error())
			return
		}
		result["hash"] = hex.EncodeToString(*txItr.GetHash())

		res, err = model.Tojson(result)
		return

	}
	if e.Error() == config.ERROR_password_fail.Error() {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	if e.Error() == config.ERROR_not_enough.Error() {
		res, err = model.Errcode(model.NotEnough)
		return
	}
	if e.Error() == config.ERROR_name_exist.Error() {
		res, err = model.Errcode(model.Exist)
		return
	}
	res, err = model.Errcode(model.Nomarl, e.Error())

	return
}

func TokenPayMore(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	var src crypto.AddressCoin
	srcAddrStr := ""
	addrItr, ok := rj.Get("srcaddress")
	if ok {
		srcAddrStr = addrItr.(string)
		if srcAddrStr != "" {
			src = crypto.AddressFromB58String(srcAddrStr)

			_, ok := keystore.FindAddress(src)
			if !ok {
				res, err = model.Errcode(model.ContentIncorrectFormat, "srcaddress")
				return
			}
		}
	}

	addrItr, ok = rj.Get("addresses")
	if !ok {
		res, err = model.Errcode(model.NoField, "addresses")
		return
	}

	bs, err := json.Marshal(addrItr)
	if err != nil {
		res, err = model.Errcode(model.TypeWrong, "addresses")
		return
	}

	addrs := make([]PayNumber, 0)
	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err = decoder.Decode(&addrs)
	if err != nil {
		res, err = model.Errcode(model.TypeWrong, "addresses")
		return
	}

	if len(addrs) <= 0 {
		res, err = model.Errcode(model.NoField, "addresses")
		return
	}

	amount := uint64(0)

	addr := make([]mining.PayNumber, 0)
	for _, one := range addrs {
		dst := crypto.AddressFromB58String(one.Address)

		if !crypto.ValidAddr(config.AddrPre, dst) {
			res, err = model.Errcode(model.ContentIncorrectFormat, "addresses")
			return
		}
		pnOne := mining.PayNumber{
			Address:      dst,
			Amount:       one.Amount,
			FrozenHeight: one.FrozenHeight,
		}
		addr = append(addr, pnOne)
		amount += one.Amount
	}

	gasItr, ok := rj.Get("gas")
	if !ok {
		res, err = model.Errcode(model.NoField, "gas")
		return
	}
	gas := uint64(gasItr.(float64))

	pwdItr, ok := rj.Get("pwd")
	if !ok {
		res, err = model.Errcode(model.NoField, "pwd")
		return
	}
	pwd := pwdItr.(string)

	comment := ""
	commentItr, ok := rj.Get("comment")
	if ok && rj.VerifyType("comment", "string") {
		comment = commentItr.(string)
	}

	txidItr, ok := rj.Get("txid")
	if !ok {
		res, err = model.Errcode(5002, "txid")
		return
	}
	txid, e := hex.DecodeString(txidItr.(string))
	if e != nil {
		res, err = model.Errcode(model.TypeWrong, "txid")
		return
	}

	txItr, e := payment.TokenPayMore(nil, src, addr, gas, pwd, comment, txid)
	if e == nil {
		result, e := utils.ChangeMap(txItr)
		if e != nil {
			res, err = model.Errcode(model.Nomarl, err.Error())
			return
		}
		result["hash"] = hex.EncodeToString(*txItr.GetHash())

		res, err = model.Tojson(result)
		return

	}
	if e.Error() == config.ERROR_password_fail.Error() {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	if e.Error() == config.ERROR_not_enough.Error() {
		res, err = model.Errcode(model.NotEnough)
		return
	}
	if e.Error() == config.ERROR_name_exist.Error() {
		res, err = model.Errcode(model.Exist)
		return
	}
	res, err = model.Errcode(model.Nomarl, e.Error())

	return
}

func PushTx(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	txidBs, e := hex.DecodeString("0400000000000000fde7187ef2ae15463eed6069f4df4ae97150c8b6f93e3c90708aad4e1e92a9fb")
	if e != nil {
		engine.Log.Info("DecodeString fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}

	payTx, e := mining.LoadTxBase(txidBs)
	if e != nil {
		engine.Log.Info("LoadTxBase fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}

	e = mining.AddTx(payTx)
	if e != nil {
		engine.Log.Info("DecodeString fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}
	res, err = model.Tojson(model.Success)
	return

	txJsonItr, ok := rj.Get("tx")
	if !ok {
		res, err = model.Errcode(model.NoField, "tx")
		return
	}

	txjson := txJsonItr.(string)

	txjsonBs, e := base64.StdEncoding.DecodeString(txjson)
	if e != nil {
		engine.Log.Info("DecodeString fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}

	txItr, e := mining.ParseTxBaseProto(0, &txjsonBs)

	if e != nil {
		engine.Log.Info("ParseTxBaseProto fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}
	engine.Log.Info("rpc transaction received %s", hex.EncodeToString(*txItr.GetHash()))
	if e := txItr.CheckSign(); e != nil {
		engine.Log.Info("transaction check fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}
	e = mining.AddTx(txItr)
	if e != nil {
		engine.Log.Info("AddTx fail:%s", e.Error())
		res, err = model.Errcode(model.Nomarl, e.Error())
		return
	}
	res, err = model.Tojson(model.Success)
	return
}

func FindBlockRange(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	startHeightItr, ok := rj.Get("startHeight")
	if !ok {
		res, err = model.Errcode(model.NoField, "startHeight")
		return
	}
	startHeight := uint64(startHeightItr.(float64))

	endHeightItr, ok := rj.Get("endHeight")
	if !ok {
		res, err = model.Errcode(model.NoField, "endHeight")
		return
	}
	endHeight := uint64(endHeightItr.(float64))

	if endHeight < startHeight {
		res, err = model.Errcode(model.NoField, "endHeight")
		return
	}

	bhvos := make([]mining.BlockHeadVO, 0, endHeight-startHeight+1)

	for i := startHeight; i <= endHeight; i++ {

		bhvo := mining.BlockHeadVO{}
		bh := mining.LoadBlockHeadByHeight(i)

		if bh == nil {
			break
		}
		bhvo.BH = bh
		bhvo.Txs = make([]mining.TxItr, 0, len(bh.Tx))

		for _, one := range bh.Tx {
			txItr, e := mining.LoadTxBase(one)

			if e != nil {
				res, err = model.Errcode(model.Nomarl, e.Error())
				return
			}
			bhvo.Txs = append(bhvo.Txs, txItr)
		}

		bhvos = append(bhvos, bhvo)
	}

	res, err = model.Tojson(bhvos)
	return

}

func FindBlockRangeProto(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	startHeightItr, ok := rj.Get("startHeight")
	if !ok {
		res, err = model.Errcode(model.NoField, "startHeight")
		return
	}
	startHeight := uint64(startHeightItr.(float64))

	endHeightItr, ok := rj.Get("endHeight")
	if !ok {
		res, err = model.Errcode(model.NoField, "endHeight")
		return
	}
	endHeight := uint64(endHeightItr.(float64))

	if endHeight < startHeight {
		res, err = model.Errcode(model.NoField, "endHeight")
		return
	}

	bhvos := make([]*[]byte, 0, endHeight-startHeight+1)

	for i := startHeight; i <= endHeight; i++ {

		bhvo := mining.BlockHeadVO{}
		bh := mining.LoadBlockHeadByHeight(i)

		if bh == nil {
			break
		}
		bhvo.BH = bh
		bhvo.Txs = make([]mining.TxItr, 0, len(bh.Tx))

		for _, one := range bh.Tx {
			txItr, e := mining.LoadTxBase(one)

			if e != nil {
				res, err = model.Errcode(model.Nomarl, e.Error())
				return
			}
			bhvo.Txs = append(bhvo.Txs, txItr)
		}
		bs, e := bhvo.Proto()
		if e != nil {
			res, err = model.Errcode(model.Nomarl, e.Error())
			return
		}
		bhvos = append(bhvos, bs)

	}

	res, err = model.Tojson(bhvos)
	return

}
