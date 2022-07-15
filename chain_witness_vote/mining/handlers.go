package mining

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"mmschainnewaccount/sqlite3_db"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/libp2parea/message_center"
	mc "github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/libp2parea/message_center/flood"
	"github.com/prestonTao/libp2parea/nodeStore"
	pgo_protos "github.com/prestonTao/libp2parea/protos/go_protos"
	"github.com/prestonTao/utils"
	"github.com/shirou/gopsutil/v3/mem"
)

func RegisteMSG() {

	message_center.Register_multicast(config.MSGID_multicast_vote_recv, MulticastVote_recv)
	message_center.Register_multicast(config.MSGID_multicast_blockhead, MulticastBlockHead_recv)
	message_center.Register_neighbor(config.MSGID_heightBlock, FindHeightBlock)
	message_center.Register_neighbor(config.MSGID_heightBlock_recv, FindHeightBlock_recv)
	message_center.Register_neighbor(config.MSGID_getStartBlockHead, GetStartBlockHead)
	message_center.Register_neighbor(config.MSGID_getStartBlockHead_recv, GetStartBlockHead_recv)
	message_center.Register_neighbor(config.MSGID_getBlockHeadVO, GetBlockHeadVO)
	message_center.Register_neighbor(config.MSGID_getBlockHeadVO_recv, GetBlockHeadVO_recv)
	message_center.Register_multicast(config.MSGID_multicast_transaction, MulticastTransaction_recv)
	message_center.Register_neighbor(config.MSGID_getUnconfirmedBlock, GetUnconfirmedBlock)
	message_center.Register_neighbor(config.MSGID_getUnconfirmedBlock_recv, GetUnconfirmedBlock_recv)
	message_center.Register_neighbor(config.MSGID_multicast_return, MulticastReturn_recv)
	message_center.Register_neighbor(config.MSGID_getblockforwitness, GetBlockForWitness)
	message_center.Register_neighbor(config.MSGID_getblockforwitness_recv, GetBlockForWitness_recv)
	message_center.Register_neighbor(config.MSGID_getDBKey_one, GetDBKeyOne)
	message_center.Register_neighbor(config.MSGID_getDBKey_one_recv, GetDBKeyOne_recv)
	message_center.Register_multicast(config.MSGID_multicast_find_witness, MulticastWitness)
	message_center.Register_p2pHE(config.MSGID_multicast_find_witness_recv, MulticastWitness_recv)
	message_center.Register_neighbor(config.MSGID_getBlockLastCurrent, GetBlockLastCurrent)
	message_center.Register_neighbor(config.MSGID_getBlockLastCurrent_recv, GetBlockLastCurrent_recv)

	message_center.Register_neighbor(config.MSGID_multicast_witness_blockhead, MulticastBlockHeadHash)
	message_center.Register_neighbor(config.MSGID_multicast_witness_blockhead_recv, MulticastBlockHeadHash_recv)
	message_center.Register_neighbor(config.MSGID_multicast_witness_blockhead_get, GetMulticastBlockHead)
	message_center.Register_neighbor(config.MSGID_multicast_witness_blockhead_get_recv, GetMulticastBlockHead_recv)

	message_center.Register_neighbor(config.MSGID_uniformity_multicast_witness_blockhead, UniformityMulticastBlockHeadHash)
	message_center.Register_neighbor(config.MSGID_uniformity_multicast_witness_blockhead_recv, UniformityMulticastBlockHeadHash_recv)
	message_center.Register_neighbor(config.MSGID_uniformity_multicast_witness_block_get, UniformityGetMulticastBlockHead)
	message_center.Register_neighbor(config.MSGID_uniformity_multicast_witness_block_get_recv, UniformityGetMulticastBlockHead_recv)
	message_center.Register_neighbor(config.MSGID_uniformity_multicast_witness_block_import, UniformityMulticastBlockImport)
	message_center.Register_neighbor(config.MSGID_uniformity_multicast_witness_block_import_recv, UniformityMulticastBlockImport_recv)

}

