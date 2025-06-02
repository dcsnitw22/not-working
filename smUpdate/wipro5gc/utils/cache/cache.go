package cache

import (
        "sync"
        "time"
        "github.com/benbjohnson/clock"
)

type ItemId     interface{}
// Cache allows caching work items (like session ids) with a timestamp. An item is
// considered ready to process if the timestamp has expired.
type WorkCache interface {
        // GetItem fetches and returns all ready items in the cache
        GetItem() []ItemId

        CheckItem(item ItemId) bool
        // AddItem inserts a new item or overwrites an existing item in cache
        AddItem(item ItemId, delay time.Duration)
        // RemoveItem removes an item from the cache
        RemoveItem(item ItemId)
        // Get Number of items in cache
        GetNumItems() int
}

type Cache struct {
        clock clock.Clock
        lock  sync.Mutex
        cache map[ItemId]time.Time
}

// New Cache returns a new basic WorkCache with the provided clock
func NewCache(clock clock.Clock) WorkCache {
        cache := make(map[ItemId]time.Time)
        return &Cache{cache: cache, clock: clock}
}

func (c *Cache) GetItem() []ItemId {
        // Lock the cache and defer until processing is done on cache
        c.lock.Lock()
        defer c.lock.Unlock()

        // Get current time
        now := c.clock.Now()

        // Get all items that have expired
        var items []ItemId
        for i, t := range c.cache {
                if t.Before(now) {
                        items = append(items, i)
                        delete(c.cache, i)
                }
        }
        return items
}

func (c *Cache) CheckItem(item ItemId) bool {
        // Lock the cache and defer until processing is done on cache
        c.lock.Lock()
        defer c.lock.Unlock()

        _, ok := c.cache[item]
        return ok
}

func (c *Cache) AddItem(item ItemId, delay time.Duration) {
        // Lock the cache and defer until processing is done on cache
        c.lock.Lock()
        defer c.lock.Unlock()

        // Add item in the cache
        c.cache[item] = c.clock.Now().Add(delay)
}

func (c *Cache) RemoveItem(item ItemId) {
        // Lock the cache and defer until processing is done on cache
        c.lock.Lock()
        defer c.lock.Unlock()

        // Remove item in the cache
        delete(c.cache, item)
}

func (c *Cache) GetNumItems() int {
        return len(c.cache)
}

