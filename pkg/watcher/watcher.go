package watcher

import (
	"github.com/stokkelol/lightbot/pkg/cache"
	"time"
)

// Watcher is a simple struct that holds the last timestamp
type Watcher struct {
	cache *cache.Cache
}

// NewWatcher creates a new watcher
func NewWatcher() *Watcher {
	return &Watcher{cache: cache.NewCache()}
}

// Run starts the watcher
func (w *Watcher) Run() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		w.SetLastTimestamp(time.Now())
	}
}

// SetLastTimestamp sets the last timestamp
func (w *Watcher) SetLastTimestamp(t time.Time) {
	w.cache.SetLastTimestamp(t.Unix())
}
