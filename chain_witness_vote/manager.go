package chain_witness_vote

import (
	"errors"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/light"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/chain_witness_vote/startblock"
	"mmschainnewaccount/config"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/hyahm/golog"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
)

func Register() error {
	golog.InitLogger("logs/randeHash.txt", 0, true)
	golog.Infof("start %s", "log")

	engine.Log.Info("CPUNUM :%d", config.CPUNUM)

	go func() {
		for {
			engine.Log.Info("NumGoroutine:%d", runtime.NumGoroutine())
			time.Sleep(time.Minute)

		}
	}()

	err := utils.StartSystemTime()
	if err != nil {
		return err
	}

	config.ParseInitFlag()

	if config.InitNode {
		os.RemoveAll(filepath.Join(config.Wallet_path))
	}

	utils.CheckCreateDir(filepath.Join(config.Wallet_path))

	err = db.InitDB(config.DB_path, config.DB_path_temp)
	if err != nil {
		panic(err)
	}

	bhvo := mining.LoadStartBlock()
	if bhvo == nil {
		config.DB_is_null = true
	}

	mining.RegisteMSG()

	if config.InitNode {
		bhvo, err = startblock.BuildFirstBlock()
		if err != nil {
			return err
		}
		engine.Log.Info("create initiation block build chain")

		config.StartBlockHash = bhvo.BH.Hash

		mining.BuildFirstChain(bhvo)
		mining.SetHighestBlock(config.Mining_block_start_height)
		mining.GetLongChain().SyncBlockFinish = true

		mining.GetLongChain().WitnessChain.BuildMiningTime()
		return nil
	}

	if config.LoadNode {

		bhvo := mining.LoadStartBlock()
		if bhvo == nil {
			return errors.New("")
		}
		engine.Log.Info("load db initiation block build chain")
		config.StartBlockHash = bhvo.BH.Hash

		mining.BuildFirstChain(bhvo)

		mining.SetHighestBlock(db.GetHighstBlock())

		mining.FindBlockHeight()

		if err := mining.GetLongChain().LoadBlockChain(); err != nil {
			return err
		}
		mining.FinishFirstLoadBlockChain()

		return nil
	}

	if config.Model == config.Model_light {
		StartModelLight()
		return nil
	}

	bhvo = mining.LoadStartBlock()
	if bhvo == nil {
		engine.Log.Info("neighbor initiation block build chain")

		err = mining.GetFirstBlock()
		if err != nil {
			engine.Log.Error("get first block error: %s", err.Error())
			panic(err.Error())
		}

		mining.FindBlockHeight()
	} else {
		engine.Log.Info("load db initiation block build chain")
		config.StartBlockHash = bhvo.BH.Hash

		mining.BuildFirstChain(bhvo)
		mining.FindBlockHeight()
	}
	if err := mining.GetLongChain().FirstDownloadBlock(); err != nil {
		return err
	}

	engine.Log.Info("build chain success")

	mining.GetLongChain().NoticeLoadBlockForDB()
	return nil

}

func StartModelLight() {

	for {
		bhvo := mining.LoadStartBlock()
		if bhvo == nil {
			engine.Log.Info("neighbor initiation block build chain")

			err := mining.GetFirstBlock()
			if err != nil {
				engine.Log.Error("get first block error: %s", err.Error())
				panic(err.Error())
			}

			mining.FindBlockHeight()
		} else {
			engine.Log.Info("load db initiation block build chain")
			config.StartBlockHash = bhvo.BH.Hash

			mining.BuildFirstChain(bhvo)
			mining.FindBlockHeight()
		}
		light.SyncBlock()
	}
}
