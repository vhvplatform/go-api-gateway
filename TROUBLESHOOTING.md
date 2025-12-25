# API Gateway Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the API Gateway service.

## Table of Contents

- [General Issues](#general-issues)
- [Performance Issues](#performance-issues)
- [Authentication & Authorization](#authentication--authorization)
- [Rate Limiting Issues](#rate-limiting-issues)
- [Circuit Breaker Issues](#circuit-breaker-issues)
- [Connection Issues](#connection-issues)
- [Memory & Resource Issues](#memory--resource-issues)
- [Monitoring & Debugging](#monitoring--debugging)

---

## General Issues

### Gateway Won't Start

**Symptom**: Gateway exits immediately or fails to start

**Possible Causes & Solutions**:

1. **Missing Environment Variables**
   ```bash
   # Check required environment variables
   echo $API_GATEWAY_PORT
   echo $AUTH_SERVICE_URL
   echo $USER_SERVICE_URL
   echo $TENANT_SERVICE_URL
   ```
   
   **Solution**: Set all required environment variables. See [Configuration](README.md#configuration) for complete list.

2. **Port Already in Use**
   ```bash
   # Check if port 8080 is already in use
   lsof -i :8080
   # or
   netstat -tuln | grep 8080
   ```
   
   **Solution**: 
   - Kill the process using the port: `kill -9 <PID>`
   - Or change `API_GATEWAY_PORT` to a different port

3. **Invalid Configuration**
   ```bash
   # Check gateway logs for configuration errors
   docker logs api-gateway
   ```
   
   **Solution**: Verify all configuration values are valid (URLs, ports, secrets, etc.)

### Health Check Failing

**Symptom**: `/health` endpoint returns 503 or unhealthy status

**Diagnosis**:
```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected healthy response:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "auth-service": "healthy",
    "user-service": "healthy",
    "tenant-service": "healthy"
  }
}
```

**Solutions**:

1. **Backend Service Unreachable**
   - Verify backend services are running
   - Check network connectivity to backend services
   - Verify service URLs in environment variables

2. **DNS Resolution Issues**
   ```bash
   # Test DNS resolution
   nslookup auth-service
   ping auth-service
   ```
   
   **Solution**: Update `/etc/hosts` or DNS configuration

---

## Performance Issues

### High Latency

**Symptom**: Requests take longer than expected

**Diagnosis**:
```bash
# Check request duration metrics
curl http://localhost:8080/metrics | grep api_gateway_request_duration

# Check active requests
curl http://localhost:8080/metrics | grep api_gateway_active_requests
```

**Solutions**:

1. **Timeout Configuration Too High**
   - Default timeout is 30 seconds
   - Reduce timeout in `internal/middleware/timeout.go` if needed

2. **Backend Service Slow**
   - Check backend service performance
   - Review distributed traces in Jaeger
   - Enable circuit breaker to fail fast

3. **Database Connection Pool Exhausted**
   - Increase backend service connection pool size
   - Monitor active database connections

4. **No Caching**
   ```bash
   # Verify Redis is connected
   redis-cli ping
   ```
   
   **Solution**: 
   - Set `REDIS_URL` environment variable
   - Verify Redis is accessible from gateway

### High CPU Usage

**Symptom**: Gateway pods/containers consuming excessive CPU

**Diagnosis**:
```bash
# Check CPU usage
docker stats api-gateway

# Profile the application
go tool pprof http://localhost:8080/debug/pprof/profile
```

**Solutions**:

1. **Too Many Requests**
   - Increase `RATE_LIMIT_RPS` is too high
   - Add more gateway replicas
   - Enable horizontal pod autoscaling

2. **Inefficient Logging**
   - Reduce log level in production
   - Disable debug logging

3. **Memory Leaks**
   - Update to latest version
   - Monitor memory usage over time
   - Check for goroutine leaks: `http://localhost:8080/debug/pprof/goroutine`

---

## Authentication & Authorization

### 401 Unauthorized Errors

**Symptom**: Valid JWT tokens being rejected

**Diagnosis**:
```bash
# Test token validation
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/v1/users

# Decode JWT to inspect (use jwt.io or jwt-cli)
jwt decode YOUR_TOKEN
```

**Solutions**:

1. **Wrong JWT Secret**
   - Verify `JWT_SECRET` matches auth service secret
   - Check for whitespace or special characters in secret

2. **Token Expired**
   - Check token `exp` claim
   - Request new token via `/api/v1/auth/refresh`

3. **Token Not in Header**
   - Ensure header format: `Authorization: Bearer <token>`
   - Check for typos in header name

4. **Token Blacklisted**
   - Check Redis for blacklisted tokens
   - User may have logged out

### 403 Forbidden Errors

**Symptom**: Authenticated user denied access

**Diagnosis**:
```bash
# Check user roles in JWT
jwt decode YOUR_TOKEN | jq '.roles'
```

**Solutions**:

1. **Insufficient Permissions**
   - Verify user has required role
   - Check RBAC configuration in backend services

2. **Wrong Tenant Context**
   - Verify `tenant_id` in JWT matches resource tenant
   - Check multi-tenancy configuration

---

## Rate Limiting Issues

### Frequent 429 Too Many Requests

**Symptom**: Legitimate requests being rate limited

**Diagnosis**:
```bash
# Check rate limit metrics
curl http://localhost:8080/metrics | grep rate_limit

# Current configuration
echo "RPS: $RATE_LIMIT_RPS"
echo "Burst: $RATE_LIMIT_BURST"
```

**Solutions**:

1. **Rate Limit Too Low**
   ```bash
   # Increase limits
   export RATE_LIMIT_RPS=500
   export RATE_LIMIT_BURST=1000
   ```

2. **Shared IP (NAT/Proxy)**
   - Multiple users behind same IP
   - Consider per-tenant rate limiting instead of per-IP
   - Implement user-based rate limiting

3. **Aggressive Polling**
   - Review client application behavior
   - Implement exponential backoff on client side
   - Add `Retry-After` header handling

### Rate Limiter Memory Growth

**Symptom**: Gateway memory usage continuously increasing

**Diagnosis**:
```bash
# Check number of active limiters
curl http://localhost:8080/metrics | grep rate_limit_active_limiters

# Monitor memory over time
docker stats api-gateway
```

**Solutions**:

1. **Cleanup Not Running**
   - Verify cleanup goroutine is active
   - Check logs for cleanup errors
   - Default cleanup interval: 10 minutes

2. **Too Many Unique IPs**
   - Each unique IP creates a limiter
   - Consider shorter cleanup interval
   - Reduce inactivity timeout from 10 minutes

---

## Circuit Breaker Issues

### Circuit Stuck Open

**Symptom**: Circuit breaker remains open even after backend recovers

**Diagnosis**:
```bash
# Check circuit breaker state
curl http://localhost:8080/metrics | grep circuit_breaker_state

# Check backend service health
curl http://backend-service:port/health
```

**Solutions**:

1. **Timeout Too Long**
   - Default timeout: 30 seconds
   - Reduce timeout in `internal/circuitbreaker/breaker.go`

2. **Half-Open Requests Failing**
   - Backend still unstable
   - Check backend logs for errors
   - Increase `MaxRequests` in half-open state

3. **Manual Reset Required**
   ```bash
   # Restart gateway to reset all breakers
   kubectl rollout restart deployment/api-gateway
   ```

### Circuit Opening Too Quickly

**Symptom**: Circuit breaker trips on minor issues

**Diagnosis**:
```bash
# Check failure threshold
# Current: 60% failure rate, minimum 3 requests
```

**Solutions**:

1. **Threshold Too Sensitive**
   - Edit `internal/circuitbreaker/breaker.go`
   - Increase failure ratio threshold (e.g., 0.6 → 0.8)
   - Increase minimum request count (e.g., 3 → 10)

2. **Timeout Too Short**
   - Increase request timeout
   - Backend service may need more time

---

## Connection Issues

### gRPC Connection Refused

**Symptom**: Cannot connect to backend gRPC services

**Diagnosis**:
```bash
# Test gRPC connectivity
grpcurl -plaintext auth-service:50051 list

# Check if service is listening
telnet auth-service 50051
```

**Solutions**:

1. **Service Not Running**
   - Start backend service
   - Check service logs

2. **Wrong Port or Host**
   - Verify environment variables:
     ```bash
     echo $AUTH_SERVICE_URL
     echo $USER_SERVICE_URL
     echo $TENANT_SERVICE_URL
     ```

3. **Firewall/Security Group**
   - Check firewall rules
   - Verify security group allows traffic on gRPC ports

4. **TLS/SSL Issues**
   - Gateway expects plaintext gRPC by default
   - If backend uses TLS, update client configuration

### Redis Connection Failed

**Symptom**: Cannot connect to Redis cache

**Diagnosis**:
```bash
# Test Redis connection
redis-cli -u $REDIS_URL ping

# Check Redis is running
docker ps | grep redis
```

**Solutions**:

1. **Redis Not Running**
   ```bash
   docker start redis
   # or
   redis-server
   ```

2. **Wrong Connection URL**
   ```bash
   # Correct format
   export REDIS_URL=redis://localhost:6379/0
   # With password
   export REDIS_URL=redis://:password@localhost:6379/0
   ```

3. **Redis Out of Memory**
   ```bash
   redis-cli INFO memory
   ```
   
   **Solution**: 
   - Increase Redis memory limit
   - Configure eviction policy
   - Clear old keys

---

## Memory & Resource Issues

### Memory Leak

**Symptom**: Memory usage grows continuously

**Diagnosis**:
```bash
# Take heap dump
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Analyze with pprof
go tool pprof heap.prof
```

**Solutions**:

1. **Rate Limiter Not Cleaning Up**
   - Verify cleanup goroutine runs
   - Check cleanup interval (default: 10 minutes)

2. **Goroutine Leak**
   ```bash
   # Check goroutine count
   curl http://localhost:8080/debug/pprof/goroutine?debug=1
   ```
   
   **Solution**: Update to latest version with goroutine fixes

3. **Large Response Caching**
   - Reduce Redis TTL
   - Limit cacheable response size
   - Implement cache size limits

### Container OOMKilled

**Symptom**: Container crashes with Out Of Memory

**Diagnosis**:
```bash
# Check container events
kubectl describe pod api-gateway-xxx
docker inspect api-gateway
```

**Solutions**:

1. **Insufficient Memory Limit**
   ```yaml
   # Increase memory limit in deployment
   resources:
     limits:
       memory: "1Gi"  # Increase from 512Mi
     requests:
       memory: "512Mi"
   ```

2. **Too Many Concurrent Requests**
   - Increase rate limiting
   - Add more replicas
   - Implement connection limits

---

## Monitoring & Debugging

### Metrics Not Available

**Symptom**: `/metrics` endpoint returns 404 or no data

**Solutions**:

1. **Metrics Disabled**
   ```bash
   # Enable metrics
   export ENABLE_METRICS=true
   ```

2. **Prometheus Not Scraping**
   - Check Prometheus configuration
   - Verify service discovery
   - Check network connectivity

### Traces Not in Jaeger

**Symptom**: No traces appearing in Jaeger UI

**Solutions**:

1. **Tracing Disabled**
   ```bash
   export ENABLE_TRACING=true
   export JAEGER_URL=http://jaeger:14268/api/traces
   ```

2. **Wrong Jaeger URL**
   - Verify Jaeger collector endpoint
   - Check Jaeger logs for errors

3. **Sampling Rate Too Low**
   - Default sampling: 100% (all requests traced)
   - Check tracing configuration

### Logs Not Showing

**Symptom**: Missing or incomplete logs

**Solutions**:

1. **Wrong Log Level**
   - Set appropriate log level
   - Default: INFO

2. **Logs Going to Wrong Output**
   - Gateway logs to stdout/stderr by default
   - Check container logs: `docker logs api-gateway`

3. **Log Aggregation Issues**
   - Check Logstash/Fluentd configuration
   - Verify Elasticsearch is receiving logs

---

## Debug Mode

Enable debug mode for detailed troubleshooting:

```bash
# Environment variables for debugging
export LOG_LEVEL=debug
export ENABLE_TRACING=true
export ENABLE_METRICS=true

# Access debug endpoints
curl http://localhost:8080/debug/pprof/
curl http://localhost:8080/debug/pprof/goroutine?debug=1
curl http://localhost:8080/debug/pprof/heap?debug=1
```

---

## Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| `dial tcp: lookup auth-service: no such host` | DNS resolution failure | Check service name and DNS configuration |
| `context deadline exceeded` | Request timeout | Increase timeout or check backend performance |
| `connection refused` | Service not running or wrong port | Verify service is up and port is correct |
| `invalid or expired token` | JWT validation failed | Check JWT secret and token expiration |
| `circuit breaker is open` | Backend service unhealthy | Wait for circuit to close or fix backend |
| `rate limit exceeded` | Too many requests | Reduce request rate or increase limits |
| `missing go.sum entry` | Dependency issue | Run `go mod tidy` |

---

## Getting Help

If you can't resolve the issue:

1. **Check Logs**: 
   ```bash
   docker logs api-gateway --tail=100 --follow
   kubectl logs -f deployment/api-gateway
   ```

2. **Collect Diagnostics**:
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Metrics
   curl http://localhost:8080/metrics
   
   # Goroutine dump
   curl http://localhost:8080/debug/pprof/goroutine?debug=2
   ```

3. **Check Monitoring**:
   - Prometheus metrics
   - Jaeger traces
   - Grafana dashboards
   - Elasticsearch logs

4. **Create GitHub Issue**:
   - Include gateway version
   - Relevant logs
   - Configuration (redact secrets!)
   - Steps to reproduce

---

## Performance Tuning Tips

1. **Enable Redis Caching**
   - Significantly reduces backend load
   - Configure appropriate TTL

2. **Tune Rate Limits**
   - Set per-tenant limits if possible
   - Monitor actual usage patterns

3. **Optimize Circuit Breaker**
   - Tune thresholds for your services
   - Monitor failure rates

4. **Connection Pooling**
   - gRPC connections are pooled by default
   - Tune pool size if needed

5. **Horizontal Scaling**
   - Run multiple gateway instances
   - Use load balancer for distribution

6. **Resource Limits**
   - Set appropriate CPU/memory limits
   - Enable autoscaling

---

## Maintenance Tasks

### Restart Gateway with Zero Downtime

```bash
# Kubernetes
kubectl rollout restart deployment/api-gateway
kubectl rollout status deployment/api-gateway

# Docker Compose
docker-compose up -d --no-deps --build api-gateway
```

### Clear Redis Cache

```bash
# Clear all cached data
redis-cli -u $REDIS_URL FLUSHDB

# Clear specific pattern
redis-cli -u $REDIS_URL --scan --pattern "cache:*" | xargs redis-cli DEL
```

### Rotate JWT Secret

1. Generate new secret
2. Update both Auth Service and API Gateway
3. Rolling restart (users will need to re-authenticate)

### Update Gateway Version

```bash
# Pull latest image
docker pull ghcr.io/vhvcorp/go-api-gateway:latest

# Rolling update
kubectl set image deployment/api-gateway \
  api-gateway=ghcr.io/vhvcorp/go-api-gateway:latest

# Monitor rollout
kubectl rollout status deployment/api-gateway
```

---

For more information, see:
- [README.md](README.md) - General documentation
- [Architecture Diagrams](docs/diagrams/) - System architecture
- [GitHub Issues](https://github.com/vhvcorp/go-api-gateway/issues) - Known issues
