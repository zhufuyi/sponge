## stat

Statistics on system and process cpu and memory information, alarm notification support.

<br>

### Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/stat"

    l, _ := zap.NewDevelopment()
    stat.Init(
        stat.WithLog(l),
        stat.WithPrintInterval(time.Minute),
        stat.WithEnableAlarm(stat.WithCPUThreshold(0.9), stat.WithMemoryThreshold(0.85)), // invalid if it is windows
    )
```
