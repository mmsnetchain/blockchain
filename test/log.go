package main

import (
	"github.com/prestonTao/libp2parea/engine"
	"time"
)

func main() {
	engine.NLog.Debug(engine.LOG_file, "%s", "nihao")
	engine.NLog.Error(engine.LOG_file, "%s", "")
	time.Sleep(time.Second)
}
