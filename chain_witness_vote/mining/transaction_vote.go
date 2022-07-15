package mining

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/chain_witness_vote/mining/name"
	"mmschainnewaccount/config"
	"mmschainnewaccount/protos/go_protos"

	"github.com/prestonTao/keystore"
	"github.com/prestonTao/keystore/crypto"
	"github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"golang.org/x/crypto/ed25519"
)

const (
	type_addr = uint16(1)
	type_name = uint16(2)

	VOTE_TYPE_community = 1
	VOTE_TYPE_vote      = 2
	VOTE_TYPE_light     = 3
)

type VoteAddress []byte

func (this *VoteAddress) GetAddress() crypto.AddressCoin {
	if len(*this) <= 2 {
		return nil
	}
	addrType := utils.BytesToUint16((*this)[:2])
	switch addrType {
	case type_addr:
		return crypto.AddressCoin((*this)[2:])
	case type_name:
		index := utils.BytesToUint16((*this)[2:4])
		nameStr := string((*this)[4:])
		nameinfo := name.FindNameToNet(nameStr)
		if nameinfo != nil && len(nameinfo.AddrCoins) >= int(index) {
			addrCoin := nameinfo.AddrCoins[index]
			return addrCoin
		}
	}
	return nil
}

func (this *VoteAddress) B58String() string {
	addr := this.GetAddress()
	if addr == nil {
		return ""
	}
	return addr.B58String()
}

func NewVoteAddressByName(name string, index uint16) VoteAddress {
	if name == "" {
		return nil
	}

	buf := bytes.NewBuffer(utils.Uint16ToBytes(type_name))

	buf.Write(utils.Uint16ToBytes(index))

	buf.WriteString(name)

	return VoteAddress(buf.Bytes())
}

func NewVoteAddressByAddr(addr crypto.AddressCoin) VoteAddress {
	if addr == nil {
		return nil
	}

	buf := bytes.NewBuffer(utils.Uint16ToBytes(type_addr))

	buf.Write(addr)

	return VoteAddress(buf.Bytes())
}

type Tx_vote_in struct {
	TxBase
	Vote     crypto.AddressCoin `json:"v"`
	VoteType uint16             `json:"vt"`
}

type Tx_vote_in_VO struct {
	TxBaseVO
	Vote     string `json:"vote"`
	VoteType uint16 `json:"votetype"`
}

func (this *Tx_vote_in) GetVOJSON() interface{} {
	return Tx_vote_in_VO{
		TxBaseVO: this.TxBase.ConversionVO(),
		Vote:     this.Vote.B58String(),
		VoteType: this.VoteType,
	}
}

func (this *Tx_vote_in) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_vote_in)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_vote_in) Proto() (*[]byte, error) {
	vins := make([]*go_protos.Vin, 0)
	for _, one := range this.Vin {
		vins = append(vins, &go_protos.Vin{

			Puk:   one.Puk,
			Sign:  one.Sign,
			Nonce: one.Nonce.Bytes(),
		})
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

	txPay := go_protos.TxVoteIn{
		TxBase:   &txBase,
		Vote:     this.Vote,
		VoteType: uint32(this.VoteType),
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *Tx_vote_in) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)

	buf.Write(this.Vote)
	buf.Write(utils.Uint16ToBytes(this.VoteType))
	*bs = buf.Bytes()
	return bs
}

func (this *Tx_vote_in) GetWaitSign(vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Vote...)
	*signDst = append(*signDst, utils.Uint16ToBytes(this.VoteType)...)

	return signDst

}

func (this *Tx_vote_in) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Vote...)
	*signDst = append(*signDst, utils.Uint16ToBytes(this.VoteType)...)

	sign := keystore.Sign(*key, *signDst)

	return &sign

}

func (this *Tx_vote_in) CheckSign() error {

	if this.Vin == nil || len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}
	if this.Vout_total != 1 {
		return config.ERROR_pay_vout_too_much
	}
	if this.Vout[0].Value <= 0 {
		return config.ERROR_amount_zero
	}

	switch this.VoteType {
	case VOTE_TYPE_community:

	case VOTE_TYPE_vote:
	case VOTE_TYPE_light:
		if this.Vote != nil && len(this.Vote) > 0 {
			return config.ERROR_deposit_light_vote
		}
	}

	for i, one := range this.Vin {

		signDst := this.GetSignSerialize(nil, uint64(i))

		*signDst = append(*signDst, this.Vote...)
		*signDst = append(*signDst, utils.Uint16ToBytes(this.VoteType)...)

		puk := ed25519.PublicKey(one.Puk)

		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*signDst))
		}
		if !ed25519.Verify(puk, *signDst, one.Sign) {
			return config.ERROR_sign_fail
		}

	}

	return nil
}

