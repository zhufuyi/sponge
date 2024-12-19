## frontend

Embed front-end web static files in gin and add routing.

<br>

### Example of use

```go
import "github.com/go-dev-frame/sponge/pkg/gin/frontend"

//go:embed user
var staticFS embed.FS

func setFrontendRouter(r *gin.Engine) error {
	var (
		isUseEmbedFS   = true
		htmlDir        = "user/home"
		configFile     = "user/home/config.js"
		modifyConfigFn = func(content []byte) []byte {
			// modify config code
			return content
		}
	)

	err := frontend.New(staticFS, isUseEmbedFs, htmlDir, configFile, modifyConfigFn).SetRouter(r)
	if err != nil {
		panic(err)
	}
}
```

Note: in the above example, `user` is the directory where the front-end is located, the static file index.html is in the `user/home` directory. If isUseEmbedFS is false and apiBaseUrl is set in the configuration file, cross-host access is supported.
