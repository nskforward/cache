package cache

import (
	"io"
	"os"
	"sync"
	"time"
)

type FileEntry struct {
	path   string
	meta   map[string]any
	mx     sync.RWMutex
	ticker *time.Ticker
	size   int64
}

func NewFileEntry(path string) *FileEntry {
	return &FileEntry{
		path: path,
	}
}

func (entry *FileEntry) SetTTL(ttl time.Duration) *FileEntry {
	entry.mx.Lock()
	defer entry.mx.Unlock()
	if entry.ticker == nil {
		entry.ticker = time.NewTicker(ttl)
	} else {
		entry.ticker.Reset(ttl)
	}
	return entry
}

func (en *FileEntry) SetMeta(name string, value any) *FileEntry {
	en.mx.Lock()
	defer en.mx.Unlock()
	if en.meta == nil {
		en.meta = make(map[string]any)
	}
	en.meta[name] = value
	return en
}

func (en *FileEntry) GetMeta(name string) any {
	en.mx.RLock()
	defer en.mx.RUnlock()
	if en.meta == nil {
		return nil
	}
	return en.meta[name]
}

func (en *FileEntry) Size() int64 {
	en.mx.RLock()
	defer en.mx.RUnlock()
	return en.size
}

func (en *FileEntry) Remove() error {
	en.mx.Lock()
	defer en.mx.Unlock()
	return os.Remove(en.path)
}

func (en *FileEntry) Expire() {
	en.mx.Lock()
	defer en.mx.Unlock()
	if en.ticker != nil {
		en.ticker.Stop()
	}
}

func (en *FileEntry) Write(r io.Reader) error {
	tmpPath := en.path + ".tmp"
	en.mx.Lock()
	defer en.mx.Unlock()
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer f.Close()
	size, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	en.size = size
	return os.Rename(tmpPath, en.path)
}

func (en *FileEntry) Read(w io.Writer) error {
	en.mx.RLock()
	defer en.mx.RUnlock()
	f, err := os.Open(en.path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}
