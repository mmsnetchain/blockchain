package config

const (
	block_start_str                = "startblock"
	History                        = "1_"
	BlockHeight                    = "2_"
	Name                           = "name_"
	Block_Highest                  = "HighestBlock"
	LEVELDB_Head_history_balance   = "h_b_"
	WitnessName                    = "3_"
	WitnessAddr                    = "4_"
	AlreadyUsed_tx                 = "5_"
	TokenInfo                      = "6_"
	TokenPublishTxid               = "7_"
	DB_PRE_Tx_Not_Import           = "8_"
	DB_spaces_mining_addr          = "9_"
	DB_community_addr_hash         = "10_"
	DBKEY_tx_blockhash             = "11_"
	DBKEY_addr_value               = "12_"
	DBKEY_community_reward         = "13_"
	DBKEY_community_addr           = "14_"
	DBKEY_vote_addr                = "15_"
	DBKEY_deposit_witness_addr     = "16_"
	DBKEY_deposit_light_addr       = "17_"
	DBKEY_deposit_community_addr   = "18_"
	DBKEY_deposit_light_vote_value = "19_"
	DBKEY_addr_nonce               = "20_"
	DBKEY_token_addr_value         = "21_"
)

var (
	Key_block_start = []byte(block_start_str)

	DBKEY_zset_frozen_height          = []byte{22}
	DBKEY_zset_frozen_time            = []byte{23}
	DBKEY_zset_frozen_height_children = []byte{24}
	DBKEY_zset_frozen_time_children   = []byte{25}
	DBKEY_addr_frozen_value           = []byte{26}
	DBKEY_addr_value_big              = []byte{27}
)

func BuildTxNotImport(txid []byte) []byte {
	return append([]byte(DB_PRE_Tx_Not_Import), txid...)
}

func BuildCommunityAddrStartHeight(addr []byte) []byte {
	return append([]byte(DB_community_addr_hash), addr...)
}

func BuildTxToBlockHash(txid []byte) []byte {
	return append([]byte(DBKEY_tx_blockhash), txid...)
}

func BuildAddrValue(addr []byte) []byte {
	return append([]byte(DBKEY_addr_value), addr...)
}

func BuildDBKeyCommunityAddrFrozen(addr []byte) *[]byte {
	key := append([]byte(DBKEY_community_reward), addr...)
	return &key
}

func BuildDBKeyCommunityAddr(addr []byte) *[]byte {
	key := append([]byte(DBKEY_community_addr), addr...)
	return &key
}

func BuildDBKeyVoteAddr(addr []byte) *[]byte {
	key := append([]byte(DBKEY_vote_addr), addr...)
	return &key
}

func BuildDBKeyDepositWitnessAddr(addr []byte) *[]byte {
	key := append([]byte(DBKEY_deposit_witness_addr), addr...)
	return &key
}

func BuildDBKeyDepositLightAddr(addr []byte) *[]byte {
	key := append([]byte(DBKEY_deposit_light_addr), addr...)
	return &key
}

func BuildDBKeyDepositCommunityAddr(addr []byte) *[]byte {
	key := append([]byte(DBKEY_deposit_community_addr), addr...)
	return &key
}

func BuildDBKeyDepositLightVoteValue(lightAddr, communityAddr []byte) *[]byte {
	keyBs := []byte(DBKEY_deposit_light_vote_value)
	key := make([]byte, 0, len(keyBs)+len(lightAddr)+len(communityAddr))
	key = append(keyBs, lightAddr...)
	key = append(key, communityAddr...)
	return &key

}

func BuildDBKeyAddrNonce(addr []byte) *[]byte {
	key := append([]byte(DBKEY_addr_nonce), addr...)
	return &key
}

func BuildDBKeyTokenAddrValue(token, addr []byte) *[]byte {
	keyBs := []byte(DBKEY_token_addr_value)
	key := make([]byte, 0, len(keyBs)+len(token)+len(addr))
	key = append(keyBs, token...)
	key = append(key, addr...)
	return &key
}

func BuildAddrFrozen(addr []byte) []byte {
	return append(DBKEY_addr_frozen_value, addr...)
}

func BuildAddrValueBig(addr []byte) []byte {
	return append(DBKEY_addr_value_big, addr...)
}
