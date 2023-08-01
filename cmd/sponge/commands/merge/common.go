package merge

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

var (
	internalECodeDir   = "internal/ecode"
	internalRoutersDir = "internal/routers"
	internalHandlerDir = "internal/handler"
	internalServiceDir = "internal/service"

	backupDir     = os.TempDir() + gofile.GetPathDelimiter() + "sponge_merge_backup_code" + gofile.GetPathDelimiter()
	splitLineMark = []byte(`// ---------- Do not delete or move this split line, this is the merge code marker ----------`)
)

// type mergeDataFunc func(subData1 []byte, subData2 []byte) []byte
type parseCodeFunc func(date []byte) []code

type code struct {
	key   string
	value string
}

func runMerge(dir string, dt string, isLineCode bool, parseCode parseCodeFunc) {
	files, err := parseFiles(dir)
	if err != nil {
		fmt.Println("Warring:", err)
	}

	for _, file := range files {
		successFile, err := runMergeCode(file, dt, isLineCode, parseCode)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if successFile != "" {
			fmt.Printf("merge code to '%s' successfully.\n", successFile)
		}
	}
}

func parseFiles(filePath string) ([]string, error) {
	result, err := gobash.Exec("ls", filePath+"/*go.gen.*")
	if err != nil {
		return nil, err
	}

	return strings.Split(string(result), "\n"), nil
}

func runMergeCode(file string, dtStr string, isLineCode bool, parseCode parseCodeFunc) (string, error) {
	if file == "" {
		return "", nil
	}

	ss := strings.Split(file, ".gen.")
	oldFile := ss[0]

	data1, err := os.ReadFile(oldFile)
	if err != nil {
		return "", err
	}

	data2, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	count1 := bytes.Count(data1, splitLineMark)
	count2 := bytes.Count(data2, splitLineMark)
	if count1 != count2 {
		return "", fmt.Errorf("merge code mark mismatch, please merge codes manually, file = %s", file)
	}

	var data []byte

	if count1 == 0 {
		data = mergeData(data1, data2, isLineCode, parseCode)
	} else {
		data1ss := bytes.Split(data1, splitLineMark)
		data2ss := bytes.Split(data2, splitLineMark)

		for index, subData1 := range data1ss {
			subData := mergeData(subData1, data2ss[index], isLineCode, parseCode)
			data = append(data, subData...)
			if index < len(data1ss)-1 {
				data = append(data, splitLineMark...)
			}
		}
	}

	if len(data1) > len(data) {
		return "", fmt.Errorf("to avoid replacing logical code, please merge codes manually, file = %s", file)
	}

	if len(data1) == len(data) {
		return "", os.Remove(file)
	}

	err = saveFile(oldFile, data, dtStr)
	if err != nil {
		return "", err
	}

	return oldFile, os.Remove(file)
}

func mergeData(subData1 []byte, subData2 []byte, isLineCode bool, parseCode parseCodeFunc) []byte {
	c1 := parseCode(subData1)
	c2 := parseCode(subData2)
	var addCode, mark []byte
	if isLineCode {
		addCode, mark = compareCode(c1, c2)
	} else {
		addCode, mark = compareCode2(c1, c2, subData2)
	}
	return mergeCode(subData1, addCode, mark)
}

func compareCode(oldCode []code, newCode []code) ([]byte, []byte) {
	var addCode []byte
	var mark []byte

	for _, code1 := range newCode {
		isEqual := false
		for _, code2 := range oldCode {
			if code1.key == code2.key {
				isEqual = true
				mark = []byte(code2.value)
				break
			}
		}
		if !isEqual {
			addCode = append(addCode, []byte(code1.value)...)
		}
	}

	return addCode, mark
}

func compareCode2(oldCode []code, newCode []code, data []byte) ([]byte, []byte) {
	var addCode []byte
	var mark []byte

	for _, code1 := range newCode {
		isEqual := false
		for _, code2 := range oldCode {
			if code1.key == code2.key {
				isEqual = true
				mark = []byte(code2.value)
				break
			}
		}
		if !isEqual {
			_, name := getTmplKey(code1.value)
			comment := getComment(name, string(data))
			addCode = append(addCode, []byte("\n\n"+comment+"\n"+code1.value)...)
		}
	}

	return addCode, mark
}

func mergeCode(oldCode []byte, addCode []byte, mark []byte) []byte {
	if len(addCode) == 0 {
		return oldCode
	}

	var data []byte

	ss := bytes.SplitN(oldCode, mark, 2)
	if len(ss) != 2 {
		return oldCode
	}
	data = append(ss[0], mark...)
	data = append(data, addCode...)
	data = append(data, ss[1]...)

	return data
}

func saveFile(file string, data []byte, dtStr string) error {
	bkDir := backupDir + dtStr
	_ = os.MkdirAll(bkDir, 0744)
	_, _ = gobash.Exec("cp", file, bkDir+gofile.GetPathDelimiter()+gofile.GetFilename(file))

	return os.WriteFile(file, data, 0766)
}

// -------------------------------------------------------------------------------------------

var (
	regStr1 = `func \(h[\w\W]*?\n}`
	reg1    = regexp.MustCompile(regStr1)

	regStr2 = `func \((.*?)\) (.*?)\(`
	reg2    = regexp.MustCompile(regStr2)
)

func parseTmplCode(data []byte) []code {
	var codes []code
	str := string(data)
	matches := reg1.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		for _, v := range match {
			key, _ := getTmplKey(v)
			codes = append(codes, code{
				key:   key,
				value: v,
			})
		}
	}
	return codes
}

func getTmplKey(str string) (string, string) {
	matches := reg2.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return "", ""
	}

	if len(matches[0]) != 3 {
		return "", ""
	}

	return matches[0][0], matches[0][2]
}

func getComment(name string, str string) string {
	//regStr := `//(.*?)` + name + `[\w\W]*?\nfunc`
	regStr := `//( ?)` + name + `[\w\W]*?\nfunc`
	reg := regexp.MustCompile(regStr)
	match := reg.FindAllString(str, -1)
	if len(match) == 0 {
		return ""
	}
	return strings.ReplaceAll(match[0], "\nfunc", "")
}
