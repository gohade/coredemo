package main

import (
	"coredemo/framework"
	"coredemo/framework/middleware"
	"net/http"
	"time"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Timeout(1 * time.Second))

	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}
