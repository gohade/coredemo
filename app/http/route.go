package http

import (
	"github.com/gohade/hade/app/http/module/demo"
	"github.com/gohade/hade/framework/gin"
	"github.com/gohade/hade/framework/middleware"
	"github.com/gohade/hade/framework/middleware/static"
)

// Routes 绑定业务层路由
func Routes(r *gin.Engine) {

	r.Use(static.Serve("/", static.LocalFile("./dist", false)))
	r.Use(middleware.Trace())
	demo.Register(r)
}
