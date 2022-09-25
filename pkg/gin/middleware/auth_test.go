package middleware

import (
	"fmt"
	"github.com/zhufuyi/sponge/pkg/utils"
	"io"
	"net/http"
	"testing"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gohttp"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/jwt"
)

var (
	uid  = "123"
	role = "admin"
)

func runAuthHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	jwt.Init()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Cors())

	tokenFun := func(c *gin.Context) {
		token, _ := jwt.GenerateToken(uid, role)
		fmt.Println("token =", token)
		response.Success(c, token)
	}

	userFun := func(c *gin.Context) {
		response.Success(c, "hello "+uid)
	}

	r.GET("/token", tokenFun)
	r.GET("/user/:id", Auth(), userFun)       // 需要鉴权
	r.GET("/admin/:id", AuthAdmin(), userFun) // 需要鉴权

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	return requestAddr
}

func TestAuth(t *testing.T) {
	role = ""
	requestAddr := runAuthHTTPServer()

	// 获取token
	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/token")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)

	// 正确的请求
	authorization := fmt.Sprintf("Bearer %s", token)
	val, err := getUser(requestAddr, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// 错误的 authorization
	val, err = getUser(requestAddr, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// 错误的 authorization
	val, err = getUser(requestAddr, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// 需要管理员访问权限
	val, err = getAdmin(requestAddr, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func TestAdminAuth(t *testing.T) {
	requestAddr := runAuthHTTPServer()

	// 获取token
	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/token")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)

	// 正确请求
	authorization := fmt.Sprintf("Bearer %s", token)
	val, err := getAdmin(requestAddr, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// 错误的 authorization
	val, err = getAdmin(requestAddr, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// 错误的 authorization
	val, err = getAdmin(requestAddr, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func getUser(requestAddr string, authorization string) (string, error) {
	client := &http.Client{}
	url := requestAddr + "/user/" + uid
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Authorization", authorization)
	if err != nil {
		return "", err
	}
	response, _ := client.Do(reqest)
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getAdmin(requestAddr string, authorization string) (string, error) {
	client := &http.Client{}
	url := requestAddr + "/admin/" + uid
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Authorization", authorization)
	if err != nil {
		return "", err
	}
	response, _ := client.Do(reqest)
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
