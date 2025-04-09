package api

import (
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
)

type Registry map[entity.ProviderName]repository.RateAPI

func (r Registry) Register(provider entity.ProviderName, rateAPI repository.RateAPI) {
	_, ok := r[provider]
	if !ok {
		r[provider] = rateAPI
	}
}
