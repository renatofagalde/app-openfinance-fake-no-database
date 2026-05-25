package output

import "bootstrap/internal/domain"

type ConfigStorage interface {
	Load() (*domain.Config, error)
	Save(*domain.Config) error
}
