package gache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex

	defaultExpireDuration time.Duration
	items                 map[string]*Item
}

type Item struct {
	Data   interface{}
	Expire time.Time
}

func (c *Cache) clear() {
	for range time.Tick(time.Minute) {
		c.Lock()
		var deleteKeys []string
		for key, i := range c.items {
			if i.Expire.Before(time.Now()) {
				deleteKeys = append(deleteKeys, key)
			}
		}
		for idx := range deleteKeys {
			delete(c.items, deleteKeys[idx])
		}
		c.Unlock()
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	i, ok := c.items[key]
	c.RUnlock()
	if !ok {
		return nil, false
	}

	if i.Expire.After(time.Now()) {
		return i.Data, true
	}

	c.Lock()
	defer c.Unlock()
	i, ok = c.items[key]
	if ok {
		if i.Expire.After(time.Now()) {
			return i.Data, true
		}
		delete(c.items, key)
	}
	return nil, false
}

func (c *Cache) Set(key string, val interface{}, expireDuration ...time.Duration) {
	expire := time.Now()
	if len(expireDuration) > 0 && expireDuration[0] > 0 {
		expire = expire.Add(expireDuration[0])
	} else {
		expire = expire.Add(c.defaultExpireDuration)
	}

	c.Lock()
	c.items[key] = &Item{Data: val, Expire: expire}
	c.Unlock()
}

func NewCacheNoExpire() *Cache {
	cache := &Cache{
		defaultExpireDuration: -1,
		items:                 map[string]*Item{},
	}
	go cache.clear()
	return cache
}

func NewCacheWithExpire(defaultExpireDuration time.Duration) *Cache {
	if defaultExpireDuration < 0 {
		defaultExpireDuration = -1
	}
	cache := &Cache{
		defaultExpireDuration: defaultExpireDuration,
		items:                 map[string]*Item{},
	}

	go cache.clear()
	return cache
}
