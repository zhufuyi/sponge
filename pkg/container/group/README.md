## group

连接etcd服务客户端。

### 使用示例

```go
    type foo struct {
        bar string
    }
    
    gr := group.NewGroup(func () interface{} {
        return &foo{"hello"}
    })

	fmt.Println(gr.Get(*foo).bar)
```
