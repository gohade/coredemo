package main

import (
	"coredemo/framework"
	"time"
)

func registerRouter(core *framework.Core) {
	core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))
}
