package wallet

import (
	"mmschainnewaccount/chain_witness_vote"

	"github.com/astaxie/beego"
)

type Account struct {
	beego.Controller
}

func (this *Account) GetInfo() {

	this.Data["CheckKey"] = chain_witness_vote.CheckKey()

	this.TplName = "wallet/index.tpl"
}
