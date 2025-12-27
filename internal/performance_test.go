package internal

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/vhvplatform/go-api-gateway/internal/cache"
	"github.com/vhvplatform/go-api-gateway/internal/circuitbreaker"
	"github.com/vhvplatform/go-api-gateway/internal/middleware"
)

// BenchmarkRateLimiter measures rate limiter performance
func BenchmarkRateLimiter(b *testing.B) {
	rl := middleware.NewRateLimiter(1000.0, 2000)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start cleanup goroutine
	go rl.CleanupLimiters(ctx)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter := rl.GetLimiter("test-key")
			_ = limiter.Allow()
		}
	})
}

// BenchmarkRateLimiterWithCleanup measures cleanup overhead
func BenchmarkRateLimiterWithCleanup(b *testing.B) {
	rl := middleware.NewRateLimiter(1000.0, 2000)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go rl.CleanupLimiters(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate different keys using strconv
		key := "key-" + strconv.Itoa(i%100)
		limiter := rl.GetLimiter(key)
		_ = limiter.Allow()
	}
}

// BenchmarkCircuitBreaker measures circuit breaker performance
func BenchmarkCircuitBreaker(b *testing.B) {
	cb := circuitbreaker.NewCircuitBreaker()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cb.Execute("test-service", func() (interface{}, error) {
				return "success", nil
			})
		}
	})
}

// BenchmarkCircuitBreakerMultipleServices measures overhead with multiple services
func BenchmarkCircuitBreakerMultipleServices(b *testing.B) {
	cb := circuitbreaker.NewCircuitBreaker()
	services := []string{"auth", "user", "tenant"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			service := services[i%len(services)]
			_, _ = cb.Execute(service, func() (interface{}, error) {
				return "success", nil
			})
			i++
		}
	})
}

// BenchmarkCacheOperations measures cache get/set performance
// Note: This benchmark requires a running Redis instance
func BenchmarkCacheOperations(b *testing.B) {
	// Skip if no Redis URL is configured
	redisURL := "redis://localhost:6379/0"
	cacheClient, err := cache.NewCache(redisURL)
	if err != nil {
		b.Skip("Redis not available for benchmark")
		return
	}
	defer cacheClient.Close()

	ctx := context.Background()
	testData := map[string]interface{}{
		"key": "value",
		"num": 12345,
	}

	b.ResetTimer()
	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cacheClient.Set(ctx, "bench-key", testData, 1*time.Minute)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Pre-populate cache
		_ = cacheClient.Set(ctx, "bench-key", testData, 1*time.Minute)

		b.ResetTimer()
		var result map[string]interface{}
		for i := 0; i < b.N; i++ {
			_ = cacheClient.Get(ctx, "bench-key", &result)
		}
	})
}

// BenchmarkContextPropagation measures context propagation overhead
func BenchmarkContextPropagation(b *testing.B) {
	ctx := context.Background()
	cb := circuitbreaker.NewCircuitBreaker()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cb.ExecuteContext(ctx, "test-service", func() (interface{}, error) {
				return "success", nil
			})
		}
	})
}