type BlockForWitness struct {
	GroupHeight uint64
	Addr        crypto.AddressCoin
}

func (this *BlockForWitness) Proto() (*[]byte, error) {
	bwp := go_protos.BlockForWitness{
		GroupHeight: this.GroupHeight,
		Addr:        this.Addr,
	}
	bs, err := bwp.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, nil

}

func ParseBlockForWitness(bs *[]byte) (*BlockForWitness, error) {
	if bs == nil {
		return nil, nil
	}
	bwp := new(go_protos.BlockForWitness)
	err := proto.Unmarshal(*bs, bwp)
	if err != nil {
		return nil, err
	}
	bw := &BlockForWitness{
		GroupHeight: bwp.GroupHeight,
		Addr:        bwp.Addr,
	}
	return bw, nil
}

func GetBlockForWitness(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	bs := []byte{}
	if message.Body.Content == nil {
		engine.Log.Warn("GetBlockForWitness message.Body.Content is nil")
		message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)
		return
	}

	bfw, err := ParseBlockForWitness(message.Body.Content)
	if err != nil {
		engine.Log.Warn("GetBlockForWitness decoder error: %s", err.Error())
		message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)
		return
	}

	witnessGroup := GetLongChain().WitnessChain.witnessGroup
	for witnessGroup.Height > bfw.GroupHeight && witnessGroup.PreGroup != nil {
		witnessGroup = witnessGroup.PreGroup
	}
	for witnessGroup.Height < bfw.GroupHeight && witnessGroup.NextGroup != nil {
		witnessGroup = witnessGroup.NextGroup
	}
	if witnessGroup.Height != bfw.GroupHeight {
		engine.Log.Warn("GetBlockForWitness not find group height")
		message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)
		return
	}

	for _, one := range witnessGroup.Witness {
		if !bytes.Equal(*one.Addr, bfw.Addr) {
			continue
		}

		if one.Block == nil {
			message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)
		} else {
			bh, tx, err := one.Block.LoadTxs()
			if err != nil {
				message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)
			}
			bhvo := CreateBlockHeadVO(nil, bh, *tx)
			newbs, err := bhvo.Proto()
			if err != nil {
				engine.Log.Warn("GetBlockForWitness bhvo encoding json error: %s", err.Error())
			}
			bs = *newbs
			engine.Log.Info("GetBlockForWitness bhvo encoding Success")
			message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)
		}
		break
	}
	engine.Log.Warn("GetBlockForWitness not find this witness")
	message_center.SendNeighborReplyMsg(message, config.MSGID_getblockforwitness_recv, &bs, msg.Session)

}

func GetBlockForWitness_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(config.CLASS_wallet_getblockforwitness, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func MulticastReturn_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	bs := []byte("ok")

	flood.ResponseWait(config.CLASS_wallet_broadcast_return, utils.Bytes2string(message.Body.Hash), &bs)
}

func MulticastVote_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

}

func MulticastBlockHead_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	chain := GetLongChain()
	if chain == nil {
		return
	}
	if !chain.SyncBlockFinish {
		return
	}

	bhVO, err := ParseBlockHeadVOProto(message.Body.Content)
	if err != nil {

		engine.Log.Warn("Parse block broadcast error: %s", err.Error())
		return
	}

	if !bhVO.Verify(bhVO.StaretBlockHash) {

		return
	}

	exist, err := db.LevelDB.CheckHashExist(bhVO.BH.Hash)
	if err != nil {
		engine.Log.Warn("this block exist error:%s", err.Error())
		return
	}
	if exist {
		engine.Log.Warn("this block exist")
		return
	}

	bhVO.FromBroadcast = true

	bhVO.BH.BuildBlockHash()

	if config.Model == config.Model_light {
		GetLongChain().Temp = bhVO

		_, err := SaveBlockHead(bhVO)
		if err != nil {
			engine.Log.Warn("save block error %s", err.Error())
			return
		}
		return
	}

	go forks.AddBlockHead(bhVO)

}

