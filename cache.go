package cache

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

// Item represents an item stored in the cache.
type Item struct {
	Value      interface{}
	Expiration int64
}

// Expired returns true if the item has expired.
func (i *Item) Expired() bool {
	if i.Expiration == NoExpiration.Nanoseconds() || i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

var (
	// NoExpiration is the duration for items without an expiration.
	NoExpiration time.Duration = -1
	// DefaultExpiration is the default duration for items stored in
	// the cache to expire.
	DefaultExpiration time.Duration = 0
	// DefaultCleanupInterval is the default duration for the janitor to clean
	// up expired items from the cache.
	DefaultCleanupInterval time.Duration = 0

	// ErrItemExpired is returned in Cache.Get when the item found in the cache
	// has expired.
	ErrItemExpired error = errors.New("item has expired")
	// ErrKeyNotFound is returned in Cache.Get and Cache.Delete when the
	// provided key could not be found in the cache.
	ErrKeyNotFound error = errors.New("key not found in cache")
	// ErrValueNotStored is returned in Cache.Put when the value could not be
	// stored in the cache.
	ErrValueNotStored error = errors.New("value not stored")
)

// Cache is the interface that wraps the cache.
type Cache interface {
	// Get gets a value from the cache by key. Returns the value or nil and an
	// error indicating whether the key was found.
	Get(k string) (interface{}, error)
	// Put stores an item to the cache replacing any existing item. Returns an
	// error when the item could not be stored.
	Put(k string, v interface{}, d time.Duration) error
	// Delete removes a key from the cache. Returns an error when the key could
	// not be removed.
	Delete(k string) error
	// Items copies all unexpired items in the cache to a new map and returns
	// it.
	Items() map[string]Item
	// WithExpired sets a flag to determine whether or not to include expired
	// items in the outputs of Get and Items. It returns itself so it can be
	// used fluently.
	WithExpired(e bool) Cache
	// DeleteExpired removes all expired items from the cache.
	DeleteExpired()
}

// New returns a new cache.
func New(opts ...Option) Cache {
	options := NewOptions(opts...)

	var items map[string]Item
	if len(options.Items) > 0 {
		items = options.Items
	} else {
		items = make(map[string]Item)
	}

	var ci time.Duration
	if options.CleanupInterval > 0 {
		ci = options.CleanupInterval
	} else {
		ci = DefaultCleanupInterval
	}

	c := &cache{
		opts:  options,
		items: items,
	}

	if ci > 0 {
		c.janitor = NewJanitor(ci)
		go c.janitor.Start(c)
		runtime.SetFinalizer(c, c.janitor.Stop)
	}

	return c
}

type cache struct {
	opts Options

	mu          sync.RWMutex
	items       map[string]Item
	withExpired bool
	janitor     Janitor
}

func (c *cache) Get(k string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[k]
	if !found {
		return nil, ErrKeyNotFound
	}
	if !c.withExpired && item.Expired() {
		return nil, ErrItemExpired
	}
	return item.Value, nil
}

func (c *cache) Put(k string, v interface{}, d time.Duration) error {
	var e int64
	if d == DefaultExpiration {
		d = c.opts.Expiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{v, e}
	return nil
}

func (c *cache) Delete(k string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[k]
	if !found {
		return ErrKeyNotFound
	}

	delete(c.items, k)
	return nil
}

func (c *cache) Items() map[string]Item {
	c.mu.Lock()
	defer c.mu.Unlock()

	m := make(map[string]Item)
	for k, v := range c.items {
		if !c.withExpired && v.Expired() {
			continue
		}
		m[k] = v
	}
	return m
}

func (c *cache) WithExpired(e bool) Cache {
	c.withExpired = e
	return c
}

func (c *cache) DeleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if !v.Expired() {
			continue
		}
		delete(c.items, k)
	}
}
