package cache

import (
	"time"
)

// Janitor is the interface that is responsible for cleaning the cache of
// expired items.
type Janitor interface {
	// Start runs the main janitor loop responsible for cleaning the cache of
	// expired items. It accepts the cache to watch as its first argument and
	// returns nothing.
	Start(c Cache)
	// Stop stops the janitor from cleaning the cache of expired items. It
	// accepts the cache to stop watching as its first argument and returns
	// nothing.
	Stop(c Cache)
}

// NewJanitor returns a new janitor.
func NewJanitor(ci time.Duration) Janitor {
	return &janitor{
		interval: ci,
		stop:     make(chan bool),
	}
}

type janitor struct {
	interval time.Duration
	stop     chan bool
}

func (j *janitor) Start(c Cache) {
	ticker := time.NewTicker(j.interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func (j *janitor) Stop(c Cache) {
	close(j.stop)
}
