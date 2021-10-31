package services

import (
	"github.com/pkg/errors"
	"time"
)

const (
	NoneDuration = time.Duration(-1)
)

var ErrKeyNotFound = errors.New("key not found")
var ErrTypeNotOk = errors.New("val type not ok")

func checkString(val interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}

	if str, ok := val.(string); ok {
		return str, nil
	}
	return "", ErrTypeNotOk
}
