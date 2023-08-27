### gohttp

The http request client, which only supports returning json format.

<br>

### Example of use

#### Standard CURD

Get, Delete request example.

```go
    import "github.com/zhufuyi/sponge/pkg/gohttp"

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

	result := &gohttp.StdResult{} // other structures can be defined to receive data
	err = resp.BindJSON(result)
```

<br>

Post, Put, Patch request example.

```go
    import "github.com/zhufuyi/sponge/pkg/gohttp"

	req := gohttp.Request{}
	req.SetURL("http://localhost:8080/user")
	req.SetHeaders(map[string]string{
		"Authorization": "Bearer token",
	})

	// body is a structure
    type User struct{
        Name string
        Email string
    }
    body := &User{"foo", "foo@bar.com"}
    req.SetJSONBody(body)
    // or body as json
    // req.SetBody(`{"name":"foo", "email":"foo@bar.com"}`)

	resp, err := req.Post()
	// resp, err := req.Put()
	// resp, err := req.Patch()

	result := &gohttp.StdResult{} // other structures can be defined to receive data
	err = resp.BindJSON(result)
```

<br>

#### simplified version of CRUD

No support for setting header, timeout, etc.

```go
    import "github.com/zhufuyi/sponge/pkg/gohttp"

    url := "http://localhost:8080/user"
    params := gohttp.KV{"id":123}
    result := &gohttp.StdResult{} // other structures can be defined to receive data

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
