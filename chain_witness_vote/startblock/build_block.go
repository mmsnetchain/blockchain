package startblock

import (
	"bytes"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining"
	"mmschainnewaccount/config"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/utils"
)

const ()

func BuildFirstBlock() (*mining.BlockHeadVO, error) {

	if !config.InitNode || !config.DB_is_null {

		return nil, nil
	}

	db.InitDB(config.DB_path, config.DB_path_temp)

	balanceTotal := uint64(config.Mining_coin_premining)

	txHashes := make([][]byte, 0)
	txs := make([]mining.TxItr, 0)

	reward := BuildReward(balanceTotal)
	txs = append(txs, reward)
	txHashes = append(txHashes, *reward.GetHash())

	depositIn := BuildDepositIn()
	txs = append(txs, depositIn)
	txHashes = append(txHashes, *depositIn.GetHash())

	blockHead1 := mining.BlockHead{
		Height:      config.Mining_block_start_height,
		GroupHeight: config.Mining_group_start_height,
		NTx:         uint64(len(txHashes)),
		Tx:          txHashes,
		Time:        utils.GetNow(),
		Witness:     keystore.GetCoinbase().Addr,
	}
	blockHead1.BuildMerkleRoot()
	blockHead1.BuildSign(keystore.GetCoinbase().Addr)
	blockHead1.BuildBlockHash()

	bhbs, _ := blockHead1.Proto()

	err := db.LevelDB.Save(blockHead1.Hash, bhbs)
	if err != nil {

		return nil, err
	}

	db.LevelDB.Save(config.Key_block_start, &blockHead1.Hash)

	hashExist := false

	for _, one := range txs {

		db.SaveTxToBlockHash(one.GetHash(), &blockHead1.Hash)

		bs, err := one.Proto()

		if err != nil {

			return nil, err
		}

		if one.CheckHashExist() {
			hashExist = true
		}

		err = db.LevelDB.Save(*one.GetHash(), bs)
		if err != nil {
			return nil, err
		}

	}

	db.SaveTxToBlockHash(depositIn.GetHash(), &blockHead1.Hash)

	bs, err := depositIn.Proto()

	if err != nil {
		return nil, err
	}
	err = db.LevelDB.Save(*depositIn.GetHash(), bs)
	if err != nil {
		return nil, err
	}

	bhvo := mining.BlockHeadVO{
		BH:  &blockHead1,
		Txs: txs,
	}

	if hashExist {
		return BuildFirstBlock()
	}

	return &bhvo, nil
}

func BuildReward(balanceTotal uint64) mining.TxItr {

	baseCoinAddr := keystore.GetCoinbase()

	vins := make([]*mining.Vin, 0)
	vin := mining.Vin{
		Puk:  baseCoinAddr.Puk,
		Sign: nil,
	}
	vins = append(vins, &vin)

	vouts := make([]*mining.Vout, 0)
	vouts = append(vouts, &mining.Vout{
		Value:   balanceTotal,
		Address: keystore.GetCoinbase().Addr,
	})
	base := mining.TxBase{
		Type:       config.Wallet_tx_type_mining,
		Vin_total:  1,
		Vin:        vins,
		Vout_total: 1,
		Vout:       vouts,
		LockHeight: 1,
	}
	reward := &mining.Tx_reward{
		TxBase: base,
	}

	for i, one := range reward.Vin {
		for _, key := range keystore.GetAddrAll() {

			puk, ok := keystore.GetPukByAddr(key.Addr)
			if !ok {
				return nil
			}

			if bytes.Equal(puk, one.Puk) {
				_, prk, _, err := keystore.GetKeyByAddr(key.Addr, config.Wallet_keystore_default_pwd)

				if err != nil {
					return nil
				}
				sign := reward.GetSign(&prk, uint64(i))

				reward.Vin[i].Sign = *sign
			}
		}
	}
	reward.BuildHash()
	return reward
}

func BuildDepositIn() mining.TxItr {
	coinbase := keystore.GetCoinbase()

	txin := mining.Tx_deposit_in{
		Puk: coinbase.Puk,
	}

	vouts := make([]*mining.Vout, 0)
	vouts = append(vouts, &mining.Vout{
		Value:   config.Mining_deposit,
		Address: coinbase.Addr,
	})
	depositInBase := mining.TxBase{
		Type: config.Wallet_tx_type_deposit_in,

		Vout_total: uint64(len(vouts)),
		Vout:       vouts,
		LockHeight: 1,

		Payload: []byte("first_witness"),
	}
	txin.TxBase = depositInBase
	txin.BuildHash()
	return &txin
}
