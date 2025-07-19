package entity

import (
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

type Symbol string

func (s Symbol) ToUpper() Symbol {
	return Symbol(strings.ToUpper(string(s)))
}

func (s Symbol) ToLower() Symbol {
	return Symbol(strings.ToLower(string(s)))
}

func (s Symbol) String() string {
	return string(s)
}

func (s Symbol) ToPair() Pair {
	parts := strings.Split(s.String(), "_")
	return Pair{
		From:   CurrencyISO(parts[0]).ToUpper(),
		To:     CurrencyISO(parts[1]).ToUpper(),
		Symbol: s,
	}
}

func NewSymbol(from, to CurrencyISO) Symbol {
	return Symbol(strings.ToUpper(fmt.Sprintf("%s_%s", from, to)))
}

type Pair struct {
	From   CurrencyISO
	To     CurrencyISO
	Symbol Symbol
}

func (p *Pair) Join(glue string) string {
	return fmt.Sprint(p.From, glue, p.To)
}
