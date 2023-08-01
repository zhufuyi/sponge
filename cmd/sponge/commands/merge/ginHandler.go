package merge

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var httpECodeMark = "errcode.NewError("
var routerMark = "c.setSinglePath("

// GinHandlerCode merge the gin handler code
func GinHandlerCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gin-handler",
		Short: "Merge the gin handler code",
		Long: `merge the gin handler code.

Examples:
  # merge gin handler code
  sponge merge gin-handler
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dt := time.Now().Format("20060102T150405")

			runMerge(internalECodeDir, dt, true, parseFromECode)
			runMerge(internalRoutersDir, dt, true, parseFromRouters)
			runMerge(internalHandlerDir, dt, false, parseFromHandler)

			return nil
		},
	}

	return cmd
}

func parseFromECode(date []byte) []code {
	var codes []code
	buf := bufio.NewReader(bytes.NewReader(date))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, httpECodeMark) {
			name := getECodeMarkName(line)
			if name != "" {
				codes = append(codes, code{
					key:   name,
					value: line,
				})
			}
		}
	}

	return codes
}

func getECodeMarkName(str string) string {
	ss := strings.SplitN(str, httpECodeMark, 2)
	name := strings.Replace(ss[0], " ", "", -1)
	name = strings.Replace(name, "=", "", -1)
	return strings.Replace(name, "	", " ", -1)
}

// ------------------------------------------------------------------------------------------

func parseFromRouters(date []byte) []code {
	var codes []code
	buf := bufio.NewReader(bytes.NewReader(date))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, routerMark) {
			name := getRoutersMarkName(line)
			if name != "" {
				codes = append(codes, code{
					key:   name,
					value: line,
				})
			}
		}
	}

	return codes
}

func getRoutersMarkName(str string) string {
	str = strings.Replace(str, " ", "", -1)
	ss := strings.SplitN(str, ",", 3)
	if len(ss) != 3 {
		return ""
	}

	sss := strings.Split(ss[0], routerMark)
	if len(sss) != 2 {
		return ""
	}
	method := strings.Replace(sss[1], "(", "", -1)
	method = strings.Replace(method, "\"", "", -1)

	router := strings.Replace(ss[1], "\"", "", -1)

	return method + "-->" + router
}

//func runMergeHandler(dir string, dt string) {
//	files, err := parseFiles(dir)
//	if err != nil {
//		fmt.Println("Warring:", err)
//	}
//
//	for _, file := range files {
//		successFile, err := runMergeCode(file, dt, parseTmplCode)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		if successFile != "" {
//			fmt.Printf("merge code to '%s' successfully.\n", successFile)
//		}
//	}
//}
//
//func mergeHandlerData(subData1 []byte, subData2 []byte) []byte {
//	r1 := parseHandler(subData1)
//	r2 := parseHandler(subData2)
//	addCode, mark := compareCode2(r1, r2, subData2)
//	return mergeCode(subData1, addCode, mark)
//}

func parseFromHandler(date []byte) []code {
	return parseTmplCode(date)
}
