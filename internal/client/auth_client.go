package client

import (
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient handles communication with auth service
type AuthClient struct {
	conn   *grpc.ClientConn
	log    *logger.Logger
	// client proto.AuthServiceClient // Uncomment when proto is generated
}

// NewAuthClient creates a new auth client
func NewAuthClient(serviceURL string, log *logger.Logger) *AuthClient {
	conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to auth service", "error", err, "url", serviceURL)
		// For now, continue without panic to allow graceful degradation
	}

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
