package rpc

import (
	"encoding/hex"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/libp2parea/virtual_node"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/cloud_reward/client"
	"mmschainnewaccount/cloud_space/fs"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"net/http"
	"strings"
)

func SpacesMiningIn(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	var addr *crypto.AddressCoin
	addrItr, ok := rj.Get("address")
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
		res, err = model.Errcode(5002, "amount")
		return
	}
	amount := uint64(amountItr.(float64))
	if amount < config.Mining_name_deposit_min {
		res, err = model.Errcode(model.Nomarl, config.ERROR_name_deposit.Error())
		return
	}

	gasItr, ok := rj.Get("gas")
	if !ok {
		res, err = model.Errcode(5002, "gas")
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
		res, err = model.Errcode(5002, "pwd")
		return
	}
	pwd := pwdItr.(string)

	nameItr, ok := rj.Get("name")
	if !ok {
		res, err = model.Errcode(5002, "name")
		return
	}
	name := nameItr.(string)

	if name == "" {
		res, err = model.Errcode(5002, "name")
		return
	}
	if strings.Contains(name, ".") || strings.Contains(name, " ") {
		res, err = model.Errcode(5002, "name")
		return
	}

	ids := make([]nodeStore.AddressNet, 0)
	netIdsItr, ok := rj.Get("netids")
	if ok {
		netIds := netIdsItr.([]interface{})
		for _, one := range netIds {
			netidOne := one.(string)
			idOne := nodeStore.AddressFromB58String(netidOne)
			ids = append(ids, idOne)
		}
	}

	coins := make([]crypto.AddressCoin, 0)
	addrcoinsItr, ok := rj.Get("addrcoins")
	if ok {
		addrcoins := addrcoinsItr.([]interface{})
		for _, one := range addrcoins {
			addrcoinOne := one.(string)
			idOne := crypto.AddressFromB58String(addrcoinOne)
			coins = append(coins, idOne)
		}
	}

	comment := ""
	commentItr, ok := rj.Get("comment")
	if ok && rj.VerifyType("comment", "string") {
		comment = commentItr.(string)
	}

	txpay, err := tx_name_in.NameIn(nil, addr, amount, gas, frozenHeight, pwd, comment, name, ids, coins)
	if err == nil {

		result, e := utils.ChangeMap(txpay)
		if e != nil {
			res, err = model.Errcode(model.Nomarl, err.Error())
			return
		}
		result["hash"] = hex.EncodeToString(*txpay.GetHash())

		res, err = model.Tojson(result)

		return
	}
	if err.Error() == config.ERROR_password_fail.Error() {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	if err.Error() == config.ERROR_not_enough.Error() {
		res, err = model.Errcode(model.NotEnough)
		return
	}
	if err.Error() == config.ERROR_name_exist.Error() {
		res, err = model.Errcode(model.Exist)
		return
	}
	res, err = model.Errcode(model.Nomarl, err.Error())

	return
}

func SpacesMiningOut(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	var addr *crypto.AddressCoin
	addrItr, ok := rj.Get("address")
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

	gasItr, ok := rj.Get("gas")
	if !ok {
		res, err = model.Errcode(5002, "gas")
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
		res, err = model.Errcode(5002, "pwd")
		return
	}
	pwd := pwdItr.(string)

	nameItr, ok := rj.Get("name")
	if !ok {
		res, err = model.Errcode(5002, "name")
		return
	}
	name := nameItr.(string)

	if name == "" {
		res, err = model.Errcode(5002, "name")
		return
	}
	if strings.Contains(name, ".") || strings.Contains(name, " ") {
		res, err = model.Errcode(5002, "name")
		return
	}

	comment := ""
	commentItr, ok := rj.Get("comment")
	if ok && rj.VerifyType("comment", "string") {
		comment = commentItr.(string)
	}

	txpay, err := tx_name_out.NameOut(nil, addr, 0, gas, frozenHeight, pwd, comment, name)
	if err == nil {

		result, e := utils.ChangeMap(txpay)
		if e != nil {
			res, err = model.Errcode(model.Nomarl, err.Error())
			return
		}
		result["hash"] = hex.EncodeToString(*txpay.GetHash())

		res, err = model.Tojson(result)

		return
	}
	if err.Error() == config.ERROR_password_fail.Error() {
		res, err = model.Errcode(model.FailPwd)
		return
	}
	if err.Error() == config.ERROR_not_enough.Error() {
		res, err = model.Errcode(model.NotEnough)
		return
	}
	if err.Error() == config.ERROR_name_not_exist.Error() {
		res, err = model.Errcode(model.NotExist)
		return
	}
	res, err = model.Errcode(model.Nomarl, err.Error())
	return

}

func SetCloudSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	nItr, ok := rj.Get("n")
	if !ok {
		res, err = model.Errcode(5002, "n")
		return
	}
	n := uint64(nItr.(float64))
	virtual_node.SetupVnodeNumber(n)
	res, err = model.Tojson("success")
	return
}

type VnodeinfoVO struct {
	Nid   string `json:"nid"`
	Index uint64 `json:"index"`
	Vid   string `json:"vid"`
}

func GetCloudSpaceList(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	list := fs.GetSpaceList()
	total := fs.GetSpaceSize()
	useSize := fs.GetUseSpaceSize()
	result := make(map[string]interface{}, 0)
	result["Code"] = 0
	result["SpaceList"] = list
	result["TotalSize"] = total
	result["UseSize"] = useSize
	res, err = model.Tojson(result)
	return
}

func AddCloudSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	nItr, ok := rj.Get("n")
	if !ok {
		res, err = model.Errcode(model.NoField, "n")
		return
	}
	n := uint64(nItr.(float64))

	absPath := ""

	absPathItr, ok := rj.Get("absPath")
	if ok {
		absPath = absPathItr.(string)
	}

	utils.Go(func() {
		fs.AddSpace(absPath, n)
		client.SendOnlineHeartBeat()
	})

	res, err = model.Tojson("success")
	return
}

func DelCloudSpaceSize(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	nItr, ok := rj.Get("n")
	if !ok {
		res, err = model.Errcode(model.NoField, "n")
		return
	}
	n := uint64(nItr.(float64))

	fs.DelSpace(n)

	res, err = model.Tojson("success")
	return
}

func DelCloudSpaceOne(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	dbpathItr, ok := rj.Get("dbpath")
	if !ok {
		res, err = model.Errcode(model.NoField, "dbpath")
		return
	}

	fs.DelSpaceForDbPath(dbpathItr.(string))

	res, err = model.Tojson("success")
	return
}
