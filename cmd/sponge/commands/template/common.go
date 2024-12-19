// Package template provides commands to generate custom code.
package template

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-dev-frame/sponge/pkg/gofile"
)

func parseFields(jsonFile string) (map[string]interface{}, error) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	return m, err
}

func mergeFields(m1 map[string]interface{}, m2 map[string]interface{}) (map[string]interface{}, error) {
	if m2 == nil {
		return m1, nil
	}

	for k, v := range m2 {
		if _, ok := m1[k]; ok {
			return nil, fmt.Errorf("'%s' is a reserved field, please change it to another name", k)
		}
		m1[k] = v
	}
	return m1, nil
}

func listTemplateFiles(builder *strings.Builder, files []string) {
	builder.WriteString("\nTemplate files:\n")
	for _, file := range files {
		builder.WriteString("    " + file + "\n")
	}
}

func listFields(builder *strings.Builder, fields map[string]interface{}) {
	jsonData, err := json.MarshalIndent(fields, "    ", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	builder.WriteString("    " + string(jsonData) + "\n\n")
}

var regPackage = regexp.MustCompile(`(?m)^package\s+([a-zA-Z0-9._]+);`)

func copyProtoFileToDir(protoFile string, targetDir string) error {
	data, err := os.ReadFile(protoFile)
	if err != nil {
		return err
	}

	protoPackage := ""
	matches := regPackage.FindStringSubmatch(string(data))
	if len(matches) == 2 {
		protoPackage = matches[1]
	}

	if protoPackage == "" {
		return fmt.Errorf("package not found in %s", protoFile)
	}
	if targetDir == "" {
		return fmt.Errorf("target directory not specified")
	}
	dir := targetDir + "/" + strings.ReplaceAll(protoPackage, ".", "/")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	newFile := dir + "/" + gofile.GetFilename(protoFile)
	if gofile.IsExists(newFile) {
		return nil
	}
	return os.WriteFile(newFile, data, 0644)
}
