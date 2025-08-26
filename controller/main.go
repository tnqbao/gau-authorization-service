package controller

import (
	"github.com/tnqbao/gau-authorization-service/config"
	"github.com/tnqbao/gau-authorization-service/infra"
	"github.com/tnqbao/gau-authorization-service/provider"
	"github.com/tnqbao/gau-authorization-service/repository"
)

type Controller struct {
	Repository     *repository.Repository
	Infrastructure *infra.Infra
	Config         *config.Config
	Provider       *provider.Provider
}

func NewController(cfg *config.Config, infra *infra.Infra, repo *repository.Repository) *Controller {
	provide := provider.InitProvider(cfg.EnvConfig)
	return &Controller{
		Repository:     repo,
		Infrastructure: infra,
		Config:         cfg,
		Provider:       provide,
	}
}
