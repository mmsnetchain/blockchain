package mining

import (
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"sync"
)

func init() {

	utils.Go(loopAddWitness)
}

var newWitnessAddrNets = make(chan *[]*crypto.AddressCoin, 100)
var addrnetWhilelist = new(sync.Map)

func loopAddWitness() {
	for one := range newWitnessAddrNets {
		chain := GetLongChain()
		if !chain.SyncBlockFinish {
			continue
		}
		lookupAddrs := new(go_protos.RepeatedBytes)
		for _, one := range *one {
			if findWitnessAddrNet(one) == nil {
				lookupAddrs.Bss = append(lookupAddrs.Bss, *one)
			}
		}
		continue
		LookupWitnessAddrNet(lookupAddrs)
	}
}

func addWitnessAddrNet(addrCoin *crypto.AddressCoin, addrNet *nodeStore.AddressNet) {
	addrnetWhilelist.Store(utils.Bytes2string(*addrCoin), addrNet)
}

func findWitnessAddrNet(addrCoin *crypto.AddressCoin) *nodeStore.AddressNet {
	value, ok := addrnetWhilelist.Load(utils.Bytes2string(*addrCoin))
	if ok {
		if value != nil {
			addrNet := value.(*nodeStore.AddressNet)
			return addrNet
		}
	}
	return nil
}

func AddWitnessAddrNets(addrs []*crypto.AddressCoin) {
	select {
	case newWitnessAddrNets <- &addrs:
	default:
	}
}

func LookupWitnessAddrNet(lookupAddrs *go_protos.RepeatedBytes) {
	if lookupAddrs == nil || len(lookupAddrs.Bss) <= 0 {
		return
	}

	bs, err := lookupAddrs.Marshal()
	if err != nil {
		engine.Log.Error("LookupWitnessAddrNet error:%s", err.Error())
		return
	}

	ok := message_center.SendMulticastMsg(config.MSGID_multicast_find_witness, &bs)
	if !ok {
		return
	}

}
