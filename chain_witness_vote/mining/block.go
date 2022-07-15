package mining

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"

	"github.com/gogo/protobuf/proto"
)

const ()

type BlockHead struct {
	Hash              []byte             `json:"H"`
	Height            uint64             `json:"Ht"`
	GroupHeight       uint64             `json:"GH"`
	GroupHeightGrowth uint64             `json:"GHG"`
	Previousblockhash []byte             `json:"Pbh"`
	Nextblockhash     []byte             `json:"Nbh"`
	NTx               uint64             `json:"NTx"`
	MerkleRoot        []byte             `json:"M"`
	Tx                [][]byte           `json:"Tx"`
	Time              int64              `json:"T"`
	Witness           crypto.AddressCoin `json:"W"`
	Sign              []byte             `json:"s"`
}

func (this *BlockHead) BuildMerkleRoot() {
	this.MerkleRoot = utils.BuildMerkleRoot(this.Tx)
}

func (this *BlockHead) Serialize() *[]byte {
	length := 0
	for _, one := range this.Tx {
		length += len(one)
	}
	length += 8 + 8 + 8 + 8 + len(this.Previousblockhash) + len(this.MerkleRoot) + len(this.Witness)
	if this.GroupHeightGrowth != 0 {
		length += 8
	}
	bs := make([]byte, 0, length)

	bs = append(bs, utils.Uint64ToBytes(this.Height)...)
	bs = append(bs, utils.Uint64ToBytes(this.GroupHeight)...)
	if this.GroupHeightGrowth != 0 {
		bs = append(bs, utils.Uint64ToBytes(this.GroupHeightGrowth)...)
	}
	bs = append(bs, this.Previousblockhash...)
	bs = append(bs, utils.Uint64ToBytes(this.NTx)...)
	bs = append(bs, this.MerkleRoot...)
	for _, one := range this.Tx {
		bs = append(bs, one...)
	}
	bs = append(bs, utils.Uint64ToBytes(uint64(this.Time))...)
	bs = append(bs, this.Witness...)
	return &bs

}

func (this *BlockHead) BuildSign(key crypto.AddressCoin) {

	bs := this.Serialize()

	_, prk, _, err := keystore.GetKeyByAddr(key, config.Wallet_keystore_default_pwd)
	if err != nil {
		return
	}
	signBs := keystore.Sign(prk, *bs)

	this.Sign = signBs
}

func (this *BlockHead) CheckBlockHead(puk []byte) bool {
	if this.Height < config.Mining_block_start_height_jump {
		return true
	}

	bs := this.Serialize()
	pkey := ed25519.PublicKey(puk)
	if !ed25519.Verify(pkey, *bs, this.Sign) {
		return false
	}

	old := this.Hash
	this.BuildBlockHash()
	if !bytes.Equal(old, this.Hash) {
		return false
	}
	if this.Height <= 1 {
		return true
	}

	return true

}

func (this *BlockHead) BuildBlockHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}

	buf := bytes.NewBuffer(*this.Serialize())

	bs := buf.Bytes()

	this.Hash = utils.Hash_SHA3_256(bs)
}

func (this *BlockHead) CheckHashExist() (bool, error) {
	return db.LevelDB.CheckHashExist(this.Hash)
}

func (this *BlockHead) Proto() (*[]byte, error) {
	bhp := go_protos.BlockHead{
		Hash:              this.Hash,
		Height:            this.Height,
		GroupHeight:       this.GroupHeight,
		GroupHeightGrowth: this.GroupHeightGrowth,
		Previousblockhash: this.Previousblockhash,
		Nextblockhash:     this.Nextblockhash,
		NTx:               this.NTx,
		MerkleRoot:        this.MerkleRoot,
		Tx:                this.Tx,
		Time:              this.Time,
		Witness:           this.Witness,
		Sign:              this.Sign,
	}
	bs, err := bhp.Marshal()
	return &bs, err
}

