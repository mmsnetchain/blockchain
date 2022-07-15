package main

import (
	"github.com/prestonTao/libp2parea/config"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/boot"
	_ "mmschainnewaccount/boot"
	"mmschainnewaccount/chain_witness_vote"

	"mmschainnewaccount/web"
	"mmschainnewaccount/web/routers"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/astaxie/beego"
)

func main() {

	boot.Step()
	StartUP()

}

func pprofMem() {

	runtime.GOMAXPROCS(1)
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	runtime.MemProfileRate = 512 * 1024

	time.Sleep(time.Minute * 5)
	stopMemProfile("mem.prof")
}

func stopMemProfile(memProfile string) {
	f, err := os.Create(memProfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create mem profile output file: %s", err)
		return
	}
	if err = pprof.WriteHeapProfile(f); err != nil {
		fmt.Fprintf(os.Stderr, "Can not write %s: %s", memProfile, err)
	}
	f.Close()
}

func StartUP() {

	nodeStore.NodeSelf.Addr = config.Init_LocalIP
	nodeStore.NodeSelf.TcpPort = config.Init_LocalPort
	nodeStore.NodeSelf.IsSuper = true

	errInt := nodeStore.InitNodeStore()
	if errInt == 1 {
		panic("")
	}

	core.StartEngine()

	web.Start()

	go func() {

		err := chain_witness_vote.Register()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {

		routers.RegisterRpc()
	}()

	go func() {
		beego.Run()
	}()
	<-utils.GetStopService()

}
