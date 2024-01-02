// Package frontend embeds the frontend static file and adds routing.
package frontend

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Note: the front-end static files must all be in the same directory, i.e. htmlDir

// when customizing customAddr, read static resources and save them to this directory.
var targetDir = "frontend"

// FrontEnd is the frontend router configuration
type FrontEnd struct {
	// directory where index.html is located, e.g. web/html.
	htmlDir string
	// defaultAddr is the default address of the back-end api request, e.g. http://localhost:8080.

	defaultAddr string
	// custom address，the default is defaultAddr, which is used to configure whether cross-host access is possible.
	customAddr string
	// configuration file in which js requests the address of the back-end api request, e.g. config.js.
	addrConfigFile string

	// embed.FS static resources.
	staticFS embed.FS
}

// New create a new frontend
func New(htmlDir string, defaultAddr string, customAddr string, addrConfigFile string, staticFS embed.FS) *FrontEnd {
	if defaultAddr == "" {
		panic("default address cannot be empty")
	}
	if customAddr == "" {
		customAddr = defaultAddr
	}
	htmlDir = strings.Trim(htmlDir, "/")
	return &FrontEnd{
		htmlDir:        htmlDir,
		defaultAddr:    defaultAddr,
		customAddr:     customAddr,
		addrConfigFile: addrConfigFile,
		staticFS:       staticFS,
	}
}

// SetRouter set frontend router
func (f *FrontEnd) SetRouter(r *gin.Engine) error {
	routerPath := fmt.Sprintf("/%s/index.html", f.htmlDir)
	// solve vue using history route 404 problem
	r.NoRoute(browserRefreshFS(f.staticFS, routerPath))

	relativePath := fmt.Sprintf("/%s/*filepath", f.htmlDir)

	// use embed file
	if f.isUseEmbedFS() {
		r.GET(relativePath, func(c *gin.Context) {
			staticServer := http.FileServer(http.FS(f.staticFS))
			staticServer.ServeHTTP(c.Writer, c.Request)
		})
		return nil
	}

	// use local file
	err := f.saveFSToLocal()
	if err != nil {
		return err
	}
	r.GET(relativePath, func(c *gin.Context) {
		localFileDir := filepath.Join(targetDir, f.htmlDir)
		filePath := c.Param("filepath")
		c.File(localFileDir + filePath)
	})

	return nil
}

func (f *FrontEnd) isUseEmbedFS() bool {
	return f.defaultAddr == f.customAddr
}

func (f *FrontEnd) saveFSToLocal() error {
	_ = os.RemoveAll(filepath.Join(targetDir, f.htmlDir))
	time.Sleep(time.Millisecond * 10)

	// Walk through the embedded filesystem
	return fs.WalkDir(f.staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Create the corresponding directory structure locally
		localPath := filepath.Join(targetDir, path)
		if d.IsDir() {
			err := os.MkdirAll(localPath, 0755)
			if err != nil {
				return err
			}
		} else {
			// Read the file from the embedded filesystem
			content, err := fs.ReadFile(f.staticFS, path)
			if err != nil {
				return err
			}

			// replace config address
			if path == f.addrConfigFile {
				content = bytes.ReplaceAll(content, []byte(f.defaultAddr), []byte(f.customAddr))
			}

			// Save the content to the local file
			err = os.WriteFile(localPath, content, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// solve vue using history route 404 problem, for embed.FS
func browserRefreshFS(efs embed.FS, path string) func(c *gin.Context) {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := efs.ReadFile(path)
			if err != nil {
				c.Writer.WriteHeader(404)
				_, _ = c.Writer.WriteString("Not Found")
				return
			}
			c.Writer.WriteHeader(200)
			c.Writer.Header().Add("Accept", "text/html")
			_, _ = c.Writer.Write(content)
			c.Writer.Flush()
		}
	}
}

// AutoOpenBrowser auto open browser
func AutoOpenBrowser(visitURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, visitURL)
	return exec.Command(cmd, args...).Start()
}
