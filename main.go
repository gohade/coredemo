package main

import (
	"coredemo/framework"
	"net/http"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		Handler: core,
	}
	server.ListenAndServe()
}
