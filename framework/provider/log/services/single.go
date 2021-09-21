package services

import (
	"github.com/gohade/hade/framework/util"
	"os"
	"path/filepath"

	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"

	"github.com/pkg/errors"
)

type HadeSingleLog struct {
	HadeLog

	folder string
	file   string
	fd     *os.File
}

// NewHadeSingleLog params sequence: level, ctxFielder, Formatter, map[string]interface(folder/file)
func NewHadeSingleLog(params ...interface{}) (interface{}, error) {

	level := params[0].(contract.LogLevel)
	ctxFielder := params[1].(contract.CtxFielder)
	formatter := params[2].(contract.Formatter)

	c := params[3].(framework.Container)

	appService := c.MustMake(contract.AppKey).(contract.App)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	log := &HadeSingleLog{}
	log.SetLevel(level)
	log.SetCxtFielder(ctxFielder)
	log.SetFormatter(formatter)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}
	log.folder = folder
	if !util.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	log.file = "hade.log"
	if configService.IsExist("log.file") {
		log.file = configService.GetString("log.file")
	}

	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}

	log.SetOutput(fd)
	log.c = c

	return log, nil
}

func (l *HadeSingleLog) SetFile(file string) {
	l.file = file
}

func (l *HadeSingleLog) SetFolder(folder string) {
	l.folder = folder
}
