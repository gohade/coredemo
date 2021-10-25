package orm

import (
	"context"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

// HadeApp 代表hade框架的App实现
type HadeGorm struct {
	container framework.Container // 服务容器
	dbs       map[string]*gorm.DB // key为dsn, value为gorm.DB（连接池）

	lock *sync.RWMutex
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

func (app *HadeGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
	logger := app.container.MustMake(contract.LogKey).(contract.Log)

	config := GetBaseConfig(app.container)
	logService := app.container.MustMake(contract.LogKey).(contract.Log)
	ormLogger := NewOrmLogger(logService)
	config.GormConfig = &gorm.Config{}
	config.GormConfig.Logger = ormLogger

	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	app.lock.RLock()
	if db, ok := app.dbs[config.Dsn]; ok {
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	app.lock.Lock()
	defer app.lock.Unlock()

	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config.GormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config.GormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config.GormConfig)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config.GormConfig)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config.GormConfig)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}
	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(liftTime)
		}
	}

	if err != nil {
		app.dbs[config.Dsn] = db
	}

	return db, err
}
