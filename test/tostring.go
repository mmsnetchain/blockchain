package main

import (
	"fmt"
	"mmschainnewaccount/cache"
)

func main() {
	b58 := cache.To58String([]byte("ok"))
	fmt.Println(b58)
}
