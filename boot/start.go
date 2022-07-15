package boot

import (
	"mmschainnewaccount/chain_witness_vote"
	"mmschainnewaccount/cloud_space"
	"mmschainnewaccount/config"
	"mmschainnewaccount/proxyhttp"
	"mmschainnewaccount/sharebox"
	"mmschainnewaccount/web"
	"mmschainnewaccount/web/routers"

	"github.com/astaxie/beego"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/libp2parea"
	pconfig "github.com/prestonTao/libp2parea/config"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

func StartUP(passwd string) {

	config.Step()
	if passwd != "" {
		config.Wallet_keystore_default_pwd = passwd
	}

	err := keystore.Load(config.KeystoreFileAbsPath, config.AddrPre)
	if err != nil {
		keystore.CreateKeystore(config.KeystoreFileAbsPath, config.AddrPre, config.Wallet_keystore_default_pwd)
	}

	pconfig.SQLITE3DB_path = "msgcache.db"
	if config.NetType == config.NetType_release {
		pconfig.NetType = pconfig.NetType_release
	} else {
		pconfig.NetType = pconfig.NetType_test
	}

	libp2parea.Start(config.ParseInitFlag(), pconfig.Init_LocalIP, pconfig.Init_LocalPort,
		config.KeystoreFileAbsPath, config.AddrPre, config.Wallet_keystore_default_pwd)

	web.Start()

	go func() {

		routers.RegisterWallet()
		err := chain_witness_vote.Register()
		if err != nil {
			engine.Log.Error(err.Error())
		}
	}()

	routers.RegisterRpc()
	go func() {
	}()

	routers.RegisterAnonymousNet()
	go func() {
		proxyhttp.Register()
	}()

	routers.RegisterSharebox()
	go func() {
		sharebox.RegsterStore()

	}()

	go func() {

		cloud_space.RegsterCloudSpace()

	}()

	go func() {
		beego.Run()
	}()
	<-utils.GetStopService()

}
