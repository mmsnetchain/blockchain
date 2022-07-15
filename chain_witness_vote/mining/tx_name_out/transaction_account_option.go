package tx_name_out

import (
	"github.com/prestonTao/keystore/crypto"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
)

func NameOut(srcAddr, addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string, name string) (mining.TxItr, error) {

	txItr, err := mining.GetLongChain().GetBalance().BuildOtherTx(config.Wallet_tx_type_account_cancel,
		srcAddr, addr, 0, gas, frozenHeight, pwd, comment, name)
	if err != nil {

	} else {

	}
	return txItr, err
}
