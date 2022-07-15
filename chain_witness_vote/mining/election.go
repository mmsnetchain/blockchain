package mining

import (
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
)

type Election struct {
	GroupHeight uint64                `json:"groupheight"`
	Addr        *nodeStore.AddressNet `json:"addr"`
	Time        int64                 `json:"time"`
}

func NewElection(addr *nodeStore.AddressNet) *Election {
	return &Election{
		Addr: addr,

		Time: utils.GetNow(),
	}
}

type BallotTicket struct {
	Addr    *crypto.AddressCoin `json:"addr"`
	Puk     []byte              `json:"puk"`
	Sign    []byte              `json:"sign"`
	Witness *crypto.AddressCoin `json:"witness"`
	Deposit []byte              `json:"deposit"`
}

func MulticastBallotTicket(deposit *[]byte, addr *crypto.AddressCoin) {

}
