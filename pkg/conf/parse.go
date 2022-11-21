package conf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Parse configuration files to struct, including yaml, toml, json, etc., and turn on listening for configuration file changes if fs is not empty
func Parse(configFile string, obj interface{}, fs ...func()) error {
	confFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	ext := strings.TrimLeft(path.Ext(filename), ".")
	filename = strings.ReplaceAll(filename, "."+ext, "") // excluding suffix names

	viper.AddConfigPath(filePathStr) // path
	viper.SetConfigName(filename)    // file name
	viper.SetConfigType(ext)         // get the configuration type from the file name
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(obj)
	if err != nil {
		return err
	}

	if len(fs) > 0 {
		watchConfig(obj, fs...)
	}

	return nil
}

// listening for profile updates
func watchConfig(obj interface{}, fs ...func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.Unmarshal(obj)
		if err != nil {
			fmt.Println("viper.Unmarshal error: ", err)
		} else {
			for _, f := range fs {
				f()
			}
		}
	})
}

// Show print configuration information (remove sensitive information)
func Show(obj interface{}, keywords ...string) string {
	var out string

	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		fmt.Println("json.MarshalIndent error: ", err)
		return ""
	}

	buf := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		keywords = append(keywords, `"dsn"`, `"password"`)

		out += replacePWD(line, keywords...)
	}

	return out
}

// replace password
func replacePWD(line string, keywords ...string) string {
	for _, keyword := range keywords {
		if strings.Contains(line, keyword) {
			index := strings.Index(line, keyword)
			if strings.Contains(line, "@") && strings.Contains(line, ":") {
				return replaceDSN(line)
			}
			return fmt.Sprintf("%s: \"******\",\n", line[:index+len(keyword)])
		}
	}

	return line
}

// replace dsn's password
func replaceDSN(str string) string {
	mysqlPWD := []byte(str)
	start, end := 0, 0
	for k, v := range mysqlPWD {
		if v == ':' {
			start = k
		}
		if v == '@' {
			end = k
			break
		}
	}

	if start >= end {
		return str
	}

	return fmt.Sprintf("%s******%s", mysqlPWD[:start+1], mysqlPWD[end:])
}
