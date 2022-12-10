## stat

Statistics on system and process cpu and memory information, alarm notification support.

<br>

### Example of use

```go
	l, _ := zap.NewDevelopment()
    stat.Init(
        WithLog(l),
        WithPrintInterval(time.Minute),
        WithEnableAlarm(WithCPUThreshold(0.9), WithMemoryThreshold(0.85)), // invalid if it is windows
    )
```
