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

// AuthClient handles communication with auth service
type AuthClient struct {
	conn *grpc.ClientConn
	log  *logger.Logger
	// client proto.AuthServiceClient // Uncomment when proto is generated
}

// NewAuthClient creates a new auth client with retry logic
func NewAuthClient(serviceURL string, log *logger.Logger) *AuthClient {
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
		log.Error("Failed to connect to auth service", zap.Error(err), zap.String("url", serviceURL))
		// Return client with nil connection for graceful degradation
		return &AuthClient{
			conn: nil,
			log:  log,
		}
	}

	log.Info("Successfully connected to auth service", zap.String("url", serviceURL))
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
