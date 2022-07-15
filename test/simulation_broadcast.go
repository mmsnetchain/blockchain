package main

import (
	"bytes"
	"fmt"
	"github.com/prestonTao/libp2parea/nodeStore"
	"math/big"
	"mmschainnewaccount/config"
	gconfig "mmschainnewaccount/config"
	"strconv"
)

func main() {

	Run2()
	fmt.Println("=============================")

}

func BuildIds() []*NodeOne {

	ids := []string{
		"8cWMZz7FvkcRr3EA1VbzGZWdMKR7Uy4VYB2bEyzqdb29",

		"8NjRXxa6KKSv9NtP2KnwzRAQGuMnqARKznCUa1U9LmR1",
		"x8sXmm9HKqyit7Pad5U5Y4o8WFPx9s3aNxQCmJfZWNk",
		"D3BbKJJyaQhWiT6eKiSryX22ZdX8AfduE1sb4Bp4PQYQ",
		"BqfZe98Sb8iiajFKHf9ndSuyapd5jXrkUSoDbCLFsfuW",
		"AhSUaV5hJ6U5Pgf7fwimpZPK7cmMwmpfmcv1fXxKJHwx",
	}

	nodes := make([]*NodeOne, 0)

	for n := 0; n < len(ids); n++ {
		nodeOne := &NodeOne{
			Nodes: make([]nodeStore.AddressNet, 0),
		}
		fmt.Println(n+1, "", ids[n])
		index := n

		idMH := nodeStore.AddressFromB58String(ids[index])
		idsm := nodeStore.NewIds(idMH, gconfig.NodeIDLevel)
		nodeOne.Self = idMH
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
			nodeOne.Nodes = append(nodeOne.Nodes, idOne)

			fmt.Println("    --", idOne.B58String())
		}
		fmt.Println()
		nodes = append(nodes, nodeOne)
	}

	return nodes

}

func Run2() {
	nodes := BuildIdsTwo()
	PrintNodes(nodes)

	broadcast(nodes)
}

func BuildIdsTwo() []*NodeOne {

	ids := []string{
		"2w5QBfujmLTAvesJRyRpxZFj4D4PJTEbhDVQJt1kbDmk",
		"8eWgP57sAKepA7h4FJV3ig4KekgNqE3exvhLLWtgDP57",
		"Bsyuy8Cpg5VWi69axQKaU6pLbHkWHffCDjcQEFJC1qEr",
		"4V2haHmFRdS5hp9VELnnVNkKwHmoUD41TSgUys28pskj",
		"5sok55G8osoiUJDbPLFpJbZUoowtxeWMHVyK9CvCS45g",
	}

	nodes := make([]*NodeOne, 0)

	for n := 0; n < len(ids); n++ {
		nodeOne := &NodeOne{
			Nodes: make([]nodeStore.AddressNet, 0),
		}
		idMH := nodeStore.AddressFromB58String(ids[n])

		nodeOne.Self = idMH
		nodes = append(nodes, nodeOne)
	}

	node0 := nodes[0]
	node0.Nodes = append(node0.Nodes, nodes[2].Self)
	node0.Nodes = append(node0.Nodes, nodes[3].Self)

	node1 := nodes[1]
	node1.Nodes = append(node1.Nodes, nodes[4].Self)

	node2 := nodes[2]
	node2.Nodes = append(node2.Nodes, nodes[0].Self)
	node2.Nodes = append(node2.Nodes, nodes[3].Self)

	node3 := nodes[3]
	node3.Nodes = append(node3.Nodes, nodes[0].Self)
	node3.Nodes = append(node3.Nodes, nodes[2].Self)

	return nodes
}

func PrintNodes(nodes []*NodeOne) {
	for _, one := range nodes {
		fmt.Println(one.Self.B58String(), one.Msg)
		for _, two := range one.Nodes {
			fmt.Println("  ", two.B58String())
		}
	}
}

func broadcast(nodes []*NodeOne) {
	for i, one := range nodes {

		for _, temp := range nodes {
			temp.Msg = false
		}
		fmt.Println("" + strconv.Itoa(i+1) + "")
		one.Msg = true

		for j, _ := range one.Nodes {
			tree := one.Nodes[j]
			for x, _ := range nodes {
				temp := nodes[x]
				if bytes.Equal(tree, temp.Self) {
					fmt.Println("src", tree.B58String(), "->", temp.Self.B58String())
					temp.Msg = true
					loopSend(temp, tree, nodes)
					break
				}
			}
		}

		PrintNodes(nodes)
	}
}
func loopSend(self *NodeOne, src nodeStore.AddressNet, nodes []*NodeOne) {
	farNodes := self.GetIdsForFar(src)
	for j, _ := range farNodes {
		tree := farNodes[j]
		for x, _ := range nodes {
			temp := nodes[x]
			if temp.Msg {
				continue
			}
			if bytes.Equal(tree, temp.Self) {
				fmt.Println("src", tree.B58String(), "->", temp.Self.B58String())
				temp.Msg = true
				loopSend(temp, tree, nodes)
				break
			}
		}
	}

}

