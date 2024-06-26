package patch

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/krand"
)

const (
	httpType = "http"
	grpcType = "grpc"
)

// ModifyDuplicateNumCommand modify duplicate numbers
func ModifyDuplicateNumCommand() *cobra.Command {
	var (
		dir string
	)

	cmd := &cobra.Command{
		Use:   "modify-dup-num",
		Short: "Modify duplicate numbers",
		Long: `modify duplicate numbers 

Examples:
  # modify duplicate numbers
  sponge patch modify-dup-num --dir=internal/ecode

`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := listErrCodeFiles(dir)
			if err != nil {
				return err
			}

			nsi, err := parseFiles(files)
			if err != nil {
				return err
			}

			msg, err := nsi.modifyHTTPDuplicateNum()
			if err != nil {
				return err
			}
			if msg != "" {
				fmt.Println("modify http duplicate numbers: ", msg)
			}
			msg, err = nsi.modifyGRPCDuplicateNum()
			if err != nil {
				return err
			}
			if msg != "" {
				fmt.Println("modify grpc duplicate numbers: ", msg)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "internal/ecode", "input directory")

	return cmd
}

func listErrCodeFiles(dir string) ([]string, error) {
	files, err := gofile.ListFiles(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("not found files")
	}

	filterFiles := []string{}
	for _, file := range files {
		if strings.Contains(file, "systemCode.go") ||
			strings.Contains(file, "systemCode_http.go") ||
			strings.Contains(file, "systemCode_rpc.go") {
			continue
		}
		filterFiles = append(filterFiles, file)
	}

	return filterFiles, nil
}

type numInfo struct {
	Name string
	Num  int
	Str  string
}

type numbersInfo struct {
	httpNumInfo         map[string]map[string]numInfo // map[file]map[code]numInfo
	httpDuplicationNums map[int][]string

	grpcNumInfo         map[string]map[string]numInfo // map[file]map[code]numInfo
	grpcDuplicationNums map[int][]string
}

func (r *numbersInfo) modifyHTTPDuplicateNum() (string, error) {
	msg := ""
	duplicateNums := []string{}

	if len(r.httpDuplicationNums) == 0 {
		return msg, nil
	}

	numMap := map[int]struct{}{}
	for num := range r.httpDuplicationNums {
		numMap[num] = struct{}{}
	}

	for num, fs := range r.httpDuplicationNums {
		if len(fs) <= 1 {
			continue
		}

		fs = sortFiles(fs)
		for i, file := range fs {
			if i == 0 {
				continue
			}
			for _, ni := range r.httpNumInfo[file] {
				newNum := genNewNum(numMap)
				if ni.Num == num {
					_, err := updateFile(file, newNum, ni)
					if err != nil {
						return msg, err
					}
					duplicateNums = append(duplicateNums, fmt.Sprintf("%d --> %d", ni.Num, newNum))
				}
			}
		}
	}

	if len(duplicateNums) == 0 {
		return msg, nil
	}
	return strings.Join(duplicateNums, ", "), nil
}

func (r *numbersInfo) modifyGRPCDuplicateNum() (string, error) {
	msg := ""
	duplicateNums := []string{}

	if len(r.grpcDuplicationNums) == 0 {
		return msg, nil
	}

	numMap := map[int]struct{}{}
	for num := range r.grpcDuplicationNums {
		numMap[num] = struct{}{}
	}

	for num, fs := range r.grpcDuplicationNums {
		if len(fs) <= 1 {
			continue
		}

		fs = sortFiles(fs)
		for i, file := range fs {
			if i == 0 {
				continue
			}
			for _, ni := range r.grpcNumInfo[file] {
				newNum := genNewNum(numMap)
				if ni.Num == num {
					_, err := updateFile(file, newNum, ni)
					if err != nil {
						return msg, err
					}
					duplicateNums = append(duplicateNums, fmt.Sprintf("%d --> %d", ni.Num, newNum))
				}
			}
		}
	}

	if len(duplicateNums) == 0 {
		return msg, nil
	}
	return strings.Join(duplicateNums, ", "), nil
}

func parseFiles(files []string) (*numbersInfo, error) {
	nsi := &numbersInfo{
		httpNumInfo:         map[string]map[string]numInfo{},
		httpDuplicationNums: map[int][]string{},
		grpcNumInfo:         map[string]map[string]numInfo{},
		grpcDuplicationNums: map[int][]string{},
	}

	for _, file := range files {
		result, err := parseNumberInfo(file)
		if err != nil {
			return nsi, err
		}
		if result == nil {
			continue
		}
		if result.errCodeType == httpType {
			for _, num := range result.nums {
				if fs, ok := nsi.httpDuplicationNums[num]; ok {
					fs = append(fs, file)
					nsi.httpDuplicationNums[num] = fs
				} else {
					nsi.httpDuplicationNums[num] = []string{file}
				}
			}
			nsi.httpNumInfo[file] = result.ni
		} else if result.errCodeType == grpcType {
			for _, num := range result.nums {
				if fs, ok := nsi.grpcDuplicationNums[num]; ok {
					fs = append(fs, file)
					nsi.grpcDuplicationNums[num] = fs
				} else {
					nsi.grpcDuplicationNums[num] = []string{file}
				}
			}
			nsi.grpcNumInfo[file] = result.ni
		}
	}

	return nsi, nil
}

type parseResult struct {
	ni          map[string]numInfo
	nums        []int
	errCodeType string
}

func parseNumberInfo(file string) (*parseResult, error) {
	errCodeType := ""
	ni := map[string]numInfo{}
	nums := []int{}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dataStr := string(data)
	if strings.Contains(dataStr, "errcode.NewError") {
		errCodeType = httpType
	} else if strings.Contains(dataStr, "errcode.NewRPCStatus") {
		errCodeType = grpcType
	}

	if errCodeType == "" {
		return nil, nil
	}

	var regStr string
	if errCodeType == httpType {
		regStr = `=[ ]*?errcode.HCode\(([\w\W]*?)\)\n`
	} else if errCodeType == grpcType {
		regStr = `=[ ]*?errcode.RCode\(([\w\W]*?)\)\n`
	}

	reg := regexp.MustCompile(regStr)
	allSubMatch := reg.FindAllStringSubmatch(dataStr, -1)
	if len(allSubMatch) == 0 {
		return nil, nil
	}
	names := []string{}

	for _, match := range allSubMatch {
		for i, v := range match {
			if i == 1 {
				names = append(names, v)
			}
		}
	}

	for _, name := range names {
		regStr = name + `[ ]*?=[ ]*?([\d]+)`
		reg = regexp.MustCompile(regStr)
		allSubMatch = reg.FindAllStringSubmatch(dataStr, -1)
		for _, match := range allSubMatch {
			if len(match) == 2 {
				num, _ := strconv.Atoi(match[1])
				nums = append(nums, num)
				ni[name] = numInfo{Name: name, Num: num, Str: match[0]}
			}
		}
	}

	return &parseResult{
		ni:          ni,
		nums:        nums,
		errCodeType: errCodeType,
	}, nil
}

func getFileCreateTime(file string) int64 {
	fi, err := os.Stat(file)
	if err != nil {
		return 0
	}

	return fi.ModTime().Unix()
}

func updateFile(file string, newNum int, ni numInfo) (numInfo, error) {
	strTmp := ni.Str
	ni.Num = newNum
	ni.Str = replaceNum(strTmp, ni.Num)

	data, err := os.ReadFile(file)
	if err != nil {
		return ni, err
	}
	data = bytes.ReplaceAll(data, []byte(strTmp), []byte(ni.Str))

	err = os.WriteFile(file, data, 0766)
	if err != nil {
		return ni, err
	}
	return ni, nil
}

func replaceNum(str string, newNum int) string {
	regStr := `([\w\W]*?=[ ]*?)[\d]+`
	reg := regexp.MustCompile(regStr)
	allSubMatch := reg.FindAllStringSubmatch(str, -1)
	for _, match := range allSubMatch {
		if len(match) == 2 {
			str = match[1] + fmt.Sprintf("%d", newNum)
		}
	}
	return str
}

type fileInfo struct {
	file        string
	createdTime int64
}

func sortFiles(files []string) []string {
	fis := []*fileInfo{}

	for _, file := range files {
		fis = append(fis, &fileInfo{
			file:        file,
			createdTime: getFileCreateTime(file),
		})
	}

	sort.Slice(fis, func(i, j int) bool {
		return fis[i].createdTime < fis[j].createdTime
	})

	var sFiles []string
	for _, fi := range fis {
		sFiles = append(sFiles, fi.file)
	}
	return sFiles
}

func genNewNum(numMap map[int]struct{}) int {
	max := 1000000
	count := 0
	for {
		count++
		newNum := krand.Int(99)
		if _, ok := numMap[newNum]; !ok {
			numMap[newNum] = struct{}{}
			return newNum
		}
		if count > max {
			break
		}
	}
	return 1
}
