package jy2struct

import (
	"bytes"
	"errors"
	"os"
	"strings"
)

// Args  参数
type Args struct {
	Format    string // 文档格式，json或yaml
	Data      string // json或yaml内容
	InputFile string // 文件
	Name      string // 结构体名称
	SubStruct bool   // 子结构体是否分开
	Tags      string // 字段tag，多个tag用逗号分隔

	tags          []string
	convertFloats bool
	parser        Parser
}

func (j *Args) checkValid() error {
	switch j.Format {
	case "json":
		j.parser = ParseJSON
		j.convertFloats = true
	case "yaml":
		j.parser = ParseYaml
	default:
		return errors.New("format must be json or yaml")
	}

	j.tags = []string{j.Format}
	tags := strings.Split(j.Tags, ",")
	for _, tag := range tags {
		if tag == j.Format || tag == "" {
			continue
		}
		j.tags = append(j.tags, tag)
	}

	if j.Name == "" {
		j.Name = "GenerateName"
	}

	return nil
}

// Covert json或yaml转go struct
func Covert(args *Args) (string, error) {
	err := args.checkValid()
	if err != nil {
		return "", err
	}

	var data []byte
	if args.Data != "" {
		data = []byte(args.Data)
	} else {
		// 读取文件
		data, err = os.ReadFile(args.InputFile)
		if err != nil {
			return "", err
		}
	}

	input := bytes.NewReader(data)

	output, err := jyParse(input, args.parser, args.Name, "main", args.tags, args.SubStruct, args.convertFloats)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
