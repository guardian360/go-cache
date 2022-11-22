package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJanitorStart(t *testing.T) {
	j := &janitor{
		interval: time.Millisecond,
		stop:     make(chan bool),
	}
	c := New()
	_ = c.Put("foo", "bar", 3*time.Millisecond)

	go j.Start(c)
	<-time.After(5 * time.Millisecond)

	assert.NotContains(t, c.WithExpired(true).Items(), "foo")
}

func TestJanitorStop(t *testing.T) {
	j := &janitor{
		interval: time.Millisecond,
		stop:     make(chan bool),
	}
	c := New()
	go j.Start(c)

	j.Stop(c)

	assert.False(t, <-j.stop)
}
