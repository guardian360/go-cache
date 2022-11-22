package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMiss(t *testing.T) {
	c := New()

	v, err := c.Get("foo")

	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestGetHit(t *testing.T) {
	i := map[string]Item{"foo": {"bar", DefaultExpiration.Nanoseconds()}}
	c := New(Items(i))

	v, err := c.Get("foo")

	assert.NoError(t, err)
	assert.Equal(t, "bar", v)
}

func TestGetNoExpiration(t *testing.T) {
	i := map[string]Item{"foo": {"bar", NoExpiration.Nanoseconds()}}
	c := New(Items(i))

	<-time.After(5 * time.Millisecond)
	v, err := c.Get("foo")

	assert.NoError(t, err)
	assert.Equal(t, "bar", v)
}

func TestGetExpired(t *testing.T) {
	d := time.Now().Add(1 * time.Millisecond).UnixNano()
	i := map[string]Item{"foo": {"bar", d}}
	c := New(Items(i))

	<-time.After(5 * time.Millisecond)
	v, err := c.Get("foo")

	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestGetValid(t *testing.T) {
	d := time.Now().Add(5 * time.Millisecond).UnixNano()
	i := map[string]Item{"foo": {"bar", d}}
	c := New(Items(i))

	<-time.After(1 * time.Millisecond)
	v, err := c.Get("foo")

	assert.NoError(t, err)
	assert.Equal(t, "bar", v)
}

func TestPut(t *testing.T) {
	c := New()
	err := c.Put("foo", "bar", 0)
	assert.NoError(t, err)
}

func TestDeleteMiss(t *testing.T) {
	c := New()
	err := c.Delete("foo")
	assert.Error(t, err)
}

func TestDeleteHit(t *testing.T) {
	i := map[string]Item{"foo": {"bar", 0}}
	c := New(Items(i))

	err := c.Delete("foo")

	assert.NoError(t, err)
}

func TestItems(t *testing.T) {
	i := map[string]Item{
		"DefaultExpiration": {"foo", DefaultExpiration.Nanoseconds()},
		"NoExpiration":      {"bar", NoExpiration.Nanoseconds()},
		"Expired":           {"baz", time.Now().Add(-time.Hour).UnixNano()},
	}
	c := New(Items(i))

	assert.Contains(t, c.Items(), "DefaultExpiration")
	assert.Contains(t, c.Items(), "NoExpiration")
	assert.NotContains(t, c.Items(), "Expired")
}

func TestDefaultExpiration(t *testing.T) {
	c := New(Expiration(1 * time.Millisecond))
	_ = c.Put("foo", "bar", 0)

	<-time.After(5 * time.Millisecond)
	v, err := c.Get("foo")

	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestCleanupInterval(t *testing.T) {
	c := New(CleanupInterval(time.Millisecond))
	_ = c.Put("foo", "bar", 3*time.Millisecond)

	<-time.After(5 * time.Millisecond)

	assert.NotContains(t, c.WithExpired(true).Items(), "foo")
}
