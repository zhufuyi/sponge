## jwt

token生成和验证。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/jwt

<br>

## 使用示例

```go
	jwt.Init(
		jwt.WithSigningKey("123456"),   // 密钥
		jwt.WithExpire(time.Hour), // 过期时间
		// jwt.WithSigningMethod(jwt.HS512), // 加密方法，默认是HS256，可以设置为HS384、HS512
	)

	uid := "123"
	// 生成token
	token, err := jwt.GenerateToken(uid)

    // 验证token
	v, err := jwt.VerifyToken(token)
	if v.Uid != uid{
	    return
	}
```