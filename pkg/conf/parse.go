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

// Parse 解析配置文件到struct，包括yaml、toml、json等文件，如果fs不为空，开启监听配置文件变化
func Parse(configFile string, obj interface{}, fs ...func()) error {
	confFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	ext := strings.TrimLeft(path.Ext(filename), ".")
	filename = strings.ReplaceAll(filename, "."+ext, "") // 不包括后缀名

	viper.AddConfigPath(filePathStr) // 路径
	viper.SetConfigName(filename)    // 名称
	viper.SetConfigType(ext)         // 从文件名中获取配置类型
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

// 监听配置文件更新
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

// Show 打印配置信息(去掉敏感信息)
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

// 替换密码
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

// 替换dsn的密码
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
