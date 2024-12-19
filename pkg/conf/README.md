## conf

Parsing yaml, json, toml configuration files to go struct.

<br>

### Example of use

```go
    import "github.com/go-dev-frame/sponge/pkg/conf"

    // Way 1: No listening configuration file
    config := &App{}
    err := conf.Parse("test.yml", config)

    // Way 2: Enable listening configuration file
    config := &App{}
    reloads  := []func(){
        func() {
            fmt.Println("close and reconnect mysql")
            fmt.Println("close and reconnect redis")
        },
    }
    err := conf.Parse("test.yml", config, reloads...)
```
