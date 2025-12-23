package client

import (
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient handles communication with user service
type UserClient struct {
	conn   *grpc.ClientConn
	log    *logger.Logger
	// client proto.UserServiceClient // Uncomment when proto is generated
}

// NewUserClient creates a new user client
func NewUserClient(serviceURL string, log *logger.Logger) *UserClient {
	conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to user service", "error", err, "url", serviceURL)
	}

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
