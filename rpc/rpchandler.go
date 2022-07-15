package rpc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/token"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"mmschainnewaccount/rpc/networks"
	"mmschainnewaccount/rpc/sharebox"
	"net/http"
	"path/filepath"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/keystore/kstore"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

type serverHandler func(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) ([]byte, error)

var rpcHandler = map[string]serverHandler{
	"getinfo":               handleGetInfo,
	"getnewaddress":         handleGetNewAddress,
	"listaccounts":          handleListAccounts,
	"getaccount":            handleGetAccount,
	"validateaddress":       handleValidateAddress,
	"import":                Import,
	"export":                Export,
	"sendtoaddress":         sendToAddress,
	"sendtoaddressmore":     sendToAddressmore,
	"depositin":             depositIn,
	"depositout":            depositOut,
	"votein":                voteIn,
	"voteout":               voteOut,
	"updatepwd":             UpdatePwd,
	"createkeystore":        CreateKeystore,
	"namesin":               NameIn,
	"namesout":              NameOut,
	"getnames":              GetNames,
	"findname":              FindName,
	"gettransactionhistory": GetTransactionHistoty,
	"getwitnessinfo":        GetWitnessInfo,
	"getcandidatelist":      GetCandidateList,
	"getcommunitylist":      GetCommunityList,
	"getvotelist":           GetVoteList,
	"findtx":                FindTx,
	"findblock":             FindBlock,
	"getcommunityreward":    GetCommunityReward,
	"sendcommunityreward":   SendCommunityReward,
	"tokenpublish":          TokenPublish,
	"tokenpay":              TokenPay,
	"tokenpaymore":          TokenPayMore,
	"pushtx":                PushTx,
	"getblocksrange":        FindBlockRange,
	"getblocksrangeproto":   FindBlockRangeProto,
	"findvalue":             GetValueForKey,
	"getnodetotal":          GetNodeTotal,
	"getnonce":              GetNonce,

	"getnetworkinfo": networks.NetworkInfo,

	"getsharefolderlist":       sharebox.ShareFolderList,
	"addlocalsharefolder":      sharebox.AddLocalShareFoler,
	"dellocalsharefolder":      sharebox.DelLocalShareFoler,
	"getremotesharefolderlist": sharebox.GetRemoteShareFolderList,

	"getminingspacelist": GetCloudSpaceList,
	"addminingspacesize": AddCloudSpaceSize,
	"delminingspacesize": DelCloudSpaceSize,
	"delminingspaceone":  DelCloudSpaceOne,

	"stopservice": StopService,
}

func StopService(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	utils.StopService()
	res, err = model.Tojson("success")
	return
}

func handleGetInfo(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	value, valuef, valuelockup := mining.FindBalanceValue()

	tbs := token.FindTokenBalanceForAll()
	tbVOs := make([]token.TokenBalanceVO, 0)
	for i, one := range tbs {
		tbs[i].TokenId = one.TokenId
		tbVO := token.TokenBalanceVO{
			TokenId:       hex.EncodeToString([]byte(one.TokenId)),
			Name:          one.Name,
			Symbol:        one.Symbol,
			Supply:        one.Supply,
			Balance:       one.Balance,
			BalanceFrozen: one.BalanceFrozen,
			BalanceLockup: one.BalanceLockup,
		}
		tbVOs = append(tbVOs, tbVO)
	}

	currentBlock := uint64(0)
	startBlock := uint64(0)
	heightBlock := uint64(0)
	pulledStates := uint64(0)
	startBlockTime := uint64(0)

	chain := mining.GetLongChain()
	if chain != nil {
		currentBlock = chain.GetCurrentBlock()
		startBlock = chain.GetStartingBlock()
		heightBlock = mining.GetHighestBlock()
		pulledStates = chain.GetPulledStates()
		startBlockTime = chain.GetStartBlockTime()
	}

	info := model.Getinfo{
		Netid:          []byte(config.AddrPre),
		TotalAmount:    config.Mining_coin_total,
		Balance:        value,
		BalanceFrozen:  valuef,
		BalanceLockup:  valuelockup,
		Testnet:        true,
		Blocks:         currentBlock,
		Group:          0,
		StartingBlock:  startBlock,
		StartBlockTime: startBlockTime,
		HighestBlock:   heightBlock,
		CurrentBlock:   currentBlock,
		PulledStates:   pulledStates,
		BlockTime:      config.Mining_block_time,
		LightNode:      config.Mining_light_min,
		CommunityNode:  config.Mining_vote,
		WitnessNode:    config.Mining_deposit,
		NameDepositMin: config.Mining_name_deposit_min,
		AddrPre:        config.AddrPre,
		TokenBalance:   tbVOs,
	}
	res, err := model.Tojson(info)
	return res, err
}

