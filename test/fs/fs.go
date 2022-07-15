package main

import (
	"fmt"
	"github.com/prestonTao/libp2parea/virtual_node"
	"mmschainnewaccount/store/fs"
)

func main() {

	fiTable := fs.FileindexSelf{
		Name:   "123",
		Vid:    virtual_node.AddressNetExtend([]byte("haoayou")),
		FileId: virtual_node.AddressNetExtend([]byte("nihao")),
		Value:  []byte("buhao"),
	}
	err := fiTable.Add(&fiTable)

	fmt.Println("end", err)
}
