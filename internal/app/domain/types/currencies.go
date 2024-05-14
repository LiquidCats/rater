package types

import (
	"strings"
)

type CurrencyNameString string

func (s CurrencyNameString) Upper() string {
	return strings.ToUpper(string(s))
}

func (s CurrencyNameString) Lower() string {
	return strings.ToLower(string(s))
}

func (s CurrencyNameString) Transform(transformerFunc func(c CurrencyNameString) string) string {
	return transformerFunc(s)
}

type BaseCurrency = CurrencyNameString
type QuoteCurrency = CurrencyNameString
