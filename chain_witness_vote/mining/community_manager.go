package mining

import (
	"time"
)

const communityRewardCountUnconfirmedHeight = uint64(20)

var communityRewardHeightStart = uint64(0)

func init() {

}
func loopSelectTx() {
	for {

		chain := GetLongChain()
		if chain == nil {
			time.Sleep(time.Minute * 10)
			continue
		}
		if !chain.SyncBlockFinish {
			time.Sleep(time.Minute * 10)
			continue
		}
		time.Sleep(time.Minute * 10)
	}
}
