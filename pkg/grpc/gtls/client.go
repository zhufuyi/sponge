package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

// GetClientTLSCredentialsByCA two-way authentication via CA-issued root certificate
func GetClientTLSCredentialsByCA(serverName string, caFile string, certFile string, keyFile string) (credentials.TransportCredentials, error) {
	// read and parse the information from the certificate file to obtain the certificate public key, key pair
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// create an empty CertPool
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	// attempts to parse the incoming PEM-encoded certificate. If the parsing is successful it will be added to the CertPool for later use
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
	}

	// building TLS-based TransportCredentials options
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, // set up a certificate chain that allows the inclusion of one or more
		ServerName:   serverName,              // requirement to verify the client's certificate
		RootCAs:      certPool,
	})

	return c, err
}

// GetClientTLSCredentials TLS encryption
func GetClientTLSCredentials(serverName string, certFile string) (credentials.TransportCredentials, error) {
	c, err := credentials.NewClientTLSFromFile(certFile, serverName)
	if err != nil {
		return nil, err
	}

	return c, err
}
