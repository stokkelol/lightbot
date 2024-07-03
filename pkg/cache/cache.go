package cache

import (
	"sync"
	"time"
)

// Cache is a simple struct that holds the last timestamp
type Cache struct {
	mu            sync.RWMutex
	lastTimestamp time.Time
}

// NewCache creates a new cache
func NewCache() *Cache {
	return &Cache{}
}

// SetLastTimestamp sets the last timestamp
func (c *Cache) SetLastTimestamp(t time.Time) {
	c.mu.Lock()
	c.lastTimestamp = t
	c.mu.Unlock()
}

// GetLastTimestamp returns the last timestamp
func (c *Cache) GetLastTimestamp() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.lastTimestamp
}
