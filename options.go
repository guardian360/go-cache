package cache

import "time"

// Options represents the available options that can be configured while
// instantiating a new cache.
type Options struct {
	// Items represents the predefined items stored in the cache.
	Items map[string]Item
	// Expiration represents the default expiration of the cache.
	Expiration time.Duration
	// CleanupInterval represents the default interval for the janitor to clean
	// the cache of expired items.
	CleanupInterval time.Duration
}

// Option manipulates the Options passed.
type Option func(o *Options)

// Items initializes the cache with preconfigured items.
func Items(i map[string]Item) Option {
	return func(o *Options) {
		o.Items = i
	}
}

// Expiration sets the duration for items stored in the cache to expire.
func Expiration(e time.Duration) Option {
	return func(o *Options) {
		o.Expiration = e
	}
}

// CleanupInterval sets the interval for the janitor to clean the cache of
// expired items.
func CleanupInterval(ci time.Duration) Option {
	return func(o *Options) {
		o.CleanupInterval = ci
	}
}

// NewOptions returns a new Options struct.
func NewOptions(opts ...Option) Options {
	options := Options{
		Items:           make(map[string]Item),
		Expiration:      DefaultExpiration,
		CleanupInterval: DefaultCleanupInterval,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
