package main

import (
	"fmt"
	"github.com/prestonTao/libp2parea/cache_store"
	"github.com/prestonTao/libp2parea/nodeStore"
	"strconv"
)

func main() {
	go read()
	Start()
}

func Start() {

	for i := 0; i < 10; i++ {

		name := cache_store.GetAddressInName(strconv.Itoa(i))
		fmt.Println(name.SuperPeerId.GetIdStr())
	}

}

func read() {
	for one := range cache_store.OutFindName {

		id := []byte{1, 1, 1, 1, 1}
		tid := nodeStore.NewTempId(id, id)
		name := cache_store.NewName(one, []*nodeStore.TempId{tid}, id)

		cache_store.AddAddressInName(one, name)

	}

}