func handleGetNewAddress(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	if !rj.VerifyType("password", "string") {
		res, err = model.Errcode(model.TypeWrong, "password")
		return
	}
	password, ok := rj.Get("password")
	if !ok {
		res, err = model.Errcode(model.NoField, "password")
		return
	}

	addr, err := keystore.GetNewAddr(password.(string))
	if err != nil {
		if err.Error() == config.ERROR_password_fail.Error() {

			res, err = model.Errcode(model.FailPwd)
			return
		}
		res, _ = model.Errcode(model.Nomarl)
		return
	}
	getnewadress := model.GetNewAddress{Address: addr.B58String()}
	res, err = model.Tojson(getnewadress)
	return
}

type AccountVO struct {
	Index       int
	AddrCoin    string
	Value       uint64
	ValueFrozen uint64
	ValueLockup uint64
	Type        int
}

func handleListAccounts(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	tokenidStr := ""
	tokenidItr, ok := rj.Get("token_id")
	if ok {
		if !rj.VerifyType("token_id", "string") {
			res, err = model.Errcode(model.TypeWrong, "token_id")
			return
		}
		tokenidStr = tokenidItr.(string)
	}

	vos := make([]AccountVO, 0)
	for i, val := range keystore.GetAddr() {
		var ba, fba, baLockup uint64
		if tokenidStr == "" {
			ba, fba, baLockup = mining.GetBalanceForAddrSelf(val.Addr)
		} else {

		}

		vo := AccountVO{
			Index:       i,
			AddrCoin:    val.AddrStr,
			Type:        mining.GetAddrState(val.Addr),
			Value:       ba,
			ValueFrozen: fba,
			ValueLockup: baLockup,
		}
		vos = append(vos, vo)
	}
	res, err = model.Tojson(vos)
	return res, err

}

func handleGetAccount(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	addr, ok := rj.Get("address")
	if !ok {
		res, err = model.Errcode(model.NoField, "address")
		return
	}
	if !rj.VerifyType("address", "string") {
		res, err = model.Errcode(model.TypeWrong, "address")
		return
	}

	addrCoin := crypto.AddressFromB58String(addr.(string))
	ok = crypto.ValidAddr(config.AddrPre, addrCoin)
	if !ok {
		res, err = model.Errcode(model.ContentIncorrectFormat, "address")
		return
	}

	value, valueFrozen, _ := mining.GetBalanceForAddrSelf(addrCoin)

	getaccount := model.GetAccount{
		Balance:       value,
		BalanceFrozen: valueFrozen,
	}
	res, err = model.Tojson(getaccount)
	return
}

func handleValidateAddress(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	if !rj.VerifyType("address", "string") {
		res, err = model.Errcode(model.TypeWrong, "address")
		return
	}
	addr, ok := rj.Get("address")
	if !ok {
		res, err = model.Errcode(model.NoField, "address")
		return
	}
	addrCoin := crypto.AddressFromB58String(addr.(string))
	ok = crypto.ValidAddr(config.AddrPre, addrCoin)
	res, err = model.Tojson(ok)
	return
}

