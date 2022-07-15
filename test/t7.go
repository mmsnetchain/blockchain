package main

import (
	"fmt"
)

const (
	year  = 20
	cycle = 40
)

func pow(x float64, n int) float64 {
	if x == 0 {
		return 0
	}
	result := calPow(x, n)
	if n < 0 {
		result = 1 / result
	}
	return result
}
func calPow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}

	result := calPow(x, n>>1)
	result *= result

	if n&1 == 1 {
		result *= x
	}
	return result
}
func main() {
	num := float64(3000000000)

	p := float64(0)
	for i := 1; i <= cycle; i++ {
		p += pow(0.9, i)
	}
	X := num / p
	fmt.Println(X)
	perkuai := 60
	t := year * 365 * 24 * 60 * 60
	kuai := t / perkuai
	allnum := float64(0)
	for i := 1; i <= cycle; i++ {
		bi := X * pow(0.9, i)
		perbi := bi / float64(kuai)
		fmt.Printf("%d %d  %.8f  %.8f \n", i, kuai/cycle, perbi, bi)
		allnum += bi
	}
	fmt.Printf("%.8f", allnum)

}
