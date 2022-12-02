package api

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/repository"
)

type Config struct {
	PkUsername   string `envconfig:"PK_USERNAME"`
	PkPassword   string `envconfig:"PK_PASSWORD"`
	TemplatePath string `envconfig:"TEMPLATE_PATH"`
}

type Api struct {
	repo *repository.Repository
	cfg  Config
}

func NewApi() *Api {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	return &Api{repo: r, cfg: cfg}
}
