package server

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/krand"
	"github.com/zhufuyi/sponge/pkg/mysql"

	"github.com/gin-gonic/gin"
)

type TestMysqlForm struct {
	Dsn string `json:"dsn"  binding:"min=10"`
}

type kv struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ListTables list tables
func ListTables(c *gin.Context) {
	form := &TestMysqlForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, ecode.InvalidParams.WithDetails(err.Error()))
		return
	}

	db, err := mysql.Init(form.Dsn)
	if err != nil {
		response.Error(c, ecode.InternalServerError.WithDetails(err.Error()))
		return
	}

	var tables []string
	err = db.Raw("show tables").Scan(&tables).Error
	if err != nil {
		response.Error(c, ecode.InternalServerError.WithDetails(err.Error()))
		return
	}

	data := []kv{}
	for _, table := range tables {
		data = append(data, kv{
			Label: table,
			Value: table,
		})
	}

	response.Success(c, data)
}

// GenerateCodeForm generate code form
type GenerateCodeForm struct {
	Arg  string `json:"arg" binding:"min=1"`
	Path string `json:"path" binding:"min=2"`
}

// GenerateCode generate code
func GenerateCode(c *gin.Context) {
	// Allow getting the value of the request header when crossing domains
	c.Writer.Header().Set("Access-Control-Expose-Headers", "content-disposition, err-msg")

	form := &GenerateCodeForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		responseErr(c, err, ecode.InvalidParams)
		return
	}

	generateCode(c, form.Path, form.Arg)
}

func generateCode(c *gin.Context, path string, arg string) {
	out := "-" + time.Now().Format("20060102150405")
	if len(path) > 1 {
		if path[0] == '/' {
			out = path[1:] + out
		} else {
			out = path + out
		}
	}
	arg += fmt.Sprintf(" --out=%s", out)

	args := strings.Split(arg, " ")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10) // nolint
	result := gobash.RunC(ctx, "sponge", args...)
	var data []string
	for v := range result.StdOut {
		data = append(data, v)
	}
	if result.Err != nil {
		responseErr(c, result.Err, ecode.InternalServerError)
		return
	}

	zipFile := out + ".zip"
	err := CompressPathToZip(out, zipFile)
	if err != nil {
		responseErr(c, err, ecode.InternalServerError)
		return
	}

	if !gofile.IsExists(zipFile) {
		err = errors.New("no found file " + zipFile)
		responseErr(c, err, ecode.InternalServerError)
		return
	}

	c.Writer.Header().Set("content-disposition", zipFile)
	c.File(zipFile)

	params := parseCommandArgs(args)
	recordObj().set(c.ClientIP(), path, params)

	_ = os.RemoveAll(out)
	_ = os.RemoveAll(zipFile)
	if strings.Contains(path, "-pb") {
		dir, _ := filepath.Split(params.ProtobufFile)
		_ = os.RemoveAll(dir)
	}
}

// GetRecord generate run command record
func GetRecord(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		response.Out(c, ecode.InvalidParams.WithDetails("path is empty"))
		return
	}

	params := recordObj().get(c.ClientIP(), path)
	if params == nil {
		params = &parameters{Embed: true}
	}

	response.Success(c, params)
}

func responseErr(c *gin.Context, err error, ec *errcode.Error) {
	k := "err-msg"
	e := ec.WithDetails(err.Error())
	c.Writer.Header().Set(k, e.Msg())
	response.Out(c, e)
}

// UploadFiles 批量上传文件
func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, ecode.InvalidParams.WithDetails(err.Error()))
		return
	}

	if len(form.File) == 0 {
		response.Error(c, ecode.InvalidParams.WithDetails("upload file is empty"))
		return
	}

	//spongeArg, err := getFormValue(form.Value, "spongeArg")
	//if err != nil {
	//	response.Error(c, ecode.InvalidParams.WithDetails("the field 'spongeArg' cannot be empty"))
	//	return
	//}

	hadSaveFiles := []string{}
	savePath := getSavePath()
	fileType := ""
	var filePath string
	for _, files := range form.File {
		for _, file := range files {
			filename := filepath.Base(file.Filename)
			fileType = path.Ext(filename)
			if !checkNameType(filename) {
				response.Error(c, ecode.InvalidParams.WithDetails("only .proto or yaml files are allowed to be uploaded"))
				return
			}

			filePath = savePath + "/" + filename
			if checkSameFile(hadSaveFiles, filePath) {
				continue
			}
			if err = c.SaveUploadedFile(file, filePath); err != nil {
				response.Error(c, ecode.InternalServerError.WithDetails(err.Error()))
				return
			}

			hadSaveFiles = append(hadSaveFiles, filePath)
		}
	}

	if fileType == ".proto" {
		filePath = savePath + "/*.proto"
	} else {
		files, err := gofile.ListFiles(savePath)
		if err != nil {
			response.Error(c, ecode.InternalServerError.WithDetails(err.Error()))
		}
		if len(files) > 0 {
			filePath = files[0]
		}
	}

	response.Success(c, filePath)
}

func getFormValue(valueMap map[string][]string, key string) (string, error) {
	valueSlice := valueMap[key]
	if len(valueSlice) == 0 {
		return "", fmt.Errorf("form '%s' is empty", key)
	}

	return valueSlice[0], nil
}

// 清空已上传的文件
func removeFiles(files []string) {
	for _, file := range files {
		_ = os.RemoveAll(file)
	}
}

func checkNameType(filename string) bool {
	tmpStr := strings.ToLower(filename)
	if strings.TrimRight(tmpStr, ".proto") != tmpStr ||
		strings.TrimRight(tmpStr, ".yml") != tmpStr ||
		strings.TrimRight(tmpStr, ".yaml") != tmpStr {
		return true
	}

	return false
}

func checkSameFile(files []string, file string) bool {
	for _, v := range files {
		if v == file {
			return true
		}
	}
	return false
}

func getSavePath() string {
	path := gofile.GetRunPath()
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "/")
	}
	path = strings.TrimRight(path, "/") + "/proto/" + krand.String(krand.R_All, 8)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0666)
	}

	return path
}

// CompressPathToZip compressed directory to zip file
func CompressPathToZip(path, targetFile string) error {
	d, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = d.Close()
	}()
	w := zip.NewWriter(d)
	defer func() {
		_ = w.Close()
	}()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	return compress(f, "", w)
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		_ = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
