package gohttp

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/utils"
)

type myBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func runGoHTTPServer() string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	oKFun := func(c *gin.Context) {
		uid := c.Query("uid")
		fmt.Printf("request parameters: uid=%s\n", uid)
		c.JSON(200, StdResult{
			Code: 0,
			Msg:  "ok",
			Data: fmt.Sprintf("uid=%v", uid),
		})
	}
	errFun := func(c *gin.Context) {
		uid := c.Query("uid")
		fmt.Printf("request parameters: uid=%s\n", uid)
		c.JSON(401, StdResult{
			Code: 401,
			Msg:  "authorization failure",
			Data: fmt.Sprintf("uid=%v", uid),
		})
	}

	oKPFun := func(c *gin.Context) {
		var body myBody
		c.BindJSON(&body)
		fmt.Println("body data:", body)
		c.JSON(200, StdResult{
			Code: 0,
			Msg:  "ok",
			Data: body,
		})
	}
	errPFun := func(c *gin.Context) {
		var body myBody
		c.BindJSON(&body)
		fmt.Println("body data:", body)
		c.JSON(401, StdResult{
			Code: 401,
			Msg:  "authorization failure",
			Data: nil,
		})
	}

	r.GET("/get", oKFun)
	r.GET("/get_err", errFun)
	r.DELETE("/delete", oKFun)
	r.DELETE("/delete_err", errFun)
	r.POST("/post", oKPFun)
	r.POST("/post_err", errPFun)
	r.PUT("/put", oKPFun)
	r.PUT("/put_err", errPFun)
	r.PATCH("/patch", oKPFun)
	r.PATCH("/patch_err", errPFun)

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

// ------------------------------------------------------------------------------------------

func TestGetStandard(t *testing.T) {
	requestAddr := runGoHTTPServer()

	req := Request{}
	req.SetURL(requestAddr + "/get")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})
	req.SetParams(KV{
		"name": "foo",
	})

	resp, err := req.GET()
	if err != nil {
		t.Fatal(err)
	}

	result := &StdResult{}
	err = resp.BindJSON(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
}

func TestDeleteStandard(t *testing.T) {
	requestAddr := runGoHTTPServer()

	req := Request{}
	req.SetURL(requestAddr + "/delete")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})
	req.SetParams(KV{
		"uid": 123,
	})

	resp, err := req.DELETE()
	if err != nil {
		t.Fatal(err)
	}

	result := &StdResult{}
	err = resp.BindJSON(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
}

func TestPostStandard(t *testing.T) {
	requestAddr := runGoHTTPServer()

	req := Request{}
	req.SetURL(requestAddr + "/post")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})
	req.SetJSONBody(&myBody{
		Name:  "foo",
		Email: "bar@gmail.com",
	})

	resp, err := req.POST()
	if err != nil {
		t.Fatal(err)
	}

	result := &StdResult{}
	err = resp.BindJSON(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
}

func TestPutStandard(t *testing.T) {
	requestAddr := runGoHTTPServer()

	req := Request{}
	req.SetURL(requestAddr + "/put")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})
	req.SetJSONBody(&myBody{
		Name:  "foo",
		Email: "bar@gmail.com",
	})

	resp, err := req.PUT()
	if err != nil {
		t.Fatal(err)
	}

	result := &StdResult{}
	err = resp.BindJSON(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
}

func TestPatchStandard(t *testing.T) {
	requestAddr := runGoHTTPServer()

	req := Request{}
	req.SetURL(requestAddr + "/patch")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})
	req.SetJSONBody(&myBody{
		Name:  "foo",
		Email: "bar@gmail.com",
	})

	resp, err := req.PATCH()
	if err != nil {
		t.Fatal(err)
	}

	result := &StdResult{}
	err = resp.BindJSON(result)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)
}

// ------------------------------------------------------------------------------------------

