package main

import (
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/token/payment"
	"mmschainnewaccount/chain_witness_vote/mining/token/publish"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_in"
	"mmschainnewaccount/chain_witness_vote/mining/tx_name_out"

	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc"
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/prestonTao/libp2parea/engine"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	splitBlockHeight = 1
)

func init() {
	tpc := new(payment.TokenPublishController)
	tpc.ActiveVoutIndex = new(sync.Map)
	mining.RegisterTransaction(config.Wallet_tx_type_token_payment, tpc)
	mining.RegisterTransaction(config.Wallet_tx_type_token_publish, new(publish.TokenPublishController))

	mining.RegisterTransaction(config.Wallet_tx_type_account, new(tx_name_in.AccountController))
	mining.RegisterTransaction(config.Wallet_tx_type_account_cancel, new(tx_name_out.AccountController))

}

func main() {

	start()
}

func start() {
	peer1 := Peer{

		Addr:       "192.168.28.36:5081",
		AddressMax: 50000,
		RPCUser:    "test",
		RPCPwd:     "testp",
		WalletPwd:  "witness2",
		PayChan:    make(chan *Address, 10),
	}

	info := peer1.GetInfo()
	fmt.Printf("%+v", info)
	startHeight := info.StartingBlock
	endHeight := info.CurrentBlock

	for i := startHeight; i <= endHeight; i++ {
		bhvos, _ := peer1.GetBlockAndTx(i, i)

		for _, one := range bhvos {
			PrintBlock(one)
		}
	}

	return

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

type PayNumber struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

func (this *Peer) GetAccount(address string) *Account {

	paramsChild := map[string]interface{}{
		"address": address,
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

func (this *Peer) GetBlockAndTx(start, end uint64) ([]mining.BlockHeadVO, error) {

	fmt.Println("", start, end)

	paramsChild := map[string]interface{}{
		"startHeight": start,
		"endHeight":   end,
	}
	params := map[string]interface{}{
		"method": "getblocksrange",
		"params": paramsChild,
	}
	result := Post(this.Addr, this.RPCUser, this.RPCPwd, params)
	bs, err := json.Marshal(result.Result)
	if err != nil {
		fmt.Println("0", err.Error())
		return nil, err
	}

	bhvos := make([]mining.BlockHeadVO, 0)

	resultMap := make([]map[string]interface{}, 0)
	err = json.Unmarshal(bs, &resultMap)
	if err != nil {
		fmt.Println("1111", err.Error())
		return nil, err
	}
	for _, one := range resultMap {
		bhItr := one["bh"]
		bs, err := json.Marshal(bhItr)
		if err != nil {
			fmt.Println("1", err.Error())
			return nil, err
		}
		bh := new(mining.BlockHead)
		err = json.Unmarshal(bs, bh)
		if err != nil {
			fmt.Println("2", err.Error())
			return nil, err
		}

		bhvo := new(mining.BlockHeadVO)
		bhvo.BH = bh
		bhvo.Txs = make([]mining.TxItr, 0)

		txsItr := one["txs"]
		txsMap := make([]map[string]interface{}, 0)
		bs, err = json.Marshal(txsItr)
		if err != nil {
			fmt.Println("3", err.Error())
			return nil, err
		}
		err = json.Unmarshal(bs, &txsMap)
		if err != nil {
			fmt.Println("4", err.Error())
			return nil, err
		}
		for _, two := range txsMap {
			bs, _ = json.Marshal(two)
			txOne, _ := mining.ParseTxBase(0, &bs)
			bhvo.Txs = append(bhvo.Txs, txOne)
		}
		bhvos = append(bhvos, *bhvo)
	}

	return bhvos, nil
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

func PrintBlock(bhVO mining.BlockHeadVO) {
	bh := bhVO.BH
	engine.Log.Info("===============================================")
	engine.Log.Info("%d %s  %d", bh.Height, hex.EncodeToString(bh.Hash), len(bh.Tx))
	txs := make([]string, 0)
	for _, one := range bh.Tx {
		txs = append(txs, hex.EncodeToString(one))
	}
	bhvo := rpc.BlockHeadVO{
		Hash:              hex.EncodeToString(bh.Hash),
		Height:            bh.Height,
		GroupHeight:       bh.GroupHeight,
		GroupHeightGrowth: bh.GroupHeightGrowth,
		Previousblockhash: hex.EncodeToString(bh.Previousblockhash),
		Nextblockhash:     hex.EncodeToString(bh.Nextblockhash),
		NTx:               bh.NTx,
		MerkleRoot:        hex.EncodeToString(bh.MerkleRoot),
		Tx:                txs,
		Time:              bh.Time,
		Witness:           bh.Witness.B58String(),
		Sign:              hex.EncodeToString(bh.Sign),
	}
	bs, _ := json.Marshal(bhvo)
	engine.Log.Info(string(bs))

	for _, one := range bhVO.Txs {
		engine.Log.Info("-----------------------------------------------")

		txItr := one.GetVOJSON()
		bs, _ := json.Marshal(txItr)
		engine.Log.Info(string(bs))
	}

}
