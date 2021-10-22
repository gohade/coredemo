package contract

import (
	"github.com/gohade/hade/framework/provider/orm"
	"gorm.io/gorm"
)

// ORMKey 代表 ORM的服务
const ORMKey = "hade:orm"

// ORMService 表示传入的参数
type ORMService interface {
	GetDB(option ...orm.DBOption) (*gorm.DB, error)
}
