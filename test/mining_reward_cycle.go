package main

import (
	"fmt"
	"github.com/prestonTao/libp2parea/engine"
	"mmschainnewaccount/config"
	"strconv"
	"time"
)

const (
	year       = 60 * 60 * 24 * 365
	year20     = year * 20
	total      = 3000000000
	blockTime  = 60
	totalBlack = year20 / blockTime
	sysle      = 40
	reduce     = 0.9
)

func main() {

	cycleRange(790617, 790620)
	cycleRange(790621, 790624)
	return

	fmt.Println(float64(32 / 31))

	start := time.Now()

	fmt.Println("", 15190001, "", config.ClacRewardForBlockHeight(15190002))

	jiange := time.Now().Sub(start)
	fmt.Println("", jiange.Seconds(), jiange.Milliseconds(), jiange.Microseconds(), jiange.Nanoseconds())

	height := uint64(1)
	interval := uint64(1)
	total := uint64(0)
	for {
		reward := config.ClacRewardForBlockHeight(height)

		if height%3153600 == 0 {
			engine.Log.Info("%d  %d  %d  %d", height/3153600, height, reward, total)

		}

		total += reward

		if reward <= 0 {
			break
		}
		height += interval
	}

	engine.Log.Info(" %d  %d", height, total)
	jiange = time.Now().Sub(start)
	engine.Log.Info(" %d %d %d %d", jiange.Seconds(), jiange.Milliseconds(), jiange.Microseconds(), jiange.Nanoseconds())

}

func exampleold() {
	height := uint64(1)
	interval := uint64(12614400)
	for {
		reward := config.ClacRewardForBlockHeight(height)
		fmt.Println("", height, "", reward)
		if reward == 0 {
			break
		}
		height += interval
	}

	fmt.Println("", config.ClacRewardForBlockHeight(1))
	fmt.Println("", config.ClacRewardForBlockHeight(12614400))

	clacCycleForTotalAndInterval(1170000000*1e8, 4)

}

func clacCycle(first, total, intervalYear uint64) {
	count := uint64(0)
	cycle := 0
	for count < total {
		cycle++
		oneCycle := 2880 * 365 * intervalYear * first
		fmt.Printf(" %d , %d \n", cycle, oneCycle)
		count += oneCycle
		first = first / 2
		if first == 0 {
			break
		}
	}
	fmt.Printf(" %d", count)
}

func clacCycleForTotalAndInterval(total, intervalYear uint64) {
	first := total / 2 / (2880 * 365 * intervalYear)
	fmt.Printf(" %d\n", first)
	clacCycle(first, total, intervalYear)

}

func jisuan() {

	ARewardTotal := float64(total) / float64(sysle)

	tb := float64(totalBlack) / float64(sysle)

	A := float64(ARewardTotal*2) / 1.9

	a := strconv.FormatFloat(A, 'f', -1, 64)
	fmt.Println(a)

	tr := float64(0)

	fmt.Println("", strconv.FormatFloat(total, 'f', -1, 64))
	fmt.Println("20", strconv.FormatFloat(totalBlack, 'f', -1, 64))
	fmt.Println("40，", strconv.FormatFloat(tb, 'f', -1, 64))

	for i := 19; i > 0; i-- {
		total := A
		for j := 0; j < i; j++ {
			total = total / 0.9
		}
		x := total / tb
		fmt.Println("", 20-i, " ", tb, " ", x, " Token", strconv.FormatFloat(total, 'f', -1, 64))
		tr = tr + total
	}
	fmt.Println("------")

	for i := 0; i <= 20; i++ {
		total := A
		for j := 0; j < i; j++ {
			total = total * 0.9
		}
		x := total / tb
		fmt.Println("", i+20, " ", tb, " ", x, " Token", strconv.FormatFloat(total, 'f', -1, 64))
		tr = tr + total
	}
	fmt.Println("", strconv.FormatFloat(tr, 'f', -1, 64))

}

func cycleRange(start, end uint64) {
	total := uint64(0)
	for i := start; i < end+1; i++ {
		rewardOne := config.ClacRewardForBlockHeight(i)
		fmt.Println("", rewardOne)
		total += rewardOne
	}
	fmt.Println("", total)
}
