package framework

import (
	"mime/multipart"
)

// 代表请求包含的方法
type IRequest interface {
	// 请求地址url中带的参数
	// 形如: foo.com?a=1&b=bar&c[]=bar
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	// 路由匹配中带的参数
	// 形如 /book/:id
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	ParamStringSlice(key string, def []string) ([]string, bool)
	Param(key string) interface{}

	// form表单中带的参数
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	// json body
	BindJson(obj interface{}) error

	// xml body
	BindXml(obj interface{}) error

	// 其他格式
	GetRawData() ([]byte, error)

	// 基础信息
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	// header
	Headers() map[string]string
	Header(key string) (string, bool)

	// cookie
	Cookes() map[string]string
	Cookie(key string) (string, bool)
}

// // #region query url
// func (ctx *Context) QueryInt(key string, def int) (int, bool) {
// 	params := ctx.QueryAll()
// 	if vals, ok := params[key]; ok {
// 		len := len(vals)
// 		if len > 0 {
// 			intval, err := strconv.Atoi(vals[len-1])
// 			if err != nil {
// 				return def, false
// 			}
// 			return intval, true
// 		}
// 	}
// 	return def, false
// }

// func (ctx *Context) QueryString(key string, def string) string {
// 	params := ctx.QueryAll()
// 	if vals, ok := params[key]; ok {
// 		len := len(vals)
// 		if len > 0 {
// 			return vals[len-1]
// 		}
// 	}
// 	return def
// }

// func (ctx *Context) QueryArray(key string, def []string) []string {
// 	params := ctx.QueryAll()
// 	if vals, ok := params[key]; ok {
// 		return vals
// 	}
// 	return def
// }

// func (ctx *Context) QueryAll() map[string][]string {
// 	if ctx.request != nil {
// 		return map[string][]string(ctx.request.URL.Query())
// 	}
// 	return map[string][]string{}
// }

// // #endregion

// 请求地址url中带的参数
// 形如: foo.com?a=1&b=bar&c[]=bar
func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Query(key string) interface{} {
	panic("not implemented") // TODO: Implement
}

// 路由匹配中带的参数
// 形如 /book/:id
func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ParamStringSlice(key string, def []string) ([]string, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Param(key string) interface{} {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Form(key string) interface{} {
	panic("not implemented") // TODO: Implement
}

// xml body
func (ctx *Context) BindXml(obj interface{}) error {
	panic("not implemented") // TODO: Implement
}

// 其他格式
func (ctx *Context) GetRawData() ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

// 基础信息
func (ctx *Context) Uri() string {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Method() string {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Host() string {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) ClientIp() string {
	panic("not implemented") // TODO: Implement
}

// header
func (ctx *Context) Headers() map[string]string {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Header(key string) (string, bool) {
	panic("not implemented") // TODO: Implement
}

// cookie
func (ctx *Context) Cookes() map[string]string {
	panic("not implemented") // TODO: Implement
}

func (ctx *Context) Cookie(key string) (string, bool) {
	panic("not implemented") // TODO: Implement
}