func (this *Tx_vote_in) GetWitness() *crypto.AddressCoin {
	witness := crypto.BuildAddr(config.AddrPre, this.Vin[0].Puk)

	return &witness
}

func (this *Tx_vote_in) SetVoteAddr(addr crypto.AddressCoin) {

	bs, err := this.Proto()
	if err != nil {
		return
	}

	db.LevelDB.Save(*this.GetHash(), bs)
}

func (this *Tx_vote_in) GetSpend() uint64 {
	return this.Vout[0].Value + this.Gas
}

func (this *Tx_vote_in) CheckRepeatedTx(txs ...TxItr) bool {

	addrSelf := this.Vout[0].Address

	switch this.VoteType {
	case VOTE_TYPE_community:

		if GetDepositCommunityAddr(&addrSelf) > 0 {
			return false
		}

		if GetDepositLightAddr(&addrSelf) > 0 {
			return false
		}

		if GetDepositWitnessAddr(&addrSelf) > 0 {
			return false
		}

	case VOTE_TYPE_vote:

		if GetDepositCommunityAddr(&addrSelf) > 0 {
			engine.Log.Info("111111111")
			return false
		}

		if GetDepositWitnessAddr(&addrSelf) > 0 {
			engine.Log.Info("111111111")
			return false
		}

		communityAddr := GetVoteAddr(&addrSelf)
		if communityAddr != nil && len(*communityAddr) > 0 && !bytes.Equal(this.Vote, *communityAddr) {
			engine.Log.Info("111111111:%s", hex.EncodeToString(*communityAddr))
			return false
		}

	case VOTE_TYPE_light:

		if GetDepositCommunityAddr(&addrSelf) > 0 {
			return false
		}

		if GetDepositLightAddr(&addrSelf) > 0 {
			return false
		}

		if GetDepositWitnessAddr(&addrSelf) > 0 {
			return false
		}

	}

	for _, one := range txs {
		if one.Class() != config.Wallet_tx_type_vote_in {
			continue
		}
		addr := (*one.GetVout())[0].Address

		votein := one.(*Tx_vote_in)
		rule := votein.VoteType

		isSelf := bytes.Equal(addrSelf, addr)
		if !isSelf {
			continue
		}

		switch this.VoteType {
		case VOTE_TYPE_community:
			return false
		case VOTE_TYPE_vote:
			if rule == VOTE_TYPE_community {
				return false
			}
			if rule == VOTE_TYPE_vote {
				if !bytes.Equal(this.Vote, votein.Vote) {
					return false
				}
			}
		case VOTE_TYPE_light:
			return false
		}
	}
	return true
}

func (this *Tx_vote_in) CountTxItemsNew(height uint64) *TxItemCountMap {
	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vin)),
		Nonce:    make(map[string]big.Int),
	}

	totalValue := this.Gas + (*this.GetVout())[0].Value

	from := this.Vin[0].GetPukToAddr()
	itemCount.Nonce[utils.Bytes2string(*from)] = this.Vin[0].Nonce

	frozenMap := make(map[uint64]int64, 0)
	frozenMap[0] = (0 - int64(totalValue))
	itemCount.AddItems[utils.Bytes2string(*from)] = &frozenMap

	return &itemCount
}

