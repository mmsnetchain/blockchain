package main

import (
	"fmt"
	"github.com/prestonTao/libp2parea/message_center"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var count int32 = 0

func main() {
	go in()
	Start()
}

func Start() {

	var group = new(sync.WaitGroup)
	for i := 0; i < 5; i++ {
		group.Add(1)
		go func() {

			bs := message_center.WaitRequest(message_center.MSG_WAIT_http_request, strconv.Itoa(int(atomic.AddInt32(&count, 1))))
			fmt.Println("", bs)
			group.Done()
		}()

	}
	group.Wait()

}

func in() {
	for i := 0; i < 6; i++ {
		time.Sleep(time.Second * 1)
		message_center.ResponseWait(message_center.MSG_WAIT_http_request, strconv.Itoa(i), &[]byte{123})
	}
}
