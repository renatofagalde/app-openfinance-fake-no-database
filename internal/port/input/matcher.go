package input

import "github.com/renatofagalde/app-openfinance-fake-no-database/internal/domain"

// MatchResult e o que o handler precisa para responder uma requisicao mockada:
// a rota encontrada, os parametros extraidos do path e a chave que a identifica.
type MatchResult struct {
	Route  domain.Route
	Params map[string]string
	Key    string
}

// Matcher e a porta consumida pelo controller HTTP generico.
// Encapsula a busca de rotas e a leitura de consentimentos.
type Matcher interface {
	Match(method, path string) (*MatchResult, bool)
	GetConsentimento(consentId string) (*domain.Consentimento, bool)
}
