package contract

import (
	"gorm.io/gorm"
)

// ORMKey 代表 ORM的服务
const ORMKey = "hade:orm"

// ORMService 表示传入的参数
type ORMService interface {
	GetDB(option ...DBOption) (*gorm.DB, error)
}

// DBOption 代表初始化的时候的选项
type DBOption func(orm ORMService)
