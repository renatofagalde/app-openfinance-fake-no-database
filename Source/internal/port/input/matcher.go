package input

import "bootstrap/internal/domain"

type MatchResult struct {
	Route  domain.Route
	Params map[string]string
	Key    string
}

type Matcher interface {
	Match(method, path string) (*MatchResult, bool)
	GetConsentimento(consentId string) (*domain.Consentimento, bool)
}
