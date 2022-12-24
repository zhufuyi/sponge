package server

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/mysql"
	"io"
	"os"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/internal/ecode"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gin/validator"
	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// NewRouter create a router
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.Logging(middleware.WithLog(logger.Get())))
	binding.Validator = validator.Init()

	r.POST("/generate", GenerateCode)
	r.POST("/listTables", ListTables)

	return r
}

type TestMysqlForm struct {
	Dsn string `json:"dsn"  binding:"min=10"`
}

type kv struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

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
	c.Writer.Header().Set("Access-Control-Expose-Headers", "content-disposition, help-for-use, err-msg")

	form := &GenerateCodeForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		responseErr(c, err, ecode.InvalidParams)
		return
	}

	out := "-" + time.Now().Format("20060102150405")
	if len(form.Path) > 1 {
		if form.Path[0] == '/' {
			out = form.Path[1:] + out
		} else {
			out = form.Path + out
		}
	}
	form.Arg += fmt.Sprintf(" --out=%s", out)

	args := strings.Split(form.Arg, " ")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20) // nolint
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
	err = CompressPathToZip(out, zipFile)
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
	if v, ok := helperInfo[form.Path]; ok {
		c.Writer.Header().Set("help-for-use", v)
	}

	c.File(zipFile)

	_ = os.RemoveAll(out)
	_ = os.RemoveAll(zipFile)
}

func responseErr(c *gin.Context, err error, ec *errcode.Error) {
	k := "err-msg"
	e := ec.WithDetails(err.Error())
	c.Writer.Header().Set(k, e.Msg())
	response.Output(c, e.ToHTTPCode())
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

var helperInfo = map[string]string{
	"/web-http": `
Step 1 Unzip the file
Step 2 Update swagger documentation: make docs
Step 3 Compile and run services: make run
Step 4 Copy http://localhost:8080/swagger/index.html to your browser to test the CRUD api
`,
	"/web-handler": `
Step 1 Unzip the file and copy to your web service directory
Step 2 Update swagger documentation: make docs
Step 3 Compile and run services: make run
Step 4 Copy http://localhost:8080/swagger/index.html to your browser to test the CRUD api
`,
	"/web-dao": `
Unzip the file and copy to your web service directory
`,
	"/web-model": `
Unzip the file and copy to your web service directory
`,
	"/micro-rpc": `
Step 1 Unzip the file
Step 2 Generate *pb.go: make proto
Step 3 Compile and run services: make run
Step 4 Using 'Goland'' or 'VS Code' to open the internal/service/xxx_client_test.go file, test or pressure test rpc methods
`,
	"/micro-service": `
Step 1 Unzip the file and copy to your rpc service directory
Step 2 Generate *pb.go: make proto
Step 3 Compile and run services: make run
Step 4 Using 'Goland' or 'VS Code' to open the internal/service/xxx_client_test.go file, test or pressure test rpc methods
`,
	"/micro-dao": `
Unzip the file and copy to your rpc service directory
`,
	"/micro-model": `
Unzip the file and copy to your rpc service directory
`,
	"/micro-protobuf": `
Unzip the file and copy to your rpc service directory
`,
	"/web-http-pb": `
Step 1 Unzip the file
Step 2 Generate *pb.go file, generate handler template code, update swagger documentation: make proto
Step 3 Compile and run services: make run
Step 4 Copy http://localhost:8080/apis/swagger/index.html to your browser to test the api
`,
	"/web-rpc-gw-pb": `
Step 1 Unzip the file
Step 2 Generate *pb.go files, generate template code, update swagger documentation: make proto
Step 3 Compile and run services: make run
Step 4 Copy http://localhost:8080/apis/swagger/index.html to your browser to test the api
`,
	"/micro-rpc-pb": `
Step 1 Unzip the file
Step 2 Generate *pb.go file, generate service template code: make proto
Step 3 Compile and run services: make run
Step 4 Using 'Goland'' or 'VS Code' to open the internal/service/xxx_client_test.go file, test rpc methods
`,
	"/yaml-config": `
Unzip the file and copy to your service directory
`,
	"/rpc-client": `
Unzip the file and copy to your service directory
`,
}
