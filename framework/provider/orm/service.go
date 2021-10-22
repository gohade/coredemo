package orm

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
)

// HadeApp 代表hade框架的App实现
type HadeGorm struct {
	container framework.Container // 服务容器

	configPath string
	dbs        map[string]*gorm.DB // 容器服务

	lock *sync.RWMutex
}

func (hg *HadeGorm) SetConfigPath(config string) {
	hg.configPath = config
}

// DBOption 代表初始化的时候的选项
type DBOption func(orm contract.ORMService)

// WithDBConfig 从database.yaml中获取配置
func WithConfigPath(path string) DBOption {
	return func(orm contract.ORMService) {
		hadeGorm := orm.(*HadeGorm)
		hadeGorm.SetConfigPath(path)
	}
}

func WithGormConfig(config *gorm.Config) DBOption {
	return func(orm contract.ORMService) {
	}
}

// NewHadeGorm 代表启动容器
func NewHadeGorm(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	return &HadeGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil
}

func (app *HadeGorm) GetDB(option ...DBOption) (*gorm.DB, error) {
	for _, op := range option {
		op(app)
	}

	if app.configPath == "" {
		WithConfigPath("database.default")
	}

	configService := app.container.MustMake(contract.ConfigKey).(contract.Config)
	config := &Config{}
	if err := configService.Load(app.configPath, config); err != nil {
		return nil, err
	}
	if config.Dsn == "" {
		config.Dsn = config.FormatDsn()
	}

	if db, ok := app.dbs[config.Dsn]; ok {
		return db, nil
	}

	logService := app.container.MustMake(contract.LogKey).(contract.Log)
	ormLogger := NewOrmLogger(logService)
	gormConfig := &gorm.Config{
		ConnPool: nil,
		Logger:   ormLogger,
	}

	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), gormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), gormConfig)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), gormConfig)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), gormConfig)
	}
	if err != nil {
		app.dbs[config.Dsn] = db
	}
	return db, err
}
