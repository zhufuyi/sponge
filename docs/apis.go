package docs

import (
	"embed"
	"fmt"
)

//go:embed apis.swagger.json
var jsonFile embed.FS

// ApiDocs swagger json file content
var ApiDocs = []byte(``)

func init() {
	data, err := jsonFile.ReadFile("apis.swagger.json")
	if err != nil {
		fmt.Printf("\nReadFile error: %v\n\n", err)
		return
	}
	ApiDocs = data
}
