package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/httpcli"
	"github.com/go-dev-frame/sponge/pkg/jwt"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var (
	uid  = "100"
	name = "tom"

	fields = jwt.KV{"id": utils.StrToUint64(uid), "name": name, "age": 10, "foo": "bar"}

	errMsg = http.StatusText(http.StatusUnauthorized)
)

func verify(claims *jwt.Claims, tokenTail10 string, c *gin.Context) error {
	if claims.UID != uid || claims.Name != name {
		return errors.New("verify failed")
	}

	// token := getToken(claims.UID)
	// if  token[len(token)-10:] != tokenTail10 { return err }

	return nil
}

func verifyCustom(claims *jwt.CustomClaims, tokenTail10 string, c *gin.Context) error {
	err := errors.New("verify failed")

	//token, fields := getToken(id)
	// if  token[len(token)-10:] != tokenTail10 { return err }

	id, exist := claims.GetUint64("id")
	if !exist || id != fields["id"] {
		return err
	}

	name, exist := claims.GetString("name")
	if !exist || name != fields["name"] {
		return err
	}

	age, exist := claims.GetInt("age")
	if !exist || age != fields["age"] {
		return err
	}

	return nil
}

func runAuthHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	jwt.Init()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Cors())

	getUserByIDHandler := func(c *gin.Context) {
		id := c.Param("id")
		response.Success(c, id)
	}

	loginHandler := func(c *gin.Context) {
		token, _ := jwt.GenerateToken(uid, name)
		fmt.Println("token =", token)
		response.Success(c, token)
	}
	r.GET("/auth/login", loginHandler)
	r.GET("/user/:id", Auth(), getUserByIDHandler)
	r.GET("/user/toHTTPCode/:id", Auth(WithVerify(verify), WithSwitchHTTPCode()), getUserByIDHandler)

	loginCustomHandler := func(c *gin.Context) {
		token, _ := jwt.GenerateCustomToken(fields)
		fmt.Println("custom token =", token)
		response.Success(c, token)
	}
	r.GET("/auth/customLogin", loginCustomHandler)
	r.GET("/user/custom/:id", AuthCustom(verifyCustom), getUserByIDHandler)
	r.GET("/user/custom/toHTTPCode/:id", AuthCustom(verifyCustom, WithSwitchHTTPCode()), getUserByIDHandler)

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
	result := &httpcli.StdResult{}
	err := httpcli.Get(result, requestAddr+"/auth/login")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)
	authorization := fmt.Sprintf("Bearer %s", token)

	// success
	val, err := getUser(requestAddr+"/user/"+uid, authorization)
	assert.Equal(t, val["data"], uid)

	// success
	val, err = getUser(requestAddr+"/user/toHTTPCode/"+uid, authorization)
	assert.Equal(t, val["data"], uid)

	// verify name failed, return 401
	name = "notfound"
	val, err = getUser(requestAddr+"/user/toHTTPCode/"+uid, authorization)
	assert.Equal(t, val["msg"], errMsg)

	// authorization format error, missing token, return 200
	val, err = getUser(requestAddr+"/user/"+uid, "Bearer ")
	assert.Equal(t, val["msg"], errMsg)

	// authorization format error, missing Bearer, return 200
	val, err = getUser(requestAddr+"/user/"+uid, token)
	assert.Equal(t, val["msg"], errMsg)
}

func TestAuthCustom(t *testing.T) {
	requestAddr := runAuthHTTPServer()

	// get token
	result := &httpcli.StdResult{}
	err := httpcli.Get(result, requestAddr+"/auth/customLogin")
	if err != nil {
		t.Fatal(err)
	}
	token := result.Data.(string)
	authorization := fmt.Sprintf("Bearer %s", token)

	// success
	val, _ := getUserCustom(requestAddr+"/user/custom/"+uid, authorization)
	assert.Equal(t, val["data"], uid)

	// success
	val, _ = getUserCustom(requestAddr+"/user/custom/toHTTPCode/"+uid, authorization)
	assert.Equal(t, val["data"], uid)

	// verify name error, return 401
	fields["name"] = "john"
	val, _ = getUser(requestAddr+"/user/custom/toHTTPCode/"+uid, authorization)
	assert.Equal(t, val["msg"], errMsg)

	// authorization format error, missing token, return 200
	val, _ = getUserCustom(requestAddr+"/user/custom/"+uid, "Bearer ")
	assert.Equal(t, val["msg"], errMsg)

	// authorization format error, missing Bearer, return 200
	val, _ = getUserCustom(requestAddr+"/user/custom/"+uid, token)
	assert.Equal(t, val["msg"], errMsg)
}

func getUser(url string, authorization string) (gin.H, error) {
	var result = gin.H{}

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", authorization)
	if err != nil {
		return result, err
	}
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)

	return result, err
}

func getUserCustom(url string, authorization string) (gin.H, error) {
	var result = gin.H{}

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", authorization)
	if err != nil {
		return result, err
	}
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)

	return result, err
}
