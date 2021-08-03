package main

import (
	"coredemo/framework"
	"fmt"
)

func SubjectAddController(c *framework.Context) error {
	c.Json(200, "ok, SubjectAddController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	c.Json(200, "ok, SubjectListController")
	return nil
}

func SubjectDelController(c *framework.Context) error {
	c.Json(200, "ok, SubjectDelController")
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.Json(200, "ok, SubjectUpdateController")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	subjectId, _ := c.ParamInt("id", 0)
	c.Json(200, "ok, SubjectGetController:"+fmt.Sprint(subjectId))

	return nil
}

func SubjectNameController(c *framework.Context) error {
	c.Json(200, "ok, SubjectNameController")
	return nil
}
