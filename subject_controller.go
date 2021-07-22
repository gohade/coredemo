package main

import "coredemo/framework"

func SubjectAddController(c *framework.Context) error {
	c.Json(200, "ok, SubjectAddController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	c.Json(200, "ok, SubjectListController")
	return nil
}
