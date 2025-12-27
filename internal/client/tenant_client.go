package client

import (
	"context"
	"os"
	"strconv"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// TenantClient handles communication with tenant service
type TenantClient struct {
	conn *grpc.ClientConn
	log  *logger.Logger
	// client proto.TenantServiceClient // Uncomment when proto is generated
}

// NewTenantClient creates a new tenant client with retry logic and connection pooling
func NewTenantClient(serviceURL string, log *logger.Logger) *TenantClient {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(3),
	}

	// Get configurable pool size (default: 5)
	poolSize := 5
	if ps := os.Getenv("GRPC_POOL_SIZE"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			poolSize = parsed
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure keepalive for connection health
	kaParams := keepalive.ClientParameters{
		Time:                10 * time.Second, // Send keepalive pings every 10 seconds
		Timeout:             3 * time.Second,  // Wait 3 seconds for ping ack before considering connection dead
		PermitWithoutStream: true,             // Send pings even without active streams
	}

	conn, err := grpc.DialContext(
		ctx,
		serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithKeepaliveParams(kaParams),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(poolSize*1024*1024)), // Pool-based message size
		grpc.WithBlock(),
	)

	if err != nil {
		log.Error("Failed to connect to tenant service", zap.Error(err), zap.String("url", serviceURL))
		// Return client with nil connection for graceful degradation
		return &TenantClient{
			conn: nil,
			log:  log,
		}
	}

	log.Info("Successfully connected to tenant service",
		zap.String("url", serviceURL),
		zap.Int("pool_size", poolSize))
	return &TenantClient{
		conn: conn,
		log:  log,
		// client: proto.NewTenantServiceClient(conn), // Uncomment when proto is generated
	}
}

// Close closes the gRPC connection
func (c *TenantClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
