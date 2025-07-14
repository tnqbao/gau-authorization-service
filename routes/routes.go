package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-authorization-service/controller"
	"github.com/tnqbao/gau-authorization-service/middlewares"
)

func SetupRouter(ctrl *controller.Controller) *gin.Engine {
	newMiddlewares, err := middlewares.NewMiddlewares(ctrl)
	if err != nil {
		panic("Failed to initialize middlewares: " + err.Error())
	}
	if newMiddlewares == nil {
		panic("Failed to initialize middlewares")
	}
	r := gin.Default()
	// apiRoutes := r.Group("/")
	// {
	// 	apiRoutes.POST("/upload", ctrl.UploadFile)

	// }
	return r
}
