package config

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"path/filepath"
)

type HadeConfigProvider struct{}

// Register registe a new function for make a service instance
func (provider *HadeConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeConfig
}

// Boot will called when the service instantiate
func (provider *HadeConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *HadeConfigProvider) IsDefer() bool {
	return true
}

// Params define the necessary params for NewInstance
func (provider *HadeConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{envFolder, envService.All()}
}

/// Name define the name for this service
func (provider *HadeConfigProvider) Name() string {
	return contract.ConfigKey
}
