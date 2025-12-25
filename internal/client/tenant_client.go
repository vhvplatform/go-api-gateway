package client

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/vhvcorp/go-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TenantClient handles communication with tenant service
type TenantClient struct {
	conn   *grpc.ClientConn
	log    *logger.Logger
	// client proto.TenantServiceClient // Uncomment when proto is generated
}

// NewTenantClient creates a new tenant client with retry logic
func NewTenantClient(serviceURL string, log *logger.Logger) *TenantClient {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(3),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
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

	log.Info("Successfully connected to tenant service", zap.String("url", serviceURL))
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
