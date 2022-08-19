package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	tests := map[string]struct {
		items      map[string]Item
		expiration time.Duration
	}{
		"DefaultOptions": {},
		"ModifiedOptions": {
			items: map[string]Item{
				"go-cache": {
					"hello go-cache", time.Minute.Nanoseconds(),
				},
			},
			expiration: 5 * time.Minute,
		},
	}

	for set, tc := range tests {
		t.Run(set, func(t *testing.T) {
			opts := NewOptions(Items(tc.items), Expiration(tc.expiration))

			assert.Equal(t, tc.items, opts.Items)
			assert.Equal(t, tc.expiration, opts.Expiration)
		})
	}
}