func TestGet(t *testing.T) {
	requestAddr := runGoHTTPServer()

	type args struct {
		result interface{}
		url    string
		params KV
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantResult *StdResult
	}{
		{
			name: "get success",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/get",
				params: KV{"uid": 123},
			},
			wantErr: false,
			wantResult: &StdResult{
				Code: 0,
				Msg:  "ok",
				Data: "uid=123",
			},
		},
		{
			name: "get err",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/get_err",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
		{
			name: "get not found",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/notfound",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Get(tt.args.result, tt.args.url, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*StdResult).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	requestAddr := runGoHTTPServer()

	type args struct {
		result interface{}
		url    string
		params KV
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantResult *StdResult
	}{
		{
			name: "delete success",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/delete",
				params: KV{"uid": 123},
			},
			wantErr: false,
			wantResult: &StdResult{
				Code: 0,
				Msg:  "ok",
				Data: "uid=123",
			},
		},
		{
			name: "delete err",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/delete_err",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
		{
			name: "delete not found",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/notfound",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.result, tt.args.url, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*StdResult).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestPost(t *testing.T) {
	requestAddr := runGoHTTPServer()

	type args struct {
		result interface{}
		url    string
		body   interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult *StdResult
		wantErr    bool
	}{
		{
			name: "post success",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/post",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "post error",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/post_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Post(tt.args.result, tt.args.url, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*StdResult).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestPut(t *testing.T) {
	requestAddr := runGoHTTPServer()

	type args struct {
		result interface{}
		url    string
		body   interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult *StdResult
		wantErr    bool
	}{
		{
			name: "put success",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/put",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "put error",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/put_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Put(tt.args.result, tt.args.url, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*StdResult).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestPatch(t *testing.T) {
	requestAddr := runGoHTTPServer()

	type args struct {
		result interface{}
		url    string
		body   interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult *StdResult
		wantErr    bool
	}{
		{
			name: "patch success",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/patch",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "patch error",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/patch_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				result: &StdResult{},
				url:    requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &StdResult{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Patch(tt.args.result, tt.args.url, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*StdResult).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestRequest_Reset(t *testing.T) {
	req := &Request{
		method: http.MethodGet,
	}
	req.Reset()
	assert.Equal(t, "", req.method)
}

func TestRequest_Do(t *testing.T) {
	req := &Request{
		method: http.MethodGet,
		url:    "http://",
	}

	_, err := req.Do(http.MethodOptions, "")
	assert.Error(t, err)

	_, err = req.Do(http.MethodGet, map[string]interface{}{"foo": "bar"})
	assert.Error(t, err)
	_, err = req.Do(http.MethodDelete, "foo=bar")
	assert.Error(t, err)

	_, err = req.Do(http.MethodPost, &myBody{
		Name:  "foo",
		Email: "bar@gmail.com",
	})
	assert.Error(t, err)

	_, err = req.Response()
	assert.Error(t, err)

	err = requestErr(err)
	assert.Error(t, err)

	err = jsonParseErr(err)
	assert.Error(t, err)
}

func TestResponse_BodyString(t *testing.T) {
	resp := &Response{
		Response: nil,
		err:      nil,
	}

	_, err := resp.BodyString()
	assert.Error(t, err)

	resp.err = errors.New("error test")
	_, err = resp.BodyString()
	assert.Error(t, err)

	err = resp.Error()
	assert.Error(t, err)
}

func TestError(t *testing.T) {
	req := &Request{}
	req.SetParam("foo", "bar")
	req.SetParam("foo3", make(chan string))
	req.SetParams(map[string]interface{}{"foo2": "bar2"})
	req.SetBody("foo")
	req.SetTimeout(time.Second * 10)
	req.CustomRequest(func(req *http.Request, data *bytes.Buffer) {
		fmt.Println("customRequest")
	})
	req.SetURL("http://127.0.0.1:0")

	resp, err := req.pull()
	assert.Error(t, err)

	req.method = http.MethodPost
	resp, err = req.push()
	assert.Error(t, err)

	_, err = resp.ReadBody()
	assert.Error(t, err)

	err = resp.BindJSON(nil)
	assert.Error(t, err)

	err = notOKErr(resp)
	assert.Error(t, err)

	err = do(http.MethodPost, nil, "", nil)
	assert.Error(t, err)
	err = do(http.MethodPost, &StdResult{}, "http://127.0.0.1:0", nil, KV{"foo": "bar"})
	assert.Error(t, err)

	err = gDo(http.MethodGet, nil, "http://127.0.0.1:0", nil)
	assert.Error(t, err)
}
