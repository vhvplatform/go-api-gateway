package client

import (
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TenantClient handles communication with tenant service
type TenantClient struct {
	conn   *grpc.ClientConn
	log    *logger.Logger
	// client proto.TenantServiceClient // Uncomment when proto is generated
}

// NewTenantClient creates a new tenant client
func NewTenantClient(serviceURL string, log *logger.Logger) *TenantClient {
	conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to tenant service", "error", err, "url", serviceURL)
	}

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
