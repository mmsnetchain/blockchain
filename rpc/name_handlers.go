package rpc

import (
	"bytes"
	"encoding/hex"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_out"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NameIn(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

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

func NameOut(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

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

func GetNames(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	nameinfoVOs := make([]NameinfoVO, 0)

	names := name.GetNameList()
	for _, one := range names {
		nets := make([]string, 0)
		for _, two := range one.NetIds {
			nets = append(nets, two.B58String())
		}
		addrs := make([]string, 0)
		for _, two := range one.AddrCoins {
			addrs = append(addrs, two.B58String())
		}
		voOne := NameinfoVO{
			Name:           one.Name,
			NetIds:         nets,
			AddrCoins:      addrs,
			Height:         one.Height,
			NameOfValidity: one.NameOfValidity,
			Deposit:        one.Deposit,
		}
		nameinfoVOs = append(nameinfoVOs, voOne)
	}

	res, err = model.Tojson(nameinfoVOs)
	return
}

type NameinfoVO struct {
	Name           string
	NetIds         []string
	AddrCoins      []string
	Height         uint64
	NameOfValidity uint64
	Deposit        uint64
}

func FindName(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {

	nameItr, ok := rj.Get("name")
	if !ok {
		res, err = model.Errcode(5002, "name")
		return
	}
	nameStr := nameItr.(string)
	nameinfo := name.FindNameToNet(nameStr)
	if nameinfo == nil || nameinfo.CheckIsOvertime(mining.GetHighestBlock()) {
		res, err = model.Errcode(model.NotExist, nameStr)
		return
	}
	bs, err := json.Marshal(nameinfo)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, "find name formate failt 1")
		return
	}
	result := make(map[string]interface{})

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err = decoder.Decode(&result)
	if err != nil {
		res, err = model.Errcode(model.Nomarl, "find name formate failt 2")
		return
	}

	res, err = model.Tojson(result)
	return
}
