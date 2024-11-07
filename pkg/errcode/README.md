## errcode

Error codes usually include system-level error codes and business-level error codes, consisting of a total of 6 decimal digits, e.g. 200101

**Error code structure:**

| First digit                                                                                                                    | Middle three digits                  | Last two digits         |
|:-------------------------------------------------------------------------------------------------------------------------------|:-------------------------------------|:------------------------|
| `1` is http system-level error<br>`2` is http business-level error<br>`3` is grpc system-level error<br>`4` is grpc system-level error | Table or module number, range 1~1000 | Custom number, range 1~100 |

<br>

**Error code ranges:**

| Service Type | System-level Error Code Range | Business-level Error Code Range |
|:-------------|:------------------------------|:--------------------------------|
| http         | 100000 ~ 200000               | 200000 ~ 300000                 |
| grpc         | 300000 ~ 400000               | 400000 ~ 500000                 |

<br>

### Example of use

### Example of http error code usage

Web services created based on **SQL**, use the following error code:

```go
    import "github.com/zhufuyi/sponge/pkg/gin/response"

    // return error
    response.Error(c, ecode.InvalidParams)
    // rewrite error messages
    response.Error(c, ecode.InvalidParams.RewriteMsg("custom error message"))

    // convert error code to standard http status code
    response.Out(c, ecode.InvalidParams)
    // convert error code to standard http status code, and rewrite error messages
    response.Out(c, ecode.InvalidParams.RewriteMsg("custom error message"))
```

Web services created based on **Protobuf**, use the following error code:

```go
    // return error
    return nil, ecode.InvalidParams.Err()
    // rewrite error messages
    return nil, ecode.InvalidParams.Err("custom error message")

    // convert error code to standard http status code
    return nil, ecode.InvalidParams.ErrToHTTP()
    // convert error code to standard http status code, and rewrite error messages
    return nil, ecode.InvalidParams.ErrToHTTP("custom error message")
```

<br>

### Example of grpc error code usage

```go
    // return error
    return nil, ecode.StatusInvalidParams.Err()
    // rewrite error messages
    return nil, ecode.StatusInvalidParams.Err("custom error message")

    // convert error code to standard grpc status code
    return nil, ecode.StatusInvalidParams.ToRPCErr()
    // convert error code to standard grpc status code, and rewrite error messages
    return nil, ecode.StatusInvalidParams.ToRPCErr("custom error message")

    // convert error code to standard http status code
    return nil, ecode.StatusInvalidParams.ErrToHTTP()
    // convert error code to standard http status code, and rewrite error messages
    return nil, ecode.StatusInvalidParams.ErrToHTTP("custom error message")
```
