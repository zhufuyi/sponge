## handlerfunc

Commonly used public handlers.

<br>

## Example of use

```go
	r := gin.New()
	r.GET("/health", handlerfunc.CheckHealth)
	r.GET("/ping", handlerfunc.Ping)
```