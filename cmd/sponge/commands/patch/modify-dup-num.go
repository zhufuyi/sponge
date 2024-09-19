package patch

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ModifyDuplicateNumCommand modify duplicate numbers
func ModifyDuplicateNumCommand() *cobra.Command {
	var (
		dir string
	)

	cmd := &cobra.Command{
		Use:   "modify-dup-num",
		Short: "Modify duplicate numbers",
		Long: color.HiBlackString(`modify duplicate numbers 

Examples:
  # modify duplicate numbers
  sponge patch modify-dup-num --dir=internal/ecode
`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := listErrCodeFiles(dir)
			if err != nil {
				return err
			}

			count, err := checkAndModifyDuplicateNum(files)
			if err != nil {
				return err
			}
			if count > 0 {
				fmt.Println("modify duplicate num successfully.")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "internal/ecode", "input directory")

	return cmd
}

type coreInfo struct {
	name   string
	num    int
	srcStr string
	dstStr string
	file   string
}

var (
	httpNumMark = "errcode.HCode"
	grpcNumMark = "errcode.RCode"
	httpPattern = `errcode\.HCode\(([^)]+)\)`
	grpcPattern = `errcode\.RCode\(([^)]+)\)`
)

func getVariableName(data []byte, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(string(data))
	if len(match) < 2 {
		return ""
	}

	return strings.ReplaceAll(match[1], " ", "")
}

func parseNumInfo(data []byte, variableName string) coreInfo {
	var info coreInfo
	pattern := variableName + `\s*=\s*(\d+)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(string(data))
	if len(match) < 2 {
		return info
	}

	num, err := strconv.Atoi(match[1])
	if err != nil {
		return info
	}

	ss := strings.Split(match[0], "=")
	if len(ss) != 2 {
		return info
	}

	info.name = variableName
	info.num = num
	info.srcStr = match[0]
	info.dstStr = ss[0] + "= "

	return info
}

func getNumberInfos(file string) []coreInfo {
	var infos []coreInfo
	data, err := os.ReadFile(file)
	if err != nil {
		return infos
	}

	serviceGroupData := bytes.Split(data, []byte(serviceGroupSeparatorMark))
	for _, groupData := range serviceGroupData {
		pattern := ""
		if bytes.Contains(groupData, []byte(httpNumMark)) {
			pattern = httpPattern
		} else if bytes.Contains(groupData, []byte(grpcNumMark)) {
			pattern = grpcPattern
		}
		if pattern != "" {
			variableName := getVariableName(groupData, pattern)
			if variableName != "" {
				info := parseNumInfo(groupData, variableName)
				if info.name != "" {
					info.file = file
					infos = append(infos, info)
				}
			}
		}
	}

	return infos
}

func getModifyNumInfos(infos []coreInfo) ([]coreInfo, map[int]struct{}) {
	m := map[int][]coreInfo{}
	allNum := map[int]struct{}{}
	for _, info := range infos {
		allNum[info.num] = struct{}{}
		if cis, ok := m[info.num]; ok {
			m[info.num] = append(cis, info)
		} else {
			m[info.num] = []coreInfo{info}
		}
	}

	needModify := []coreInfo{}
	for _, numInfos := range m {
		if len(numInfos) > 1 {
			needModify = append(needModify, numInfos[1:]...)
		}
	}

	return needModify, allNum
}

func modifyNumberInfos(infos []coreInfo, allNum map[int]struct{}) (int, error) {
	l := len(infos)
	if l == 0 {
		return 0, nil
	}

	var nums []int
	for i := 1; i < 100; i++ {
		if _, ok := allNum[i]; !ok {
			nums = append(nums, i)
			if len(nums) == len(infos) {
				break
			}
		}
	}

	if len(nums) < l {
		for i := 0; i < l-len(nums); i++ {
			nums = append(nums, 99) // 99 is the largest number
		}
	}

	count := 0
	for i, info := range infos {
		data, err := os.ReadFile(info.file)
		if err != nil {
			return 0, err
		}

		newData := bytes.ReplaceAll(data, []byte(info.srcStr), []byte(info.dstStr+strconv.Itoa(nums[i])))

		err = os.WriteFile(info.file, newData, 0666)
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

func checkAndModifyDuplicateNum(files []string) (int, error) {
	var allInfos []coreInfo
	for _, file := range files {
		infos := getNumberInfos(file)
		if len(infos) > 0 {
			allInfos = append(allInfos, infos...)
		}
	}

	needModify, allNum := getModifyNumInfos(allInfos)

	return modifyNumberInfos(needModify, allNum)
}
