package postgres

import (
	"sync"
	"time"

	"github.com/paniccaaa/wbtech/internal/model"
)

type Cache struct {
	mu          sync.RWMutex
	data        map[model.OrderUID]model.Order
	expiration  map[model.OrderUID]time.Time
	ttl         time.Duration
	cleanupTick time.Duration
}

func newCache(ttl, cleanupTick time.Duration) *Cache {
	cache := &Cache{
		data:        make(map[model.OrderUID]model.Order),
		expiration:  make(map[model.OrderUID]time.Time),
		ttl:         ttl,
		cleanupTick: cleanupTick,
	}

	// start background goroutine
	go cache.startCleanup()

	return cache
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupTick)
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}

func (c *Cache) cleanup() {
	currentTime := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, expTime := range c.expiration {
		if currentTime.After(expTime) {
			// remove old elem
			delete(c.data, key)
			delete(c.expiration, key)
		}
	}
}

func (c *Cache) Set(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[order.OrderUID] = order
	c.expiration[order.OrderUID] = time.Now().Add(c.ttl)
}

func (c *Cache) Get(orderUID model.OrderUID) (model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.data[orderUID]
	return order, ok
}

func (c *Cache) Restore(orders []model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		// add new elem
		c.data[order.OrderUID] = order
		c.expiration[order.OrderUID] = time.Now().Add(c.ttl)
	}
}
