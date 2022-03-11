package state

import (
	"github.com/ethereum/go-ethereum/log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// sharedStorage is used to store maps of originStorage of stateObjects
type SharedStorage struct {
	poolLock   *sync.RWMutex
	shared_map map[common.Address]sync.Map
}

func NewSharedStorage() SharedStorage {
	sharedMap := make(map[common.Address]sync.Map, 1000)
	return SharedStorage{
		poolLock:   &sync.RWMutex{},
		shared_map: sharedMap,
	}
}

func (storage *SharedStorage) GetStorage(address common.Address, key common.Hash) (interface{}, bool) {
	storage.poolLock.RLock()
	storageMap, ok := storage.shared_map[address]
	storage.poolLock.RUnlock()
	if !ok {
		log.Error("can not find originStorage on:" + address.String())
		return nil, false
	}
	val, ok := storageMap.Load(key)
	log.Info("read originStorage on:" + address.String() + "key:" + key.String() + "val:" + val.(common.Hash).String())
	return storageMap.Load(key)
}

func (storage *SharedStorage) setStorage(address common.Address, key common.Hash, val common.Hash) {
	storage.poolLock.RLock()
	log.Info("write originStorage on:" + address.String() + "key:" + key.String() + "val:" + val.String())
	storageMap, ok := storage.shared_map[address]
	storage.poolLock.RUnlock()
	if !ok {
		log.Error("can not find originStorage on:" + address.String())
	}
	storageMap.Store(key, val)
}

// Check whether the storage exist in pool,
// new one if not exist, it will be fetched in stateObjects.GetCommittedState()
func (storage *SharedStorage) checkSharedStorage(address common.Address) {
	storage.poolLock.RLock()
	_, ok := storage.shared_map[address]
	storage.poolLock.RUnlock()
	log.Info("check object on:"+address.String(), "begin")
	if !ok {
		log.Info("check object on:"+address.String(), "not finish , new one")
		m := sync.Map{}
		storage.poolLock.Lock()
		storage.shared_map[address] = m
		storage.poolLock.Unlock()
	}
	log.Info("check object on:"+address.String(), "finish")
}
