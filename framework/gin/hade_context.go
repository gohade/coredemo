// Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package gin

import (
	"context"

	"github.com/gohade/hade/framework"
)

func (ctx *Context) BaseContext() context.Context {
	return ctx.Request.Context()
}

// engine实现container的绑定封装
func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

// context 实现container的几个封装

// 实现make的封装
func (ctx *Context) Make(key string) (interface{}, error) {
	return ctx.container.Make(key)
}

// 实现mustMake的封装
func (ctx *Context) MustMake(key string) interface{} {
	return ctx.container.MustMake(key)
}

// 实现makenew的封装
func (ctx *Context) MakeNew(key string, params []interface{}) (interface{}, error) {
	return ctx.container.MakeNew(key, params)
}
