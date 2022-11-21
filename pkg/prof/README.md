## prof

Wrap the official `net/http/pprof` route and add the profile io wait time route.

<br>

### Example of use

```go
	mux := http.NewServeMux()
    prof.Register(r, WithPrefix("/myServer"), WithIOWaitTime())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic("listen and serve error: " + err.Error())
    }
```
