package kernel

import (
	"github.com/gohade/hade/framework/gin"
)

// 引擎服务
type HadeKernelService struct {
	engine *gin.Engine
}

// 初始化web引擎服务实例
func NewHadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &HadeKernelService{engine: httpEngine}, nil
}

// 返回web引擎
func (s *HadeKernelService) HttpEngine() *gin.Engine {
	return s.engine
}
