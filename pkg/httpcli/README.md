### httpcli

`httcli` is a simple HTTP request client, which only supports returning json format.

<br>

### Example of use

#### Request way 1

```go
    import "github.com/zhufuyi/sponge/pkg/httpcli"

    type User struct{
        Name string
        Email string
    }

    url := "http://localhost:8080/user"
    params := httpcli.KV{"id":123}
    headers := map[string]string{"Authorization": "Bearer token"}
    body := &User{"foo", "foo@bar.com"}
    result := &httpcli.StdResult{} // other structures can be defined to receive data

    var err error

    // Get
    err = httpcli.Get(result, url)
    err = httpcli.Get(result, url, httpcli.WithParams(params))
    err = httpcli.Get(result, url, httpcli.WithParams(params), httpcli.WithHeaders(headers))

    // Delete
    err = httpcli.Delete(result, url)
    err = httpcli.Delete(result, httpcli.WithParams(params))
    err = httpcli.Delete(result, httpcli.WithParams(params), httpcli.WithHeaders(headers))

    // Post
    err = httpcli.Post(result, url, body)
    err = httpcli.Post(result, url, body, httpcli.WithParams(params))
    err = httpcli.Delete(result, httpcli.WithParams(params), httpcli.WithHeaders(headers))
    // Put
    err := httpcli.Put(result, url, body)
    // Patch
    err := httpcli.Patch(result, url, body)
```

<br>

#### Request way 2

Get, Delete request example.

```go
    import "github.com/zhufuyi/sponge/pkg/httpcli"

    url := "http://localhost:8080/user"
    headers := map[string]string{"Authorization": "Bearer token"}
    params := httpcli.KV{"id": 123}

    cli := httpcli.New().SetURL(url).SetHeaders(headers).SetParams(params)

    // Get
    resp, err := cli.GET()
    // Delete
    // resp, err := cli.Delete()

    defer resp.Body.Close()

    result := &httpcli.StdResult{} // other structures can be defined to receive data
    err = resp.BindJSON(result)
```

<br>

Post, Put, Patch request example.

```go
    import "github.com/zhufuyi/sponge/pkg/httpcli"


    type User struct{
        Name string
        Email string
    }

    body := &User{"foo", "foo@bar.com"}
    url := "http://localhost:8080/user"
    headers := map[string]string{"Authorization": "Bearer token"}

    cli := httpcli.New().SetURL(url).SetHeaders(headers).SetBody(body)

    // Post
    resp, err := cli.Post()
    // Put
    // resp, err := cli.Put()
    // Patch
    // resp, err := cli.Patch()

   defer resp.Body.Close()

    result := &httpcli.StdResult{} // other structures can be defined to receive data
    err = resp.BindJSON(result)
```
