package main

import (
	"fmt"
	"mmschainnewaccount/wallet/mining"
)

func main() {
	bt1 := mining.BallotTicket{
		Addr:        "a",
		Miner:       "hong",
		GroupHeight: 1,
	}
	mining.AddBallotTicket(&bt1)

	bt2 := mining.BallotTicket{
		Addr:        "b",
		Miner:       "hong",
		GroupHeight: 1,
	}
	mining.AddBallotTicket(&bt2)

	bt3 := mining.BallotTicket{
		Addr:        "c",
		Miner:       "hong",
		GroupHeight: 1,
	}
	mining.AddBallotTicket(&bt3)

	bt4 := mining.BallotTicket{
		Addr:        "c",
		Miner:       "hong",
		GroupHeight: 2,
	}
	mining.AddBallotTicket(&bt4)
	bt5 := mining.BallotTicket{
		Addr:        "c",
		Miner:       "hong",
		GroupHeight: 2,
	}
	mining.AddBallotTicket(&bt5)

	total := mining.FindTotal(2)

	fmt.Println(total)
}
