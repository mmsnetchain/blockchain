package config

import (
	"sync/atomic"
)

const (
	Chunk_size = 10 * 1024 * 1024
	Spacenum   = 1024 * Chunk_size
)

var (
	spaceTotal     = uint64(0)
	spaceTotalAddr = uint64(0)
)

func GetSpaceTotal() uint64 {
	return atomic.LoadUint64(&spaceTotal)
}

func SetSpaceTotal(st uint64) {
	atomic.StoreUint64(&spaceTotal, st)
}

func GetSpaceTotalAddr() uint64 {
	return atomic.LoadUint64(&spaceTotalAddr)
}

func SetSpaceTotalAddr(sta uint64) {
	atomic.StoreUint64(&spaceTotalAddr, sta)
}
