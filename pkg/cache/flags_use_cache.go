package cache

import (
	"sync"
	"time"
)

type FlagsUseCacheItem struct {
	value      string
	expiration time.Time
}

type FlagsUseCache struct {
	items map[string]FlagsUseCacheItem
	mutex sync.RWMutex
}

var flagsCache = &FlagsUseCache{
	items: make(map[string]FlagsUseCacheItem),
}

func GetFlagsUseCache() *FlagsUseCache {
	return flagsCache
}

func (c *FlagsUseCache) AddFlagKey(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = FlagsUseCacheItem{
		value:      "",
		expiration: time.Now().Add(1 * time.Hour),
	}
}

func (c *FlagsUseCache) Exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false
	}

	// Check if expired
	if time.Now().After(item.expiration) {
		return false
	}

	return true
}
