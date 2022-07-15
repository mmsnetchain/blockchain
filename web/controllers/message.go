package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/prestonTao/libp2parea"
	"github.com/prestonTao/libp2parea/persistence"
)

type MsgController struct {
	beego.Controller
}

func (this *MsgController) MsgPage() {
	fs := persistence.Friends_getall()

	this.Data["Fs"] = fs

	this.TplName = "message.tpl"

}

func (this *MsgController) GetMsg() {
	out := make(map[string]interface{})

	overtime := time.NewTicker(time.Second * 20)
	select {
	case <-overtime.C:
		out["Code"] = 1
	case <-this.Ctx.ResponseWriter.CloseNotify():

		overtime.Stop()

	case msg := <-libp2parea.MsgChannl:
		overtime.Stop()
		out["Code"] = 0
		out["Id"] = msg.Id
		out["Index"] = msg.Index
		out["Content"] = msg.Content
	}
	this.Data["json"] = out
	this.ServeJSON(true)

}

func (this *MsgController) AddFriend() {
	out := make(map[string]interface{})
	id := this.GetString("ID")
	if id == "" {
		out["Code"] = 1
		this.Data["json"] = out
		this.ServeJSON(true)
		return
	}
	err := persistence.Friends_add(id)
	if err != nil {
		out["Code"] = 1
	} else {
		out["Code"] = 0
	}

	this.Data["json"] = out
	this.ServeJSON(true)
}
