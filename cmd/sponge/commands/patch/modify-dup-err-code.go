package patch

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// ModifyDuplicateErrCodeCommand modify duplicate error codes
func ModifyDuplicateErrCodeCommand() *cobra.Command {
	var (
		dir string
	)

	cmd := &cobra.Command{
		Use:   "modify-dup-err-code",
		Short: "Modify duplicate error codes",
		Long: `modify duplicate error codes 

Examples:
  # modify duplicate error codes
  sponge patch modify-dup-err-code --dir=internal/ecode

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := listErrCodeFiles(dir)
			if err != nil {
				return err
			}

			for _, file := range files {
				ecsis, err := parseErrCodeInfo(file)
				if err != nil {
					return err
				}
				for _, ecsi := range ecsis {
					msg, err := ecsi.modifyHTTPDuplicateNum()
					if err != nil {
						return err
					}
					if msg != "" {
						fmt.Println("modify http duplicate error codes: ", msg)
					}
					msg, err = ecsi.modifyGRPCDuplicateNum()
					if err != nil {
						return err
					}
					if msg != "" {
						fmt.Println("modify grpc duplicate error codes: ", msg)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "internal/ecode", "input directory")

	return cmd
}

type eCodeInfo struct {
	Name string
	Num  int
	Str  string
}

type errCodesInfo struct {
	file string

	httpErrCodeInfo     map[string]eCodeInfo // map[name]eCodeInfo
	httpDuplicationNums map[int][]string
	httpMaxNum          int

	grpcErrCodeInfo     map[string]eCodeInfo // map[name]eCodeInfo
	grpcDuplicationNums map[int][]string
	grpcMaxNum          int
}

func (e *errCodesInfo) getHTTPMaxNum() int {
	maxNum := 0
	for num := range e.httpDuplicationNums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

func (e *errCodesInfo) getGRPCMaxNum() int {
	maxNum := 0
	for num := range e.grpcDuplicationNums {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

func (e *errCodesInfo) modifyHTTPDuplicateNum() (string, error) {
	msg := ""
	duplicateNums := []string{}

	if len(e.httpDuplicationNums) == 0 {
		return msg, nil
	}

	numMap := map[int]struct{}{}
	for num := range e.httpDuplicationNums {
		numMap[num] = struct{}{}
	}

	e.httpMaxNum = e.getHTTPMaxNum()
	for _, names := range e.httpDuplicationNums {
		if len(names) <= 1 {
			continue
		}

		for i, name := range names {
			if i == 0 {
				continue
			}

			eci := e.httpErrCodeInfo[name]
			e.httpMaxNum++
			newNum := e.httpMaxNum

			_, err := updateErrCodeFile(e.file, newNum, eci)
			if err != nil {
				return msg, err
			}
			duplicateNums = append(duplicateNums, fmt.Sprintf("%d --> %d", eci.Num, newNum))
		}
	}

	if len(duplicateNums) == 0 {
		return msg, nil
	}
	return strings.Join(duplicateNums, ", "), nil
}

func (e *errCodesInfo) modifyGRPCDuplicateNum() (string, error) {
	msg := ""
	duplicateNums := []string{}

	if len(e.grpcDuplicationNums) == 0 {
		return msg, nil
	}

	numMap := map[int]struct{}{}
	for num := range e.grpcDuplicationNums {
		numMap[num] = struct{}{}
	}

	e.grpcMaxNum = e.getGRPCMaxNum()
	for _, names := range e.grpcDuplicationNums {
		if len(names) <= 1 {
			continue
		}

		for i, name := range names {
			if i == 0 {
				continue
			}

			eci := e.grpcErrCodeInfo[name]
			e.grpcMaxNum++
			newNum := e.grpcMaxNum

			_, err := updateErrCodeFile(e.file, newNum, eci)
			if err != nil {
				return msg, err
			}
			duplicateNums = append(duplicateNums, fmt.Sprintf("%d --> %d", eci.Num, newNum))
		}
	}

	if len(duplicateNums) == 0 {
		return msg, nil
	}
	return strings.Join(duplicateNums, ", "), nil
}

func parseErrCodeInfo(file string) ([]*errCodesInfo, error) {
	errCodeType := ""
	ecsis := []*errCodesInfo{}

	data, err := os.ReadFile(file)
	if err != nil {
		return ecsis, err
	}
	dataStr := string(data)
	if strings.Contains(dataStr, "errcode.NewError") {
		errCodeType = httpType
	} else if strings.Contains(dataStr, "errcode.NewRPCStatus") {
		errCodeType = grpcType
	}

	if errCodeType == "" {
		return ecsis, nil
	}

	var regStr string
	if errCodeType == httpType {
		regStr = `(Err[\w\W]*?)[ ]*?=[ ]*?errcode.NewError\(([\w\W]*?)BaseCode\+(\d),`
	} else if errCodeType == grpcType {
		regStr = `(Status[\w\W]*?)[ ]*?=[ ]*?errcode.NewRPCStatus\(([\w\W]*?)BaseCode\+(\d),`
	}

	reg := regexp.MustCompile(regStr)
	allSubMatch := reg.FindAllStringSubmatch(dataStr, -1)
	if len(allSubMatch) == 0 {
		return ecsis, nil
	}

	groupNames := make(map[string][][]string)
	for _, match := range allSubMatch {
		if len(match) == 4 {
			gns, ok := groupNames[match[2]]
			if ok {
				gns = append(gns, match)
			} else {
				gns = [][]string{match}
			}
			groupNames[match[2]] = gns
		}
		continue
	}

	for _, gn := range groupNames {
		ecsi := &errCodesInfo{}
		eci := make(map[string]eCodeInfo)
		duplicationNums := make(map[int][]string)
		for _, match := range gn {
			if len(match) == 4 {
				num, _ := strconv.Atoi(match[3])
				if num == 0 {
					continue
				}

				if names, ok := duplicationNums[num]; ok {
					duplicationNums[num] = append(names, match[1])
				} else {
					duplicationNums[num] = []string{match[1]}
				}

				eci[match[1]] = eCodeInfo{Name: match[1], Num: num, Str: match[0]}
			}
		}
		if errCodeType == httpType {
			ecsi.httpDuplicationNums = duplicationNums
			ecsi.httpErrCodeInfo = eci
		} else if errCodeType == grpcType {
			ecsi.grpcDuplicationNums = duplicationNums
			ecsi.grpcErrCodeInfo = eci
		}
		ecsi.file = file
		ecsis = append(ecsis, ecsi)
	}

	return ecsis, nil
}

func updateErrCodeFile(file string, newNum int, eci eCodeInfo) (eCodeInfo, error) {
	strTmp := eci.Str
	oldNum := eci.Num
	eci.Str = replaceNumStr(strTmp, oldNum, newNum)
	eci.Num = newNum

	data, err := os.ReadFile(file)
	if err != nil {
		return eci, err
	}
	data = bytes.ReplaceAll(data, []byte(strTmp), []byte(eci.Str))

	err = os.WriteFile(file, data, 0766)
	if err != nil {
		return eci, err
	}
	return eci, nil
}

func replaceNumStr(str string, oldNum int, newNum int) string {
	oldNumStr := fmt.Sprintf("+%d", oldNum)
	newNumStr := fmt.Sprintf("+%d", newNum)
	return strings.ReplaceAll(str, oldNumStr, newNumStr)
}
