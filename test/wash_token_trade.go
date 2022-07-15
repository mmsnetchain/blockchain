package main

import (
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"
)

const (
	splitBlockHeight = 1
)

func main() {

	start()
}

func start() {
	peer1 := Peer{

		Addr:       "127.0.0.1:2080",
		AddressMax: 10000,
		RPCUser:    "test",
		RPCPwd:     "testp",
		WalletPwd:  "123456789",
		PayChan:    make(chan *Address, 10),
	}

	gas := uint64(1)

	totalOne, addrsOne := GetAddressList(peer1)
	log.Println("", totalOne)
	if totalOne < peer1.AddressMax {
		for i := uint64(0); i < uint64(peer1.AddressMax)-totalOne; i++ {
			peer1.CreatNewAddr()
		}
	}

	blockHeight := uint64(0)

	log.Println("")

	totalOne, addrsOne = GetAddressList(peer1)
	info := peer1.GetInfo()
	payNumber := make([]PayNumber, 0)
	for i, one := range addrsOne {
		if uint64(i) >= peer1.AddressMax {
			break
		}
		payNumberOne := PayNumber{
			Address: one.AddrCoin,
			Amount:  (info.Balance - gas) / totalOne,
		}
		payNumber = append(payNumber, payNumberOne)
	}

	peer1.TxPayMore(payNumber, gas)

	for {
		info := peer1.GetInfo()
		if info.CurrentBlock <= blockHeight {
			time.Sleep(time.Second)
			continue
		}
		blockHeight = info.CurrentBlock

		log.Println("info", info)

		log.Println("")

		_, addrsOne := GetAddressList(peer1)
		for i, _ := range addrsOne {
			if i == 0 {
				continue
			}
			if uint64(i) >= peer1.AddressMax {
				break
			}
			if addrsOne[i].Value <= 0 {
				continue
			}

			peer1.StartPay(&addrsOne[i])
		}

	}

}

type Fission struct {
	SrcAddress string
	SrcIndex   int
	DstAddress string
	DstIndex   int
	Fission    int
}

type Peer struct {
	Addr       string
	AddressMax uint64
	RPCUser    string
	RPCPwd     string
	WalletPwd  string
	PayChan    chan *Address
}

func (this *Peer) Fission() {

}

func (this *Peer) StartPay(addr *Address) {
	go func() {
		this.PayChan <- addr
		this.TxPay(addr.AddrCoin, addr.Value-1, 1, addr.AddrCoin)
		<-this.PayChan
	}()
}

func (this *Peer) GetInfo() *Info {

	params := map[string]interface{}{
		"method": "getinfo",
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return nil
	}

	info := new(Info)

	buf := bytes.NewBuffer(bs)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(info)
	return info
}

type Info struct {
	Netid          []byte `json:"netid"`
	TotalAmount    uint64 `json:"TotalAmount"`
	Balance        uint64 `json:"balance"`
	BalanceFrozen  uint64 `json:"BalanceFrozen"`
	Testnet        bool   `json:"testnet"`
	Blocks         uint64 `json:"blocks"`
	Group          uint64 `json:"group"`
	StartingBlock  uint64 `json:"StartingBlock"`
	HighestBlock   uint64 `json:"HighestBlock"`
	CurrentBlock   uint64 `json:"CurrentBlock"`
	PulledStates   uint64 `json:"PulledStates"`
	BlockTime      uint64 `json:"BlockTime"`
	LightNode      uint64 `json:"LightNode"`
	CommunityNode  uint64 `json:"CommunityNode"`
	WitnessNode    uint64 `json:"WitnessNode"`
	NameDepositMin uint64 `json:"NameDepositMin"`
	AddrPre        string `json:"AddrPre"`
}

func GetAddressList(peer Peer) (uint64, []Address) {
	fmt.Println("，")

	params := map[string]interface{}{
		"method": "listaccounts",
	}
	result := Post(peer.Addr, peer.RPCUser, peer.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return 0, nil
	}

	addrs := make([]Address, 0)

	buf := bytes.NewBuffer(bs)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(&addrs)

	return uint64(len(addrs)), addrs
}