type NodeOne struct {
	Self  nodeStore.AddressNet
	Nodes []nodeStore.AddressNet
	Msg   bool
}

func (this *NodeOne) GetIdsForFar(id nodeStore.AddressNet) []nodeStore.AddressNet {

	kl := nodeStore.NewKademlia()
	kl.Add(new(big.Int).SetBytes(this.Self))
	kl.Add(new(big.Int).SetBytes(id))

	for i, _ := range this.Nodes {
		kl.Add(new(big.Int).SetBytes(this.Nodes[i]))
	}

	list := kl.Get(new(big.Int).SetBytes(id))

	out := make([]nodeStore.AddressNet, 0)
	find := false
	for _, one := range list {

		if bytes.Equal(one.Bytes(), this.Self) {
			find = true
		} else {
			if find {

				mh := nodeStore.AddressNet(one.Bytes())
				out = append(out, mh)
			}
		}

	}

	return out
}

func BuildIdsTree() []*Node {
	ids := []string{
		"2w5QBfujmLTAvesJRyRpxZFj4D4PJTEbhDVQJt1kbDmk",
		"8eWgP57sAKepA7h4FJV3ig4KekgNqE3exvhLLWtgDP57",
		"Bsyuy8Cpg5VWi69axQKaU6pLbHkWHffCDjcQEFJC1qEr",
		"4V2haHmFRdS5hp9VELnnVNkKwHmoUD41TSgUys28pskj",
		"5sok55G8osoiUJDbPLFpJbZUoowtxeWMHVyK9CvCS45g",
	}

	nodes := make([]*Node, 0)

	for n := 0; n < len(ids); n++ {
		nodeOne := &Node{
			Nodes: make([]*Node, 0),
		}
		idMH := nodeStore.AddressFromB58String(ids[n])
		nodeOne.Self = idMH

		idBuilder := nodeStore.NewLogicNumBuider(nodeOne.Self, config.NodeIDLevel)
		nodeOne.LogicIds = idBuilder.GetNodeNetworkNum()
		nodes = append(nodes, nodeOne)
	}

	node0 := nodes[0]
	node0.Nodes = append(node0.Nodes, nodes[2])
	node0.Nodes = append(node0.Nodes, nodes[3])

	node1 := nodes[1]
	node1.Nodes = append(node1.Nodes, nodes[4])

	node2 := nodes[2]
	node2.Nodes = append(node2.Nodes, nodes[0])
	node2.Nodes = append(node2.Nodes, nodes[3])

	node3 := nodes[3]
	node3.Nodes = append(node3.Nodes, nodes[0])
	node3.Nodes = append(node3.Nodes, nodes[2])

	return nodes
}

func Select(nodes []*Node) {
	for n := 0; n < 1; n++ {
		fmt.Println("" + strconv.Itoa(n) + "")

		for _, nodeOne := range nodes {

			for _, idOne := range nodeOne.LogicIds {

				nodeOne.Send(nodeOne.Self, *idOne)

			}

		}

	}
}
func Print(nodes []*Node) {
	for _, one := range nodes {
		fmt.Println(one.Self.B58String())
		for _, two := range one.Nodes {
			fmt.Println("    ", two.Self.B58String())
		}
	}
}

type Node struct {
	Self        nodeStore.AddressNet
	Nodes       []*Node
	NodesClient []*Node
	LogicIds    []*[]byte
	Msg         bool
}

func (this *Node) Send(src, findId nodeStore.AddressNet) {
	id := this.FindNearNodeId(findId, nil, true)

	if bytes.Equal(id, this.Self) {

		this.ReturnMsg(src, this)
		return
	}
	for _, one := range this.Nodes {
		if bytes.Equal(one.Self, id) {
			one.Send(src, findId)
			return
		}
	}
}

func (this *Node) ReturnMsg(src nodeStore.AddressNet, findNode *Node) {
	id := this.FindNearNodeId(src, nil, true)
	if bytes.Equal(id, this.Self) {

		if bytes.Equal(this.Self, findNode.Self) {
			return
		}

		have := false
		for _, one := range this.Nodes {
			if bytes.Equal(one.Self, findNode.Self) {
				have = true
				break
			}
		}
		if have {
			return
		}
		fmt.Println(this.Self.B58String(), "", findNode.Self.B58String())
		return
	}
	for _, one := range this.Nodes {
		if bytes.Equal(one.Self, id) {
			one.ReturnMsg(src, findNode)
			return
		}
	}
}

