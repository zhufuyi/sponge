## benchmark

压测rpc方法，并生成报告结果。

### 使用示例

```go
func benchmarkExample() error {
	host := "127.0.0.1:8282"
	protoFile := "api/serverNameExample/v1/userExample.proto"
	// 如果压测过程中缺少第三方依赖，复制到项目的third_party目录下(不包括import路径)
	importPaths := []string{"third_party"}
	message := &serverNameV1.GetUserExampleByIDRequest{
		ID: 2,
	}

	b, err := benchmark.New(host, protoFile, "GetByID", message, 1000, importPaths...)
	if err != nil {
		return err
	}
	return b.Run()
}
```

压测完毕后，复制输出的html文件路径到浏览器查看详细的压测报告。
