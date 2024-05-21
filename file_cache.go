package cache

import (
	"path/filepath"
	"sync"
)

type FileCache struct {
	baseDir string
	items   map[string]*FileEntry
	mx      sync.RWMutex
}

func NewFileCache(baseDir string) *FileCache {
	return &FileCache{
		baseDir: baseDir,
		items:   make(map[string]*FileEntry),
	}
}

func (c *FileCache) Find(key string) *FileEntry {
	c.mx.RLock()
	defer c.mx.RUnlock()
	entry, ok := c.items[key]
	if ok {
		return entry
	}
	return nil
}

func (c *FileCache) FindOrCreate(key string) *FileEntry {
	c.mx.RLock()
	entry, ok := c.items[key]
	c.mx.RUnlock()
	if ok {
		return entry
	}
	entry = NewFileEntry(filepath.Join(c.baseDir, key))
	c.mx.Lock()
	c.items[key] = entry
	c.mx.Unlock()
	return entry
}

func (c *FileCache) RemoveOnExpiration(key string, entry *FileEntry) {
	if entry.ticker != nil {
		go func() {
			<-entry.ticker.C
			entry.Remove()
			c.mx.Lock()
			defer c.mx.Unlock()
			delete(c.items, key)
		}()
	}
}

func (c *FileCache) RemoveNow(key string) {
	entry := c.Find(key)
	if entry != nil {
		entry.Expire()
		entry.Remove()
		c.mx.Lock()
		defer c.mx.Unlock()
		delete(c.items, key)
	}
}