type Address struct {
	Index       int
	AddrCoin    string
	Value       uint64
	ValueFrozen uint64
	Type        int
}

func (this *Peer) CreatNewAddr() bool {
	fmt.Println("")

	paramsChild := map[string]interface{}{
		"password": this.WalletPwd,
	}
	params := map[string]interface{}{
		"method": "getnewaddress",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	_, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return false
	}
	return true
}

func (this *Peer) TxPay(srcAddress string, amount, gas uint64, address string) bool {

	paramsChild := map[string]interface{}{
		"srcaddress": srcAddress,
		"amount":     amount,
		"address":    address,
		"gas":        gas,
		"pwd":        this.WalletPwd,
	}
	params := map[string]interface{}{
		"method": "sendtoaddress",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	if result == nil {

		return false
	}
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return false
	}
	if result == nil {
		fmt.Println("", srcAddress, amount, gas, address)
	} else {
		if result.Code != 2000 {
			fmt.Println(result.Code, result.Message, string(bs))
		}
	}
	return true
}

func (this *Peer) PublishToken(owner, name, symbol string, supply, gas uint64) bool {

	paramsChild := map[string]interface{}{
		"owner":  owner,
		"name":   name,
		"symbol": symbol,
		"supply": supply,
		"gas":    gas,
		"pwd":    this.WalletPwd,
	}
	params := map[string]interface{}{
		"method": "tokenpublish",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return false
	}
	fmt.Println(result.Code, result.Message, string(bs))
	return true
}

func (this *Peer) TxPayMore(payNumber []PayNumber, gas uint64) bool {

	fmt.Println("", payNumber, gas)

	paramsChild := map[string]interface{}{
		"addresses": payNumber,
		"gas":       gas,
		"pwd":       this.WalletPwd,
	}
	params := map[string]interface{}{
		"method": "sendtoaddressmore",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return false
	}
	fmt.Println(result.Code, result.Message, string(bs))
	return true
}

func (this *Peer) TxTokenPayMore(tokenId string, payNumber []PayNumber, gas uint64) bool {
	fmt.Println("Token", payNumber, gas)

	paramsChild := map[string]interface{}{
		"txid":      tokenId,
		"addresses": payNumber,
		"gas":       gas,
		"pwd":       this.WalletPwd,
	}
	params := map[string]interface{}{
		"method": "tokenpaymore",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return false
	}
	fmt.Println(result.Code, result.Message, string(bs))
	return true
}

type PayNumber struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

func (this *Peer) GetAccount(address, tokenid string) *Account {

	paramsChild := map[string]interface{}{
		"address":  address,
		"token_id": tokenid,
	}
	params := map[string]interface{}{
		"method": "getaccount",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("", err.Error())
		return nil
	}

	account := new(Account)
	buf := bytes.NewBuffer(bs)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(&account)

	return account
}

type Account struct {
	Balance       uint64 `json:"Balance"`
	BalanceFrozen uint64 `json:"BalanceFrozen"`
}

func Post(addr, rpcUser, rpcPwd string, params map[string]interface{}) *PostResult {
	url := "/rpc"
	method := "POST"

	header := http.Header{"user": []string{rpcUser}, "password": []string{rpcPwd}}
	client := &http.Client{}

	bs, err := json.Marshal(params)
	req, err := http.NewRequest(method, "http:
	if err != nil {
		fmt.Println("request")
		return nil
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("", err.Error())
		return nil
	}

	var resultBs []byte

	if resp.StatusCode == 200 {
		resultBs, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("body")
			return nil
		}

		result := new(PostResult)

		buf := bytes.NewBuffer(resultBs)
		decoder := json.NewDecoder(buf)
		decoder.UseNumber()
		err = decoder.Decode(result)
		return result
	}
	fmt.Println("StatusCode:", resp.StatusCode)
	return nil
}

type PostResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func GetRandNum(n int64) int64 {
	if n == 0 {
		return 0
	}
	result, _ := crand.Int(crand.Reader, big.NewInt(int64(n)))
	return result.Int64()
}

type PayPlan struct {
	SrcIndex int
	DstIndex int
	Value    uint64
}