func (this *Node) FindNearNodeId(nodeId, outId nodeStore.AddressNet, includeSelf bool) nodeStore.AddressNet {
	kl := nodeStore.NewKademlia()
	if includeSelf {
		kl.Add(new(big.Int).SetBytes(this.Self))
	}

	for _, one := range this.Nodes {
		if bytes.Equal(one.Self, outId) {
			continue
		}
		kl.Add(new(big.Int).SetBytes(one.Self))
	}

	targetIds := kl.Get(new(big.Int).SetBytes(nodeId))
	if len(targetIds) == 0 {
		return nil
	}
	targetId := targetIds[0]
	if targetId == nil {
		return nil
	}
	mh := nodeStore.AddressNet(targetId.Bytes())
	return mh
}

func run4() {
	nodes := BuildIds4()
	Select4(nodes)
}

func BuildIds4() []*Node {
	ids := []string{

		"8eWgP57sAKepA7h4FJV3ig4KekgNqE3exvhLLWtgDP57",
		"2w5QBfujmLTAvesJRyRpxZFj4D4PJTEbhDVQJt1kbDmk",
		"Bsyuy8Cpg5VWi69axQKaU6pLbHkWHffCDjcQEFJC1qEr",
		"4V2haHmFRdS5hp9VELnnVNkKwHmoUD41TSgUys28pskj",
		"7JBa2oUeYgYSp9FUsHy6wt5WuqigJkDcsYoAntig7eTt",
		"FPFEKi4MmDi9PssqkYbYjwTbzFnwqmanFZ7fwo7DNC1x",
		"84yEekKXynEx3SSaQjEQUr5JDf6B1Fp34Kn2hBNQmNZS",
		"DNDywcPsJqsWq2gn7gH4yZg5GrAZbR5JvbpxoJDhyoAs",
	}

	nodes := make([]*Node, 0)

	for n := 0; n < len(ids); n++ {
		nodeOne := &Node{
			Nodes: make([]*Node, 0),
		}
		idMH := nodeStore.AddressFromB58String(ids[n])
		nodeOne.Self = idMH

		idBuilder := nodeStore.NewLogicNumBuider(nodeOne.Self, config.NodeIDLevel)
		nodeOne.LogicIds = idBuilder.GetNodeNetworkNum()
		nodes = append(nodes, nodeOne)
	}

	for i, _ := range nodes {
		if i == 0 {
			continue
		}
		one := nodes[i]
		one.Nodes = append(one.Nodes, nodes[0])
	}

	return nodes
}

func Select4(nodes []*Node) {
	for n := 0; n < 10; n++ {
		fmt.Println("" + strconv.Itoa(n) + "")

		for _, nodeOne := range nodes {

			for _, logicNodeOne := range nodeOne.Nodes {
				idsm := nodeStore.NewIds(nodeOne.Self, nodeStore.NodeIdLevel)
				for _, one := range append(logicNodeOne.Nodes, logicNodeOne.NodesClient...) {
					if bytes.Equal(nodeOne.Self, one.Self) {
						continue
					}
					idsm.AddId(one.Self)
				}
				ids := idsm.GetIds()
				for _, one := range ids {

					for _, findNodeOne := range nodes {
						if bytes.Equal(findNodeOne.Self, one) {
							nodeOne.AddNode(findNodeOne)
							break
						}
					}
				}
			}

		}

		Print(nodes)
	}
}

func (this *Node) AddNode(node *Node) {

	if bytes.Equal(this.Self, node.Self) {
		return
	}

	idm := nodeStore.NewIds(this.Self, nodeStore.NodeIdLevel)

	for _, one := range this.Nodes {
		idm.AddId(one.Self)
	}

	ok, removeIDs := idm.AddId(node.Self)
	if ok {
		this.Nodes = append(this.Nodes, node)
		node.NodesClient = append(node.NodesClient, this)

		for _, one := range removeIDs {
			for i, nodeOne := range this.Nodes {
				if bytes.Equal(nodeOne.Self, one) {
					temp := this.Nodes[:i]
					temp = append(temp, this.Nodes[i+1:]...)
					this.Nodes = temp

					for j, clientOne := range nodeOne.NodesClient {
						if bytes.Equal(clientOne.Self, nodeOne.Self) {
							temp := nodeOne.NodesClient[:j]
							temp = append(temp, nodeOne.NodesClient[j+1:]...)
							nodeOne.NodesClient = temp
							break
						}
					}
					break
				}
			}
		}
	}

	return

}
