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
	apiRoutes := r.Group("/api/v2/authorization", useMiddlewares.CORSMiddleware, useMiddlewares.PrivateMiddlewares)
	{
		apiRoutes.POST("/token", ctrl.CreateNewToken)
		apiRoutes.GET("/token/renew", ctrl.RenewAccessToken)
		apiRoutes.GET("/token/validate", ctrl.CheckAccessToken)
		apiRoutes.DELETE("/token", ctrl.RevokeToken)
	}

	return r
}
