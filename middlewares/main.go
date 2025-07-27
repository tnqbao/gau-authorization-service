package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-authorization-service/controller"
)

type Middlewares struct {
	CORSMiddleware     gin.HandlerFunc
	PrivateMiddlewares gin.HandlerFunc
}

func NewMiddlewares(ctrl *controller.Controller) (*Middlewares, error) {
	cors := CORSMiddleware(ctrl.Config.EnvConfig)
	if cors == nil {
		return nil, nil
	}
	private := PrivateMiddleware(ctrl.Config.EnvConfig)
	if private == nil {
		return nil, nil
	}

	return &Middlewares{
		CORSMiddleware: cors,
	}, nil
}
