package framework

// IResponse代表返回方法
type IResponse interface {
	// Json输出
	Json(obj interface{}) IResponse

	// Jsonp输出
	Jsonp(obj interface{}) IResponse

	//xml输出
	Xml(obj interface{}) IResponse

	// html输出
	Html(template string, obj interface{}) IResponse

	// string
	Text(format string, values ...interface{}) IResponse

	// 重定向
	Redirect(path string) IResponse

	// header
	SetHeader(key string, val string) IResponse

	// Cookie
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// 基础操作
	Abort()

	// 设置状态码
	SetStaus(code int) IResponse

	// 设置200状态
	SetOkStatus() IResponse
}

// Jsonp输出
func (ctx *Context) Jsonp(obj interface{}) IResponse {
	panic("not implemented") // TODO: Implement
}

//xml输出
func (ctx *Context) Xml(obj interface{}) IResponse {
	panic("not implemented") // TODO: Implement
}

// html输出
func (ctx *Context) Html(template string, obj interface{}) IResponse {
	panic("not implemented") // TODO: Implement
}

// 重定向
func (ctx *Context) Redirect(path string) IResponse {
	panic("not implemented") // TODO: Implement
}

// header
func (ctx *Context) SetHeader(key string, val string) IResponse {
	panic("not implemented") // TODO: Implement
}

// Cookie
func (ctx *Context) SetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) IResponse {
	panic("not implemented") // TODO: Implement
}

// 基础操作
func (ctx *Context) Abort() {
	panic("not implemented") // TODO: Implement
}

// 设置状态码
func (ctx *Context) SetStaus(code int) IResponse {
	panic("not implemented") // TODO: Implement
}

// 设置200状态
func (ctx *Context) SetOkStatus() IResponse {
	panic("not implemented") // TODO: Implement
}