func sendToAddress(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	var src crypto.AddressCoin
	addrItr, ok := rj.Get("srcaddress")
	if ok {
		srcaddr := addrItr.(string)
		if srcaddr != "" {
			src = crypto.AddressFromB58String(srcaddr)

			if !crypto.ValidAddr(config.AddrPre, src) {
				res, err = model.Errcode(model.ContentIncorrectFormat, "srcaddress")
				return
			}
			_, ok := keystore.FindAddress(src)
			if !ok {
				res, err = model.Errcode(model.ContentIncorrectFormat, "srcaddress")
				return
			}
		}
	}

	addrItr, ok = rj.Get("address")
	if !ok {
		res, err = model.Errcode(model.NoField, "address")
		return
	}
	addr := addrItr.(string)

	dst := crypto.AddressFromB58String(addr)
	if !crypto.ValidAddr(config.AddrPre, dst) {
		res, err = model.Errcode(model.ContentIncorrectFormat, "address")
		return
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
		res, err = model.Errcode(model.NoField, "gas")
		return
	}
	gas := uint64(gasItr.(float64))

	frozenHeight := uint64(0)
	frozenHeightItr, ok := rj.Get("frozen_height")
	if ok {
		frozenHeight = uint64(frozenHeightItr.(float64))
	}

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
	runeLength := len([]rune(comment))
	if runeLength > 1024 {
		res, err = model.Errcode(model.CommentOverLengthMax, "comment")
		return
	}
	temp := new(big.Int).Mul(big.NewInt(int64(runeLength)), big.NewInt(100000000))
	temp = new(big.Int).Div(temp, big.NewInt(1024))
	if gas < temp.Uint64() {
		res, err = model.Errcode(model.GasTooLittle, "gas")
		return
	}

	total, _ := mining.GetLongChain().GetBalance().BuildPayVinNew(&src, amount+gas)
	if total < amount+gas {

		res, err = model.Errcode(model.BalanceNotEnough)
		return
	}

	txpay, err := mining.SendToAddress(&src, &dst, amount, gas, frozenHeight, pwd, comment)

	if err != nil {

		if err.Error() == config.ERROR_password_fail.Error() {

			res, err = model.Errcode(model.FailPwd)
			return
		}

		if err.Error() == config.ERROR_amount_zero.Error() {
			res, err = model.Errcode(model.AmountIsZero, "amount")
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}

	result, err := utils.ChangeMap(txpay)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	result["hash"] = hex.EncodeToString(*txpay.GetHash())

	res, err = model.Tojson(result)

	return
}

type PayNumber struct {
	Address      string `json:"address"`
	Amount       uint64 `json:"amount"`
	FrozenHeight uint64 `json:"frozen_height"`
}

func sendToAddressmore(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

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

	value, _, _ := mining.FindBalanceValue()
	if amount+gas > value {
		res, err = model.Errcode(model.BalanceNotEnough)

		return
	}

	txpay, err := mining.SendToMoreAddress(&src, addr, gas, pwd, comment)
	if err != nil {

		if err.Error() == config.ERROR_password_fail.Error() {

			res, err = model.Errcode(model.FailPwd)
			return
		}

		if err.Error() == config.ERROR_amount_zero.Error() {
			res, err = model.Errcode(model.AmountIsZero, "amount")
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	result, err := utils.ChangeMap(txpay)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	result["hash"] = hex.EncodeToString(*txpay.GetHash())

	res, err = model.Tojson(result)
	return
}

func depositIn(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	amountItr, ok := rj.Get("amount")
	if !ok {
		res, err = model.Errcode(5002, "amount")
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

	payload := ""
	payloadItr, ok := rj.Get("payload")
	if ok {
		payload = payloadItr.(string)
	}

	heightBlock := mining.GetHighestBlock()
	if heightBlock <= config.Wallet_vote_start_height {
		res, err = model.Errcode(model.VoteNotOpen)
		return
	}

	value, _, _ := mining.FindBalanceValue()
	if amount > value {
		res, err = model.Errcode(model.BalanceNotEnough)
		return
	}

	err = mining.DepositIn(amount, gas, pwd, payload)
	if err != nil {
		if err.Error() == config.ERROR_password_fail.Error() {
			res, err = model.Errcode(model.FailPwd)
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	res, err = model.Tojson("success")
	config.SubmitDepositin = true
	return
}

func depositOut(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	addr := ""
	addrItr, ok := rj.Get("address")
	if ok {
		addr = addrItr.(string)

	}

	if addr != "" {
		dst := crypto.AddressFromB58String(addr)
		if !crypto.ValidAddr(config.AddrPre, dst) {
			res, err = model.Errcode(model.ContentIncorrectFormat, "address")
			return
		}
	}

	amount := uint64(0)
	amountItr, ok := rj.Get("amount")
	if ok {
		amount = uint64(amountItr.(float64))
		if amount < 0 {
			res, err = model.Errcode(model.AmountIsZero, "amount")
			return
		}
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

	engine.Log.Info("address:%s amount:%d gas:%d", addr, amount, gas)

	err = mining.DepositOut(addr, amount, gas, pwd)
	if err != nil {
		if err.Error() == config.ERROR_password_fail.Error() {
			res, err = model.Errcode(model.FailPwd)
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	res, err = model.Tojson("success")
	return
}

func voteIn(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	vtItr, ok := rj.Get("votetype")
	if !ok {
		res, err = model.Errcode(model.NoField, "votetype")
		return
	}
	voteType := uint16(vtItr.(float64))

	addr := ""
	addrItr, ok := rj.Get("address")
	if !ok {
		res, err = model.Errcode(model.NoField, "votetype")
		return
	}
	addr = addrItr.(string)
	if addr == "" {
		res, err = model.Errcode(model.NoField, "votetype")
		return
	}

	var witnessAddr crypto.AddressCoin
	witnessAddrItr, ok := rj.Get("witness")
	if ok {

		witnessStr := witnessAddrItr.(string)

		witnessAddr = crypto.AddressFromB58String(witnessStr)

		if witnessStr != "" {
			dst := crypto.AddressFromB58String(witnessStr)
			if !crypto.ValidAddr(config.AddrPre, dst) {
				res, err = model.Errcode(model.ContentIncorrectFormat, "witness")
				return
			}
		}
	}

	switch voteType {
	case mining.VOTE_TYPE_community:
	case mining.VOTE_TYPE_vote:
	case mining.VOTE_TYPE_light:
		witnessAddr = nil
	default:
		res, err = model.Errcode(model.Nomarl, "votetype")
		return
	}
	var dst crypto.AddressCoin
	dst = crypto.AddressFromB58String(addr)
	if !crypto.ValidAddr(config.AddrPre, dst) {
		res, err = model.Errcode(model.ContentIncorrectFormat, "address")
		return
	}

	amountItr, ok := rj.Get("amount")
	if !ok {
		res, err = model.Errcode(5002, "amount")
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

	payload := ""
	payloadItr, ok := rj.Get("payload")
	if ok {
		payload = payloadItr.(string)
	}

	heightBlock := mining.GetHighestBlock()
	if heightBlock <= config.Wallet_vote_start_height {
		res, err = model.Errcode(model.VoteNotOpen)
		return
	}

	value, _, _ := mining.FindBalanceValue()
	if amount > value {
		res, err = model.Errcode(model.BalanceNotEnough)
		return
	}

	err = mining.VoteIn(voteType, witnessAddr, dst, amount, gas, pwd, payload)
	if err != nil {
		if err.Error() == config.ERROR_password_fail.Error() {
			res, err = model.Errcode(model.FailPwd)
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	res, err = model.Tojson("success")
	return
}

func voteOut(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	vtItr, ok := rj.Get("votetype")
	if !ok {
		res, err = model.Errcode(model.NoField, "votetype")
		return
	}
	voteType := uint16(vtItr.(float64))

	addrStr := ""
	var addr crypto.AddressCoin
	addrItr, ok := rj.Get("address")
	if ok {
		addrStr = addrItr.(string)
	}
	if addrStr != "" {
		addr = crypto.AddressFromB58String(addrStr)
		if !crypto.ValidAddr(config.AddrPre, addr) {
			res, err = model.Errcode(model.ContentIncorrectFormat, "address")
			return
		}
	}

	switch voteType {
	case mining.VOTE_TYPE_community:
	case mining.VOTE_TYPE_vote:
	case mining.VOTE_TYPE_light:

		di := mining.GetLongChain().GetBalance().GetDepositVote(&addr)
		if di != nil {
			res, err = model.Errcode(model.NotDepositOutLight)
			return
		}
	default:
		res, err = model.Errcode(model.Nomarl, "votetype")
		return
	}

	amountItr, ok := rj.Get("amount")
	var amount uint64
	if !ok && voteType == mining.VOTE_TYPE_vote {
		res, err = model.Errcode(model.NoField, "amount")
		return
	}
	if ok {
		amount = uint64(amountItr.(float64))
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

	payload := ""
	payloadItr, ok := rj.Get("payload")
	if ok {
		payload = payloadItr.(string)
	}

	err = mining.VoteOut(voteType, addr, amount, gas, pwd, payload)
	if err != nil {

		if err.Error() == config.ERROR_password_fail.Error() {
			res, err = model.Errcode(model.FailPwd)
			return
		}

		if err.Error() == config.ERROR_not_enough.Error() {
			res, err = model.Errcode(model.BalanceNotEnough)
			return
		}

		if err.Error() == config.ERROR_vote_exist.Error() {
			res, err = model.Errcode(model.VoteExist)
			return
		}
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	res, err = model.Tojson("success")
	return
}

func UpdatePwd(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	oldpwdItr, ok := rj.Get("oldpwd")
	if !ok {
		res, err = model.Errcode(5002, "oldpwd")
		return
	}
	oldpwd := oldpwdItr.(string)

	pwdItr, ok := rj.Get("newpwd")
	if !ok {
		res, err = model.Errcode(5002, "newpwd")
		return
	}
	pwd := pwdItr.(string)

	ok, err = keystore.UpdatePwd(oldpwd, pwd)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}
	if !ok {

		res, err = model.Errcode(model.FailPwd, errors.New("password fail").Error())
		return
	}
	config.Wallet_keystore_default_pwd = pwd
	res, err = model.Tojson("success")
	return
}

func CreateKeystore(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	randomItr, ok := rj.Get("random")
	if !ok {
		res, err = model.Errcode(model.NoField, "random")
		return
	}

	randomItrs := randomItr.([]interface{})
	buf := bytes.NewBuffer(nil)
	for _, one := range randomItrs {
		onePoint := uint16(one.(float64))
		_, e := buf.Write(utils.Uint16ToBytes(onePoint))
		if e != nil {
			res, err = model.Errcode(model.Nomarl, e.Error())
			return
		}
	}
	if buf.Len() != 4000 {

		res, err = model.Errcode(model.Nomarl, "Random number length not equal to 2000")
		return
	}

	rand1 := buf.Bytes()[:2000]
	rand2 := buf.Bytes()[2000:]

	pwdItr, ok := rj.Get("pwd")
	if !ok {
		res, err = model.Errcode(model.NoField, "pwd")
		return
	}
	pwd := pwdItr.(string)

	err = keystore.CreateKeystoreRand(filepath.Join(config.Path_configDir, config.Core_keystore), rand1, rand2, pwd)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, err.Error())
		return
	}

	res, err = model.Tojson("success")
	return
}

func Export(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	pwdItr, ok := rj.Get("password")
	if !ok {
		res, err = model.Errcode(5002, "password")
		return
	}
	pwd := pwdItr.(string)
	rs := kstore.Export(pwd)
	di := kstore.ParseDataInfo(rs)
	if di.Code == 500 {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	res, err = model.Tojson(di.Data)
	return
}

func Import(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	pwdItr, ok := rj.Get("password")
	if !ok {
		res, err = model.Errcode(5002, "password")
		return
	}
	pwd := pwdItr.(string)
	seeds, ok := rj.Get("seed")
	if !ok {
		res, err = model.Errcode(5002, "seed")
		return
	}
	seed := seeds.(string)
	rs := kstore.Import(config.Path_configDir, pwd, seed, config.AddrPre)
	di := kstore.ParseDataInfo(rs)
	if di.Code == 500 {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	res, err = model.Tojson(di.Data)
	return
}
