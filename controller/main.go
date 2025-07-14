package controller

import (
	"github.com/tnqbao/gau-authorization-service/config"
	"github.com/tnqbao/gau-authorization-service/infra"
	"github.com/tnqbao/gau-authorization-service/repository"
)

type Controller struct {
	Repository     *repository.Repository
	Infrastructure *infra.Infra
	Config         *config.Config
}

func NewController() *Controller {
	cfg := config.NewConfig()
	repo := repository.NewRepository(cfg)
	infra := infra.NewInfra(cfg)

	return &Controller{
		Repository:     repo,
		Infrastructure: infra,
		Config:         cfg,
	}
}

func NewController(cfg *config.Config, infra *infra.Infra) *Controller {
	repo := repository.NewRepository(cfg)
	return &Controller{
		Repository:     repo,
		Infrastructure: infra,
		Config:         cfg,
	}
}