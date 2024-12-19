package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data    any
	Expires time.Time
}

type Cache struct {
	mu    sync.RWMutex
	store map[string]CacheEntry
}

var cache = Cache{
	store: make(map[string]CacheEntry),
}

func DefaultExpiration() time.Time {
	return time.Now().Add(time.Hour * 24)
}

func Set(key string, value any, expires time.Time) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.store[key] = CacheEntry{
		Data:    value,
		Expires: expires,
	}
}

func Get(key string) any {
	cache.mu.RLock()
	entry, ok := cache.store[key]
	cache.mu.RUnlock()

	if !ok {
		return nil
	}

	if entry.Expires.Before(time.Now()) {
		Invalidate(key)
		return nil
	}

	return entry.Data
}

func Invalidate(key string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	delete(cache.store, key)
}
