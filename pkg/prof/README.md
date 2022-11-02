## prof

自定义官方`net/http/pprof`路由。

<br>

### 使用示例

```go
	mux := http.NewServeMux()
	Register(mux, WithPrefix(""), WithPrefix("/myServer"))

	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: ":8080",
	}
	
    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic("listen and serve error: " + err.Error())
    }
```
