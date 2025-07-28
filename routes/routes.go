package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-authorization-service/controller"
	"github.com/tnqbao/gau-authorization-service/middlewares"
)

func SetupRouter(ctrl *controller.Controller) *gin.Engine {
	useMiddlewares, err := middlewares.NewMiddlewares(ctrl)
	if err != nil {
		panic("Failed to initialize middlewares: " + err.Error())
	}
	if useMiddlewares == nil {
		panic("Failed to initialize middlewares")
	}
	r := gin.Default()
	apiRoutes := r.Group("/api/v2/authorization", useMiddlewares.CORSMiddleware)
	{
		apiRoutes.POST("/token", useMiddlewares.PrivateMiddlewares, ctrl.CreateNewToken)
		apiRoutes.GET("/token/renew", ctrl.RenewAccessToken)
		apiRoutes.GET("/token/validate", ctrl.CheckAccessToken)
		apiRoutes.DELETE("/token", useMiddlewares.PrivateMiddlewares, ctrl.RevokeToken)
	}

	return r
}
