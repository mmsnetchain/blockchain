package main

import (
	"encoding/hex"
	"fmt"
	"github.com/prestonTao/utils"
	"time"
)

func main() {
	example2()
}

func example2() {
	bs := []byte("nihao")
	fmt.Println("start", hex.EncodeToString(bs))
	mh := utils.Multihash(bs)
	fmt.Println(hex.EncodeToString([]byte(mh)))
	fmt.Println("end", mh.Data(), hex.EncodeToString(mh.Data()))

}

func example1() {
	data := make(map[string]int)
	fmt.Println("start")
	time.Sleep(time.Second * 20)

	for i := 0; i < 1000000; i++ {
		time.Sleep(time.Nanosecond)
		one := hex.EncodeToString(utils.GetHashForDomain(time.Now().Format("2006-01-02 15:04:05.999999999")))
		data[one] = 0
	}
	fmt.Println("end", len(data))
	time.Sleep(time.Second * 60)
}
