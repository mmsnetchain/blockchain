package main

import (
	"fmt"
	"mmschainnewaccount/rpc"
)

func main() {
	fmt.Println("start...")
	rpc.RegisterRpcServer()
}
