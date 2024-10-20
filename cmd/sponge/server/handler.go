package server

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/zhufuyi/sponge/cmd/sponge/global"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/ggorm"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/krand"
	"github.com/zhufuyi/sponge/pkg/mgo"
	"github.com/zhufuyi/sponge/pkg/utils"
)

var (
	recordDirName = "sponge_record"
	saveDir       = fmt.Sprintf("%s/.%s", getSpongeDir(), recordDirName)
)

type dbInfoForm struct {
	Dsn      string `json:"dsn" binding:"required"`
	DbDriver string `json:"dbDriver"`
}

type kv struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ListDbDrivers list db drivers
func ListDbDrivers(c *gin.Context) {
	dbDrivers := []string{
		ggorm.DBDriverMysql,
		mgo.DBDriverName,
		ggorm.DBDriverPostgresql,
		ggorm.DBDriverTidb,
		ggorm.DBDriverSqlite,
	}

	data := []kv{}
	for _, driver := range dbDrivers {
		data = append(data, kv{
			Label: driver,
			Value: driver,
		})
	}

	response.Success(c, data)
}

// ListTables list tables
func ListTables(c *gin.Context) {
	form := &dbInfoForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		response.Error(c, errcode.InvalidParams.RewriteMsg(err.Error()))
		return
	}
	dbParams := strings.Split(form.Dsn, ";")
	form.Dsn = dbParams[0]
	var tables []string
	switch strings.ToLower(form.DbDriver) {
	case ggorm.DBDriverMysql, ggorm.DBDriverTidb:
		tables, err = getMysqlTables(form.Dsn)
	case ggorm.DBDriverPostgresql:
		tables, err = getPostgresqlTables(form.Dsn)
	case ggorm.DBDriverSqlite:
		tables, err = getSqliteTables(form.Dsn)
	case mgo.DBDriverName:
		tables, err = getMongodbTables(form.Dsn)
	case "":
		response.Error(c, errcode.InvalidParams.RewriteMsg("database type cannot be empty"))
		return
	default:
		response.Error(c, errcode.InvalidParams.RewriteMsg("unsupported database type: "+form.DbDriver))
		return
	}
	if err != nil {
		response.Error(c, errcode.InternalServerError.RewriteMsg(err.Error()))
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
	Arg  string `json:"arg" binding:"required"`
	Path string `json:"path" binding:"required"`
}

// GenerateCode generate code
func GenerateCode(c *gin.Context) {
	// Allow getting the value of the request header when crossing domains
	c.Writer.Header().Set("Access-Control-Expose-Headers", "content-disposition, err-msg")

	form := &GenerateCodeForm{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		responseErr(c, err, errcode.InvalidParams)
		return
	}

	handleGenerateCode(c, form.Path, form.Arg)
}

// nolint
func handleGenerateCode(c *gin.Context, outPath string, arg string) {
	out := "-" + time.Now().Format("150405")
	if len(outPath) > 1 {
		if outPath[0] == '/' {
			out = outPath[1:] + out
		} else {
			out = outPath + out
		}
	}

	args := strings.Split(arg, " ")
	params := parseCommandArgs(args)
	if params.ServerName != "" {
		out = params.ServerName + "-" + out
	} else {
		if params.ModuleName != "" {
			out = params.ModuleName + "-" + out
		}
	}
	out = global.Path
	if params.SuitedMonoRepo {
		out += "-mono-repo"
	}

	//out = os.TempDir() + gofile.GetPathDelimiter() + "sponge-generate-code" + gofile.GetPathDelimiter() + out
	args = append(args, fmt.Sprintf("--out=%s", out))
	fmt.Println(strings.Join(args, " "))

	ctx, _ := context.WithTimeout(context.Background(), time.Second*2) // nolint
	result := gobash.Run(ctx, "sponge", args...)
	for v := range result.StdOut {
		_ = v
	}
	if result.Err != nil {
		responseErr(c, result.Err, errcode.InternalServerError)
		return
	}

	//zipFile := out + ".zip"
	//err := CompressPathToZip(out, zipFile)
	//if err != nil {
	//	responseErr(c, err, errcode.InternalServerError)
	//	return
	//}
	//
	//if !gofile.IsExists(zipFile) {
	//	err = errors.New("no found file " + zipFile)
	//	responseErr(c, err, errcode.InternalServerError)
	//	return
	//}
	//
	//c.Writer.Header().Set("content-disposition", gofile.GetFilename(zipFile))
	//c.File(zipFile)

	recordObj().set(c.ClientIP(), outPath, params)

	go func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 5):
				//err := os.RemoveAll(out)
				//if err != nil {
				//	continue
				//}
				//err = os.RemoveAll(zipFile)
				//if err != nil {
				//	continue
				//}

				if params.ProtobufFile != "" && strings.Contains(params.ProtobufFile, recordDirName) {
					err := os.RemoveAll(gofile.GetFileDir(params.ProtobufFile))
					if err != nil {
						continue
					}
				}
				if params.YamlFile != "" && strings.Contains(params.YamlFile, recordDirName) {
					err := os.RemoveAll(gofile.GetFileDir(params.YamlFile))
					if err != nil {
						continue
					}
				}
				return
			}
		}
	}()
}

