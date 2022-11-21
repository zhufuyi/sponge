## response

The wrapper gin returns json data in the same format.

<br>

## Example of use

- `Output`  return a compatible http status code.
- `Success` and `Error` return a uniform status code of 200, with a custom status code in data.code

all requests return a uniform json

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

```go
    // c is *gin.Context

    // return success
    response.Success(c)
    // return success and return data
    response.Success(c, gin.H{"users":users})

    // return failure
    response.Error(c, errcode.SendEmailErr)
    // returns a failure and returns the data
    response.Error(c,  errcode.SendEmailErr, gin.H{"user":user})
```