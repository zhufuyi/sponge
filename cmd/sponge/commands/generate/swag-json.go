package generate

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

// ConvertSwagJSONCommand convert 64-bit fields type string to integer
func ConvertSwagJSONCommand(parentName string) *cobra.Command {
	var (
		jsonFile string
		isSort   bool
	)

	cmd := &cobra.Command{
		Use:   "swagger",
		Short: "Convert 64-bit fields type string to integer",
		Long: color.HiBlackString(fmt.Sprintf(`convert 64-bit fields type string to integer

Examples:
  # convent file docs/apis.swagger.json.
  sponge %s swagger

  # convent file test/swagger.json
  sponge %s swagger --file=test/swagger.json

  # convent file docs/apis.swagger.json and sort json key.
  sponge %s swagger --is-sort
`, parentName, parentName, parentName)),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if isSort {
				err = handlerJSONFormatTypeWithSortKey(jsonFile)
			} else {
				err = handlerJSONFormatType(jsonFile)
			}
			if err != nil {
				return err
			}

			fmt.Printf("convert json file successfully, out = %s\n", jsonFile)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&isSort, "is-sort", "s", false, "formatting json, json's fields are sorted in ascending")
	cmd.Flags().StringVarP(&jsonFile, "file", "f", "docs/apis.swagger.json", "input json file")

	return cmd
}

func handlerJSONFormatType(jsonFilePath string) error {
	newData, err := convertStringToInteger(jsonFilePath)
	if err != nil {
		return err
	}

	return saveJSONFile(newData, jsonFilePath)
}

func handlerJSONFormatTypeWithSortKey(jsonFilePath string) error {
	data, err := formatJSON(jsonFilePath)
	if err != nil {
		return err
	}

	newData, err := convertStringToIntegerWithSortKey(data)
	if err != nil {
		return err
	}

	return saveJSONFile(newData, jsonFilePath)
}

func convertStringToInteger(jsonFilePath string) ([]byte, error) {
	f, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	scanner := bufio.NewScanner(f)
	contents := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `"format": "uint64"`) || strings.Contains(line, `"format": "int64"`) {
			l := len(contents)
			previousLine := contents[l-1]
			if len(contents) > 0 && strings.Contains(previousLine, `"type": "string"`) {
				contents[l-1] = strings.ReplaceAll(previousLine, `"type": "string"`, `"type": "integer"`)
			}
		}
		contents = append(contents, line+"\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	newData := []byte{}
	for _, v := range contents {
		newData = append(newData, []byte(v)...)
	}

	return newData, nil
}

func formatJSON(jsonFilePath string) ([]byte, error) {
	content, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	indentedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	return indentedJSON, nil
}

func convertStringToIntegerWithSortKey(data []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	contents := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(contents) > 0 && strings.Contains(line, `"type": "string"`) {
			l := len(contents)
			previousLine := contents[l-1]
			if strings.Contains(previousLine, `"format": "uint64"`) || strings.Contains(previousLine, `"format": "int64"`) {
				if tmpLine := strings.ReplaceAll(line, `"type": "string"`, `"type": "integer"`); line != tmpLine {
					line = tmpLine
				}
			}
		}
		contents = append(contents, line+"\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	newData := []byte{}
	for _, v := range contents {
		newData = append(newData, []byte(v)...)
	}

	return newData, nil
}

func saveJSONFile(data []byte, jsonFilePath string) error {
	if gofile.IsExists(jsonFilePath) {
		tmpFile := jsonFilePath + ".tmp"
		err := os.WriteFile(tmpFile, data, 0666)
		if err != nil {
			return err
		}
		return os.Rename(tmpFile, jsonFilePath)
	}

	dir := gofile.GetFileDir(jsonFilePath)
	_ = os.MkdirAll(dir, 0766)
	return os.WriteFile(jsonFilePath, data, 0666)
}
