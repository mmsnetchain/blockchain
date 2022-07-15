package main

import (
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Config struct {
	RpcHost string
	RpcUser string
	RpcPwd  string
}

func main() {
	conf := Config{
		RpcHost: "127.0.0.1:5081",
		RpcUser: "test",
		RpcPwd:  "testp",
	}
	getinfoParams := Info{
		Method: "getinfo",
	}
	localHeight := uint64(0)
	remoteHeight := uint64(0)
	for {

		result := HttpPost(conf, getinfoParams)
		fmt.Println(string(result))

		resultVO := new(RpcResult)
		json.Unmarshal(result, resultVO)
		bs, _ := json.Marshal(resultVO.Result)

		getinfo := new(model.Getinfo)
		json.Unmarshal(bs, getinfo)

		fmt.Println("", getinfo.CurrentBlock)
		remoteHeight = getinfo.CurrentBlock

		for localHeight < remoteHeight {
			params := map[string]interface{}{
				"startHeight": localHeight + 1,
				"endHeight":   localHeight + 1,
			}
			getBlockParams := Info{
				Method: "getblocksrange",
				Params: params,
			}
			resultBs := HttpPost(conf, getBlockParams)

			resultVO := new(RpcResult)
			json.Unmarshal(resultBs, resultVO)
			bs, _ := json.Marshal(resultVO.Result)

			array := make([]interface{}, 0)
			json.Unmarshal(bs, &array)
			for _, one := range array {
				bsOne, _ := json.Marshal(one)

				bhvo, _ := mining.ParseBlockHeadVO(&bsOne)

				bhvojsonbs, _ := bhvo.Json()
				fmt.Println(string(*bhvojsonbs))

				count(bhvo)

				localHeight = bhvo.BH.Height

			}

		}

		time.Sleep(time.Second)
	}

	fmt.Println("hello")
}

type Info struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

func HttpPost(conf Config, info Info) []byte {
	jsons, _ := json.Marshal(info)
	result := string(jsons)
	jsonInfo := strings.NewReader(result)
	req, _ := http.NewRequest("POST", "http:
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("user", conf.RpcUser)
	req.Header.Add("password", conf.RpcPwd)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error create client:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error getInfo:%v", err)
	}
	return body
}

type RpcResult struct {
	JsonRpc string      `json:"jsonrpc"`
	Code    int         `json:"code"`
	Result  interface{} `json:"result"`
}

func count(bhvo *mining.BlockHeadVO) {
	for _, txOne := range bhvo.Txs {
		switch txOne.Class() {
		case config.Wallet_tx_type_pay:
			vouts := txOne.GetVout()

			for voutIndex, voutOne := range *vouts {
				ok := voutOne.CheckIsSelf()

				if !ok {
					continue
				}

				txItem := mining.TxItem{
					Addr:         &voutOne.Address,
					Value:        voutOne.Value,
					Txid:         bhvo.BH.Hash,
					OutIndex:     uint64(voutIndex),
					Height:       bhvo.BH.Height,
					LockupHeight: voutOne.FrozenHeight,
				}
				fmt.Println("txitem", txItem)
			}
		case config.Wallet_tx_type_account:
			tx := txOne.(*tx_name_in.Tx_account)
			fmt.Println("", string(tx.Account))
		}
	}
}
