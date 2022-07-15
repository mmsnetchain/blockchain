package main

import (
	"fmt"
	"mmschainnewaccount/config"
)

func main() {
	example()
}

func example() {
	ok := config.CheckAddBlacklist(3, 5)
	fmt.Println("isok:", ok)
}
