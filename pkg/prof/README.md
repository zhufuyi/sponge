## prof

Wrap the official `net/http/pprof` route and add the profile io wait time route.

<br>

### Example of use

#### sampling profile by http

```go
    import "github.com/zhufuyi/sponge/pkg/prof"

    mux := http.NewServeMux()
    prof.Register(mux, prof.WithPrefix("/myServer"), prof.WithIOWaitTime())

    httpServer := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
	
    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic("listen and serve error: " + err.Error())
    }
```

<br>

#### sampling profile by system notification signal

```go
import "github.com/zhufuyi/sponge/pkg/prof"

func WaitSign() {
	p := prof.NewProfile()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTRAP)

	for {
		v := <-signals
		switch v {
		case syscall.SIGTRAP:
			p.StartOrStop()

		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
			os.Exit(0)
		}
	}
}
```

```bash
# view the program's pid
ps -aux | grep serverName

# notification of sampling profile, default 60s, in less than 60s, if the second execution will actively stop sampling profile
kill -trap pid
```
