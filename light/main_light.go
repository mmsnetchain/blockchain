package main

import (
	"fmt"
	"github.com/prestonTao/libp2parea/config"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote"
	"mmschainnewaccount/web"
	"mmschainnewaccount/web/routers"

	"github.com/astaxie/beego"
)

func main() {
	StartUP()
	fmt.Println("end")
}

func StartUP() {

	nodeStore.NodeSelf.Addr = config.Init_LocalIP
	nodeStore.NodeSelf.TcpPort = config.Init_LocalPort
	nodeStore.NodeSelf.IsSuper = true

	errInt := nodeStore.InitNodeStore()
	if errInt == 1 {
		panic("password fail")
	}

	core.StartEngine()

	web.Start()

	go func() {

		routers.RegisterWallet()
		err := chain_witness_vote.Register()
		if err != nil {
			fmt.Println(err)
		}
	}()

	routers.RegisterRpc()
	go func() {
	}()

	go func() {
		beego.Run()
	}()
	<-utils.GetStopService()

}
