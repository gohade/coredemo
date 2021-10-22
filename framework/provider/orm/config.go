package orm

import (
	"github.com/go-sql-driver/mysql"
	"net"
	"time"
)

type Config struct {
	WriteTimeout time.Duration  `yaml:"write_timeout"`
	Loc          *time.Location `yaml:"loc"`
	Port         string         `yaml:"port"`
	ReadTimeout  time.Duration  `yaml:"read_timeout"`
	Charset      string         `yaml:"charset"`
	ParseTime    bool           `yaml:"parse_time"`
	Protocol     string         `yaml:"protocol"`
	Dsn          string         `yaml:"dsn"`
	Database     string         `yaml:"database"`
	Collation    string         `yaml:"collation"`
	Timeout      time.Duration  `yaml:"timeout"`
	Username     string         `yaml:"username"`
	Password     string         `yaml:"password"`
	Driver       string         `yaml:"driver"`
	Host         string         `yaml:"host"`
}

// FormatDsn 生成dsn
func (conf *Config) FormatDsn() string {
	driverConf := &mysql.Config{
		User:         conf.Username,
		Passwd:       conf.Password,
		Net:          conf.Protocol,
		Addr:         net.JoinHostPort(conf.Host, conf.Port),
		DBName:       conf.Database,
		Collation:    conf.Collation,
		Loc:          conf.Loc,
		Timeout:      conf.Timeout,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		ParseTime:    conf.ParseTime,
	}
	return driverConf.FormatDSN()
}
