package entity

import (
	"errors"
	"fmt"
	"strings"
)

type CurrencyISO string

func (s CurrencyISO) ToUpper() CurrencyISO {
	return CurrencyISO(strings.ToUpper(string(s)))
}

func (s CurrencyISO) ToLower() CurrencyISO {
	return CurrencyISO(strings.ToLower(string(s)))
}

func (s CurrencyISO) String() string {
	return string(s)
}

type CurrencyPairString string

func (s CurrencyPairString) ToUpper() CurrencyPairString {
	return CurrencyPairString(strings.ToUpper(string(s)))
}

func (s CurrencyPairString) ToLower() CurrencyPairString {
	return CurrencyPairString(strings.ToLower(string(s)))
}

func (s CurrencyPairString) ToPair() (Pair, error) {
	parts := strings.Split(string(s), "_")

	if len(parts) != 2 {
		return Pair{}, errors.New("invalid currency pair")
	}

	return Pair{
		From: CurrencyISO(parts[0]).ToUpper(),
		To:   CurrencyISO(parts[1]).ToLower(),
	}, nil
}

type Pair struct {
	From CurrencyISO
	To   CurrencyISO
	_    any
}

func (p Pair) Join(glue string) string {
	return fmt.Sprintf("%s%s%s", p.From.ToUpper(), glue, p.To.ToUpper())
}

func (p Pair) ToCurrencyPairString() CurrencyPairString {
	return CurrencyPairString(fmt.Sprintf("%s_%s", p.From.ToUpper(), p.To.ToUpper()))
}
