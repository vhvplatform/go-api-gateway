package client

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient handles communication with user service
type UserClient struct {
	conn *grpc.ClientConn
	log  *logger.Logger
	// client proto.UserServiceClient // Uncomment when proto is generated
}

// NewUserClient creates a new user client with retry logic
func NewUserClient(serviceURL string, log *logger.Logger) *UserClient {
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
		log.Error("Failed to connect to user service", zap.Error(err), zap.String("url", serviceURL))
		// Return client with nil connection for graceful degradation
		return &UserClient{
			conn: nil,
			log:  log,
		}
	}

	log.Info("Successfully connected to user service", zap.String("url", serviceURL))
	return &UserClient{
		conn: conn,
		log:  log,
		// client: proto.NewUserServiceClient(conn), // Uncomment when proto is generated
	}
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
