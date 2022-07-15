package main

import (
	"encoding/hex"
	"fmt"
	"github.com/prestonTao/libp2parea/nodeStore"
	"github.com/prestonTao/utils"
	"math/big"
	gconfig "mmschainnewaccount/config"
)

func main() {
	BuildIds()

	distance()

}

func BuildIds() {

	fmt.Println("---------- BuildIds ----------")

	ids := []string{

		"Bsyuy8Cpg5VWi69axQKaU6pLbHkWHffCDjcQEFJC1qEr",

		"2w5QBfujmLTAvesJRyRpxZFj4D4PJTEbhDVQJt1kbDmk",
		"D79GZyCBcNKyvc3pHco6nvnzJuTGiaYDLqFjkLG1fkcR",
		"84yEekKXynEx3SSaQjEQUr5JDf6B1Fp34Kn2hBNQmNZS",
	}

	for n := 0; n < len(ids); n++ {
		fmt.Println("", ids[n])
		index := n

		idMH := nodeStore.AddressFromB58String(ids[index])
		idsm := nodeStore.NewIds(idMH, gconfig.NodeIDLevel)
		for i, one := range ids {
			if i == index {
				continue
			}

			idMH := nodeStore.AddressFromB58String(one)
			idsm.AddId(idMH)

		}

		is := idsm.GetIds()
		for _, one := range is {

			idOne := nodeStore.AddressNet(one)

			fmt.Println("--", idOne.B58String())
		}

	}

}

func distance() {

	fmt.Println("---------- distance ----------")

	ids := []string{
		"W1aLWC4unTJZhSFc4VNLFsazAJ1PyTocV7agmteQDL3J3N",
		"W1gfVGa52yUJ4Gws4TiA9YbwGP8qCGgaYeeT8APjSiNk6U",
		"W1j9RJ1xYHaoAuRk2HGBrVA82njoxFAoctYKQMH43k8hXu",
		"W1atFt7bJ5Ubk4MXuV5GfsEYE7srWXR51exDgUEJcVr5fZ",
		"W1n9XtbLAjRsh9sr2kbwfkfy3VGenyhazbHJwrEYsnDZ8M",
	}

	index := 4

	kl := nodeStore.NewKademlia()
	for i, one := range ids {
		if i == index {
			continue
		}

		idMH, _ := utils.FromB58String(one)
		kl.Add(new(big.Int).SetBytes(idMH.Data()))

	}

	idMH, _ := utils.FromB58String(ids[index])
	is := kl.Get(new(big.Int).SetBytes(idMH.Data()))
	src := new(big.Int).SetBytes(idMH.Data())

	for _, one := range is {
		tag := new(big.Int).SetBytes(one.Bytes())
		juli := tag.Xor(tag, src)

		bs, err := utils.Encode(one.Bytes(), gconfig.HashCode)
		if err != nil {
			fmt.Println("")
			continue
		}
		mh := utils.Multihash(bs)

		fmt.Println("", mh.B58String(), "", hex.EncodeToString(juli.Bytes()))
	}

}
