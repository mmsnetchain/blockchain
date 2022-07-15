package mining

import (
	"encoding/hex"
	"errors"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"
	"runtime"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	mc "github.com/prestonTao/libp2parea/message_center"
	"github.com/prestonTao/utils"
	"golang.org/x/crypto/ed25519"
)

const ()

type TxItr interface {
	Class() uint64
	BuildHash()
	GetHash() *[]byte

	CheckLockHeight(lockHeight uint64) error

	CheckSign() error
	CheckHashExist() bool

	GetSpend() uint64
	CheckRepeatedTx(txs ...TxItr) bool

	Proto() (*[]byte, error)
	Serialize() *[]byte
	GetVin() *[]*Vin
	GetVout() *[]*Vout
	GetGas() uint64
	GetVoutSignSerialize(voutIndex uint64) *[]byte
	GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte

	GetLockHeight() uint64
	SetSign(index uint64, bs []byte) bool
	SetPayload(bs []byte)
	GetPayload() []byte
	GetVOJSON() interface{}

	CountTxItemsNew(height uint64) *TxItemCountMap
	CountTxHistory(height uint64)
}

type TxBase struct {
	Hash       []byte  `json:"h"`
	Type       uint64  `json:"t"`
	Vin_total  uint64  `json:"vin_t"`
	Vin        []*Vin  `json:"vin"`
	Vout_total uint64  `json:"vot_t"`
	Vout       []*Vout `json:"vot"`
	Gas        uint64  `json:"g"`
	LockHeight uint64  `json:"l_h"`
	Payload    []byte  `json:"p"`
	BlockHash  []byte  `json:"bh"`
}

type TxBaseVO struct {
	Hash       string    `json:"hash"`
	Type       uint64    `json:"type"`
	Vin_total  uint64    `json:"vin_total"`
	Vin        []*VinVO  `json:"vin"`
	Vout_total uint64    `json:"vout_total"`
	Vout       []*VoutVO `json:"vout"`
	Gas        uint64    `json:"gas"`
	LockHeight uint64    `json:"lock_height"`
	Payload    string    `json:"payload"`
	BlockHash  string    `json:"blockhash"`
}

func (this *TxBase) ConversionVO() TxBaseVO {
	vins := make([]*VinVO, 0)
	for _, one := range this.Vin {
		vins = append(vins, one.ConversionVO())
	}

	vouts := make([]*VoutVO, 0)
	for _, one := range this.Vout {
		vouts = append(vouts, one.ConversionVO())
	}

	return TxBaseVO{
		Hash:       hex.EncodeToString(this.Hash),
		Type:       this.Type,
		Vin_total:  this.Vin_total,
		Vin:        vins,
		Vout_total: this.Vout_total,
		Vout:       vouts,
		Gas:        this.Gas,
		LockHeight: this.LockHeight,
		Payload:    string(this.Payload),
		BlockHash:  hex.EncodeToString(this.BlockHash),
	}
}

func (this *TxBase) SetSign(index uint64, bs []byte) bool {
	this.Vin[index].Sign = bs
	return true
}

func (this *TxBase) SetBlockHash(bs []byte) {
	this.BlockHash = bs
}

func (this *TxBase) GetBlockHash() *[]byte {
	return &this.BlockHash
}

func (this *TxBase) GetLockHeight() uint64 {
	return this.LockHeight
}

func (this *TxBase) Serialize() *[]byte {
	length := 0
	var vinSs []*[]byte
	if this.Vin != nil {
		vinSs = make([]*[]byte, 0, len(this.Vin))
		for _, one := range this.Vin {
			bsOne := one.SerializeVin()
			vinSs = append(vinSs, bsOne)
			length += len(*bsOne)
		}
	}
	var voutSs []*[]byte
	if this.Vout != nil {
		voutSs = make([]*[]byte, 0, len(this.Vout))
		for _, one := range this.Vout {
			bsOne := one.Serialize()
			voutSs = append(voutSs, bsOne)
			length += len(*bsOne)
		}
	}
	length += 8 + 8 + 8 + 8 + len(this.Payload)
	bs := make([]byte, 0, length)

	bs = append(bs, utils.Uint64ToBytes(this.Type)...)
	bs = append(bs, utils.Uint64ToBytes(this.Vin_total)...)
	if vinSs != nil {
		for _, one := range vinSs {
			bs = append(bs, *one...)
		}
	}
	bs = append(bs, utils.Uint64ToBytes(this.Vout_total)...)
	if voutSs != nil {
		for _, one := range voutSs {
			bs = append(bs, *one...)
		}
	}
	bs = append(bs, utils.Uint64ToBytes(this.Gas)...)
	bs = append(bs, utils.Uint64ToBytes(this.LockHeight)...)
	bs = append(bs, this.Payload...)
	return &bs

}

