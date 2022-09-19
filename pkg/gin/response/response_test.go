package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/zhufuyi/sponge/pkg/errcode"

	"github.com/gin-gonic/gin"
)

var (
	requestAddr string
	wantCode    int
	wantData    interface{}
	wantErrInfo *errcode.Error
)

func init() {
	port, _ := getAvailablePort()
	requestAddr = fmt.Sprintf("http://localhost:%d", port)
	addr := fmt.Sprintf(":%d", port)

	r := gin.Default()
	r.GET("/hello1", func(c *gin.Context) { Output(c, wantCode, wantData) })
	r.GET("/hello2", func(c *gin.Context) { Success(c, wantData) })
	r.GET("/hello3", func(c *gin.Context) { Error(c, wantErrInfo) })

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

// 获取可用端口
func getAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}

func do(method string, url string, body interface{}) ([]byte, error) {
	var (
		resp        *http.Response
		err         error
		contentType = "application/json"
	)

	v, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodGet:
		resp, err = http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	case http.MethodPost:
		resp, err = http.Post(url, contentType, bytes.NewReader(v))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	case http.MethodDelete, http.MethodPut, http.MethodPatch:
		req, err := http.NewRequest(method, url, bytes.NewReader(v))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	default:
		return nil, fmt.Errorf("%s method not supported", method)
	}

	return io.ReadAll(resp.Body)
}

func get(url string) ([]byte, error) {
	return do(http.MethodGet, url, nil)
}

func delete(url string) ([]byte, error) {
	return do(http.MethodDelete, url, nil)
}

func post(url string, body interface{}) ([]byte, error) {
	return do(http.MethodPost, url, body)
}

func put(url string, body interface{}) ([]byte, error) {
	return do(http.MethodPut, url, body)
}

func patch(url string, body interface{}) ([]byte, error) {
	return do(http.MethodPatch, url, body)
}

// ------------------------------------------------------------------------------------------

func TestRespond(t *testing.T) {
	type args struct {
		url  string
		code int
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "respond 200",
			args: args{
				url:  requestAddr + "/hello1",
				code: http.StatusOK,
				data: gin.H{"name": "zhangsan"},
			},
			wantErr: false,
		},
		{
			name: "respond 400",
			args: args{
				url:  requestAddr + "/hello1",
				code: http.StatusBadRequest,
				data: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantCode = tt.args.code
			wantData = tt.args.data
			data, err := get(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("http.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%s", data)
			var resp = &Result{}
			err = json.Unmarshal(data, resp)
			if err != nil {
				t.Error(err)
				return
			}
			if resp.Code != wantCode {
				t.Errorf("%s, got = %v, want %v", tt.name, resp.Code, wantCode)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	type args struct {
		url  string
		code int
		data interface{}
		ei   *errcode.Error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				url:  requestAddr + "/hello2",
				code: http.StatusOK,
				data: gin.H{"name": "zhangsan"},
				ei:   errcode.Success,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantData = tt.args.data
			wantErrInfo = tt.args.ei
			data, err := get(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("http.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%s", data)
			var resp = &Result{}
			err = json.Unmarshal(data, resp)
			if err != nil {
				t.Error(err)
				return
			}
			if resp.Code != wantErrInfo.Code() && resp.Msg != wantErrInfo.Msg() {
				t.Errorf("%s, got = %v, want %v", tt.name, resp, wantErrInfo)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		url  string
		code int
		data interface{}
		ei   *errcode.Error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unauthorized",
			args: args{
				url:  requestAddr + "/hello3",
				code: http.StatusOK,
				data: nil,
				ei:   errcode.Unauthorized,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErrInfo = tt.args.ei
			data, err := get(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("http.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%s", data)
			var resp = &Result{}
			err = json.Unmarshal(data, resp)
			if err != nil {
				t.Error(err)
				return
			}
			if resp.Code != wantErrInfo.Code() && resp.Msg != wantErrInfo.Msg() {
				t.Errorf("%s, got = %v, want %v", tt.name, resp, wantErrInfo)
			}
		})
	}
}
