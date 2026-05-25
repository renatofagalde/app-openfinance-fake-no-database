package input

import "bootstrap/internal/domain"

type Admin interface {
	Config() *domain.Config
	UpsertRoute(key string, route domain.Route) error
	DeleteRoute(key string) error
	UpsertConsentimento(id string, c domain.Consentimento) error
	DeleteConsentimento(id string) error
}
