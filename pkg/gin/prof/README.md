## prof

Wrap the official `net/http/pprof` route and add the profile io wait time route.

<br>

### Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/gin/prof"

	r := gin.Default()
	prof.Register(r, WithPrefix("/myServer"), WithIOWaitTime())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	
    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic("listen and serve error: " + err.Error())
    }
```
