package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	jsoniter "github.com/json-iterator/go"
	"github.com/prestonTao/libp2parea/config"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	Path_configDir     = "conf"
	Path_config        = "config.json"
	Core_keystore      = "keystore.key"
	Core_addr_prk      = "addr_ec_prk.pem"
	Core_addr_puk      = "addr_ec_puk.pem"
	Core_addr_prk_type = "EC PRIVATE KEY"
	Core_addr_puk_type = "EC PUBLIC KEY"
)

const (
	Name_prk = "name_ec_prk.pem"
	Name_puk = "name_ec_puk.pem"
)

const (
	Store_path_dir            = "store"
	Store_path_fileinfo_self  = "self"
	Store_path_fileinfo_local = "local"
	Store_path_fileinfo_net   = "net"
	Store_path_fileinfo_cache = "cache"
	Store_path_temp           = "temp"
	Store_path_files          = "files"
	IsRemoveStore             = false

	HashCode    = utils.SHA3_256
	NodeIDLevel = 256

	Model_complete = "complete"
	Model_light    = "light"
)

var (
	WebAddr                = ""
	WebPort         uint16 = 0
	Web_path_static        = ""
	Web_path_views         = ""
	AddrPre                = ""
	NetIds                 = []byte{0}
	NetType_release        = "release"
	NetType                = "test"
	RpcServer              = false
	RPCUser                = ""
	RPCPassword            = ""
	Model                  = Model_complete

	Wallet_txitem_save_db = false
	Entry                 = []string{}
	CPUNUM                = runtime.NumCPU()
	OS                    = runtime.GOOS
)

var (
	KeystoreFileAbsPath = filepath.Join(Path_configDir, Core_keystore)

	Store_dir            string = filepath.Join(Store_path_dir)
	Store_fileinfo_self  string = filepath.Join(Store_path_dir, Store_path_fileinfo_self)
	Store_fileinfo_local string = filepath.Join(Store_path_dir, Store_path_fileinfo_local)
	Store_fileinfo_net   string = filepath.Join(Store_path_dir, Store_path_fileinfo_net)
	Store_fileinfo_cache string = filepath.Join(Store_path_dir, Store_path_fileinfo_cache)
	Store_temp           string = filepath.Join(Store_path_dir, Store_path_temp)
	Store_files          string = filepath.Join(Store_path_dir, Store_path_files)
)

func ParseConfig() {

	if OS == "windows" {
		Wallet_print_serialize_hex = false
	} else {
		Wallet_print_serialize_hex = false
	}

	ok, err := utils.PathExists(filepath.Join(Path_configDir, Path_config))
	if err != nil {

		panic("：" + err.Error())
		return
	}

	if !ok {

		cfi := new(Config)
		cfi.Port = 9981
		bs, _ := json.Marshal(cfi)

		f, err := os.OpenFile(filepath.Join(Path_configDir, Path_config), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {

			panic("：" + err.Error())
			return
		}

		_, err = f.Write(bs)
		if err != nil {

			panic("：" + err.Error())
			return
		}

		f.Close()
	}

	bs, err := ioutil.ReadFile(filepath.Join(Path_configDir, Path_config))
	if err != nil {

		panic("：" + err.Error())
		return
	}

	cfi := new(Config)

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err = decoder.Decode(cfi)
	if err != nil {

		panic("：" + err.Error())
		return
	}

	config.Init_LocalIP = cfi.IP
	config.Init_LocalPort = cfi.Port
	config.Init_GatewayPort = cfi.Port
	Web_path_static = cfi.WebStatic
	Web_path_views = cfi.WebViews
	engine.Netid = cfi.Netid
	WebAddr = cfi.WebAddr
	WebPort = cfi.WebPort
	Miner = cfi.Miner
	NetType = cfi.NetType

	AddrPre = cfi.AddrPre
	RpcServer = cfi.RpcServer
	RPCUser = cfi.RpcUser
	RPCPassword = cfi.RpcPassword
	Wallet_txitem_save_db = cfi.BalanceDB
	if cfi.Model == Model_light {
		Model = Model_light
	}
}

type Config struct {
	Netid       uint32 `json:"netid"`
	IP          string `json:"ip"`
	Port        uint16 `json:"port"`
	WebAddr     string `json:"WebAddr"`
	WebPort     uint16 `json:"WebPort"`
	WebStatic   string `json:"WebStatic"`
	WebViews    string `json:"WebViews"`
	RpcServer   bool   `json:"RpcServer"`
	RpcUser     string `json:"RpcUser"`
	RpcPassword string `json:"RpcPassword"`
	Miner       bool   `json:"miner"`
	NetType     string `json:"NetType"`
	AddrPre     string `json:"AddrPre"`
	BalanceDB   bool   `json:"balancedb"`
	Model       string `json:"model"`
}
