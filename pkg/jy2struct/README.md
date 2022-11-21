## jy2struct

A library for generating go struct code, supporting json and yaml.

<br>

### Example of use

Main setting parameters.

```go
type Args struct {
	Format    string // document format, json or yaml
	Data      string // json or yaml content
	InputFile string // file
	Name      string // name of structure
	SubStruct bool   // are sub-structures separated
	Tags      string // add additional tags, multiple tags separated by commas
}
```

<br>

Example of conversion.

```go
    // json covert to struct
    code, err := jy2struct.Covert(&jy2struct.Args{
        Format: "json",
        // InputFile: "user.json", // source from json file
        SubStruct: true,
    })

    // yaml covert to struct
    code, err := jy2struct.Covert(&jy2struct.Args{
        Format: "yaml",
        // InputFile: "user.yaml", // Source from yaml file
        SubStruct: true,
    })
```
