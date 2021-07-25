package main

import (
	"coredemo/framework"
	"time"
)

func UserLoginController(c *framework.Context) error {
	time.Sleep(5 * time.Second)
	c.Json(200, "ok, UserLoginController")
	return nil
}
