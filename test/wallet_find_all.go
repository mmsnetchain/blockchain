package main

import (
	"encoding/json"
	"fmt"
	"github.com/prestonTao/utils"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/mining/token/payment"
	"mmschainnewaccount/config"
	"sync"
)

var TagTxSubmiting string = "TxSub_"

func init() {
	tpc := new(payment.TokenPublishController)
	tpc.ActiveVoutIndex = new(sync.Map)
	mining.RegisterTransaction(config.Wallet_tx_type_account, tpc)

	err := db.InitDB(config.DB_path)
	if err != nil {
		fmt.Println("init db error:", err.Error())
		panic(err)
	}

}
func main() {
	db.PrintAll()

}

func findAll() {
	keys, values, _ := db.FindPrefixKeyAll([]byte(TagTxSubmiting))
	fmt.Println(len(keys), len(values))

	outBs := make([]byte, 0)
	for _, one := range values {

		txItr, err := mining.ParseTxBase(0, &one)
		panicError(err)
		txVO := txItr.GetVOJSON()

		txBs, err := json.Marshal(txVO)
		panicError(err)

		outBs = append(outBs, txBs...)
		outBs = append(outBs, []byte("\n")...)

	}

	utils.SaveFile("tx.txt", &outBs)
}
func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
