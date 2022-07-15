package core

import (
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"sort"
	"sync"
	"time"
)

var netSpeedMap = new(sync.Map)

func AddNodeAddrSpeed(addr nodeStore.AddressNet, speed time.Duration) {
	netSpeedMap.Store(utils.Bytes2string(addr), speed)
}

func DelNodeAddrSpeed(addr nodeStore.AddressNet) {
	netSpeedMap.Delete(utils.Bytes2string(addr))
}

func SortNetAddrForSpeed(netAddrs []nodeStore.AddressNet) []AddrNetSpeedInfo {
	addrNetSpeedASC := NewAddrNetSpeedASC(netAddrs)
	return addrNetSpeedASC.Sort()
}

type AddrNetSpeedASC struct {
	addrs []AddrNetSpeedInfo
}

func (this AddrNetSpeedASC) Len() int {
	return len(this.addrs)
}

func (this AddrNetSpeedASC) Less(i, j int) bool {

	if this.addrs[i].Speed > this.addrs[j].Speed {
		return false
	} else {
		return true
	}
}

func (this AddrNetSpeedASC) Swap(i, j int) {
	this.addrs[i], this.addrs[j] = this.addrs[j], this.addrs[i]
}

func (this AddrNetSpeedASC) Sort() []AddrNetSpeedInfo {

	sort.Stable(this)

	return this.addrs
}

func NewAddrNetSpeedASC(addrs []nodeStore.AddressNet) *AddrNetSpeedASC {
	addrNetSpeedASC := new(AddrNetSpeedASC)

	addrMap := make(map[string]nodeStore.AddressNet)
	for i, one := range addrs {
		key := utils.Bytes2string(one)
		if _, ok := addrMap[key]; ok {
			continue
		}
		info := AddrNetSpeedInfo{
			AddrNet: addrs[i],
			Speed:   0,
		}
		value, ok := netSpeedMap.Load(key)
		if ok {
			info.Speed = int64(value.(time.Duration))
		}
		addrMap[key] = addrs[i]
		addrNetSpeedASC.addrs = append(addrNetSpeedASC.addrs, info)
	}

	return addrNetSpeedASC
}

type AddrNetSpeedInfo struct {
	AddrNet nodeStore.AddressNet
	Speed   int64
}
