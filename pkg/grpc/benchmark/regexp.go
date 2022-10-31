package benchmark

import (
	"regexp"
	"strings"
)

const (
	packagePattern = `\npackage (.*);`
	servicePattern = `\nservice (\w+)`
	methodPattern  = `rpc (\w+)`
)

func getName(data []byte, pattern string) string {
	re := regexp.MustCompile(pattern)
	matchArr := re.FindStringSubmatch(string(data))
	if len(matchArr) == 2 {
		return strings.ReplaceAll(matchArr[1], " ", "")
	}
	return ""
}

func getMethodNames(data []byte, methodPattern string) []string {
	re := regexp.MustCompile(methodPattern)
	matchArr := re.FindAllStringSubmatch(string(data), -1)
	names := []string{}
	for _, arr := range matchArr {
		if len(arr) == 2 {
			names = append(names, strings.ReplaceAll(arr[1], " ", ""))
		}
	}

	return names
}

// 匹配名称，不区分大小写
func matchName(names []string, name string) string {
	out := ""
	for _, s := range names {
		if strings.EqualFold(s, name) {
			out = s
			break
		}
	}
	return out
}
