package client

import (
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// UserClient handles communication with user service
type UserClient struct {
	conn *grpc.ClientConn
	log  *logger.Logger
	// client proto.UserServiceClient // Uncomment when proto is generated
}

// NewUserClient creates a new user client with retry logic and mTLS support
func NewUserClient(serviceURL string, log *logger.Logger, tlsCfg *TLSConfig) *UserClient {
	conn, err := NewGRPCConnection(serviceURL, log, tlsCfg)
	if err != nil {
		log.Error("Failed to connect to user service", zap.Error(err), zap.String("url", serviceURL))
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
