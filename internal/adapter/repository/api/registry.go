package api

import (
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
)

type Registry map[entity.ProviderName]repository.RateApi

func (r Registry) Register(provider entity.ProviderName, rateApi repository.RateApi) {
	_, ok := r[provider]
	if !ok {
		r[provider] = rateApi
	}
}
