package input

import "github.com/renatofagalde/app-openfinance-fake-no-database/internal/domain"

// Admin e a porta consumida pelo controller administrativo.
// Permite inspecionar e mutar a configuracao em runtime.
type Admin interface {
	Config() *domain.Config
	UpsertRoute(key string, route domain.Route) error
	DeleteRoute(key string) error
	UpsertConsentimento(id string, c domain.Consentimento) error
	DeleteConsentimento(id string) error
}
