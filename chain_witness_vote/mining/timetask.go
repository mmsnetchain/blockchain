package mining

import (
	"runtime"
	"time"

	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

func (this *Witness) SyncBuildBlock(n int64) {
	goroutineId := utils.GetRandomDomain() + utils.TimeFormatToNanosecondStr()
	_, file, line, _ := runtime.Caller(0)
	engine.AddRuntime(file, line, goroutineId)
	defer engine.DelRuntime(file, line, goroutineId)

	timer := time.NewTimer(time.Second * time.Duration(n))
	select {
	case <-timer.C:
	case <-this.StopMining:

		timer.Stop()
		return
	}

	if GetLongChain().WitnessChain.witnessGroup.NextGroup == nil {
		GetLongChain().WitnessChain.BuildWitnessGroup(false, true)
		GetLongChain().WitnessChain.BuildMiningTime()
	}

	this.BuildBlock()

}
