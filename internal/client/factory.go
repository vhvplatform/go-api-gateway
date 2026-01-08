package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// TLSConfig holds the configuration for mTLS
type TLSConfig struct {
	Enabled        bool
	CACertFile     string
	ClientCertFile string
	ClientKeyFile  string
	ServerName     string // Override server name for testing/mismatched certs
}

// NewGRPCConnection creates a new gRPC connection with optional mTLS and retries
func NewGRPCConnection(target string, log *logger.Logger, tlsCfg *TLSConfig) (*grpc.ClientConn, error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(3),
	}

	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithBlock(), // Wait for connection to be established
	}

	if tlsCfg != nil && tlsCfg.Enabled {
		creds, err := loadTLSCredentials(tlsCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
		log.Info("mTLS enabled for gRPC connection", zap.String("target", target))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		log.Warn("Using insecure gRPC connection", zap.String("target", target))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", target, err)
	}

	return conn, nil
}

func loadTLSCredentials(cfg *TLSConfig) (credentials.TransportCredentials, error) {
	// Load Certificate Authority
	pemServerCA, err := os.ReadFile(cfg.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to append CA cert")
	}

	// Load Client Cert and Key
	clientCert, err := tls.LoadX509KeyPair(cfg.ClientCertFile, cfg.ClientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert/key: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		ServerName:   cfg.ServerName, 
		// If ServerName is empty, it uses the target hostname. 
		// Set InsecureSkipVerify if needed for dev (not recommended for prod)
	}

	return credentials.NewTLS(tlsConfig), nil
}
