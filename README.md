# go-cache

go-cache is an in-memory key-value store that is suitable for running
single-machine applications. It is essentially a thread-safe
`map[string]interface{}` with expiration times. Any object can be stored for a
given duration or forever and the cache can safely be used by multiple
goroutines.

## Installation

go-cache can be installed by using the `go get` command.

```bash
go get github.com/guardian360/go-cache
```

## Usage

```golang
package main

import (
	"fmt"
	"time"

	"github.com/guardian360/go-cache"
)

func main() {
	// Create a cache with a default expiration time of 5 minutes and a
	// cleanup interval of 10 minutes.
	c := cache.New(
		cache.Expiration(5 * time.Minute),
		cache.CleanupInterval(10 * time.Minute),
	)

	// Set the value of the key "foo" to "bar" with the default expiration
	// time.
	c.Put("foo", "bar", 0)

	// Set the value of the key "baz" to 42 with no expiration time.
	c.Put("baz", 42, cache.NoExpiration)

	// Get the string associated with the key "foo" from the cache.
	foo, err := c.Get("foo")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(foo)

	// Delete a key from the cache.
	c.Delete("foo")
}
```

## Reference

Documentation can be found using the `go doc` command or at [pkg.go.dev][docs].

## Credits

* [go-cache][go-cache-credits]
* [go-micro cache][go-micro-cache-credits]

[docs]: https://pkg.go.dev/github.com/guardian360/go-cache
[go-cache-credits]: https://github.com/patrickmn/go-cache
[go-micro-cache-credits]: https://github.com/asim/go-micro/tree/master/cache
