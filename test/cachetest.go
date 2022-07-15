package main

import (
	"fmt"

	"github.com/prestonTao/libp2parea/cache"
	"github.com/prestonTao/libp2parea/nodeStore"
)

func main() {
	c := make(chan int)
	core.StartEngine()
	cache.Register()
	cachedata := cache.BuildCacheData([]byte("key"), []byte("value"))
	cachedata.AddOwnId(nodeStore.NodeSelf.IdInfo.Id)
	cachedata.SetTime()
	cache.Save(cachedata)
	cache.AddSyncDataTask([]byte("key"))

	fmt.Println(cache.Get([]byte("key")))
	fmt.Printf("%+v", cache.GetCacheData([]byte("key")))
	<-c
}
