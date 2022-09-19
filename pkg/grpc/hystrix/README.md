## hystrix

### 使用示例

#### grpc client

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 熔断拦截器
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			hystrix.UnaryClientInterceptor("hello_grpc",
				hystrix.WithTimeout(time.Second*2),       // 执行command的超时时间，时间单位是ms，默认时间是1000ms
				hystrix.WithMaxConcurrentRequests(20),    // command的最大并发量，默认值是10并发量
				hystrix.WithSleepWindow(10*time.Second),  // 熔断器被打开后使用，在熔断器被打开后，根据SleepWindow设置的时间控制多久后尝试服务是否可用，默认时间为5000ms
				hystrix.WithRequestVolumeThreshold(1000), // 判断熔断开关的条件之一，统计10s（代码中写死了）内请求数量，达到这个请求数量后再根据错误率判断是否要开启熔断；
				hystrix.WithErrorPercentThreshold(25),    // 判断熔断开关的条件之一，统计错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断 默认值是50
			),
		),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	// ......
}
```