func (this *TxBase) GetVin() *[]*Vin {
	return &this.Vin
}
func (this *TxBase) GetVout() *[]*Vout {
	return &this.Vout
}

func (this *TxBase) GetGas() uint64 {
	return this.Gas
}

func (this *TxBase) GetHash() *[]byte {
	return &this.Hash
}

func (this *TxBase) Class() uint64 {
	return this.Type
}

func (this *TxBase) SetPayload(bs []byte) {
	this.Payload = bs
}

func (this *TxBase) GetPayload() []byte {
	return this.Payload
}

func (this *TxBase) GetVoutSignSerialize(voutIndex uint64) *[]byte {
	if voutIndex > uint64(len(this.Vout)) {
		return nil
	}
	vout := this.Vout[voutIndex]
	voutBs := vout.Serialize()
	bs := make([]byte, 0, len(*voutBs)+8)
	bs = append(bs, utils.Uint64ToBytes(voutIndex)...)
	bs = append(bs, *voutBs...)
	return &bs

}

func (this *TxBase) GetSignSerialize(voutBs *[]byte, vinIndex uint64) *[]byte {
	if vinIndex > uint64(len(this.Vin)) {
		return nil
	}

	voutBssLenght := 0
	voutBss := make([]*[]byte, 0, len(this.Vout))
	for _, one := range this.Vout {
		voutBsOne := one.Serialize()
		voutBss = append(voutBss, voutBsOne)
		voutBssLenght += len(*voutBsOne)
	}

	var bs []byte
	if voutBs == nil {
		bs = make([]byte, 0, 8+8+8+8+voutBssLenght+8+len(this.Payload))
	} else {
		bs = make([]byte, 0, len(*voutBs)+8+8+8+8+voutBssLenght+8+len(this.Payload))
		bs = append(bs, *voutBs...)
	}

	bs = append(bs, utils.Uint64ToBytes(this.Type)...)
	bs = append(bs, utils.Uint64ToBytes(this.Vin_total)...)
	bs = append(bs, utils.Uint64ToBytes(vinIndex)...)
	bs = append(bs, utils.Uint64ToBytes(this.Vout_total)...)
	for _, one := range voutBss {
		bs = append(bs, *one...)
	}
	bs = append(bs, utils.Uint64ToBytes(this.Gas)...)
	bs = append(bs, utils.Uint64ToBytes(this.LockHeight)...)
	bs = append(bs, this.Payload...)
	return &bs

}

func (this *TxBase) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	sign := keystore.Sign(*key, *signDst)

	return &sign
}

func (this *TxBase) GetWaitSign(vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	return signDst
}

func (this *TxBase) CheckLockHeight(lockHeight uint64) error {

	if lockHeight < config.Mining_block_start_height+config.Mining_block_start_height_jump {
		return nil
	}

	if this.GetLockHeight() < lockHeight {

		engine.Log.Warn("Failed to compare lock block height: LockHeight=%d %d %s", this.GetLockHeight(), lockHeight, hex.EncodeToString(*this.GetHash()))
		return config.ERROR_tx_lockheight
	}
	return nil
}

func (this *TxBase) CheckBase() error {
	if len(this.Vin) > 1 {
		return config.ERROR_pay_vin_too_much
	}

	for i, one := range this.Vout {
		if i != 0 && one.Value <= 0 {
			return config.ERROR_amount_zero
		}
	}

	for i, one := range this.Vin {

		sign := this.GetSignSerialize(nil, uint64(i))

		puk := ed25519.PublicKey(one.Puk)
		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*sign))
		}

		if !ed25519.Verify(puk, *sign, one.Sign) {

			return config.ERROR_sign_fail
		}

	}

	return nil
}

func (this *TxBase) CheckHashExist() bool {
	ok, _ := db.LevelDB.CheckHashExist(this.Hash)
	return ok
}

