package middleware

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
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
	r.GET("/user/:id", Auth(), userFun)
	r.GET("/admin/:id", AuthAdmin(), userFun)

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestAuth(t *testing.T) {
	role = ""
	requestAddr := runAuthHTTPServer()

	// get token
	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/token")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)

	// the right request
	authorization := fmt.Sprintf("Bearer %s", token)
	val, err := getUser(requestAddr, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getUser(requestAddr, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getUser(requestAddr, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// administrator access required
	val, err = getAdmin(requestAddr, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func TestAdminAuth(t *testing.T) {
	requestAddr := runAuthHTTPServer()

	// get token
	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/token")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)

	// the right request
	authorization := fmt.Sprintf("Bearer %s", token)
	val, err := getAdmin(requestAddr, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getAdmin(requestAddr, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
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