func (this *Tx_vote_in) CountTxHistory(height uint64) {

	hiOut := HistoryItem{
		IsIn:    false,
		Type:    this.Class(),
		InAddr:  make([]*crypto.AddressCoin, 0),
		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
	}

	addrCoin := make(map[string]bool)
	for _, vin := range this.Vin {
		addrInfo, isSelf := keystore.FindPuk(vin.Puk)

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
	for voutIndex, vout := range this.Vout {
		if voutIndex != 0 {
			continue
		}
		hiOut.OutAddr = append(hiOut.OutAddr, &vout.Address)
		hiOut.Value += vout.Value
		_, ok := keystore.FindAddress(vout.Address)
		if !ok {
			continue
		}

		if _, ok := addrCoin[utils.Bytes2string(vout.Address)]; ok {
			continue
		} else {
			addrCoin[utils.Bytes2string(vout.Address)] = false
		}

	}
	if len(hiOut.InAddr) > 0 {
		balanceHistoryManager.Add(hiOut)
	}

}

func CreateTxVoteIn(voteType uint16, witnessAddr crypto.AddressCoin, addr crypto.AddressCoin, amount, gas uint64, pwd, payload string) (*Tx_vote_in, error) {

	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()

	vins := make([]*Vin, 0)

	total, item := chain.Balance.BuildPayVinNew(nil, amount+gas)
	if total < amount+gas {

		return nil, config.ERROR_not_enough
	}

	puk, ok := keystore.GetPukByAddr(*item.Addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := Vin{

		Puk: puk,

		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}
	vins = append(vins, &vin)

	var dstAddr crypto.AddressCoin = addr
	if addr == nil {

		dstAddr = keystore.GetAddr()[0].Addr
	}

	vouts := make([]*Vout, 0)

	vout := Vout{
		Value:   amount,
		Address: dstAddr,
	}
	vouts = append(vouts, &vout)

	var txin *Tx_vote_in
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_vote_in,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
			Payload:    []byte(payload),
		}

		txin = &Tx_vote_in{
			TxBase:   base,
			Vote:     witnessAddr,
			VoteType: voteType,
		}

		for i, one := range txin.Vin {
			for _, key := range keystore.GetAddr() {

				puk, ok := keystore.GetPukByAddr(key.Addr)
				if !ok {

					return nil, config.ERROR_get_sign_data_fail
				}

				if bytes.Equal(puk, one.Puk) {
					_, prk, _, err := keystore.GetKeyByAddr(key.Addr, pwd)

					if err != nil {

						return nil, err
					}
					sign := txin.GetSign(&prk, uint64(i))

					txin.Vin[i].Sign = *sign
				}
			}
		}

		txin.BuildHash()
		if txin.CheckHashExist() {
			txin = nil
			continue
		} else {
			break
		}
	}
	chain.Balance.AddLockTx(txin)
	return txin, nil
}

type Tx_vote_out struct {
	TxBase
	Vote     crypto.AddressCoin `json:"v"`
	VoteType uint16             `json:"vt"`
}

type Tx_vote_out_VO struct {
	TxBaseVO
	Vote     string `json:"vote"`
	VoteType uint16 `json:"votetype"`
}

func (this *Tx_vote_out) GetVOJSON() interface{} {
	return Tx_vote_out_VO{
		TxBaseVO: this.TxBase.ConversionVO(),
		Vote:     this.Vote.B58String(),
		VoteType: this.VoteType,
	}
}

func (this *Tx_vote_out) BuildHash() {
	if this.Hash != nil && len(this.Hash) > 0 {
		return
	}
	bs := this.Serialize()

	id := make([]byte, 8)
	binary.PutUvarint(id, config.Wallet_tx_type_vote_out)

	this.Hash = append(id, utils.Hash_SHA3_256(*bs)...)
}

func (this *Tx_vote_out) Proto() (*[]byte, error) {
	vins := make([]*go_protos.Vin, 0)
	for _, one := range this.Vin {
		vins = append(vins, &go_protos.Vin{

			Puk:   one.Puk,
			Sign:  one.Sign,
			Nonce: one.Nonce.Bytes(),
		})
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

	txPay := go_protos.TxVoteOut{
		TxBase:   &txBase,
		Vote:     this.Vote,
		VoteType: uint32(this.VoteType),
	}

	bs, err := txPay.Marshal()
	if err != nil {
		return nil, err
	}
	return &bs, err
}

func (this *Tx_vote_out) Serialize() *[]byte {
	bs := this.TxBase.Serialize()
	buf := bytes.NewBuffer(*bs)
	buf.Write(this.Vote)
	buf.Write(utils.Uint16ToBytes(this.VoteType))
	*bs = buf.Bytes()
	return bs
}

func (this *Tx_vote_out) GetWaitSign(vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Vote...)
	*signDst = append(*signDst, utils.Uint16ToBytes(this.VoteType)...)

	return signDst

}

func (this *Tx_vote_out) GetSign(key *ed25519.PrivateKey, vinIndex uint64) *[]byte {

	signDst := this.GetSignSerialize(nil, vinIndex)

	*signDst = append(*signDst, this.Vote...)
	*signDst = append(*signDst, utils.Uint16ToBytes(this.VoteType)...)

	sign := keystore.Sign(*key, *signDst)

	return &sign

}

