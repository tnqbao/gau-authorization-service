package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-upload-service/controller"
)

func SetupRouter(ctrl *controller.Controller) *gin.Engine {
	r := gin.Default()
	// apiRoutes := r.Group("/")
	// {
	// 	apiRoutes.POST("/upload", ctrl.UploadFile)

	// }
	return r
}