package main

import (
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/prestonTao/keystore/crypto"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	payMore()
}

func payMore() {
	bs, err := ioutil.ReadFile("config.json")
	if err != nil {

		panic("：" + err.Error())
		return
	}
	cfi := new(Config)

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err = decoder.Decode(cfi)
	if err != nil {

		panic("：" + err.Error())
		return
	}

	payNumbers := readTextFile("addrs.txt")

	payCount := uint64(0)
	for i, _ := range payNumbers {
		one := payNumbers[i]
		payNumbers[i].Amount = one.Amount / 10

		payCount += one.Amount / 10
	}
	fmt.Println(":", len(payNumbers), ":", payCount)

	input := bufio.NewScanner(os.Stdin)
	fmt.Print(":")
	input.Scan()
	pwd := input.Text()
	fmt.Println("：", pwd)

	pageNumber := 500

	for i := 0; i <= len(payNumbers)/pageNumber; i++ {
		start := i * pageNumber
		end := (i + 1) * pageNumber
		if end > len(payNumbers) {
			end = len(payNumbers)
		}
		params := map[string]interface{}{
			"addresses": payNumbers[start:end],
			"gas":       config.Wallet_tx_gas_min,
			"pwd":       pwd,
		}
		payMores := Info{
			Method: "sendtoaddressmore",
			Params: params,
		}

		log.Println("，start:" + strconv.Itoa(start) + " end:" + strconv.Itoa(end))
		for _, one := range payNumbers[start:end] {
			log.Println("：", one.Address, ":", one.Amount)
		}
		resultBs := HttpPost(*cfi, payMores)
		resultVO := new(RpcResult)
		json.Unmarshal(resultBs, resultVO)

		if resultVO.Code != model.Success {
			log.Println("start:" + strconv.Itoa(start) + " end:" + strconv.Itoa(end))
			break
		} else {
		}
		log.Println(string(resultBs))
		time.Sleep(time.Second * 20)
	}

}

type PayNumber struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
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
	req.Header.Add("password", conf.RpcPassword)
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

func readTextFile(path string) []PayNumber {
	index := 1
	pns := make([]PayNumber, 0)
	file, _ := os.Open(path)
	buf := bufio.NewReader(file)
	for {

		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		fmt.Println(string(line))

		strs := strings.Split(string(line), " ")
		if len(strs) != 2 {
			panic("" + strconv.Itoa(index) + "")
		}
		addr := crypto.AddressFromB58String(strs[0])
		if addr == nil {
			panic("" + strconv.Itoa(index) + "，")
		}
		ok := crypto.ValidAddr("MMS", addr)
		if !ok {
			panic("" + strconv.Itoa(index) + "，")
		}

		amout, err := strconv.Atoi(strs[1])
		if err != nil {
			panic("" + strconv.Itoa(index) + ",")
		}
		pnOne := PayNumber{
			Address: strs[0],
			Amount:  uint64(amout),
		}
		pns = append(pns, pnOne)
		index++
	}
	file.Close()
	return pns
}

type Config struct {
	Netid       uint32 `json:"netid"`
	IP          string `json:"ip"`
	Port        uint16 `json:"port"`
	WebAddr     string `json:"WebAddr"`
	WebPort     uint16 `json:"WebPort"`
	WebStatic   string `json:"WebStatic"`
	WebViews    string `json:"WebViews"`
	RpcServer   bool   `json:"RpcServer"`
	RpcUser     string `json:"RpcUser"`
	RpcPassword string `json:"RpcPassword"`
	Miner       bool   `json:"miner"`
	NetType     string `json:"NetType"`
	AddrPre     string `json:"AddrPre"`
}
