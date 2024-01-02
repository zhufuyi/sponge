## frontend

Embed front-end web static files in gin and add routing.

<br>

### Example of use

```go1
import "github.com/zhufuyi/sponge/pkg/gin/frontend"

//go:embed user
var staticFS embed.FS

func setFrontendRouter(r *gin.Engine) error {
	var (
		// index.html file path, also the routing of access
		htmlPath        = "user/home"
		// file setting
		addrConfigFile = "user/home/config.js"

		// addr setting, the set address is in the addrConfigFile file
		defaultAddr    = "http://localhost:8080"
		// if cross-service is required, fill in the address of the server where the service is deployed here, e.g. http://192.168.3.37:8080
		customAddr     = ""
	)

	return frontend.New(htmlPath, defaultAddr, customAddr, addrConfigFile, staticFS).SetRouter(r)
}
```

Note: in the above example, `user` is the directory where the front-end is located, the static file index.html is in the `user/home` directory, if customAddr is empty and the default address is `http://localhost:8080`, then the access to the index.html routing address is `http:// localhost:8080/user/home`. If you set customAddr to `http://192.168.3.37:8080`, the index.html routing address will be `http://192.168.3.37:8080/user/home`.
