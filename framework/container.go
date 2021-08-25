package framework

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// 绑定一个服务提供者，如果关键字已经存在，会进行替换操作，但是会进行
	Bind(provider ServiceProvider) error
	// 关键字是否已经绑定服务提供者
	IsBind(key string) bool

	// 根据关键字获取一个服务，
	Make(key string) (interface{}, error)
	// MustMake to get a service provider which will not return error
	// if error，panic it
	// If you use it, make sure it will not panic, or you can cover it
	MustMake(key string) interface{}
	// MakeNew to new a service
	// The service Must not be singlton, it with params to make a new service
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// HadeContainer is instance of Container
type HadeContainer struct {
	Container
	// providers store many provider
	providers map[string]ServiceProvider
	// instance store instances
	instances map[string]interface{}
	// lock for container for change bind
	lock sync.RWMutex
}

// NewHadeContainer is new instance
func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

func (hade *HadeContainer) PrintList() []string {
	ret := []string{}
	for _, provider := range hade.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind make relationship between provider and contract
func (hade *HadeContainer) Bind(provider ServiceProvider) error {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	key := provider.Name()

	hade.providers[key] = provider

	// if provider is not defer
	if provider.IsDefer() == false {
		if err := provider.Boot(hade); err != nil {
			return err
		}
		params := provider.Params()
		method := provider.Register(hade)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		hade.instances[key] = instance
	}
	return nil
}

func (hade *HadeContainer) IsBind(key string) bool {
	return hade.findServiceProvider(key) != nil
}

func (hade *HadeContainer) findServiceProvider(key string) ServiceProvider {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	if sp, ok := hade.providers[key]; ok {
		return sp
	}
	return nil
}

func (hade *HadeContainer) Make(key string) (interface{}, error) {
	return hade.make(key, nil, false)
}

func (hade *HadeContainer) MustMake(key string) interface{} {
	serv, err := hade.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

func (hade *HadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return hade.make(key, params, true)
}

func (hade *HadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	hade.lock.RLock()
	defer hade.lock.RUnlock()
	// check has Register
	sp := hade.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		// force new a
		if err := sp.Boot(hade); err != nil {
			return nil, err
		}
		if params == nil {
			params = sp.Params()
		}
		method := sp.Register(hade)
		ins, err := method(params...)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		return ins, nil
	}

	// not force New
	// bool ins
	if ins, ok := hade.instances[key]; ok {
		return ins, nil
	}

	if err := sp.Boot(hade); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params()
	}
	method := sp.Register(hade)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	hade.instances[key] = ins
	return ins, nil
}
