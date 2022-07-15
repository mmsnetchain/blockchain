package config

import (
	"flag"
	"math/big"
	"path/filepath"
	"time"
)

const (
	Version_0 = 0
	Version_1 = 1
)

const (
	Wallet_path          = "wallet"
	Wallet_path_prkName  = "ec_prk.pem"
	Wallet_path_pukName  = "ec_puk.pem"
	Wallet_seed          = "seed_key.json"
	Wallet_addr_puk_type = "EC PUBLIC KEY"
)

const (
	Wallet_tx_type_start          = 0
	Wallet_tx_type_mining         = 1
	Wallet_tx_type_deposit_in     = 2
	Wallet_tx_type_deposit_out    = 3
	Wallet_tx_type_pay            = 4
	Wallet_tx_type_account        = 5
	Wallet_tx_type_account_cancel = 6
	Wallet_tx_type_vote_in        = 7
	Wallet_tx_type_vote_out       = 8

	Wallet_tx_type_token_publish = 10
	Wallet_tx_type_token_payment = 11

	Wallet_tx_type_spaces_mining_in  = 12
	Wallet_tx_type_spaces_mining_out = 13
	Wallet_tx_type_spaces_use_in     = 14
	Wallet_tx_type_spaces_use_out    = 15
	Wallet_tx_type_voting_reward     = 16

	Wallet_tx_type_end = 100
)

const (
	Mining_coin_total            = 13 * 10000 * 10000 * 1e8
	Mining_coin_premining        = 242164181942014
	Mining_coin_rest             = Mining_coin_total - Mining_coin_premining
	Mining_block_cycle           = 180
	Mining_block_time            = 10
	Mining_block_start_height    = 1
	Mining_group_start_height    = Mining_block_start_height
	Mining_block_hash_count      = 100
	Mining_group_min             = 1
	Mining_group_max             = 3
	Mining_deposit               = uint64(100000 * 1e8)
	Mining_vote                  = uint64(10000 * 1e8)
	Mining_light_min             = uint64(10 * 1e8)
	Mining_name_deposit_min      = uint64(1 * 1e8)
	Mining_community_reward_time = 60
	Wallet_community_reward_max  = 500
	Mining_pay_vin_max           = 100
	Mining_pay_vout_max          = 10000
	Witness_backup_min           = Mining_group_min
	Witness_backup_max           = 31
	Witness_backup_reward_max    = 99

	Block_size_max       = 1024 * 1024 * 8
	Wallet_tx_lockHeight = 30

	Block_confirm                         = 6
	Wallet_balance_history                = 10
	Wallet_sync_block_interval_time       = time.Second / 30
	Wallet_sync_block_interval_time_relax = time.Second / 300
	Witness_token_supply_min              = 1

	Wallet_frozen_time_min = 1600000000

	Wallet_vote_start_height   = 0
	Wallet_not_build_block_max = time.Hour * 24 * 3
	Wallet_sync_block_timeout  = 10
	Wallet_multicas_block_time = 20

	Wallet_Memory_percentage_max = 90
	Wallet_addr_tx_count_max     = 128
	Wallet_tx_gas_min            = uint64(10 * 1e8)
)

const (
	DB_name = "data"
	DB_temp = "temp"
)

var (
	DB_path                     = filepath.Join(Wallet_path, DB_name)
	DB_path_temp                = filepath.Join(Wallet_path, DB_temp)
	Miner                       = false
	InitNode                    = false
	LoadNode                    = false
	DB_is_null                  = false
	Wallet_keystore_default_pwd = "xhy19liu21@"

	SubmitDepositin            = false
	AlreadyMining              = false
	EnableCache                = true
	StartBlockHash             = []byte{}
	Wallet_print_serialize_hex = false

	Witness_backup_group            = 5
	Witness_backup_group_overheight = uint64(100000)
	Witness_backup_group_new        = 60
)

func ParseInitFlag() bool {
	if InitNode {
		return true
	}
	for _, param := range flag.Args() {
		switch param {
		case "init":
			InitNode = true
			Model = Model_complete
			return true
		case "load":
			LoadNode = true
			Model = Model_complete
			return true
		}
	}
	return false
}

func ClacRewardForBlockHeight(height uint64) uint64 {
	height += BlockRewardHeightOffset

	heightBig := big.NewInt(int64(height))
	num31 := big.NewInt(31)
	num347 := big.NewInt(347000)
	res := big.NewInt(28600000000)
	if height < 490000*31+1 {
		temp := new(big.Int).Div(heightBig, num31)
		temp = new(big.Int).Mul(temp, num347)
		temp = new(big.Int).Add(res, temp)
		temp = new(big.Int).Div(temp, num31)
		return temp.Uint64()
	} else {
		num49 := big.NewInt(490000)
		num95 := big.NewInt(95)
		num100 := big.NewInt(100)
		num2 := big.NewInt(20000)
		base := new(big.Int).Mul(num49, num347)
		base = new(big.Int).Add(res, base)
		res = new(big.Int).Mul(num95, base)

		temp1 := new(big.Int).Mul(num49, num31)
		temp1 = new(big.Int).Sub(heightBig, temp1)
		temp2 := new(big.Int).Mul(num2, num31)
		temp := new(big.Int).Div(temp1, temp2)

		div1 := new(big.Int).Exp(num95, temp, nil)
		div2 := new(big.Int).Exp(num100, temp, nil)

		res = new(big.Int).Mul(res, div1)
		res = new(big.Int).Div(res, num100)
		temp = new(big.Int).Div(res, div2)

		temp = new(big.Int).Div(temp, num31)

		if temp.Cmp(big.NewInt(0)) == -1 {
			return 0
		}

		return temp.Uint64()
	}

}

func CheckAddBlacklist(witnessCount, total uint64) bool {
	maxHeight := uint64(Wallet_not_build_block_max / (Mining_block_time * time.Second))
	if witnessCount < Mining_group_max {
		maxHeight = maxHeight / witnessCount
	} else {
		maxHeight = maxHeight / Mining_group_max
	}
	if maxHeight < total {
		return true
	}
	return false
}
