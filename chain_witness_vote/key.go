package chain_witness_vote

import (
	"bytes"
	"errors"
	"github.com/prestonTao/utils"
	"io/ioutil"
	"mmschainnewaccount/config"
	"path/filepath"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var keyLock = new(sync.RWMutex)
var key *Key

func CheckKey() bool {
	ok := false
	keyLock.RLock()
	if key != nil {
		ok = true
	}
	keyLock.RUnlock()
	return ok
}

func LoadKeyStore() (*KeyStore, error) {
	key, err := LoadKeyStoreToLocal()
	if err != nil {

		return nil, err
	}
	return key, err
}

type Keys struct {
	Seeds []SeedKey
}

func (this *Keys) GetAddrs() []string {
	addr := make([]string, 0)
	for _, one := range this.Seeds {
		addr = append(addr, one.GetAddrs()...)
	}
	return addr
}

type SeedKey struct {
	Seed []byte
	Keys []Key
}

func (this *SeedKey) GetAddrs() []string {
	addr := make([]string, 0)
	for _, one := range this.Keys {
		addr = append(addr, one.Addr)
	}
	return addr
}

type Key struct {
	Index int64
	Puk   string
	Addr  string
}

type KeyStore struct {
	Seed  []SeedKeyStore
	Alone []AloneKeyStore
}

type SeedKeyStore struct {
	Key    []byte
	Indexs []int64
}

type AloneKeyStore struct {
	Prk string
}

func (this *KeyStore) Save() error {
	bs, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return utils.SaveFile(filepath.Join(config.Wallet_path, config.Wallet_seed), &bs)
}

func LoadKeyStoreToLocal() (*KeyStore, error) {

	ok, err := utils.PathExists(filepath.Join(config.Wallet_path, config.Wallet_seed))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("There is no key file locally")
	}

	bs, err := ioutil.ReadFile(filepath.Join(config.Wallet_path, config.Wallet_seed))
	if err != nil {
		return nil, err
	}
	keyStore := new(KeyStore)

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()
	err = decoder.Decode(keyStore)
	return keyStore, err

}
