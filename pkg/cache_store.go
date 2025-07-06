package pkg

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type CacheStore struct {
	items    map[uuid.UUID]*CacheItem
	mutex    sync.RWMutex
	stopChan chan struct{}
}

func NewCacheStore(cleanupInterval time.Duration) *CacheStore {
	store := &CacheStore{
		items:    make(map[uuid.UUID]*CacheItem),
		stopChan: make(chan struct{}),
	}

	go store.startCleanupTicker(cleanupInterval)

	return store
}

func (store *CacheStore) Set(item *CacheItem) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.items[item.ID] = item
}

func (store *CacheStore) Get(id uuid.UUID) (*CacheItem, bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	item, found := store.items[id]
	if found && item.IsExpired() {
		return nil, false
	}
	return item, found
}

func (store *CacheStore) Delete(id uuid.UUID) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	delete(store.items, id)
}

func (store *CacheStore) startCleanupTicker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			store.cleanupExpiredItems()
		case <-store.stopChan:
			return
		}
	}
}

func (store *CacheStore) cleanupExpiredItems() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	for id, item := range store.items {
		if item.IsExpired() {
			delete(store.items, id)
		}
	}
}

func (store *CacheStore) StopCleanup() {
	close(store.stopChan)
}
