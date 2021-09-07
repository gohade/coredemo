package contract

import (
	"github.com/gohade/hade/framework/gin"
)

const KernelKey = "hade:kernel"

// Kernel 接口提供框架最核心的结构
type Kernel interface {
	// HttpEngine 提供gin的Engine结构
	HttpEngine() *gin.Engine
}
