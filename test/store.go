package main

import (
	"fmt"
	"mmschainnewaccount/store/fs"
)

func main() {
	ch := make(chan int)
	s := fs.NewSpace()
	go s.Init()
	free := s.FreeSpace()
	fmt.Println(free)

	<-ch
}