func ParseBlockHeadProto(bs *[]byte) (*BlockHead, error) {
	if bs == nil {
		return nil, nil
	}
	bhp := new(go_protos.BlockHead)
	err := proto.Unmarshal(*bs, bhp)
	if err != nil {
		return nil, err
	}
	bh := BlockHead{
		Hash:              bhp.Hash,
		Height:            bhp.Height,
		GroupHeight:       bhp.GroupHeight,
		GroupHeightGrowth: bhp.GroupHeightGrowth,
		Previousblockhash: bhp.Previousblockhash,
		Nextblockhash:     bhp.Nextblockhash,
		NTx:               bhp.NTx,
		MerkleRoot:        bhp.MerkleRoot,
		Tx:                bhp.Tx,
		Time:              bhp.Time,
		Witness:           bhp.Witness,
		Sign:              bhp.Sign,
	}
	return &bh, nil
}

type BlockHeadVO struct {
	FromBroadcast   bool       `json:"-"`
	StaretBlockHash []byte     `json:"sbh"`
	BH              *BlockHead `json:"bh"`
	Txs             []TxItr    `json:"txs"`
}

func (this *BlockHeadVO) Proto() (*[]byte, error) {
	bh := go_protos.BlockHead{
		Hash:              this.BH.Hash,
		Height:            this.BH.Height,
		GroupHeight:       this.BH.GroupHeight,
		GroupHeightGrowth: this.BH.GroupHeightGrowth,
		Previousblockhash: this.BH.Previousblockhash,
		Nextblockhash:     this.BH.Nextblockhash,
		NTx:               this.BH.NTx,
		MerkleRoot:        this.BH.MerkleRoot,
		Tx:                this.BH.Tx,
		Time:              this.BH.Time,
		Witness:           this.BH.Witness,
		Sign:              this.BH.Sign,
	}

	bhat := go_protos.BlockHeadAndTxs{
		StaretBlockHash: this.StaretBlockHash,
		Bh:              &bh,
		TxBs:            make([][]byte, 0),
	}
	for _, one := range this.Txs {
		bs, err := one.Proto()
		if err != nil {
			return nil, err
		}
		bhat.TxBs = append(bhat.TxBs, *bs)
	}
	bs, err := bhat.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, nil

}

const startBlockHashLength = 6

func (this *BlockHeadVO) Verify(sbh []byte) bool {
	if sbh == nil || len(sbh) < startBlockHashLength {
		engine.Log.Info("Illegal block start block hash")
		return false
	}
	if !bytes.Equal(sbh, config.StartBlockHash) {
		engine.Log.Info("Illegal block start block hash %s", hex.EncodeToString(sbh[:startBlockHashLength]))
		return false
	}
	for i, one := range this.Txs {
		one.BuildHash()
		if !bytes.Equal(this.BH.Tx[i], *one.GetHash()) {

			engine.Log.Info("Illegal block")
			return false
		}
	}
	return true

}

func CreateBlockHeadVO(sbh []byte, bh *BlockHead, txs []TxItr) *BlockHeadVO {
	bhvo := BlockHeadVO{
		StaretBlockHash: sbh,
		BH:              bh,
		Txs:             txs,
	}
	return &bhvo
}

type BlockHeadVOParse struct {
	StaretBlockHash []byte        `json:"sbh"`
	BH              *BlockHead    `json:"bh"`
	Txs             []interface{} `json:"txs"`
}

func ParseBlockHeadVOProto(bs *[]byte) (*BlockHeadVO, error) {
	if bs == nil {
		return nil, nil
	}
	bhatp := new(go_protos.BlockHeadAndTxs)
	err := proto.Unmarshal(*bs, bhatp)
	if err != nil {
		return nil, err
	}

	txs := make([]TxItr, 0)
	for i, one := range bhatp.TxBs {
		tx, err := ParseTxBaseProto(ParseTxClass(bhatp.Bh.Tx[i]), &one)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	bh := BlockHead{
		Hash:              bhatp.Bh.Hash,
		Height:            bhatp.Bh.Height,
		GroupHeight:       bhatp.Bh.GroupHeight,
		GroupHeightGrowth: bhatp.Bh.GroupHeightGrowth,
		Previousblockhash: bhatp.Bh.Previousblockhash,
		Nextblockhash:     bhatp.Bh.Nextblockhash,
		NTx:               bhatp.Bh.NTx,
		MerkleRoot:        bhatp.Bh.MerkleRoot,
		Tx:                bhatp.Bh.Tx,
		Time:              bhatp.Bh.Time,
		Witness:           bhatp.Bh.Witness,
		Sign:              bhatp.Bh.Sign,
	}

	return CreateBlockHeadVO(bhatp.StaretBlockHash, &bh, txs), nil

}
