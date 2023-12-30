package config

import (
	"github.com/gin-gonic/gin"
	"seg-red-file/internal/app/common"
	"seg-red-file/internal/app/controller"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(common.GlobalErrorHandler())

	v1 := r.Group("/api/v1")

	fileCtrl := controller.NewFileController()
	fileCtrl.RegisterRoutes(v1)

	return r
}
