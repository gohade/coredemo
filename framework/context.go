package framework

import "net/http"

type Context struct {
	r *http.Request
	w http.ResponseWriter
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.r
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.w
}
