package db

import (
	"os"
	"sync"

	"github.com/prestonTao/utils"
)

var Once_ConnLevelDB sync.Once

var LevelDB *utils.LedisDB
var LevelTempDB *utils.LedisDB

func InitDB(name, tempDBName string) (err error) {
	Once_ConnLevelDB.Do(func() {
		LevelDB, err = utils.CreateLedisDB(name)
		if err != nil {
			return
		}
		os.RemoveAll(tempDBName)
		LevelTempDB, err = utils.CreateLedisDB(tempDBName)
		if err != nil {
			return
		}
		return
	})
	return
}
