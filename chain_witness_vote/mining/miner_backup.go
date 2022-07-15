package mining

import (
	"bytes"
	"github.com/prestonTao/utils"
	"sync"
)

var groupMinersLock = new(sync.RWMutex)
var groupMiners = make([]utils.Multihash, 0)

type BackupMiners struct {
	Time   int64
	Miners []BackupMiner
}

type BackupMiner struct {
	Miner utils.Multihash
	Count uint64
}

func AddGroupBackupMiner(miners ...*utils.Multihash) {
	groupMinersLock.Lock()

	for _, one := range miners {
		find := false
		for _, two := range groupMiners {
			if bytes.Equal(*one, two) {
				find = true
				break
			}
		}
		if find {
			continue
		}
		groupMiners = append(groupMiners, *one)
	}
	groupMinersLock.Unlock()
}

func TotalBackupMiner() (n uint64) {
	groupMinersLock.RLock()
	n = uint64(len(groupMiners))
	groupMinersLock.RUnlock()
	return
}

func RemoveGroupBackupMiner(miners ...utils.Multihash) {
	newMiners := make([]utils.Multihash, 0)
	groupMinersLock.Lock()
	for _, one := range groupMiners {
		find := false
		for _, two := range miners {
			if bytes.Equal(one, two) {
				find = true
				break
			}
		}
		if find {
		} else {
			newMiners = append(newMiners, one)
		}
	}
	groupMiners = newMiners
	groupMinersLock.Unlock()
}
