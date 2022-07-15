package config

import (
	"errors"
	"strconv"
)

const (
	ERROR_fail = 5005
)

var (
	ERROR_chain_sysn_block_fail    = errors.New("sync fail,not find block")
	ERROR_chain_sync_block_timeout = errors.New("sync block timeout")
	ERROR_wait_msg_timeout         = errors.New("wait message timeout")

	ERROR_deposit_witness    = errors.New("deposit shoud be:" + strconv.Itoa(int(Mining_deposit)))
	ERROR_deposit_not_exist  = errors.New("deposit not exist")
	ERROR_deposit_exist      = errors.New("deposit exist")
	ERROR_deposit_light_vote = errors.New("tx vote feild error")

	ERROR_password_fail            = errors.New("password fail")
	ERROR_not_enough               = errors.New("balance is not enough")
	ERROR_token_not_enough         = errors.New("token balance is not enough")
	ERROR_public_key_not_exist     = errors.New("not find public key")
	ERROR_amount_zero              = errors.New("Transfer amount cannot be 0")
	ERROR_tx_not_exist             = errors.New("Transaction not found")
	ERROR_tx_format_fail           = errors.New("Error parsing transaction")
	ERROR_tx_Repetitive_vin        = errors.New("Duplicate VIN in transaction")
	ERROR_tx_is_use                = errors.New("Transaction has been used")
	ERROR_tx_fail                  = errors.New("Transaction error")
	ERROR_tx_lockheight            = errors.New("Lock height error")
	ERROR_tx_frozenheight          = errors.New("frozen height error")
	ERROR_public_and_addr_notMatch = errors.New("The public key and address do not match")
	ERROR_sign_fail                = errors.New("Signature error")
	ERROR_vote_exist               = errors.New("Vote already exists")
	ERROR_pay_vin_too_much         = errors.New("vin too much")
	ERROR_pay_vout_too_much        = errors.New("vout too much")

	ERROR_name_deposit = errors.New("Domain name deposit is required at least" + strconv.Itoa(int(Mining_name_deposit_min/1e8)))

	ERROR_name_not_self      = errors.New("Domain name does not belong to itself")
	ERROR_name_exist         = errors.New("Domain name already exists")
	ERROR_name_not_exist     = errors.New("Domain name does not exist")
	ERROR_get_sign_data_fail = errors.New("Error getting signed data")
	ERROR_params_not_enough  = errors.New("params not enough")
	ERROR_params_fail        = errors.New("params fail")
	ERROR_token_min_fail     = errors.New("params token min fail")

	ERROR_get_node_conn_fail = errors.New("get node conn fail")

	ERROR_get_reward_count_sync     = errors.New("get reward count sync")
	ERROR_vote_reward_addr_disunity = errors.New("vote reward address disunity")
	ERROR_pay_nonce_is_nil          = errors.New("nonce is nil")
	ERROR_addr_value_big            = errors.New("The balance of the address is too large")
)
