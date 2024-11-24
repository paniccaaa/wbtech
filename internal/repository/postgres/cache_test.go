package postgres

import (
	"reflect"
	"testing"
	"time"

	"github.com/paniccaaa/wbtech/internal/model"
)

func Test_newCache(t *testing.T) {
	tests := []struct {
		name        string
		ttl         time.Duration
		cleanupTick time.Duration
		want        *Cache
	}{
		{
			name:        "Create cache with valid TTL and cleanup tick",
			ttl:         1 * time.Second,
			cleanupTick: 500 * time.Millisecond,
			want: &Cache{
				ttl:         1 * time.Second,
				cleanupTick: 500 * time.Millisecond,
				data:        make(map[model.OrderUID]model.Order),
				expiration:  make(map[model.OrderUID]time.Time),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCache(tt.ttl, tt.cleanupTick)
			if got.ttl != tt.want.ttl || got.cleanupTick != tt.want.cleanupTick || reflect.DeepEqual(got.data, tt.want.data) == false {
				t.Errorf("newCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_startCleanup(t *testing.T) {
	cache := newCache(500*time.Millisecond, 100*time.Millisecond)
	order := model.Order{OrderUID: "order123"}
	cache.Set(order)

	t.Run("Start cleanup and verify expired orders are removed", func(t *testing.T) {
		go cache.startCleanup()
		time.Sleep(600 * time.Millisecond) // Wait for cleanup to run
		_, ok := cache.Get(order.OrderUID)
		if ok {
			t.Errorf("Cache.startCleanup() did not remove expired order")
		}
	})
}

func TestCache_cleanup(t *testing.T) {
	cache := newCache(500*time.Millisecond, 100*time.Millisecond)
	order := model.Order{OrderUID: "order123"}
	cache.Set(order)

	t.Run("Cleanup expired orders manually", func(t *testing.T) {
		time.Sleep(600 * time.Millisecond)
		cache.cleanup()
		_, ok := cache.Get(order.OrderUID)
		if ok {
			t.Errorf("Cache.cleanup() did not remove expired order")
		}
	})
}

func TestCache_Set(t *testing.T) {
	cache := newCache(1*time.Second, 500*time.Millisecond)
	order := model.Order{OrderUID: "order123"}

	t.Run("Set order and verify it exists", func(t *testing.T) {
		cache.Set(order)
		got, ok := cache.Get(order.OrderUID)
		if !ok || !reflect.DeepEqual(got, order) {
			t.Errorf("Cache.Set() = %v, %v; want %v, true", got, ok, order)
		}
	})
}

func TestCache_Get(t *testing.T) {
	cache := newCache(1*time.Second, 500*time.Millisecond)
	order := model.Order{OrderUID: "order123"}
	cache.Set(order)

	t.Run("Get existing order", func(t *testing.T) {
		got, ok := cache.Get(order.OrderUID)
		if !ok || !reflect.DeepEqual(got, order) {
			t.Errorf("Cache.Get() = %v, %v; want %v, true", got, ok, order)
		}
	})

	t.Run("Get non-existing order", func(t *testing.T) {
		_, ok := cache.Get("nonexistent")
		if ok {
			t.Errorf("Cache.Get() = true; want false")
		}
	})
}

func TestCache_Restore(t *testing.T) {
	cache := newCache(1*time.Second, 500*time.Millisecond)
	orders := []model.Order{
		{OrderUID: "order123"},
		{OrderUID: "order456"},
	}

	t.Run("Restore multiple orders and verify they exist", func(t *testing.T) {
		cache.Restore(orders)
		for _, order := range orders {
			got, ok := cache.Get(order.OrderUID)
			if !ok || !reflect.DeepEqual(got, order) {
				t.Errorf("Cache.Restore() = %v, %v; want %v, true", got, ok, order)
			}
		}
	})
}
