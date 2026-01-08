package client

import (
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TenantClient handles communication with tenant service
type TenantClient struct {
	conn *grpc.ClientConn
	log  *logger.Logger
	// client proto.TenantServiceClient // Uncomment when proto is generated
}

// NewTenantClient creates a new tenant client with retry logic and mTLS support
func NewTenantClient(serviceURL string, log *logger.Logger, tlsCfg *TLSConfig) *TenantClient {
	conn, err := NewGRPCConnection(serviceURL, log, tlsCfg)
	if err != nil {
		log.Error("Failed to connect to tenant service", zap.Error(err), zap.String("url", serviceURL))
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
