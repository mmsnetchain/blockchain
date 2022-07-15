package main

import (
	"fmt"
	"mmschainnewaccount/wallet/mining"
	"time"
)

func main() {
	example2()
}

func example2() {
	bss := make([][]byte, 0)
	for i := 0; i < 9500000; i++ {
		id := []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks")
		bss = append(bss, id)
	}
	fmt.Println("end")
	time.Sleep(time.Minute)
}

func example1() {
	bhs := make([]mining.BlockHead, 0)
	for i := 0; i < 9500000; i++ {
		bh := mining.BlockHead{
			Hash:              []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
			Height:            1,
			GroupHeight:       1,
			Previousblockhash: []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
			Nextblockhash:     []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
			NTx:               1,
			MerkleRoot:        []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
			Tx: [][]byte{[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
				[]byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks")},
			Time:        1,
			BackupMiner: []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
			DepositId:   []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
			Witness:     []byte("12FRzz2xrVtEm9cwzgFArrLE7VA7ks"),
		}
		bhs = append(bhs, bh)
	}
	fmt.Println("ok")
	time.Sleep(time.Minute)
	fmt.Println("end")
}
