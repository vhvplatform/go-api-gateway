package cache

import (
"context"
"encoding/json"
"time"

"github.com/redis/go-redis/v9"
)

// Cache provides caching functionality using Redis
type Cache struct {
client *redis.Client
}

// NewCache creates a new cache instance
func NewCache(redisURL string) (*Cache, error) {
opt, err := redis.ParseURL(redisURL)
if err != nil {
return nil, err
}

client := redis.NewClient(opt)

// Test connection
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := client.Ping(ctx).Err(); err != nil {
return nil, err
}

return &Cache{client: client}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
val, err := c.client.Get(ctx, key).Result()
if err != nil {
return err
}
return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
data, err := json.Marshal(value)
if err != nil {
return err
}
return c.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes a value from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
return c.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (c *Cache) Close() error {
return c.client.Close()
}
