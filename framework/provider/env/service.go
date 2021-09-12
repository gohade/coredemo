package env

import (
	"bufio"
	"bytes"
	"github.com/gohade/hade/framework/contract"
	"io"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type HadeEnv struct {
	folder string // represent env folder

	maps map[string]string
}

// NewHadeEnv have two params: folder and env
// for example: NewHadeEnv("/envfolder/")
// It will read file: /envfolder/.env
// The file have format XXX=XXX
func NewHadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewHadeEnv param error")
	}

	folder := params[0].(string)

	hadeEnv := &HadeEnv{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// parse .env
	file := path.Join(folder, ".env")
	_, err := os.Stat(file)
	if err == nil {
		fi, err := os.Open(file)
		if err == nil {
			defer fi.Close()

			br := bufio.NewReader(fi)
			for {
				line, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				s := bytes.SplitN(line, []byte{'='}, 2)
				if len(s) < 2 {
					continue
				}
				key := string(s[0])
				val := string(s[1])
				hadeEnv.maps[key] = val
			}
		}
	}

	// 获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		hadeEnv.maps[pair[0]] = pair[1]
	}
	return hadeEnv, nil
}

// AppEnv get current environment
func (en *HadeEnv) AppEnv() string {
	return en.Get("APP_ENV")
}

// IsExist check setting is exist
func (en *HadeEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// Get environment setting, if not exist, return ""
func (en *HadeEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

// All return all settings
func (en *HadeEnv) All() map[string]string {
	return en.maps
}
