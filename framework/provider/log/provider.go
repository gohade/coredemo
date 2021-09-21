package log

import (
	"github.com/gohade/hade/framework/provider/log/formatter"
	"strings"

	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/provider/log/services"
)

type HadeLogServiceProvider struct {
	framework.ServiceProvider

	driver string // driver

	// common config for log
	Formatter  contract.Formatter
	Level      contract.LogLevel
	CtxFielder contract.CtxFielder
}

// Register registe a new function for make a service instance
func (l *HadeLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	tcs, err := c.Make(contract.ConfigKey)
	if err != nil {
		return services.NewHadeConsoleLog
	}

	cs := tcs.(contract.Config)

	l.driver = strings.ToLower(cs.GetString("log.driver"))

	switch l.driver {
	case "single":
		return services.NewHadeSingleLog
	case "rotate":
		return services.NewHadeRotateLog
	case "console":
		return services.NewHadeConsoleLog
	default:
		return services.NewHadeConsoleLog
	}
}

// Boot will called when the service instantiate
func (l *HadeLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (l *HadeLogServiceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (l *HadeLogServiceProvider) Params(c framework.Container) []interface{} {
	// param sequence: level, ctxFielder, Formatter, map[string]string(folder/file)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	// 设置参数formatter
	if l.Formatter == nil {
		l.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				l.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				l.Formatter = formatter.TextFormatter
			}
		}
	}

	if l.Level == contract.UnknownLevel {
		l.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			l.Level = logLevel(configService.GetString("log.level"))
		}
	}

	return []interface{}{l.Level, l.CtxFielder, l.Formatter, c}
}

/// Name define the name for this service
func (l *HadeLogServiceProvider) Name() string {
	return contract.LogKey
}

// logLevel get level from string
func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