func FindHeightBlock(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	dataBuf := bytes.NewBuffer([]byte{})
	binary.Write(dataBuf, binary.LittleEndian, forks.GetLongChain().GetStartingBlock())
	binary.Write(dataBuf, binary.LittleEndian, forks.GetLongChain().GetCurrentBlock())
	bs := dataBuf.Bytes()

	message_center.SendNeighborReplyMsg(message, config.MSGID_heightBlock_recv, &bs, msg.Session)

}

func FindHeightBlock_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(mc.CLASS_findHeightBlock, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func GetStartBlockHead(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	bhash, err := db.LevelDB.Find(config.Key_block_start)
	if err != nil {
		return
	}

	chainInfo := ChainInfo{
		StartBlockHash: *bhash,
		HightBlock:     forks.GetHighestBlock(),
	}

	bs, err := chainInfo.Proto()
	if err != nil {
		return
	}

	message_center.SendNeighborReplyMsg(message, config.MSGID_getStartBlockHead_recv, &bs, msg.Session)

}

type ChainInfo struct {
	StartBlockHash []byte
	HightBlock     uint64
}

func (this *ChainInfo) Proto() ([]byte, error) {
	cip := go_protos.ChainInfo{
		StartBlockHash: this.StartBlockHash,
		HightBlock:     this.HightBlock,
	}
	return cip.Marshal()
}

func ParseChainInfo(bs *[]byte) (*ChainInfo, error) {
	if bs == nil {
		return nil, nil
	}
	cip := new(go_protos.ChainInfo)
	err := proto.Unmarshal(*bs, cip)
	if err != nil {
		return nil, err
	}

	ci := ChainInfo{
		StartBlockHash: cip.StartBlockHash,
		HightBlock:     cip.HightBlock,
	}

	return &ci, nil
}

func GetStartBlockHead_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(mc.CLASS_getBlockHead, utils.Bytes2string(message.Body.Hash), message.Body.Content)

}

func GetBlockHeadVO(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	utils.SetTimeToken(config.TIMETOKEN_GetTransaction, config.Wallet_sync_block_interval_time)
	utils.SetTimeToken(config.TIMETOKEN_GetTransactionRelax, config.Wallet_sync_block_interval_time_relax)

	chain := GetLongChain()
	if chain == nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockHeadVO_recv, nil, msg.Session)
		engine.Log.Info("return nil message")
		return
	}

	if !chain.SyncBlockFinish {
		utils.GetTimeToken(config.TIMETOKEN_GetTransaction, true)

	} else if CheckNameStore() {
		utils.GetTimeToken(config.TIMETOKEN_GetTransaction, true)

	} else if chain.WitnessChain.FindWitness(keystore.GetCoinbase().Addr) {
		utils.GetTimeToken(config.TIMETOKEN_GetTransaction, true)

	} else {
		utils.GetTimeToken(config.TIMETOKEN_GetTransactionRelax, true)

	}

	memInfo, _ := mem.VirtualMemory()
	if memInfo.UsedPercent > 92 {
		time.Sleep(time.Second)
	}

	netid := nodeStore.AddressNet([]byte(msg.Session.GetName()))

	if message.Body.Content == nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockHeadVO_recv, nil, msg.Session)
		engine.Log.Info("return nil message")
		return
	}

	bhvo := new(BlockHeadVO)

	bh, err := LoadBlockHeadByHash(message.Body.Content)
	if err != nil {

		message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockHeadVO_recv, nil, msg.Session)
		engine.Log.Info("return nil message")
		return
	} else {

		bhvo.BH = bh
		bhvo.Txs = make([]TxItr, 0, len(bh.Tx))
		for _, one := range bh.Tx {

			txOne, err := LoadTxBase(one)

			if err != nil {
				message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockHeadVO_recv, nil, msg.Session)
				engine.Log.Info("return error message")
				return
			}
			bhvo.Txs = append(bhvo.Txs, txOne)
		}
	}
	if bhvo.BH.Nextblockhash == nil {

		if GetHighestBlock() > bhvo.BH.Height+1 {
			engine.Log.Info("neighbor %s find next block %d hash nil. hight:%d", netid.B58String(), bhvo.BH.Height, GetHighestBlock())
		}

		tempGroup := GetLongChain().WitnessChain.witnessGroup
		for tempGroup != nil {

			if tempGroup.Height < bhvo.BH.GroupHeight {
				break
			}
			if tempGroup.Height > bhvo.BH.GroupHeight {
				tempGroup = tempGroup.PreGroup
				continue
			}
			for _, one := range tempGroup.Witness {
				if one.Block == nil {
					continue
				}
				if one.Block.Height == bhvo.BH.Height {
					if one.Block.NextBlock != nil {
						engine.Log.Info("neighbor %s find next block %d hash nil.", netid.B58String(), bhvo.BH.Height)
						bhvo.BH.Nextblockhash = one.Block.NextBlock.Id
					}
					tempGroup = nil
					break
				}
			}
			break
		}
	}

	bs, err := bhvo.Proto()
	if err != nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockHeadVO_recv, nil, msg.Session)
		engine.Log.Info("return json fialt message")
		return
	}
	err = message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockHeadVO_recv, bs, msg.Session)
	if err != nil {
		engine.Log.Info("returning query transaction or block message Error: %s", err.Error())
	} else {

	}
}

