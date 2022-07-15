package mining

import (
	_ "github.com/prestonTao/libp2parea/engine"
	"github.com/prestonTao/utils"
	"sync"
)

var TxCache *Cache

func init() {
	TxCache = &Cache{
		txitemLock:         new(sync.RWMutex),
		txitem:             utils.NewCache(40000),
		txCacheLock:        new(sync.RWMutex),
		txCache:            utils.NewCache(40000 * 2),
		blockHeadCacheLock: new(sync.RWMutex),
		blockHeadCache:     utils.NewCache(100),
	}
}

type Cache struct {
	txitemLock         *sync.RWMutex
	txitem             *utils.Cache
	txCacheLock        *sync.RWMutex
	txCache            *utils.Cache
	blockHeadCacheLock *sync.RWMutex
	blockHeadCache     *utils.Cache
}

func (this *Cache) AddTxInTxItem(keybs []byte, txItr TxItr) {

	key := utils.Bytes2string(keybs)
	this.txitemLock.Lock()
	_, ok := this.txitem.Get(key)
	if !ok {
		this.txitem.Add(key, txItr)
	}
	this.txitemLock.Unlock()

}

func (this *Cache) FindTxInCache(keybs []byte) (TxItr, bool) {

	key := utils.Bytes2string(keybs)
	this.txCacheLock.RLock()
	v, ok := this.txCache.Get(key)
	this.txCacheLock.RUnlock()
	if !ok {
		this.txitemLock.RLock()
		v, ok = this.txitem.Get(key)
		this.txitemLock.RUnlock()
	}
	if ok {
		tx := v.(TxItr)
		return tx, ok
	}
	return nil, ok
}

func (this *Cache) AddTxInCache(keybs []byte, txItr TxItr) {
	if txItr == nil {
		return
	}
	key := utils.Bytes2string(keybs)
	this.txCacheLock.Lock()
	_, ok := this.txCache.Get(key)
	if !ok {
		this.txCache.Add(key, txItr)
	}
	this.txCacheLock.Unlock()
}

func (this *Cache) FlashTxInCache(keybs []byte, txItr TxItr) {
	if txItr == nil {
		return
	}
	key := utils.Bytes2string(keybs)
	this.txCacheLock.Lock()
	this.txCache.Add(key, txItr)
	this.txCacheLock.Unlock()
}

func (this *Cache) AddBlockHeadCache(keybs []byte, bh *BlockHead) {
	key := utils.Bytes2string(keybs)
	this.blockHeadCacheLock.Lock()
	_, ok := this.blockHeadCache.Get(key)
	if !ok {
		this.blockHeadCache.Add(key, bh)
	}
	this.blockHeadCacheLock.Unlock()
}

func (this *Cache) FindBlockHeadCache(keybs []byte) (*BlockHead, bool) {
	key := utils.Bytes2string(keybs)
	this.blockHeadCacheLock.RLock()
	v, ok := this.blockHeadCache.Get(key)
	this.blockHeadCacheLock.RUnlock()
	if ok {
		tx := v.(*BlockHead)
		return tx, ok
	}
	return nil, ok
}

func (this *Cache) FlashBlockHeadCache(keybs []byte, bh *BlockHead) {
	key := utils.Bytes2string(keybs)
	this.blockHeadCacheLock.RLock()
	v, ok := this.blockHeadCache.Get(key)
	this.blockHeadCacheLock.RUnlock()
	if !ok {
		this.AddBlockHeadCache(keybs, bh)

		return
	}

	bhOld := v.(*BlockHead)
	bhOld.Nextblockhash = bh.Nextblockhash
}
