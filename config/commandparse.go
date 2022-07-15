package config

import (
	"flag"
)

var Input_pstrName = flag.String("name", "gerry", "input ur name")
var Input_piAge = flag.Int("age", 20, "input ur age")
var Input_flagvar int

var (
	Init        = flag.String("init", "", "(：genesis.json)")
	Conf        = flag.String("conf", "conf/config.json", "(：conf/config.json)")
	Port        = flag.Int("port", 9811, "(：9811)")
	NetId       = flag.Int("netid", 20, "id(：20)")
	Ip          = flag.String("ip", "0.0.0.0", "IP(：0.0.0.0)")
	webAddr     = flag.String("webaddr", "0.0.0.0", "webIP(：0.0.0.0)")
	webPort     = flag.Int("webport", 2080, "web(：2080)")
	WebStatic   = flag.String("westatic", "", "web")
	WebViews    = flag.String("webviews", "", "web Views")
	DataDir     = flag.String("datadir", "", "")
	DbCache     = flag.Int("dbcache", 25, "，（MB）（：25）")
	TimeOut     = flag.Int("timeout", 0, "，")
	rpcServer   = flag.Bool("rpcserver", false, "JSON-RPC true/false(：false)")
	RpcUser     = flag.String("rpcuser", "", "JSON-RPC ")
	RpcPassword = flag.String("rpcpassword", "", "JSON-RPC ")
	WalletPwd   = flag.String("walletpwd", Wallet_keystore_default_pwd, "")
	classpath   = flag.String("classpath", "", "jar")
	Load        = flag.String("load", "", "")
)

func Step() {

	flag.Parse()
	parseParam()
	ParseConfig()

}

func StepDLL() {

	parseParam()
	ParseConfig()

}

func parseParam() {
	flag.VisitAll(func(v *flag.Flag) {
		switch v.Name {
		case "port":

		case "netid":

		case "webaddr":

		case "webport":

		case "westatic":

		case "webviews":

		case "datadir":

		case "dbcache":

		case "timeout":

		case "rpcserver":

		case "rpcuser":

		case "rpcpassword":

		case "walletpwd":
			Wallet_keystore_default_pwd = *WalletPwd

		}
	})

}
