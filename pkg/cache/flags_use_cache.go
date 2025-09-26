package cache

import (
	"log"
	"sync"
	"time"
)

type FlagsUseCacheItem struct {
	value      string
	expiration time.Time
}

type command struct {
	action string
	key    string
	done   chan struct{}
}

type FlagsUseCache struct {
	items    map[string]FlagsUseCacheItem
	mutex    sync.RWMutex
	commands chan command
}

var flagsCache = NewFlagsUseCache()

func NewFlagsUseCache() *FlagsUseCache {
	cache := &FlagsUseCache{
		items:    make(map[string]FlagsUseCacheItem),
		commands: make(chan command, 1000),
	}
	go cache.processor()
	return cache
}

func GetFlagsUseCache() *FlagsUseCache {
	return flagsCache
}

func (c *FlagsUseCache) processor() {
	for cmd := range c.commands {
		c.mutex.Lock()
		c.items[cmd.key] = FlagsUseCacheItem{
			value:      "",
			expiration: time.Now().Add(1 * time.Hour),
		}
		c.mutex.Unlock()

		if cmd.done != nil {
			close(cmd.done)
		}
	}
}

func (c *FlagsUseCache) AddFlagKey(key string) {
	select {
	case c.commands <- command{action: "add", key: key}:
	default:
		log.Printf("Cache command queue full, skipping key: %s", key)
	}
}

func (c *FlagsUseCache) Exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exists := c.items[key]
	if !exists {
		return false
	}
	if time.Now().After(item.expiration) {
		return false
	}
	return true
}
