package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/gohttp"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/utils"
)

var (
	uid    = "123"
	role   = "admin"
	fields = jwt.KV{"id": 1, "foo": "bar"}
)

func verify(claims *jwt.Claims, tokenTail10 string, c *gin.Context) error {
	if claims.UID != uid || claims.Role != role {
		return errors.New("verify failed")
	}

	// token := getToken(claims.UID)
	// if  token[len(token)-10:] != tokenTail10 { return err }

	return nil
}

func verifyCustom(claims *jwt.CustomClaims, tokenTail10 string, c *gin.Context) error {
	err := errors.New("verify failed")

	id, exist := claims.Get("id")
	if !exist {
		return err
	}
	foo, exist := claims.Get("foo")
	if !exist {
		return err
	}
	if int(id.(float64)) != fields["id"].(int) || foo.(string) != fields["foo"].(string) {
		return err
	}

	// token := getToken(id)
	// if  token[len(token)-10:] != tokenTail10 { return err }

	return nil
}

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
	customTokenFun := func(c *gin.Context) {
		token, _ := jwt.GenerateCustomToken(fields)
		fmt.Println("token custom =", token)
		response.Success(c, token)
	}

	userFun := func(c *gin.Context) {
		response.Success(c, "hello "+uid)
	}

	r.GET("/token", tokenFun)
	r.GET("/user/:id", Auth(), userFun)
	r.GET("/user2/:id", Auth(WithVerify(verify), WithSwitchHTTPCode()), userFun)

	r.GET("/token/custom", customTokenFun)
	r.GET("/user/custom", AuthCustom(verifyCustom), userFun)

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
	val, err := getUser(requestAddr+"/user/"+uid, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	val, err = getUser(requestAddr+"/user2/"+uid, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// verify error
	role = "foobar"
	val, err = getUser(requestAddr+"/user2/"+uid, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getUser(requestAddr+"/user/"+uid, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getUser(requestAddr+"/user/"+uid, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func TestAuthCustom(t *testing.T) {
	requestAddr := runAuthHTTPServer()

	// get token
	result := &gohttp.StdResult{}
	err := gohttp.Get(result, requestAddr+"/token/custom")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)

	url := requestAddr + "/user/custom"

	// the right request
	authorization := fmt.Sprintf("Bearer %s", token)
	val, err := getUserCustom(url, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// verify error
	fields["foo"] = "bar2"
	val, err = getUser(url, authorization)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getUserCustom(url, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)

	// wrong authorization
	val, err = getUserCustom(url, token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func getUser(url string, authorization string) (string, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", authorization)
	if err != nil {
		return "", err
	}
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getUserCustom(url string, authorization string) (string, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", authorization)
	if err != nil {
		return "", err
	}
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
