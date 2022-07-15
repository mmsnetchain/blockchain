package controllers

import (
	"github.com/astaxie/beego"
)

type AccountController struct {
	beego.Controller
}

type Getinfo struct {
	Balance float64 `json:"balance"`
	Testnet bool    `json:"testnet"`
	Blocks  uint64  `json:"blocks"`
}

type Result struct {
	Jsonrpc string  `json:"jsonrpc"`
	Code    int     `json:"code"`
	Result  Getinfo `json:"result"`
}

func (o *AccountController) Getinfo() {
	var ob Getinfo
	var result Result
	result.Jsonrpc = "2.0"
	result.Code = 2000
	result.Result = ob
	o.Data["json"] = result
	o.ServeJSON()
}

func (o *AccountController) GetNewAddress() {
	var ob Getinfo
	var result Result
	result.Jsonrpc = "2.0"
	result.Code = 2000
	result.Result = ob
	o.Data["json"] = result
	o.ServeJSON()
}

func (o *AccountController) ListAccounts() {
	var ob Getinfo
	var result Result
	result.Jsonrpc = "2.0"
	result.Code = 2000
	result.Result = ob
	o.Data["json"] = result
	o.ServeJSON()
}

func (o *AccountController) GetAccount() {
	var ob Getinfo
	var result Result
	result.Jsonrpc = "2.0"
	result.Code = 2000
	result.Result = ob
	o.Data["json"] = result
	o.ServeJSON()
}
