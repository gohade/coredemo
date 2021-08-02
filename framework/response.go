package framework

import "encoding/json"

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

// #region response

func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

// #endregion
