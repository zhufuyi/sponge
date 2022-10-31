## handlerfunc

常用公共的handler。

<br>

## 使用示例

```go
	r := gin.New()
	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)
```