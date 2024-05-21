package cache

import (
	"sync"
)

type MemCache struct {
	items map[string]*MemEntry
	mx    sync.RWMutex
}

func NewMemCache() *MemCache {
	return &MemCache{
		items: make(map[string]*MemEntry),
	}
}

func (c *MemCache) Find(key string) *MemEntry {
	c.mx.RLock()
	entry, ok := c.items[key]
	c.mx.RUnlock()
	if ok {
		return entry
	}
	return nil
}

func (c *MemCache) FindOrCreate(key string) *MemEntry {
	c.mx.RLock()
	entry, ok := c.items[key]
	c.mx.RUnlock()
	if ok {
		return entry
	}
	entry = &MemEntry{}
	c.mx.Lock()
	c.items[key] = entry
	c.mx.Unlock()
	return entry
}

func (c *MemCache) RemoveOnExpiration(key string, entry *MemEntry) {
	if entry.ticker != nil {
		go func() {
			<-entry.ticker.C
			c.mx.Lock()
			defer c.mx.Unlock()
			delete(c.items, key)
		}()
	}
}

func (c *MemCache) RemoveNow(key string) {
	entry := c.Find(key)
	if entry != nil {
		entry.Expire()
		c.mx.Lock()
		defer c.mx.Unlock()
		delete(c.items, key)
	}
}
