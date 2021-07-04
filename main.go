package main

import (
	"coredemo/framework"
	"net/http"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
	}
	server.ListenAndServe()
}
