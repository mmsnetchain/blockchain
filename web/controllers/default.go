package controllers

import (
	"fmt"
	gconfig "mmschainnewaccount/config"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/prestonTao/libp2parea/cache_store"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/nodeStore"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {

	this.Ctx.Redirect(http.StatusFound, "/self/test")
}

func (this *MainController) Test() {

	this.Data["Ip"] = nodeStore.NodeSelf.Addr

	this.Data["RootExist"] = cache_store.Root.Exist

	this.Data["IsSuper"] = nodeStore.NodeSelf.IsSuper
	if nodeStore.SuperPeerId != nil {
		this.Data["SuperId"] = nodeStore.SuperPeerId.B58String()
	} else {
		this.Data["SuperId"] = ""
	}

	this.Data["ID"] = nodeStore.NodeSelf.IdInfo.Id.B58String()

	ids := nodeStore.GetLogicNodes()
	idsStr := make([]string, 0)
	for _, one := range ids {
		idsStr = append(idsStr, one.B58String())
	}
	this.Data["ids"] = idsStr

	names := cache_store.Debug_GetAllName()
	this.Data["names"] = names

	this.TplName = "test.tpl"
}

func (this *MainController) SendMeg() {
	id := this.GetString("id")
	recvId := nodeStore.AddressFromB58String(id)

	content := []byte(this.GetString("content"))

	message_center.SendP2pMsgHE(gconfig.MSGID_TextMsg, &recvId, &content)

	out := make(map[string]interface{})
	out["Code"] = 0
	this.Data["json"] = out
	this.ServeJSON(true)
}

func (this *MainController) AgentToo() {
	fmt.Println("")
}

func (this *MainController) BtTest() {

	out := make(map[string]interface{})
	out["Code"] = 0
	this.Data["json"] = out
	this.ServeJSON(true)
}
