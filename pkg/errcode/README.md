## errcode

Error codes usually include system-level error codes and business-level error codes, consisting of a total of 5 decimal digits, e.g. 20101

| First digit                                                                          | Middle two digits | Last two digits |
|:-------------------------------------------------------------------------------------|:-------|:-------|
| For http error codes, 2 indicates a business-level error (1 is a system-level error) | Service Module Code | Specific error codes |
| For grpc error codes, 4 indicates a business-level error (3 is a system-level error) | Service Module Code | Specific error codes |

- Error levels occupy one digit: 1 (http) and 3 (grpc) indicate system-level errors, 2 (http) and 4 (grpc) indicate business-level errors, usually caused by illegal user operations.
- Double-digit service modules: A large system usually has no more than two service modules; if it exceeds that, it's time to split the system.
- Error codes take up two digits: prevents a module from being customised with too many error codes, which are not well maintained later.

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
