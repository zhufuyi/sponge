package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"google.golang.org/grpc/credentials"
)

// GetServerTLSCredentialsByCA two-way authentication via CA-issued root certificate
func GetServerTLSCredentialsByCA(caFile string, certFile string, keyFile string) (credentials.TransportCredentials, error) {
	//read and parse the information from the certificate file to obtain the certificate public key, key pair
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

	//attempts to parse the incoming PEM-encoded certificate. If the parsing is successful it will be added to the CertPool for later use
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("certPool.AppendCertsFromPEM err")
	}

	//building TLS-based TransportCredentials options
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},        // set up a certificate chain that allows the inclusion of one or more
		ClientAuth:   tls.RequireAndVerifyClientCert, // requirement to verify the client's certificate
		ClientCAs:    certPool,                       // set the set of root certificates and use the mode set in ClientAuth for verification
	})

	return c, err
}

// GetServerTLSCredentials server-side authentication
func GetServerTLSCredentials(certFile string, keyFile string) (credentials.TransportCredentials, error) {
	c, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return c, err
}
