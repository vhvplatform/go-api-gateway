package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/ristretto"
)

// Cache provides caching functionality using Ristretto (in-memory)
type Cache struct {
	client *ristretto.Cache[string, any]
}

// NewCache creates a new Ristretto cache instance
// maxCost: maximum cost of cache (approx memory usage in bytes if cost=1 means 1 byte, or count)
// numCounters: should be 10x the number of keys
func NewCache(maxCost int64, numCounters int64) (*Cache, error) {
	config := &ristretto.Config[string, any]{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: 64, // recommend 64
	}

	cache, err := ristretto.NewCache(config)
	if err != nil {
		return nil, err
	}

	return &Cache{client: cache}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	val, found := c.client.Get(key)
	if !found {
		// Mimic Redis nil error if needed, or just return custom error
		return context.DeadlineExceeded // Or a custom ErrCacheMiss
	}

	// If val is []byte (stored as JSON)
	if data, ok := val.([]byte); ok {
		return json.Unmarshal(data, dest)
	}

	// If stored as object directly (Ristretto can do that, but let's stick to JSON for consistency with previous interface just in case)
	// But actually Ristretto stores interface{}, so we could store the object directly.
	// However, the interface asks for `dest interface{}` pointer usually for Unmarshal.
	// To be safe and compatible with "Unmarshal" style usage:

	// Check if dest is a pointer
	// For now, let's assume we store []byte to transparently support the "Unmarshal" pattern used before
	// or we can try to cast back if we stored the object.
	// Let's stick to []byte storage for consistency with the previous Redis implementation which used JSON.

	return json.Unmarshal(val.([]byte), dest)
}

// Set stores a value in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// Ristretto cost: 1 per item or length of data. Let's use length of data for MaxCost to work as size limit.
	cost := int64(len(data))
	c.client.SetWithTTL(key, data, cost, ttl)
	return nil
}

// Delete removes a value from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.client.Del(key)
	return nil
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	_, found := c.client.Get(key)
	return found, nil
}

// Close closes the cache (Ristretto implementation uses Close)
func (c *Cache) Close() error {
	c.client.Close()
	return nil
}
