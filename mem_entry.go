package cache

import (
	"sync"
	"time"
)

type MemEntry struct {
	mx     sync.RWMutex
	ticker *time.Ticker
	value  any
	meta   map[string]any
}

func NewMemEntry() *MemEntry {
	return &MemEntry{}
}

func (en *MemEntry) SetMeta(name string, value any) *MemEntry {
	en.mx.Lock()
	defer en.mx.Unlock()
	if en.meta == nil {
		en.meta = make(map[string]any)
	}
	en.meta[name] = value
	return en
}

func (en *MemEntry) SetValue(value any) *MemEntry {
	en.mx.Lock()
	defer en.mx.Unlock()
	en.value = value
	return en
}

func (en *MemEntry) GetValue() any {
	en.mx.RLock()
	defer en.mx.RUnlock()
	return en.value
}

func (en *MemEntry) GetMeta(name string) any {
	en.mx.RLock()
	defer en.mx.RUnlock()
	if en.meta == nil {
		return nil
	}
	return en.meta[name]
}

func (en *MemEntry) SetTTL(ttl time.Duration) *MemEntry {
	en.mx.Lock()
	defer en.mx.Unlock()
	if en.ticker != nil {
		en.ticker.Reset(ttl)
	} else {
		en.ticker = time.NewTicker(ttl)
	}
	return en
}

func (en *MemEntry) Expire() {
	en.mx.Lock()
	defer en.mx.Unlock()
	if en.ticker != nil {
		en.ticker.Stop()
	}
}
