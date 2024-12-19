## jwt

Generate and parse token based on [jwt](https://github.com/golang-jwt/jwt) library.

<br>

## Example of use

Example 1: common fields jwt

```go
    import "github.com/go-dev-frame/sponge/pkg/jwt"

	jwt.Init(
		// jwt.WithSigningKey("123456"),   // key
		// jwt.WithExpire(time.Hour), // expiry time
		// jwt.WithSigningMethod(jwt.HS512), // encryption method, default is HS256, can be set to HS384, HS512
	)

	uid := "123"
	name := "admin"

	// generate token
	token, err := jwt.GenerateToken(uid, name)
	// handle err

	// parse token
	claims, err := jwt.ParseToken(token)
	// handle err

	// verify
	if claims.Uid != uid || claims.Name != name {
		print("verify failed")
	    return
	}
```

<br>

Example 2: custom fields jwt

```go
    import "github.com/go-dev-frame/sponge/pkg/jwt"

	jwt.Init(
		// jwt.WithSigningKey("123456"),   // key
		// jwt.WithExpire(time.Hour), // expiry time
		// jwt.WithSigningMethod(jwt.HS512), // encryption method, default is HS256, can be set to HS384, HS512
	)

	fields := jwt.KV{"id": 123, "foo": "bar"}

	// generate token
	token, err := jwt.GenerateCustomToken(fields)
	// handle err

	// parse token
	claims, err := jwt.ParseCustomToken(token)
	// handle err

	// verify
	id, isExist1 := claims.Get("id")
	foo, isExist2 := claims.Get("foo")
	if !isExist1 || !isExist2 || int(id.(float64)) != fields["id"].(int) || foo.(string) != fields["foo"].(string) {
	    print("verify failed")
	    return
	}
```