func GetBlockHeadVO_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(mc.CLASS_getTransaction, utils.Bytes2string(message.Body.Hash), message.Body.Content)

}

func MulticastTransaction_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	if config.Model == config.Model_light {
		return
	}

	txbase, err := ParseTxBaseProto(0, message.Body.Content)
	if err != nil {

		engine.Log.Warn("Broadcast transaction format error %s", err.Error())
		return
	}
	engine.Log.Debug("Broadcast transaction received %s", hex.EncodeToString(*txbase.GetHash()))

	txbase.BuildHash()

	chain := GetLongChain()
	if chain == nil {
		engine.Log.Info("chain is nil")
		return
	}
	if !chain.SyncBlockFinish {
		engine.Log.Info("chain is not SyncBlockFinish")
		return
	}

	if txbase.GetGas() < config.Wallet_tx_gas_min {
		engine.Log.Info("gas is little")
		return
	}

	checkTxQueue <- txbase

}

func GetUnconfirmedBlock(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	utils.SetTimeToken(config.TIMETOKEN_GetUnconfirmedBlock, time.Second)
	utils.GetTimeToken(config.TIMETOKEN_GetUnconfirmedBlock, true)

	height := utils.BytesToUint64(*message.Body.Content)

	engine.Log.Info("Get the unconfirmed block, the height of this synchronization block %d", height)

	witnessGroup := GetLongChain().WitnessChain.witnessGroup

	group := witnessGroup
	for {
		group = group.PreGroup
		if group.BlockGroup != nil {
			break
		}
	}
	block := group.BlockGroup.Blocks[0]
	for i := 0; i < config.Mining_group_max*5; i++ {
		if block == nil || block.Height == height {
			break
		}
		if block.Height > height {
			block = block.PreBlock
		}
		if block.Height < height {
			block = block.NextBlock
		}
	}

	var bs []byte
	bhvos := make([]*BlockHeadVO, 0)
	var err error

	for block != nil {
		bh, txs, e := block.LoadTxs()
		if e != nil {
			err = e
			break
		}
		bhvo := &BlockHeadVO{BH: bh, Txs: *txs}
		bhvos = append(bhvos, bhvo)
		block = block.NextBlock
	}
	if err == nil {
		for {
			if witnessGroup == nil {
				break
			}

			for _, one := range witnessGroup.Witness {
				if one.Block == nil {
					continue
				}

				bh, txs, e := one.Block.LoadTxs()
				if e != nil {
					err = e
					break
				}
				bhvo := &BlockHeadVO{BH: bh, Txs: *txs}
				bhvos = append(bhvos, bhvo)

			}
			if err != nil {
				break
			}
			witnessGroup = witnessGroup.NextGroup
		}
	}
	if err == nil {

		rbsp := go_protos.RepeatedBytes{
			Bss: make([][]byte, 0),
		}
		for _, one := range bhvos {
			bsOne, err := one.Proto()
			if err != nil {
				return
			}
			rbsp.Bss = append(rbsp.Bss, *bsOne)
		}
		bsOne, err := rbsp.Marshal()

		if err == nil {
			bs = bsOne
		}
	}

	err = message_center.SendNeighborReplyMsg(message, config.MSGID_getUnconfirmedBlock_recv, &bs, msg.Session)
	if err != nil {
		engine.Log.Info("returning query transaction or block message Error %s", err.Error())
	}
}

