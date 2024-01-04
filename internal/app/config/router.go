package config

import (
	"github.com/gin-gonic/gin"
	"seg-red-file/internal/app/common"
	"seg-red-file/internal/app/controller"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(common.GlobalErrorHandler())
	r.NoRoute(common.HandleNoRoute())
	v1 := r.Group("/api/v1")

	controller.NewFileController(v1)

	return r
}
