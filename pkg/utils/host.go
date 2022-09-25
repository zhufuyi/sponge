package utils

import (
	"fmt"
	"net"
	"os"
)

// GetHostname 获取主机名
func GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	return name
}

// GetLocalHTTPAddrPairs 获取可用的http服务端和请求地址
func GetLocalHTTPAddrPairs() (string, string) {
	port, err := GetAvailablePort()
	if err != nil {
		fmt.Printf("GetAvailablePort error: %v\n", err)
		return "", ""
	}
	serverAddr := fmt.Sprintf(":%d", port)
	requestAddr := fmt.Sprintf("http://localhost:%d", port)
	return serverAddr, requestAddr
}

// GetAvailablePort 获取可用端口
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}
