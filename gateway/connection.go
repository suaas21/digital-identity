package main

import (
	"crypto/x509"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"path"
	"time"
)

func InitFabricConnection() (*FabricClient, error) {
	// Load connection parameters from environment
	channelName := envOrDefault("CHANNEL_NAME", "mychannel")
	chaincodeName := envOrDefault("CHAINCODE_NAME", "identity")
	mspID := envOrDefault("MSP_ID", "Org1MSP")

	// Load crypto material
	certPath := envOrDefault("CERT_PATH", "/etc/secret-volume/certPath")
	keyPath := envOrDefault("KEY_PATH", "/etc/secret-volume/keyPath")
	tlsCertPath := envOrDefault("TLS_CERT_PATH", "/etc/secret-volume/tlsCertPath")
	peerEndpoint := envOrDefault("PEER_ENDPOINT", "test-network-org1-peer1-peer.localho.st:443")

	// Create gRPC client connection
	clientConnection, err := newGrpcConnection(tlsCertPath, peerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	id, err := newIdentity(certPath, mspID)
	if err != nil {
		return nil, err
	}
	sign, err := newSign(keyPath)
	if err != nil {
		return nil, err
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	return &FabricClient{
		gateway:  gw,
		contract: contract,
	}, nil
}

func newGrpcConnection(tlsCertPath, peerEndpoint string) (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, peerEndpoint)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity(certPath, mspId string) (*identity.X509Identity, error) {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		return nil, err
	}

	return identity.NewX509Identity(mspId, certificate)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign(keyPath string) (identity.Sign, error) {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key directory: %w", err)
	}

	if len(files) <= 0 {
		return nil, fmt.Errorf("no file name found to load")
	}

	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

func envOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
