## gohttp

http请求客户端，只支持返回json格式。

<br>

## 使用示例

### 标准CURD

Get、Delete请求示例

```go
	req := gohttp.Request{}
	req.SetURL("http://localhost:8080/user")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})
	req.SetParams(gohttp.KV{
		"id": 123,
	})

	resp, err := req.GET()
	// resp, err := req.Delete()

	result := &gohttp.StdResult{} // 可以定义其他结构体接收数据
	err = resp.BindJSON(result)
```

<br>

Post、Put、Patch请求示例

```go
	req := gohttp.Request{}
	req.SetURL("http://localhost:8080/user")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})

	// body为结构体
    type User struct{
        Name string
        Email string
    }
    body := &User{"foo", "foo@bar.com"}
    req.SetJSONBody(body)
    // 或者 body为json
    // req.SetBody(`{"name":"foo", "email":"foo@bar.com"}`)

	resp, err := req.Post()
	// resp, err := req.Put()
	// resp, err := req.Patch()

	result := &gohttp.StdResult{} // 可以定义其他结构体接收数据
	err = resp.BindJSON(result)
```

<br>

### 简化版CRUD

不支持设置header、超时等

```go
    url := "http://localhost:8080/user"
    params := gohttp.KV{"id":123}
    result := &gohttp.StdResult{} // 可以定义其他结构体接收数据

    // Get
    err := gohttp.Get(result, url)
    err := gohttp.Get(result, url, params)

    // Delete
    err := gohttp.Delete(result, url)
    err := gohttp.Delete(result, url, params)

    type User struct{
        Name string
        Email string
    }
    body := &User{"foo", "foo@bar.com"}

    // Post
    err := gohttp.Post(result, url, body)
    // Put
    err := gohttp.Put(result, url, body)
    // Patch
    err := gohttp.Patch(result, url, body)
```

