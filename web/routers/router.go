package routers

import (
	"mmschainnewaccount/rpc"
	"mmschainnewaccount/web/controllers"
	"mmschainnewaccount/web/controllers/anonymousnet"
	"mmschainnewaccount/web/controllers/sharebox"
	"mmschainnewaccount/web/controllers/store"
	"mmschainnewaccount/web/controllers/wallet"
	"sync"

	"github.com/astaxie/beego"
)

var routerLock = new(sync.RWMutex)

func Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) (app *beego.App) {
	routerLock.Lock()
	app = beego.Router(rootpath, c, mappingMethods...)
	routerLock.Unlock()
	return
}

func Start() {
	Router("/", &controllers.MainController{}, "get:Test")

	Router("/self/test", &controllers.MainController{}, "get:Test")

}

func RegisterStore() {
	Router("/store/getlist", &store.Index{}, "get:GetList")
	Router("/store/addfile", &store.Index{}, "post:AddFile")
	Router("/store/addcryptfile", &store.Index{}, "post:AddCryptFile")
	Router("/store/:hash", &store.Index{}, "get:GetFile")
	Router("/store/:down/:hash", &store.Index{}, "get:GetFile")
}

func RegisterSharebox() {
	Router("/sharebox/page", &sharebox.Index{}, "get:Index")
	Router("/sharebox/getlist", &sharebox.Index{}, "get:GetList")
	Router("/sharebox/addfile", &sharebox.Index{}, "post:AddFile")
	Router("/sharebox/:hash", &sharebox.Index{}, "get:GetFile")
}

func RegisterWallet() {

	Router("/self/getinfo", &wallet.Index{}, "post:Getinfo")
	Router("/self/block", &wallet.Index{}, "post:Block")
	Router("/self/witnesslist", &wallet.Index{}, "post:GetWitnessList")
}
func RegisterRpc() {
	Router("/rpc", &rpc.Bind{}, "post:Index")
}

func RegisterAnonymousNet() {
	Router("/*", &anonymousnet.MainController{}, "*:Agent")
	Router("/:urls", &anonymousnet.MainController{}, "get:AgentToo")
}
