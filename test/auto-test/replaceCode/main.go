package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	paramFile := os.Args[1]
	targetFile := os.Args[2]
	fmt.Printf("parameter file: %s, procces file:%s\n", paramFile, targetFile)
	data, err := os.ReadFile(paramFile)
	if err != nil {
		panic(err)
	}

	ss := strings.Split(string(data), "|-|-|-|-|-|")
	if len(ss) == 0 || len(ss)%2 == 1 {
		panic(fmt.Sprintf(`%sThe file content format does not meet the specification, example:
This is the original content|-|-|-|-|-|This is the replacement content.
`, paramFile))
	}

	data, err = os.ReadFile(targetFile)
	if err != nil {
		panic(err)
	}

	var newData = string(data)
	count, total := 0, len(ss)/2
	for i := 0; i < len(ss); i += 2 {
		srcStr := ss[i]
		dstStr := ss[i+1]
		if strings.Contains(newData, srcStr) {
			count++
			newData = strings.ReplaceAll(newData, srcStr, dstStr)
		}
	}

	if count > 0 {
		err = os.WriteFile(targetFile, []byte(newData), 0666)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Successfully replaced the %d group of strings, for a total of %d groups.\n", count, total)
		return
	}

	time.Sleep(time.Millisecond * 100)
}
