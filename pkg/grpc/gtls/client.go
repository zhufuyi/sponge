package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

// GetClientTLSCredentialsByCA 通过CA颁发的根证书，双向认证
func GetClientTLSCredentialsByCA(serverName string, caFile string, certFile string, keyFile string) (credentials.TransportCredentials, error) {
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	//尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
	}

	//构建基于 TLS 的 TransportCredentials 选项
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //设置证书链，允许包含一个或多个
		ServerName:   serverName,              // 要求必须校验客户端的证书
		RootCAs:      certPool,
	})

	return c, err
}

// GetClientTLSCredentials TLS加密
func GetClientTLSCredentials(serverName string, certFile string) (credentials.TransportCredentials, error) {
	c, err := credentials.NewClientTLSFromFile(certFile, serverName)
	if err != nil {
		return nil, err
	}

	return c, err
}
