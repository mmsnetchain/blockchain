package main

import (
	"fmt"
	"time"
)

func main() {
	plant_10s()
}

func plant_10s() {
	blockTime := int64(10)
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-12-03 10:20:00", time.Local)
	fmt.Println("", err, startTime.Format("2006-01-02 15:04:05"))

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-12-08 00:00:00", time.Local)
	fmt.Println("", err, startTime.Format("2006-01-02 15:04:05"))

	jiange := endTime.Unix() - startTime.Unix()
	fmt.Println(" ", jiange)

	lostHeight := jiange / blockTime

	engHeight := plant_20s()

	fmt.Println(" ", engHeight-lostHeight)

}

func plant_20s() int64 {

	startHeight := int64(381960)
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-12-01 11:12:05", time.Local)
	fmt.Println("", err, startTime.Format("2006-01-02 15:04:05"))

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2020-12-08 00:00:00", time.Local)
	fmt.Println("", err, startTime.Format("2006-01-02 15:04:05"))

	jiange := endTime.Unix() - startTime.Unix()
	fmt.Println(" ", jiange)

	blockTime := int64(20)
	lostHeight := jiange / blockTime

	fmt.Println(" ", startHeight+lostHeight)

	return startHeight + lostHeight
}
