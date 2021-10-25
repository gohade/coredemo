package orm

import (
	"context"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"gorm.io/gorm"
)

func GetBaseConfig(c framework.Container) *contract.DBConfig {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.DBConfig{}
	err := configService.Load("database", config)
	if err != nil {
		logService.Error(context.Background(), "parse database config error", nil)
		return nil
	}
	return config
}

// WithConfigPath 加载配置文件地址
func WithConfigPath(configPath string) contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}
}

// WithGormConfig 表示自行配置Gorm的配置信息
func WithGormConfig(gormConfig *gorm.Config) contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		if gormConfig.Logger == nil {
			gormConfig.Logger = config.GormConfig.Logger
		}
		config.GormConfig = gormConfig
		return nil
	}
}