func (this *Tx_vote_out) CheckSign() error {

	if this.Vin == nil || len(this.Vin) != 1 {
		return config.ERROR_pay_vin_too_much
	}
	if len(this.Vin[0].Nonce.Bytes()) == 0 {

		return config.ERROR_pay_nonce_is_nil
	}
	if this.Vout_total != 1 {
		return config.ERROR_pay_vout_too_much
	}
	if this.Vout[0].Value <= 0 {
		return config.ERROR_amount_zero
	}

	switch this.VoteType {
	case VOTE_TYPE_community:

	case VOTE_TYPE_vote:
	case VOTE_TYPE_light:
		if this.Vote != nil && len(this.Vote) > 0 {
			return config.ERROR_deposit_light_vote
		}

	}

	for i, one := range this.Vin {

		signDst := this.GetSignSerialize(nil, uint64(i))

		*signDst = append(*signDst, this.Vote...)
		*signDst = append(*signDst, utils.Uint16ToBytes(this.VoteType)...)

		puk := ed25519.PublicKey(one.Puk)

		if config.Wallet_print_serialize_hex {
			engine.Log.Info("sign serialize:%s", hex.EncodeToString(*signDst))
		}
		if !ed25519.Verify(puk, *signDst, one.Sign) {
			return config.ERROR_sign_fail
		}

	}

	return nil
}

func (this *Tx_vote_out) GetWitness() *crypto.AddressCoin {
	witness := crypto.BuildAddr(config.AddrPre, this.Vin[0].Puk)

	return &witness
}

func (this *Tx_vote_out) GetSpend() uint64 {
	return this.Gas
}

func (this *Tx_vote_out) CheckRepeatedTx(txs ...TxItr) bool {
	superiorAddr := this.Vote
	subordinateAddr := this.Vin[0].GetPukToAddr()

	switch this.VoteType {
	case VOTE_TYPE_community:

		if GetDepositCommunityAddr(subordinateAddr) <= 0 {
			return false
		}
	case VOTE_TYPE_vote:

		if value := GetDepositLightVoteValue(subordinateAddr, &superiorAddr); value <= 0 || this.Vout[0].Value > value {
			engine.Log.Info(":%d %d", value, this.Vout[0].Value)
			return false
		}
	case VOTE_TYPE_light:

		if GetDepositLightAddr(subordinateAddr) <= 0 {

			return false
		}

		engine.Log.Info("%v", subordinateAddr)
		superiorAddr := GetVoteAddr(subordinateAddr)
		if superiorAddr != nil && len(*superiorAddr) > 0 {
			if value := GetDepositLightVoteValue(subordinateAddr, superiorAddr); value > 0 {
				engine.Log.Info(":%d %d", value, this.Vout[0].Value)
				return false
			}
		}
	}

	voteOutTotal := uint64(0)

	for _, one := range txs {
		if one.Class() != config.Wallet_tx_type_vote_out {
			continue
		}
		addr := (*one.GetVout())[0].Address

		voteout := one.(*Tx_vote_out)
		rule := voteout.VoteType

		isSelf := bytes.Equal(*subordinateAddr, addr)
		if !isSelf {
			continue
		}

		switch this.VoteType {
		case VOTE_TYPE_community:

		case VOTE_TYPE_vote:

			if rule == VOTE_TYPE_vote {
				value := GetDepositLightVoteValue(subordinateAddr, &superiorAddr)
				if value < this.Vout[0].Value+voteOutTotal {
					engine.Log.Info(":%d %d", value, this.Vout[0].Value, voteOutTotal)
					return false
				} else {
					voteOutTotal += this.Vout[0].Value
				}

			}
		case VOTE_TYPE_light:

		}
	}
	return true
}

func (this *Tx_vote_out) CountTxItemsNew(height uint64) *TxItemCountMap {
	itemCount := TxItemCountMap{
		AddItems: make(map[string]*map[uint64]int64, len(this.Vout)+len(this.Vin)),
		Nonce:    make(map[string]big.Int),
	}

	totalValue := this.Gas
	for _, vout := range this.Vout {

		frozenMap, ok := itemCount.AddItems[utils.Bytes2string(vout.Address)]
		if ok {
			oldValue, ok := (*frozenMap)[vout.FrozenHeight]
			if ok {
				oldValue += int64(vout.Value)
				(*frozenMap)[vout.FrozenHeight] = oldValue
			} else {
				(*frozenMap)[vout.FrozenHeight] = int64(vout.Value)
			}
		} else {
			frozenMap := make(map[uint64]int64, 0)
			frozenMap[vout.FrozenHeight] = int64(vout.Value)
			itemCount.AddItems[utils.Bytes2string(vout.Address)] = &frozenMap
		}
	}

	from := this.Vin[0].GetPukToAddr()
	itemCount.Nonce[utils.Bytes2string(*from)] = this.Vin[0].Nonce
	frozenMap, ok := itemCount.AddItems[utils.Bytes2string(*from)]
	if ok {
		oldValue, ok := (*frozenMap)[0]
		if ok {
			oldValue -= int64(totalValue)
			(*frozenMap)[0] = oldValue
		} else {
			(*frozenMap)[0] = (0 - int64(totalValue))
		}
	} else {
		frozenMap := make(map[uint64]int64, 0)
		frozenMap[0] = (0 - int64(totalValue))
		itemCount.AddItems[utils.Bytes2string(*from)] = &frozenMap
	}
	return &itemCount
}

