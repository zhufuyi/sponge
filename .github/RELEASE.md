## Change log

1. Adjust errcode package, support custom error message.

example:

```go
//http code
    // code(10003) and message
    ecode.InvalidParams.Err("custom error message")
    // code(400) and message
    ecode.InvalidParams.ErrToHTTP("custom error message")

// grpc code
    // code(30003) and message
    ecode.StatusInvalidParams.Err("custom error message")
    // code(3) and message
    ecode.StatusInvalidParams.ToRPCErr("custom error message")
    // code(30003) and message, use in grpc-gateway
    ecode.StatusInvalidParams.ErrToHTTP("custom error message")
```

2. Optimize print log, grpc support custom marshal data.

3. Adjust some code.
