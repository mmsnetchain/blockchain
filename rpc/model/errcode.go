package model

import (
	"fmt"
	"strconv"
)

const (
	Success                = 2000
	NoMethod               = 4001
	TypeWrong              = 5001
	NoField                = 5002
	Nomarl                 = 5003
	Timeout                = 5004
	Exist                  = 5005
	FailPwd                = 5006
	NotExist               = 5007
	NotEnough              = 5008
	ContentIncorrectFormat = 5009
	AmountIsZero           = 5010
	RuleField              = 5011
	BalanceNotEnough       = 5012
	VoteExist              = 5013
	VoteNotOpen            = 5014
	RewardCountSync        = 5015
	CommentOverLengthMax   = 5016
	GasTooLittle           = 5017
	NotDepositOutLight     = 5018
)

var codes = map[int]string{
	NoMethod:               "no method",
	TypeWrong:              "type wrong",
	NoField:                "no field",
	Nomarl:                 "",
	Timeout:                "timeout",
	Exist:                  "exist",
	FailPwd:                "fail password",
	NotExist:               "not exist",
	NotEnough:              "not enough",
	ContentIncorrectFormat: "",
	AmountIsZero:           "",
	RuleField:              "",
	BalanceNotEnough:       "BalanceNotEnough",
	VoteExist:              "VoteExist",
	VoteNotOpen:            "VoteNotOpen",
	RewardCountSync:        "reward sync execution",
	CommentOverLengthMax:   "comment over length max",
	GasTooLittle:           "gas too little",
	NotDepositOutLight:     "Not Deposit Out Light",
}

func Errcode(code int, p ...string) (res []byte, err error) {
	res = []byte(strconv.Itoa(code))
	c, ok := codes[code]
	if ok {
		if len(p) > 0 {
			if c == "" {
				err = fmt.Errorf("%s", p[0])
			} else {
				err = fmt.Errorf("%s: %s", p[0], c)
			}
		} else {
			err = fmt.Errorf("%s", c)
		}
	}
	return
}