func (this *TxBase) Proto() (*[]byte, error) {
	vins := make([]*go_protos.Vin, 0)
	for _, one := range this.Vin {
		vinOne := &go_protos.Vin{

			Puk:  one.Puk,
			Sign: one.Sign,
		}

		vinOne.Nonce = one.Nonce.Bytes()

		vins = append(vins, vinOne)
	}
	vouts := make([]*go_protos.Vout, 0)
	for _, one := range this.Vout {
		vouts = append(vouts, &go_protos.Vout{
			Value:        one.Value,
			Address:      one.Address,
			FrozenHeight: one.FrozenHeight,
		})
	}
	txBase := go_protos.TxBase{
		Hash:       this.Hash,
		Type:       this.Type,
		VinTotal:   this.Vin_total,
		Vin:        vins,
		VoutTotal:  this.Vout_total,
		Vout:       vouts,
		Gas:        this.Gas,
		LockHeight: this.LockHeight,
		Payload:    this.Payload,
		BlockHash:  this.BlockHash,
	}

	txPay := go_protos.TxPay{
		TxBase: &txBase,
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *TxBase) CountTxHistory(height uint64) {

	hiOut := HistoryItem{
		IsIn:    false,
		Type:    this.Class(),
		InAddr:  make([]*crypto.AddressCoin, 0),
		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
	}

	hiIn := HistoryItem{
		IsIn:    true,
		Type:    this.Class(),
		InAddr:  make([]*crypto.AddressCoin, 0),
		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
	}

	if this.Class() == config.Wallet_tx_type_pay {
		hiOut.Payload = this.Payload
		hiIn.Payload = this.Payload
	}

	addrCoin := make(map[string]bool)
	for _, vin := range this.Vin {
		addrInfo, isSelf := keystore.FindPuk(vin.Puk)
		hiIn.InAddr = append(hiIn.InAddr, &addrInfo.Addr)
		if !isSelf {
			continue
		}
		if _, ok := addrCoin[utils.Bytes2string(addrInfo.Addr)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(addrInfo.Addr)] = false
		}
		hiOut.InAddr = append(hiOut.InAddr, &addrInfo.Addr)
	}

	addrCoin = make(map[string]bool)
	for _, vout := range this.Vout {
		hiOut.OutAddr = append(hiOut.OutAddr, &vout.Address)
		hiOut.Value += vout.Value
		_, ok := keystore.FindAddress(vout.Address)
		if !ok {
			continue
		}
		hiIn.Value += vout.Value
		if _, ok := addrCoin[utils.Bytes2string(vout.Address)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(vout.Address)] = false
		}
		hiIn.OutAddr = append(hiIn.OutAddr, &vout.Address)
	}
	if len(hiOut.InAddr) > 0 {
		balanceHistoryManager.Add(hiOut)
	}
	if len(hiIn.OutAddr) > 0 {
		balanceHistoryManager.Add(hiIn)
	}
}

func CleanZeroVouts(vs *[]*Vout) []*Vout {
	vouts := make([]*Vout, 0)
	for _, one := range *vs {
		if one.Value == 0 {
			continue
		}
		vouts = append(vouts, one)
	}
	return vouts
}

func MergeVouts(vs *[]*Vout) []*Vout {
	voutMap := make(map[string]*Vout)
	for i, one := range *vs {
		if one.Value == 0 {
			continue
		}
		v, ok := voutMap[utils.Bytes2string(one.Address)+strconv.Itoa(int(one.FrozenHeight))]
		if ok {
			v.Value = v.Value + one.Value
			continue
		}
		voutMap[utils.Bytes2string(one.Address)+strconv.Itoa(int(one.FrozenHeight))] = (*vs)[i]
	}
	vouts := make([]*Vout, 0)
	for _, v := range voutMap {
		vouts = append(vouts, v)
	}
	return vouts
}

func (this *TxBase) MergeVout() {
	this.Vout = MergeVouts(this.GetVout())
	this.Vout_total = uint64(len(this.Vout))
}

func (this *TxBase) CleanZeroVout() {
	this.Vout = CleanZeroVouts(this.GetVout())
	this.Vout_total = uint64(len(this.Vout))
}

