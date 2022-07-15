package name

import (
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"sync"

	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/nodeStore"

	"github.com/gogo/protobuf/proto"
)

var NameOfValidity = uint64(0)
var names = new(sync.Map)

func init() {

	NameOfValidity = 60 * 60 * 24 * 365 / config.Mining_block_time

}

func AddName(name Nameinfo) {
	names.Store(string(name.Name), name)
}

func DelName(name []byte) {
	names.Delete(string(name))
}

func FindName(name string) *Nameinfo {
	itr, ok := names.Load(name)
	if !ok {
		return nil
	}
	nameinfo := itr.(Nameinfo)
	return &nameinfo
}

func GetNameList() []Nameinfo {
	name := make([]Nameinfo, 0)
	names.Range(func(k, v interface{}) bool {
		nameOne := v.(Nameinfo)
		nameOne.NameOfValidity = nameOne.Height + NameOfValidity
		name = append(name, nameOne)
		return true
	})
	return name
}

type Nameinfo struct {
	Name           string
	Txid           []byte
	NetIds         []nodeStore.AddressNet
	AddrCoins      []crypto.AddressCoin
	Height         uint64
	NameOfValidity uint64
	Deposit        uint64
}

func (this *Nameinfo) Proto() ([]byte, error) {
	netids := make([][]byte, 0)
	for _, one := range this.NetIds {
		netids = append(netids, one)
	}
	addrCoins := make([][]byte, 0)
	for _, one := range this.AddrCoins {
		addrCoins = append(addrCoins, one)
	}
	nip := go_protos.Nameinfo{
		Name:           this.Name,
		Txid:           this.Txid,
		NetIds:         netids,
		AddrCoins:      addrCoins,
		Height:         this.Height,
		NameOfValidity: this.NameOfValidity,
		Deposit:        this.Deposit,
	}
	return nip.Marshal()
}

func ParseNameinfo(bs []byte) (*Nameinfo, error) {
	nip := new(go_protos.Nameinfo)
	err := proto.Unmarshal(bs, nip)
	if err != nil {
		return nil, err
	}

	netids := make([]nodeStore.AddressNet, 0)
	for _, one := range nip.NetIds {
		netids = append(netids, one)
	}
	addrCoins := make([]crypto.AddressCoin, 0)
	for _, one := range nip.AddrCoins {
		addrCoins = append(addrCoins, one)
	}
	nameinfo := Nameinfo{
		Name:           nip.Name,
		Txid:           nip.Txid,
		NetIds:         netids,
		AddrCoins:      addrCoins,
		Height:         nip.Height,
		NameOfValidity: nip.NameOfValidity,
		Deposit:        nip.Deposit,
	}
	return &nameinfo, nil
}

func (this *Nameinfo) CheckIsOvertime(height uint64) bool {
	if (this.Height + NameOfValidity) > height {
		return false
	}
	return true
}
