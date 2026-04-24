package state

import "sync"

var (
	mu    sync.RWMutex
	cache map[string]int8
)

func SetCache(key string, value int8) {
	mu.Lock()
	defer mu.Unlock()
	if cache == nil {
		cache = make(map[string]int8)
	}
	cache[key] = value
}

func GetCache(key string) int8 {
	mu.RLock()
	defer mu.RUnlock()
	return cache[key]
}