func (this *Tx_vote_out) CountTxHistory(height uint64) {

	hiIn := HistoryItem{
		IsIn:    true,
		Type:    this.Class(),
		InAddr:  make([]*crypto.AddressCoin, 0),
		OutAddr: make([]*crypto.AddressCoin, 0),

		Txid:   *this.GetHash(),
		Height: height,
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

	}

	addrCoin = make(map[string]bool)
	for _, vout := range this.Vout {

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

	if len(hiIn.OutAddr) > 0 {
		balanceHistoryManager.Add(hiIn)
	}
}

func CreateTxVoteOut(voteType uint16, addr crypto.AddressCoin, amount, gas uint64, pwd string, payload string) (*Tx_vote_out, error) {
	engine.Log.Info("start CreateTxVoteOut")
	chain := forks.GetLongChain()

	currentHeight := chain.GetCurrentBlock()

	vins := make([]*Vin, 0)

	var di *DepositInfo
	switch voteType {
	case VOTE_TYPE_community:
		di = chain.Balance.GetDepositCommunity(&addr)
		amount = di.Value
	case VOTE_TYPE_vote:
		di = chain.Balance.GetDepositVote(&addr)
		if di == nil {
			return nil, config.ERROR_deposit_not_exist
		}
		if amount > di.Value {
			amount = di.Value
		}
	case VOTE_TYPE_light:
		di = chain.Balance.GetDepositLight(&addr)
		amount = di.Value
	default:
		return nil, config.ERROR_tx_not_exist
	}

	puk, ok := keystore.GetPukByAddr(addr)
	if !ok {
		return nil, config.ERROR_public_key_not_exist
	}

	totalAll, item := chain.Balance.BuildPayVinNew(&addr, gas)
	if totalAll < gas {

		return nil, config.ERROR_not_enough
	}
	nonce := chain.GetBalance().FindNonce(item.Addr)
	vin := Vin{
		Puk:   puk,
		Nonce: *new(big.Int).Add(&nonce, big.NewInt(1)),
	}
	vins = append(vins, &vin)

	var dstAddr crypto.AddressCoin = addr
	if addr == nil {

		dstAddr = keystore.GetAddr()[0].Addr
	}

	engine.Log.Info(":%d", amount)

	vouts := make([]*Vout, 0)

	vout2 := Vout{
		Value:   amount,
		Address: dstAddr,
	}
	vouts = append(vouts, &vout2)

	var txout *Tx_vote_out
	for i := uint64(0); i < 10000; i++ {

		base := TxBase{
			Type:       config.Wallet_tx_type_vote_out,
			Vin_total:  uint64(len(vins)),
			Vin:        vins,
			Vout_total: uint64(len(vouts)),
			Vout:       vouts,
			Gas:        gas,
			LockHeight: currentHeight + config.Wallet_tx_lockHeight + i,
		}

		txout = &Tx_vote_out{
			TxBase:   base,
			Vote:     di.WitnessAddr,
			VoteType: voteType,
		}

		for i, one := range txout.Vin {
			for _, key := range keystore.GetAddr() {
				puk, ok := keystore.GetPukByAddr(key.Addr)
				if !ok {

					return nil, config.ERROR_public_key_not_exist
				}

				if bytes.Equal(puk, one.Puk) {
					_, prk, _, err := keystore.GetKeyByAddr(key.Addr, pwd)

					if err != nil {

						return nil, err
					}
					sign := txout.GetSign(&prk, uint64(i))

					txout.Vin[i].Sign = *sign
				}
			}
		}

		txout.BuildHash()
		if txout.CheckHashExist() {
			txout = nil
			continue
		} else {
			break
		}
	}
	chain.Balance.AddLockTx(txout)
	engine.Log.Info("end CreateTxVoteOut %s", hex.EncodeToString(*txout.GetHash()))
	return txout, nil
}
