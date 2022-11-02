## stat

统计系统和进程的cpu和内存信息。

<br>

### 使用示例

```go
	l, _ := zap.NewDevelopment()
    stat.Init(
        WithLog(l),
        WithPrintInterval(time.Minute),
    )
```
