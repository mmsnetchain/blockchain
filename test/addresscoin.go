package main

import (
	"fmt"
	"github.com/prestonTao/keystore/crypto"
	"time"
)

var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func main() {
	example3()
}

func example3() {

	pre := "SELF"

	addrCoinStr := "SELFEzLXgGoCeb1tygrRt2qYTQ7MfTu7Sg8gr5"
	addrCoin := crypto.AddressFromB58String(addrCoinStr)

	fmt.Println(crypto.ValidAddr(pre, addrCoin))

	start := time.Now()
	for i := 0; i < 1000000; i++ {
		crypto.ValidAddr(pre, addrCoin)
	}
	fmt.Println("", time.Now().Sub(start))

}

func example2() {

	addrCoinStr := "SELFEzLXgGoCeb1tygrRt2qYTQ7MfTu7Sg8gr5"
	addrCoin := crypto.AddressFromB58String(addrCoinStr)

	start := time.Now()
	for i := 0; i < 10000; i++ {
		addrCoin.B58String()
	}
	fmt.Println("1", time.Now().Sub(start))

	start = time.Now()
	for i := 0; i < 10000; i++ {
		base58.Encode(addrCoin)
	}
	fmt.Println("2", time.Now().Sub(start))

}

func example1() {
	pre := ""

	addr := crypto.BuildAddr(pre, []byte("fdafjkdlfajkldajf"))

	fmt.Println(crypto.ValidAddr(pre, addr))

	addrStr := addr.B58String()
	fmt.Println(addrStr)

	addr = crypto.AddressFromB58String(addrStr)
	fmt.Println(addr.B58String())

	ok := crypto.ValidAddr(pre, addr)
	fmt.Println(ok)
}
