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

    configPath string
    dbs        map[string]*gorm.DB // 容器服务
    gormConfig *gorm.Config        // gorm的配置文件，可以修改

    lock *sync.RWMutex
}

// WithDBConfig 从database.yaml中获取配置
func WithConfigPath(path string) contract.DBOption {
    return func(orm contract.ORMService) {
        hadeGorm := orm.(*HadeGorm)
        hadeGorm.configPath = path
    }
}

func WithGormConfig(config *gorm.Config) contract.DBOption {
    return func(orm contract.ORMService) {
        hadeGorm := orm.(*HadeGorm)
        hadeGorm.gormConfig = config
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

func (app *HadeGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
    for _, op := range option {
        op(app)
    }

    if app.configPath == "" {
        WithConfigPath("database.default")
    }

    configService := app.container.MustMake(contract.ConfigKey).(contract.Config)
    logger := app.container.MustMake(contract.LogKey).(contract.Log)

    config := GetBaseConfig(app.container)
    if err := configService.Load(app.configPath, config); err != nil {
        return nil, err
    }
    if config.Dsn == "" {
        dsn, err := config.FormatDsn()
        if err != nil {
            return nil, err
        }
        config.Dsn = dsn
    }

    if db, ok := app.dbs[config.Dsn]; ok {
        return db, nil
    }

    logService := app.container.MustMake(contract.LogKey).(contract.Log)
    ormLogger := NewOrmLogger(logService)
    gormConfig := app.gormConfig
    if gormConfig == nil {
        gormConfig = &gorm.Config{}
    }
    if gormConfig.Logger == nil {
        gormConfig.Logger = ormLogger
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
