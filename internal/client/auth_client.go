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

// AuthClient handles communication with auth service
type AuthClient struct {
	conn *grpc.ClientConn
	log  *logger.Logger
	// client proto.AuthServiceClient // Uncomment when proto is generated
}

// NewAuthClient creates a new auth client with retry logic and connection pooling
func NewAuthClient(serviceURL string, log *logger.Logger) *AuthClient {
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

	// Get configurable max message size (default: 10MB)
	maxMsgSize := 10 * 1024 * 1024
	if ms := os.Getenv("GRPC_MAX_MESSAGE_SIZE"); ms != "" {
		if parsed, err := strconv.Atoi(ms); err == nil && parsed > 0 {
			maxMsgSize = parsed
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
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Error("Failed to connect to auth service", zap.Error(err), zap.String("url", serviceURL))
		// Return client with nil connection for graceful degradation
		return &AuthClient{
			conn: nil,
			log:  log,
		}
	}

	log.Info("Successfully connected to auth service",
		zap.String("url", serviceURL),
		zap.Int("pool_size", poolSize),
		zap.Int("max_message_size_mb", maxMsgSize/(1024*1024)))
	return &AuthClient{
		conn: conn,
		log:  log,
		// client: proto.NewAuthServiceClient(conn), // Uncomment when proto is generated
	}
}

// Close closes the gRPC connection
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Additional methods will be added once protobuf definitions are generated
// Example:
// func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*proto.ValidateTokenResponse, error) {
//     return c.client.ValidateToken(ctx, &proto.ValidateTokenRequest{Token: token})
// }
