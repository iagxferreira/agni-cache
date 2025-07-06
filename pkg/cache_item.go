package pkg

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type CacheItem struct {
	ID         uuid.UUID
	Value      []byte
	Expiration int64
	mutex      sync.RWMutex
}

func NewCacheItem(value []byte, ttl time.Duration) *CacheItem {
	return &CacheItem{
		ID:         uuid.New(),
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (item *CacheItem) SetValue(value []byte) {
	item.mutex.Lock()
	defer item.mutex.Unlock()
	item.Value = value
}

func (ci *CacheItem) GetValue() []byte {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()
	return ci.Value
}

func (ci *CacheItem) IsExpired() bool {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()
	return time.Now().UnixNano() > ci.Expiration
}
