package token

import (
	"bytes"
	"mmschainnewaccount/chain_witness_vote/db"
	"mmschainnewaccount/config"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func SaveTokenInfo(tokenid []byte, name, symbol string, supply uint64) error {
	tokeninfo := TokenInfo{
		Txid:   tokenid,
		Name:   name,
		Symbol: symbol,
		Supply: supply,
	}
	bs, err := json.Marshal(tokeninfo)
	if err != nil {
		return err
	}
	return db.LevelTempDB.Save(BuildKeyForPublishToken(tokenid), &bs)
}

func FindTokenInfo(tokenid []byte) (*TokenInfo, error) {
	bs, err := db.LevelTempDB.Find(BuildKeyForPublishToken(tokenid))
	if err != nil {
		return nil, err
	}
	tokeninfo := new(TokenInfo)
	buf := bytes.NewBuffer(*bs)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(tokeninfo)
	return tokeninfo, err
}

func BuildKeyForPublishToken(txid []byte) []byte {
	return append([]byte(config.TokenInfo), txid...)
}

type TokenInfo struct {
	Txid   []byte
	Name   string
	Symbol string
	Supply uint64
}
