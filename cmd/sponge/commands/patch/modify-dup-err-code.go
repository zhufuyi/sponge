package patch

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
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
		Long:  "Modify duplicate error codes.",
		Example: color.HiBlackString(`  # Modify duplicate error codes
  sponge patch modify-dup-err-code --dir=internal/ecode`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := listErrCodeFiles(dir)
			if err != nil {
				return err
			}

			var total int
			for _, file := range files {
				count, err := checkAndModifyDuplicateErrCode(file)
				if err != nil {
					return err
				}
				total += count
			}
			if total > 0 {
				fmt.Println("modify duplicate error codes successfully.")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "internal/ecode", "input directory")

	return cmd
}

type eCodeInfo struct {
	Name   string
	Num    int
	Str    string
	DstStr string
}

var (
	serviceGroupSeparatorMark = "// ---------- Do not delete or move this split line, this is the merge code marker ----------"
	defineHTTPErrCodeMark     = "errcode.NewError("
	defineGRPCErrCodeMark     = "errcode.NewRPCStatus("
)

func parseErrCodeInfo(line string) eCodeInfo {
	ci := eCodeInfo{}

	pattern := `(\w+)\s*=\s*\w+\.(.*?)\((.*?),`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(line)
	if len(match) < 4 {
		return ci
	}
	baseCodeStr := match[3]

	index := strings.Index(line, baseCodeStr)
	if index < 0 {
		return ci
	}
	srcStr := line[:index] + baseCodeStr

	ss := strings.Split(baseCodeStr, "+")
	if len(ss) != 2 {
		return ci
	}
	num, _ := strconv.Atoi(strings.TrimSpace(ss[1]))

	ci.Name = match[1]
	ci.Num = num
	ci.Str = srcStr
	ci.DstStr = line[:index] + ss[0] + "+"

	return ci
}

func getModifyCodeInfos(codes []eCodeInfo) ([]eCodeInfo, int) {
	maxCode := 0
	m := map[int][]eCodeInfo{}

	for _, ci := range codes {
		if ci.Num > maxCode {
			maxCode = ci.Num
		}

		if cis, ok := m[ci.Num]; ok {
			m[ci.Num] = append(cis, ci)
		} else {
			m[ci.Num] = []eCodeInfo{ci}
		}
	}

	needModify := []eCodeInfo{}
	for _, infos := range m {
		if len(infos) > 1 {
			needModify = append(needModify, infos[1:]...)
		}
	}

	return needModify, maxCode
}

func modifyErrCode(data []byte, infos []eCodeInfo, maxCode int) []byte {
	for _, info := range infos {
		maxCode++
		data = bytes.ReplaceAll(data, []byte(info.Str), []byte(info.DstStr+strconv.Itoa(maxCode)))
	}
	return data
}

func getDuplicateErrCodeInfo(data []byte) ([]eCodeInfo, int) {
	cis := []eCodeInfo{}

	buf := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if !strings.Contains(line, defineHTTPErrCodeMark) && !strings.Contains(line, defineGRPCErrCodeMark) {
			continue
		}

		ci := parseErrCodeInfo(line)
		if ci.Name != "" {
			cis = append(cis, ci)
		}
	}

	return getModifyCodeInfos(cis)
}

func checkAndModifyDuplicateErrCode(file string) (int, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return 0, err
	}

	serviceGroupData := bytes.Split(data, []byte(serviceGroupSeparatorMark))
	var fileContent [][]byte
	var count int
	for _, groupData := range serviceGroupData {
		ecis, maxCode := getDuplicateErrCodeInfo(groupData)
		fileContent = append(fileContent, modifyErrCode(groupData, ecis, maxCode))
		count += len(ecis)
	}

	data = bytes.Join(fileContent, []byte(serviceGroupSeparatorMark))
	err = os.WriteFile(file, data, 0666)
	return count, err
}
