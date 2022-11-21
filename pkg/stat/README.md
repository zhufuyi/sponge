## stat

Statistics on system and process cpu and memory information.

<br>

### Example of use

```go
	l, _ := zap.NewDevelopment()
    stat.Init(
        WithLog(l),
        WithPrintInterval(time.Minute),
    )
```
