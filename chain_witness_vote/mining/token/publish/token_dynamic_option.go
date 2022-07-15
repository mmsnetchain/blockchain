package publish

import (
	"github.com/prestonTao/keystore/crypto"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
)

const (
	Wallet_tx_class = config.Wallet_tx_type_token_publish
)

func PublishToken(srcAddr, addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string,
	name, symbol string, supply uint64, owner crypto.AddressCoin) (mining.TxItr, error) {

	txItr, err := mining.GetLongChain().GetBalance().BuildOtherTx(Wallet_tx_class,
		srcAddr, addr, amount, gas, frozenHeight, pwd, comment, name, symbol, supply, owner)
	if err != nil {

	} else {

	}
	return txItr, err
}
