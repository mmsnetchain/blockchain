package mining

import (
	"github.com/prestonTao/utils"
	"sync"
)

var bytePool = sync.Pool{
	New: func() interface{} {

		buf := utils.NewBufferByte(0)
		return &buf
	},
}
