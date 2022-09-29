package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zhufuyi/sponge/internal/serverNameExample/routers"
	"github.com/zhufuyi/sponge/pkg/app"

	"github.com/gin-gonic/gin"
)

var _ app.IServer = (*httpServer)(nil)

type httpServer struct {
	addr   string
	server *http.Server
}

// Start http service
func (s *httpServer) Start() error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen server error: %v", err)
	}
	return nil
}

// Stop http service
func (s *httpServer) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
	return s.server.Shutdown(ctx)
}

// String comment
func (s *httpServer) String() string {
	return "http service, addr = " + s.addr
}

// NewHTTPServer creates a new web server
func NewHTTPServer(addr string, opts ...HTTPOption) app.IServer {
	o := defaultHTTPOptions()
	o.apply(opts...)

	if o.isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := routers.NewRouter()
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    o.readTimeout,
		WriteTimeout:   o.writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return &httpServer{
		addr:   addr,
		server: server,
	}
}
