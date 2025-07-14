package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-authorization-service/controller"
)

type Middlewares struct {
	CORSMiddleware gin.HandlerFunc
}

func NewMiddlewares(ctrl *controller.Controller) (*Middlewares, error) {
	cors := CORSMiddleware(ctrl.Config.EnvConfig)
	return &Middlewares{
		CORSMiddleware: cors,
	}, nil
}
