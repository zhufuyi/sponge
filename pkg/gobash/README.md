## gobash

在go环境中执行命令、脚本、可执行文件，日志实时输出。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gobash

<br>

## 使用示例

### Run

Run执行命令，可以主动结束命令，实时返回日志和错误信息，推荐使用

```go

    command := "for i in $(seq 1 5); do echo 'test cmd' $i;sleep 1; done"
    ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) // 超时控制
	
    // 执行
    result := Run(ctx, command)
    // 实时输出日志和错误信息
    for v := range result.StdOut {
        fmt.Printf(v)
    }
    if result.Err != nil {
        fmt.Println("exec command failed,", result.Err.Error())
    }
```

<br>

### Exec

Exec 适合执行单条非阻塞命令，输出标准和错误日志，但日志输出不是实时，注：如果执行命令永久阻塞，会造成协程泄露

```go
    command := "for i in $(seq 1 5); do echo 'test cmd' $i;sleep 1; done"
    out, err := gobash.Exec(command)
    if err != nil {
        return
    }
    fmt.Println(string(out))
```
