package name

import (
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"

	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

func FindNameToNet(name string) *Nameinfo {
	dbKey := append([]byte(config.Name), []byte(name)...)
	txbs, err := db.LevelTempDB.Find(dbKey)
	if err != nil {
		return nil
	}
	if txbs == nil || len(*txbs) <= 0 {
		return nil
	}
	nameinfo, err := ParseNameinfo(*txbs)
	if err != nil {
		return nil
	}
	return nameinfo
}

func FindNameToNetRandOne(name string, height uint64) *nodeStore.AddressNet {
	nameinfo := FindNameToNet(name)
	if nameinfo == nil {
		return nil
	}
	if nameinfo.CheckIsOvertime(height) {
		return nil
	}
	addr := nameinfo.NetIds[utils.GetRandNum(int64(len(nameinfo.NetIds)))]
	return &addr
}
