package cache

import "sync"

// Cache is a simple struct that holds the last timestamp
type Cache struct {
	mu            sync.RWMutex
	lastTimestamp int64
}

// NewCache creates a new cache
func NewCache() *Cache {
	return &Cache{}
}

// SetLastTimestamp sets the last timestamp
func (c *Cache) SetLastTimestamp(t int64) {
	c.mu.Lock()
	c.lastTimestamp = t
	c.mu.Unlock()
}
