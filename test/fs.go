package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func main() {

	example3()

}

func example3() {

}

func inition() {

	var err error
	db, err = leveldb.OpenFile("db", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func example() {

	for i := 0; i < 2097152; i++ {
		key := []byte(strconv.Itoa(i))
		code := sha256.Sum256(key)
		value := make([]byte, 0)
		for j := 0; j < 16; j++ {
			value = append(value, code[:]...)
		}
		fmt.Println(len(value))
		db.Put(key, value, nil)
	}

}

func example2() {
	for i := 0; i < 2097152; i++ {
		key := sha256.Sum256([]byte(strconv.Itoa(i)))
		db.Put(key[:], randStr(), nil)
	}
}

func randStr() []byte {
	n := 512
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}
