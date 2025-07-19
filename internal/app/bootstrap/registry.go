package bootstrap

import (
	"context"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/rotisserie/eris"
)

type Registry map[entity.ProviderName]repository.RateAPI

func NewRegistry(ctx context.Context, cfg configs.Config, db repository.ProviderDatabase) (Registry, error) {
	providers, err := db.GetAllProviders(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "get providers for registry")
	}

	r := make(Registry, len(providers))

	for _, providerEntry := range providers {
		providerName := entity.ProviderName(providerEntry.Name)
		r[providerName] = api.ProviderFactory(cfg, providerName)
	}

	return r, nil
}

func (r Registry) Register(provider entity.ProviderName, rateAPI repository.RateAPI) {
	_, ok := r[provider]
	if !ok {
		r[provider] = rateAPI
	}
}
