package common

import "github.com/gin-gonic/gin"

type ModuleChild struct {
	Route   string
	Method  string
	Handles []gin.HandlerFunc
}

type ModuleOption struct {
	Name      string
	ChildList []ModuleChild
}

type ModuleOptionList []ModuleOption