func GetUnconfirmedBlock_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(mc.CLASS_getUnconfirmedBlock, utils.Bytes2string(message.Body.Hash), message.Body.Content)

}

func GetDBKeyOne(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	utils.SetTimeToken(config.TIMETOKEN_GetTransaction, config.Wallet_sync_block_interval_time)
	utils.GetTimeToken(config.TIMETOKEN_GetTransaction, true)

	if message.Body.Content == nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getDBKey_one_recv, nil, msg.Session)
		return
	}

	bs, err := db.LevelDB.Find(*message.Body.Content)
	if err != nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getDBKey_one_recv, nil, msg.Session)
		return
	}

	err = message_center.SendNeighborReplyMsg(message, config.MSGID_getDBKey_one_recv, bs, msg.Session)
	if err != nil {
		engine.Log.Info("returning query transaction message Error: %s", err.Error())
	}
}

func GetDBKeyOne_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	flood.ResponseWait(mc.CLASS_getTransaction, utils.Bytes2string(message.Body.Hash), message.Body.Content)

}

func MulticastWitness(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	chain := GetLongChain()
	if chain == nil {
		return
	}

	IsBackup := chain.WitnessChain.FindWitness(keystore.GetCoinbase().Addr)
	if !IsBackup {
		return
	}

}

func MulticastWitness_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {

}

func GetBlockLastCurrent(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	if config.Model == config.Model_light {
		currentHeight := GetLongChain().GetCurrentBlock()
		bhash := LoadBlockHashByHeight(currentHeight)
		if bhash == nil {
			message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockLastCurrent_recv, nil, msg.Session)
			return
		}
		bs, err := db.LevelDB.Find(*bhash)
		if err != nil {
			message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockLastCurrent_recv, nil, msg.Session)
			return
		}
		err = message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockLastCurrent_recv, bs, msg.Session)
		if err != nil {
			engine.Log.Info("returning GetBlockLastCurrent Error %s", err.Error())
		}
		return
	}

	_, block := GetLongChain().GetLastBlock()

	bh, err := block.Load()
	if err != nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockLastCurrent_recv, nil, msg.Session)
		return
	}
	bs, err := bh.Proto()
	if err != nil {
		message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockLastCurrent_recv, nil, msg.Session)
		return
	}
	err = message_center.SendNeighborReplyMsg(message, config.MSGID_getBlockLastCurrent_recv, bs, msg.Session)
	if err != nil {
		engine.Log.Info("returning GetBlockLastCurrent Error %s", err.Error())
	}
}

