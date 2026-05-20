package output

import "github.com/renatofagalde/app-openfinance-fake-no-database/internal/domain"

// ConfigStorage e a porta de saida para persistir e ler a configuracao.
// Implementacoes podem usar arquivo JSON, redis, banco, etc.
type ConfigStorage interface {
	Load() (*domain.Config, error)
	Save(*domain.Config) error
}
