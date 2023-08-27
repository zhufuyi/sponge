## jwt

token generation and validation.

<br>

## Example of use

```go
    import "github.com/zhufuyi/sponge/pkg/gwt"

	jwt.Init(
		jwt.WithSigningKey("123456"),   // key
		jwt.WithExpire(time.Hour), // expiry time
		// jwt.WithSigningMethod(jwt.HS512), // encryption method, default is HS256, can be set to HS384, HS512
	)

	uid := "123"
	// generate token
	token, err := jwt.GenerateToken(uid)

    // verify token
	v, err := jwt.VerifyToken(token)
	if v.Uid != uid{
	    return
	}
```
