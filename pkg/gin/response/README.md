## response

封装gin返回json数据插件。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/response

<br>

## 使用示例

`Output`函数返回兼容http状态码

`Success`和`Error`统一返回状态码200，在data.code自定义状态码

所有请求统一返回json

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

```go
    // c是*gin.Context

    // 返回成功
    response.Success(c)
    // 返回成功，并返回数据
    response.Success(c, gin.H{"users":users})

    // 返回失败
    response.Error(c, errcode.SendEmailErr)
    // 返回失败，并返回数据
    response.Error(c,  errcode.SendEmailErr, gin.H{"user":user})
```