package tx_spaces_mining_in

import (
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/nodeStore"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"
)

func SpacesMiningIn(srcAddr, addr *crypto.AddressCoin, amount, gas, frozenHeight uint64, pwd, comment string,
	name string, netIds []nodeStore.AddressNet, addrCoins []crypto.AddressCoin) (mining.TxItr, error) {

	txItr, err := mining.GetLongChain().GetBalance().BuildOtherTx(config.Wallet_tx_type_spaces_mining_in,
		srcAddr, addr, amount, gas, frozenHeight, pwd, comment, name, netIds, addrCoins)
	if err != nil {

	} else {

	}
	return txItr, err
}