// GetRecord generate run command record
func GetRecord(c *gin.Context) {
	pathParam := c.Param("path")
	if pathParam == "" {
		response.Out(c, errcode.InvalidParams.RewriteMsg("path param is empty"))
		return
	}

	params := recordObj().get(c.ClientIP(), pathParam)
	if params == nil {
		params = &parameters{Embed: true}
	}

	response.Success(c, params)
}

func responseErr(c *gin.Context, err error, ec *errcode.Error) {
	k := "err-msg"
	e := ec.RewriteMsg(err.Error())
	c.Writer.Header().Set(k, e.Msg())
	response.Out(c, e)
}

// UploadFiles batch files upload
func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, errcode.InvalidParams.RewriteMsg(err.Error()))
		return
	}

	if len(form.File) == 0 {
		response.Error(c, errcode.InvalidParams.RewriteMsg("upload file is empty"))
		return
	}

	//spongeArg, err := getFormValue(form.Value, "spongeArg")
	//if err != nil {
	//	response.Error(c, errcode.InvalidParams.RewriteMsg("the field 'spongeArg' cannot be empty"))
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
			if !checkFileType(fileType) {
				response.Error(c, errcode.InvalidParams.RewriteMsg("only .proto or yaml files are allowed to be uploaded"))
				return
			}

			filePath = savePath + "/" + filename
			if checkSameFile(hadSaveFiles, filePath) {
				continue
			}
			if err = c.SaveUploadedFile(file, filePath); err != nil {
				response.Error(c, errcode.InternalServerError.RewriteMsg(err.Error()))
				return
			}

			hadSaveFiles = append(hadSaveFiles, filePath)
		}
	}

	if fileType == ".proto" {
		filePath = savePath + "/*.proto"
	}

	response.Success(c, filePath)
}

//func getFormValue(valueMap map[string][]string, key string) (string, error) {
//	valueSlice := valueMap[key]
//	if len(valueSlice) == 0 {
//		return "", fmt.Errorf("form '%s' is empty", key)
//	}
//
//	return valueSlice[0], nil
//}

func checkFileType(typeName string) bool {
	switch typeName {
	case ".proto", ".yml", ".yaml":
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
	var dir = saveDir
	if gofile.IsWindows() {
		dir = strings.ReplaceAll(saveDir, "\\", "/")
	}
	dir += "/" + krand.String(krand.R_All, 8)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0766)
	}
	return dir
}

// CompressPathToZip compressed directory to zip file
func CompressPathToZip(outPath, targetFile string) error {
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

	f, err := os.Open(outPath)
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

func getSpongeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("can't get home directory'")
		return ""
	}

	return dir
}

func getMysqlTables(dsn string) ([]string, error) {
	dsn = utils.AdaptiveMysqlDsn(dsn)
	db, err := ggorm.InitMysql(dsn)
	if err != nil {
		return nil, err
	}
	defer ggorm.CloseSQLDB(db)

	var tables []string
	err = db.Raw("show tables").Scan(&tables).Error
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func getPostgresqlTables(dsn string) ([]string, error) {
	dsn = utils.AdaptivePostgresqlDsn(dsn)
	db, err := ggorm.InitPostgresql(dsn)
	if err != nil {
		return nil, err
	}
	defer ggorm.CloseSQLDB(db)

	var tables []string
	err = db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", "public").Scan(&tables).Error
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func getSqliteTables(dbFile string) ([]string, error) {
	if !gofile.IsExists(dbFile) {
		return nil, fmt.Errorf("sqlite db file %s not found in local host", dbFile)
	}

	db, err := ggorm.InitSqlite(dbFile)
	if err != nil {
		return nil, err
	}
	defer ggorm.CloseSQLDB(db)

	var tables []string
	err = db.Raw("select name from sqlite_master where type = ?", "table").Scan(&tables).Error
	if err != nil {
		return nil, err
	}

	filteredTables := []string{}
	for _, table := range tables {
		if table == "sqlite_sequence" {
			continue
		}
		filteredTables = append(filteredTables, table)
	}

	return filteredTables, nil
}

func getMongodbTables(dsn string) ([]string, error) {
	dsn = utils.AdaptiveMongodbDsn(dsn)
	db, err := mgo.Init(dsn)
	if err != nil {
		return nil, err
	}
	defer mgo.Close(db) //nolint

	tables, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	if len(tables) == 0 {
		u, _ := url.Parse(dsn)
		return nil, fmt.Errorf("mongodb db %s has no tables", strings.TrimLeft(u.Path, "/"))
	}

	return tables, nil
}
