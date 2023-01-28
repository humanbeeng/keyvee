package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string]string
}

func New() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, found := c.data[string(key)]

	if !found {
		return []byte(val), fmt.Errorf("not found %s", key)
	}

	return []byte(val), nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, found := c.data[string(key)]
	return found
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	go func() {
		<-time.After(ttl)
		delete(c.data, string(key))
	}()
	c.data[string(key)] = string(value)
	return nil
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}
