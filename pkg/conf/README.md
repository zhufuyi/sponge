## conf

Parsing yaml, json, toml configuration files to go struct.

<br>

### Example of use

```go
    // Way 1: No listening profile
	conf := &App{}
	err := Parse("test.yml", conf)

    // Way 2: Enable listening profile
	conf := &App{}
	fs := []func(){
		func() {
			fmt.Println("Listening for updates to the configuration file")
		},
	}
	err := Parse("test.yml", conf, fs...)
```