func GetBlockLastCurrent_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	flood.ResponseWait(mc.CLASS_getBlockLastCurrent, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

var syncHashLock = new(sync.Mutex)
var blockheadHashMap = make(map[string]int)

func MulticastBlockHeadHash(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	if !forks.GetLongChain().SyncBlockFinish {

		message_center.SendNeighborReplyMsg(message, config.MSGID_multicast_witness_blockhead_recv, nil, msg.Session)
		return
	}

	var bhVO *BlockHeadVO
	success := false
	syncHashLock.Lock()
	_, err := new(sqlite3_db.MessageCache).FindByHash(*message.Body.Content)
	if err == nil {

		message_center.SendNeighborReplyMsg(message, config.MSGID_multicast_witness_blockhead_recv, nil, msg.Session)
	} else {

		addrNet := nodeStore.AddressNet(msg.Session.GetName())
		newmsg, _ := message_center.SendNeighborMsg(config.MSGID_multicast_witness_blockhead_get, &addrNet, message.Body.Content)
		bs, err := flood.WaitRequest(config.CLASS_witness_get_blockhead, utils.Bytes2string(newmsg.Body.Hash), int64(4))
		if err != nil {

			engine.Log.Warn("Timeout receiving broadcast reply message %s %s", addrNet.B58String(), hex.EncodeToString(newmsg.Body.Hash))

		} else {

			mmp := new(pgo_protos.MessageMulticast)
			err = proto.Unmarshal(*bs, mmp)
			if err != nil {
				engine.Log.Error("proto unmarshal error %s", err.Error())
			} else {

				bhVOmessage, err := message_center.ParserMessageProto(mmp.Head, mmp.Body, 0)
				if err != nil {
					engine.Log.Error("proto unmarshal error %s", err.Error())
				} else {

					err = bhVOmessage.ParserContentProto()
					if err != nil {
						engine.Log.Error("proto unmarshal error %s", err.Error())
					} else {

						bhVO, err = ParseBlockHeadVOProto(bhVOmessage.Body.Content)
						if err != nil {
							engine.Log.Warn("Parse block broadcast error: %s", err.Error())
						} else {

							message_center.SendNeighborReplyMsg(message, config.MSGID_multicast_witness_blockhead_recv, nil, msg.Session)

							err = new(sqlite3_db.MessageCache).Add(*message.Body.Content, mmp.Head, mmp.Body)
							if err != nil {
								engine.Log.Error(err.Error())
							} else {

								success = true
							}
						}
					}
				}
			}
		}
	}
	syncHashLock.Unlock()

	if !success {

		return
	}

	whiltlistNodes := nodeStore.GetWhiltListNodes()
	err = message_center.BroadcastsAll(1, config.MSGID_multicast_witness_blockhead, whiltlistNodes, nil, nil, message.Body.Content)
	if err != nil {

		return
	}

	if !bhVO.Verify(bhVO.StaretBlockHash) {

		return
	}

	exist, err := db.LevelDB.CheckHashExist(bhVO.BH.Hash)
	if err != nil {
		engine.Log.Warn("this block exist error:%s", err.Error())
		return
	}
	if exist {
		engine.Log.Warn("this block exist")
		return
	}

	bhVO.FromBroadcast = true

	bhVO.BH.BuildBlockHash()
	go forks.AddBlockHead(bhVO)

	go MulticastBlock(*bhVO)

}

func MulticastBlockHeadHash_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	flood.ResponseWait(config.CLASS_wallet_broadcast_return, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func GetMulticastBlockHead(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	messageCache, err := new(sqlite3_db.MessageCache).FindByHash(*message.Body.Content)
	if err != nil {

		engine.Log.Error("find message hash error:%s", err.Error())
		return
	}

	mmp := pgo_protos.MessageMulticast{
		Head: messageCache.Head,
		Body: messageCache.Body,
	}

	content, err := mmp.Marshal()
	if err != nil {

		engine.Log.Error(err.Error())
		return
	}

	message_center.SendNeighborReplyMsg(message, config.MSGID_multicast_witness_blockhead_get_recv, &content, msg.Session)

}

func GetMulticastBlockHead_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	flood.ResponseWait(config.CLASS_witness_get_blockhead, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

var syncUniformityHashLock = new(sync.Mutex)

func UniformityMulticastBlockHeadHash(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	if !forks.GetLongChain().SyncBlockFinish {

		message_center.SendNeighborReplyMsg(message, config.MSGID_uniformity_multicast_witness_blockhead_recv, nil, msg.Session)
		return
	}
	var bhVO *BlockHeadVO
	success := false
	syncUniformityHashLock.Lock()
	_, err := new(sqlite3_db.MessageCache).FindByHash(*message.Body.Content)
	if err == nil {

		message_center.SendNeighborReplyMsg(message, config.MSGID_uniformity_multicast_witness_blockhead_recv, nil, msg.Session)
	} else {

		addrNet := nodeStore.AddressNet(msg.Session.GetName())

		bs, err := message_center.SendNeighborWithReplyMsg(config.MSGID_uniformity_multicast_witness_block_get, &addrNet, message.Body.Content,
			config.CLASS_uniformity_witness_get_blockhead, config.Wallet_sync_block_timeout)

		if err != nil {

			engine.Log.Warn("Timeout receiving broadcast reply message %s", addrNet.B58String())

		} else {

			mmp := new(pgo_protos.MessageMulticast)
			err = proto.Unmarshal(*bs, mmp)
			if err != nil {
				engine.Log.Error("proto unmarshal error %s", err.Error())
			} else {

				bhVOmessage, err := message_center.ParserMessageProto(mmp.Head, mmp.Body, 0)
				if err != nil {
					engine.Log.Error("proto unmarshal error %s", err.Error())
				} else {

					err = bhVOmessage.ParserContentProto()
					if err != nil {
						engine.Log.Error("proto unmarshal error %s", err.Error())
					} else {

						bhVO, err = ParseBlockHeadVOProto(bhVOmessage.Body.Content)
						if err != nil {
							engine.Log.Warn("Parse block broadcast error: %s", err.Error())
						} else {

							message_center.SendNeighborReplyMsg(message, config.MSGID_uniformity_multicast_witness_blockhead_recv, nil, msg.Session)

							err = new(sqlite3_db.MessageCache).Add(*message.Body.Content, mmp.Head, mmp.Body)
							if err != nil {
								engine.Log.Error(err.Error())
							} else {

								success = true
							}
						}
					}
				}
			}
		}
	}
	syncUniformityHashLock.Unlock()

	if !success {

		return
	}

	if !bhVO.Verify(bhVO.StaretBlockHash) {

		return
	}

	exist, err := db.LevelDB.CheckHashExist(bhVO.BH.Hash)
	if err != nil {
		engine.Log.Warn("this block exist error:%s", err.Error())
		return
	}
	if exist {
		engine.Log.Warn("this block exist")
		return
	}

	bhVO.FromBroadcast = true
	bhVO.BH.BuildBlockHash()

	AddBlockToCache(bhVO)

}

func UniformityMulticastBlockHeadHash_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	flood.ResponseWait(config.CLASS_uniformity_witness_multicas_blockhead, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func UniformityGetMulticastBlockHead(c engine.Controller, msg engine.Packet, message *message_center.Message) {

	messageCache, err := new(sqlite3_db.MessageCache).FindByHash(*message.Body.Content)
	if err != nil {

		engine.Log.Error("find message hash error:%s", err.Error())
		return
	}

	mmp := pgo_protos.MessageMulticast{
		Head: messageCache.Head,
		Body: messageCache.Body,
	}

	content, err := mmp.Marshal()
	if err != nil {

		engine.Log.Error(err.Error())
		return
	}

	message_center.SendNeighborReplyMsg(message, config.MSGID_uniformity_multicast_witness_block_get_recv, &content, msg.Session)

}

func UniformityGetMulticastBlockHead_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	flood.ResponseWait(config.CLASS_uniformity_witness_get_blockhead, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}

func UniformityMulticastBlockImport(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	message_center.SendNeighborReplyMsg(message, config.MSGID_uniformity_multicast_witness_block_import_recv, nil, msg.Session)
	if !forks.GetLongChain().SyncBlockFinish {

		return
	}

	ImportBlockByCache(message.Body.Content)

}

func UniformityMulticastBlockImport_recv(c engine.Controller, msg engine.Packet, message *message_center.Message) {
	flood.ResponseWait(config.CLASS_uniformity_witness_multicas_block_import, utils.Bytes2string(message.Body.Hash), message.Body.Content)
}