func ParseTxBaseProto(txtype uint64, bs *[]byte) (TxItr, error) {
	if bs == nil {
		return nil, nil
	}

	if txtype == 0 {

		txPay := go_protos.TxPay{}
		err := proto.Unmarshal(*bs, &txPay)
		if err != nil {
			return nil, err
		}
		txtype = txPay.TxBase.Type

	}

	var tx interface{}

	switch txtype {
	case config.Wallet_tx_type_mining:

		txProto := new(go_protos.TxReward)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}

		if txProto.TxBase.Type != config.Wallet_tx_type_mining {
			return nil, errors.New("tx type error")
		}

		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {

			vins = append(vins, &Vin{

				Puk:  one.Puk,
				Sign: one.Sign,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_reward{
			TxBase: txBase,
		}

	case config.Wallet_tx_type_deposit_in:

		txProto := new(go_protos.TxDepositIn)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}
		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {
			nonce := new(big.Int).SetBytes(one.Nonce)
			vins = append(vins, &Vin{

				Puk:   one.Puk,
				Sign:  one.Sign,
				Nonce: *nonce,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_deposit_in{
			TxBase: txBase,
			Puk:    txProto.Puk,
		}

	case config.Wallet_tx_type_deposit_out:

		txProto := new(go_protos.TxDepositOut)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}
		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {
			nonce := new(big.Int).SetBytes(one.Nonce)
			vins = append(vins, &Vin{

				Puk:   one.Puk,
				Sign:  one.Sign,
				Nonce: *nonce,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_deposit_out{
			TxBase: txBase,
		}
	case config.Wallet_tx_type_pay:
		txProto := new(go_protos.TxPay)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}
		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {
			nonce := new(big.Int).SetBytes(one.Nonce)
			vins = append(vins, &Vin{

				Puk:   one.Puk,
				Sign:  one.Sign,
				Nonce: *nonce,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_Pay{
			TxBase: txBase,
		}
	case config.Wallet_tx_type_vote_in:
		txProto := new(go_protos.TxVoteIn)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}
		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {
			nonce := new(big.Int).SetBytes(one.Nonce)
			vins = append(vins, &Vin{

				Puk:   one.Puk,
				Sign:  one.Sign,
				Nonce: *nonce,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_vote_in{
			TxBase:   txBase,
			Vote:     txProto.Vote,
			VoteType: uint16(txProto.VoteType),
		}

	case config.Wallet_tx_type_vote_out:
		txProto := new(go_protos.TxVoteOut)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}
		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {
			nonce := new(big.Int).SetBytes(one.Nonce)
			vins = append(vins, &Vin{

				Puk:   one.Puk,
				Sign:  one.Sign,
				Nonce: *nonce,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_vote_out{
			TxBase:   txBase,
			Vote:     txProto.Vote,
			VoteType: uint16(txProto.VoteType),
		}
	case config.Wallet_tx_type_voting_reward:
		txProto := new(go_protos.TxVoteReward)
		err := proto.Unmarshal(*bs, txProto)
		if err != nil {
			return nil, err
		}
		vins := make([]*Vin, 0, len(txProto.TxBase.Vin))
		for _, one := range txProto.TxBase.Vin {
			nonce := new(big.Int).SetBytes(one.Nonce)
			vins = append(vins, &Vin{

				Puk:   one.Puk,
				Sign:  one.Sign,
				Nonce: *nonce,
			})
		}
		vouts := make([]*Vout, 0, len(txProto.TxBase.Vout))
		for _, one := range txProto.TxBase.Vout {
			vouts = append(vouts, &Vout{
				Value:        one.Value,
				Address:      one.Address,
				FrozenHeight: one.FrozenHeight,
			})
		}
		txBase := TxBase{}
		txBase.Hash = txProto.TxBase.Hash
		txBase.Type = txProto.TxBase.Type
		txBase.Vin_total = txProto.TxBase.VinTotal
		txBase.Vin = vins
		txBase.Vout_total = txProto.TxBase.VoutTotal
		txBase.Vout = vouts
		txBase.Gas = txProto.TxBase.Gas
		txBase.LockHeight = txProto.TxBase.LockHeight
		txBase.Payload = txProto.TxBase.Payload
		txBase.BlockHash = txProto.TxBase.BlockHash
		tx = &Tx_Vote_Reward{
			TxBase:      txBase,
			StartHeight: txProto.StartHeight,
			EndHeight:   txProto.EndHeight,
		}

	default:
		tx = GetNewTransaction(txtype, bs)
		if tx == nil {

			return nil, errors.New("Unknown transaction type")
		}
	}

	return tx.(TxItr), nil
}

type Vin struct {
	Puk   []byte  `json:"puk"`
	Nonce big.Int `json:"n"`

	PukIsSelf int                `json:"-"`
	PukToAddr crypto.AddressCoin `json:"-"`

	Sign []byte `json:"sign"`
}

type VinVO struct {
	Puk   string `json:"puk"`
	Nonce string `json:"n"`
	Sign  string `json:"sign"`
}

func (this *Vin) ConversionVO() *VinVO {
	vinvo := &VinVO{

		Puk: hex.EncodeToString(this.Puk),

		Sign: hex.EncodeToString(this.Sign),
	}

	vinvo.Nonce = this.Nonce.Text(10)

	return vinvo
}

func (this *Vin) CheckIsSelf() bool {

	if this.PukIsSelf == 0 {

		_, ok := keystore.FindPuk(this.Puk)
		if ok {
			this.PukIsSelf = 2

		} else {
			this.PukIsSelf = 1
		}
	}

	if this.PukIsSelf == 1 {
		return false
	} else {
		return true
	}

}

func (this *Vin) GetPukToAddr() *crypto.AddressCoin {
	if this.PukToAddr == nil || len(this.PukToAddr) <= 0 {
		addr := crypto.BuildAddr(config.AddrPre, this.Puk)
		this.PukToAddr = addr
	}
	return &this.PukToAddr
}

func (this *Vin) SerializeVin() *[]byte {
	bs := make([]byte, 0, len(this.Puk)+len(this.Sign))

	bs = append(bs, this.Puk...)
	bs = append(bs, this.Sign...)
	return &bs

}

func (this *Vin) ValidateAddr() (*crypto.AddressCoin, bool) {
	addr := crypto.BuildAddr(config.AddrPre, this.Puk)

	_, ok := keystore.FindAddress(addr)
	if !ok {
		return &addr, false
	}

	return &addr, true
}

type Vout struct {
	Value        uint64             `json:"value"`
	Address      crypto.AddressCoin `json:"address"`
	FrozenHeight uint64             `json:"frozen_height"`
	AddrIsSelf   int                `json:"-"`
	AddrStr      string             `json:"-"`
}

type VoutVO struct {
	Value        uint64 `json:"value"`
	Address      string `json:"address"`
	FrozenHeight uint64 `json:"frozen_height"`
}

func (this *Vout) ConversionVO() *VoutVO {
	return &VoutVO{
		Value:        this.Value,
		Address:      this.Address.B58String(),
		FrozenHeight: this.FrozenHeight,
	}
}

func (this *Vout) CheckIsSelf() bool {

	if this.AddrIsSelf == 0 {
		_, ok := keystore.FindAddress(this.Address)
		if ok {
			this.AddrIsSelf = 2
		} else {
			this.AddrIsSelf = 1
		}
	}
	if this.AddrIsSelf == 1 {
		return false
	} else {
		return true
	}
}

func (this *Vout) GetAddrStr() string {
	if this.AddrStr == "" {
		this.AddrStr = this.Address.B58String()
	}
	return this.AddrStr
}

func (this *Vout) Serialize() *[]byte {
	bs := make([]byte, 0, len(this.Address)+8+8)
	bs = append(bs, utils.Uint64ToBytes(this.Value)...)
	bs = append(bs, this.Address...)
	bs = append(bs, utils.Uint64ToBytes(this.FrozenHeight)...)

	return &bs
}

func MulticastTx(txItr TxItr) {

	utils.Go(func() {
		goroutineId := utils.GetRandomDomain() + utils.TimeFormatToNanosecondStr()
		_, file, line, _ := runtime.Caller(0)
		engine.AddRuntime(file, line, goroutineId)
		defer engine.DelRuntime(file, line, goroutineId)

		bs, err := txItr.Proto()
		if err != nil {

			return
		}
		mc.SendMulticastMsg(config.MSGID_multicast_transaction, bs)
	})

}

func ParseTxClass(txid []byte) uint64 {

	classBs := txid[:8]

	return utils.BytesToUint64(classBs)
}
