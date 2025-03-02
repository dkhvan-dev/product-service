package graph

import (
	"github.com/dkhvan-dev/product-service/internal/config"
)

type Resolver struct {
	config config.AppConfig
}

func NewResolver(cfg config.AppConfig) *Resolver {
	return &Resolver{
		config: cfg,
	}
}
