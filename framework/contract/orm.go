package contract

import (
	"github.com/go-sql-driver/mysql"
	"github.com/gohade/hade/framework"
	"gorm.io/gorm"
	"net"
	"strconv"
	"time"
)

// ORMKey 代表 ORM的服务
const ORMKey = "hade:orm"

// ORMService 表示传入的参数
type ORMService interface {
	GetDB(option ...DBOption) (*gorm.DB, error)
}

// DBOption 代表初始化的时候的选项
type DBOption func(container framework.Container, config *DBConfig) error

type DBConfig struct {
	WriteTimeout    string `yaml:"write_timeout"`
	Loc             string `yaml:"loc"`
	Port            int    `yaml:"port"`
	ReadTimeout     string `yaml:"read_timeout"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parse_time"`
	Protocol        string `yaml:"protocol"`
	Dsn             string `yaml:"dsn"`
	Database        string `yaml:"database"`
	Collation       string `yaml:"collation"`
	Timeout         string `yaml:"timeout"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	ConnMaxIdle     int    `yaml:"conn_max_idle"`
	ConnMaxOpen     int    `yaml:"conn_max_open"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`
	ConnMaxIdletime string `yaml:"conn_max_idletime"`

	GormConfig *gorm.Config
}

// FormatDsn 生成dsn
func (conf *DBConfig) FormatDsn() (string, error) {
	port := strconv.Itoa(conf.Port)
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		return "", err
	}
	readTimeout, err := time.ParseDuration(conf.ReadTimeout)
	if err != nil {
		return "", err
	}
	writeTimeout, err := time.ParseDuration(conf.WriteTimeout)
	if err != nil {
		return "", err
	}
	location, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return "", err
	}
	driverConf := &mysql.Config{
		User:         conf.Username,
		Passwd:       conf.Password,
		Net:          conf.Protocol,
		Addr:         net.JoinHostPort(conf.Host, port),
		DBName:       conf.Database,
		Collation:    conf.Collation,
		Loc:          location,
		Timeout:      timeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ParseTime:    conf.ParseTime,
	}
	return driverConf.FormatDSN(), nil
}
