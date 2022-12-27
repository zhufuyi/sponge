package server

import (
	"fmt"
	"net/http"
	"time"
)

// RunHTTPServer run http server
func RunHTTPServer(addr string) {
	initRecord()

	router := NewRouter()
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("listen server error: %v", err))
	}
}
