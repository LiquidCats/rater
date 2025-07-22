package bootstrap

import "github.com/LiquidCats/rater/internal/app/domain/entity"

type PairRegistry interface {
	GetPair(symbol entity.Symbol) (entity.Pair, bool)
	GetAllPairs() []entity.Pair
}
type ProviderRegistry interface {
	GetProvider(name entity.ProviderName) (entity.Provider, bool)
	GetAllProviders() []entity.Provider
}
