package services

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/patrickmn/go-cache"
)

var (
	c *Cache
)

// Cache describe the cache client
type Cache struct {
	c      *cache.Cache
	enable bool
}

// CacheInterface describe a cache client interface
type CacheInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, data interface{})
}

// NewCache create a new cache client
func NewCache(enable bool) *Cache {
	if c != nil && c.enable == enable {
		return c
	}

	exp := 10 * time.Second
	c = &Cache{
		c:      cache.New(exp, 1*time.Minute),
		enable: enable,
	}

	if !c.enable {
		log.Println("Cache is disabled")
	} else {
		log.Debugf("Set up cache with default expiration at %s", exp)
	}

	return c
}

// Get return stored data at given key
func (c *Cache) Get(key string) (interface{}, bool) {
	if c.enable {
		log.Debugln("Data stored in cache")
		return c.c.Get(key)
	}
	return nil, false
}

// Set store data at given key
func (c *Cache) Set(key string, data interface{}) {
	if c.enable {
		c.c.Set(key, data, cache.DefaultExpiration)
	}
}
