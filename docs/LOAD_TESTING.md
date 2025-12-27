# Load Testing and Performance Benchmarking

This guide describes how to conduct load testing and performance benchmarking for the API Gateway.

## Benchmarking

### Running Benchmarks

The project includes built-in benchmarks for critical performance paths:

```bash
# Run all benchmarks
go test -bench=. -benchmem ./internal/

# Run specific benchmark
go test -bench=BenchmarkRateLimiter -benchmem ./internal/

# Run with CPU profiling
go test -bench=. -benchmem -cpuprofile=cpu.prof ./internal/
go tool pprof cpu.prof

# Run with memory profiling
go test -bench=. -benchmem -memprofile=mem.prof ./internal/
go tool pprof mem.prof
```

### Benchmark Coverage

Current benchmarks cover:

1. **Rate Limiter Performance**
   - Single key throughput
   - Multiple keys with cleanup
   - Parallel access

2. **Circuit Breaker Performance**
   - Single service execution
   - Multiple services
   - Context propagation overhead

3. **Cache Operations**
   - Set operations
   - Get operations
   - Requires Redis running locally

### Interpreting Results

Good benchmark results should show:
- Rate limiter: < 1000 ns/op for single key access
- Circuit breaker: < 500 ns/op for successful execution
- Cache operations: < 1ms for local Redis

## Load Testing

### Prerequisites

Install load testing tools:

```bash
# Install hey (HTTP load generator)
go install github.com/rakyll/hey@latest

# Or use Apache Bench
sudo apt-get install apache2-utils

# Or use k6
brew install k6  # macOS
# or download from https://k6.io/
```

### Basic Load Test

```bash
# Start the gateway
make run

# Or with Docker
docker-compose up -d api-gateway

# Run load test with hey (1000 requests, 50 concurrent)
hey -n 1000 -c 50 http://localhost:8080/health

# With authentication
hey -n 1000 -c 50 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/users
```

### Apache Bench Example

```bash
# 10000 requests with 100 concurrent connections
ab -n 10000 -c 100 http://localhost:8080/health

# With authentication and POST
ab -n 1000 -c 10 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -p payload.json \
  http://localhost:8080/api/v1/users
```

### Advanced k6 Script

Create `load-test.js`:

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up to 20 users
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 100 }, // Ramp to 100 users
    { duration: '1m', target: 100 },  // Stay at 100 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests < 500ms
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
  },
};

export default function () {
  // Test health endpoint
  let res = http.get('http://localhost:8080/health');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 200ms': (r) => r.timings.duration < 200,
  });
  
  sleep(1);
}
```

Run with:
```bash
k6 run load-test.js
```

### Rate Limiting Test

Test rate limiting behavior:

```bash
# Fast requests to trigger rate limiting
hey -n 1000 -c 100 -q 200 http://localhost:8080/api/v1/users

# Check metrics for rate limit rejections
curl http://localhost:8080/metrics | grep rate_limiter_rejected
```

## Monitoring During Load Tests

### Watch Key Metrics

```bash
# Watch metrics endpoint
watch -n 1 'curl -s http://localhost:8080/metrics | grep -E "(active_requests|duration|rate_limiter|circuit_breaker)"'

# Monitor goroutines
watch -n 5 'curl -s http://localhost:8080/debug/pprof/goroutine?debug=1 | head -1'

# Monitor memory
watch -n 5 'curl -s http://localhost:8080/debug/pprof/heap?debug=1 | head -10'
```

### CPU Profiling During Load

```bash
# Start load test in background
hey -n 10000 -c 100 http://localhost:8080/api/v1/users &

# Capture 30-second CPU profile
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# Analyze
go tool pprof cpu.prof
```

### Memory Profiling

```bash
# Before load test
curl http://localhost:8080/debug/pprof/heap > heap-before.prof

# During load test
hey -n 10000 -c 100 http://localhost:8080/api/v1/users

# After load test
curl http://localhost:8080/debug/pprof/heap > heap-after.prof

# Compare
go tool pprof -base=heap-before.prof heap-after.prof
```

## Performance Targets

### Response Time Targets

- **Health check**: < 10ms (p95)
- **Authenticated requests**: < 100ms (p95) without backend
- **Backend proxied requests**: < 500ms (p95)
- **With circuit breaker open**: < 5ms (immediate failure)

### Throughput Targets

- **Without rate limiting**: > 1000 rps per instance
- **With rate limiting (100 rps)**: Stable at configured limit
- **Circuit breaker overhead**: < 5% latency increase

### Resource Usage Targets

- **Memory**: < 512MB under normal load
- **CPU**: < 50% per core at 1000 rps
- **Goroutines**: < 200 under normal load
- **Rate limiters**: Automatic cleanup keeps count < 1000

## Continuous Performance Testing

### Automated Performance Checks

Add to CI/CD pipeline:

```bash
# In .github/workflows/performance.yml or similar

# Run benchmarks with baseline comparison
go test -bench=. -benchmem ./internal/ > new-bench.txt

# Compare with baseline (if exists)
if [ -f baseline-bench.txt ]; then
  benchcmp baseline-bench.txt new-bench.txt
fi

# Store new baseline
cp new-bench.txt baseline-bench.txt
```

### Performance Regression Detection

Set up alerts for:
- Benchmark performance degradation > 10%
- Memory usage increase > 20%
- Response time p95 increase > 50ms
- Rate of 5xx errors > 1%

## Troubleshooting Performance Issues

### High Latency

1. Check backend service response times
2. Review circuit breaker state
3. Check cache hit rate
4. Monitor database connection pool

### High CPU Usage

1. Check log level (should be warn/error in production)
2. Review rate limiting configuration
3. Profile with pprof
4. Check for hot paths in code

### Memory Growth

1. Monitor rate limiter active count
2. Check goroutine count
3. Review cache size
4. Check for connection leaks

### Rate Limiting Issues

1. Monitor rejection metrics
2. Review configuration vs actual load
3. Check for IP clustering (NAT/proxy)
4. Consider per-tenant limiting

## References

- [Prometheus Metrics](http://localhost:8080/metrics)
- [pprof Profiling](http://localhost:8080/debug/pprof/)
- [TROUBLESHOOTING.md](../TROUBLESHOOTING.md)
- [go-infrastructure Performance Standards](https://github.com/vhvplatform/go-infrastructure)
