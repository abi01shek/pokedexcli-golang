package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	expTime time.Time
	val     []byte
}

type Cache struct {
	cacheMap   map[string]cacheEntry
	mu         *sync.RWMutex
	expiration time.Duration
	ticker     *time.Ticker
}

// NewCache create a new empty cache with default values
func NewCache(expiration time.Duration) *Cache {
	c := new(Cache)
	c.cacheMap = make(map[string]cacheEntry)
	c.mu = new(sync.RWMutex)
	c.expiration = expiration
	c.ticker = time.NewTicker(expiration)
	return c
}

// Add: add a new string and byte array to the cache
func (c *Cache) Add(key string, val []byte) {
	thisCe := new(cacheEntry)
	thisCe.val = val
	thisCe.expTime = time.Now().Add(c.expiration)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = *thisCe
}

// Get: get the byte array for a specific string
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if ce, exists := c.cacheMap[key]; exists {
		// if current time is after the expiration time
		if time.Now().After(ce.expTime) {
			delete(c.cacheMap, key)
			return nil, false
		}
		return ce.val, true
	}
	return nil, false
}

// ReadLoop: loops around the cache and clears out expired data
// forever loop with a blocking ticker channel
func (c *Cache) ReadLoop() {
	for {
		<-c.ticker.C // block for ticker channel
		// clean up the caches
		c.mu.Lock()
		for k := range c.cacheMap {
			if time.Now().After(c.cacheMap[k].expTime) {
				delete(c.cacheMap, k)
			}
		}
		c.mu.Unlock()
		//time.Sleep(c.expiration)
	}
}
