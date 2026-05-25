package filestore

import (
	"bootstrap/internal/domain"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// JSONStore persiste a Config em um unico arquivo JSON em disco.
// Implementa output.ConfigStorage.
//
// As operacoes de escrita sao atomicas via tmpfile + rename, evitando
// arquivos parcialmente escritos em caso de crash.
type JSONStore struct {
	path string
	mu   sync.Mutex
}

// NewJSONStore retorna um ponteiro pronto para uso.
// Pasta do arquivo e criada se nao existir.
func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

// Load le o arquivo JSON e devolve uma Config preenchida.
// Se o arquivo nao existir, retorna uma Config vazia (nao e erro).
func (s *JSONStore) Load() (*domain.Config, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return emptyConfig(), nil
		}
		return nil, err
	}
	cfg := &domain.Config{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	if cfg.Routes == nil {
		cfg.Routes = map[string]domain.Route{}
	}
	if cfg.Consentimentos == nil {
		cfg.Consentimentos = map[string]domain.Consentimento{}
	}
	if cfg.Version == 0 {
		cfg.Version = 1
	}
	return cfg, nil
}

// Save escreve a Config com escrita atomica.
func (s *JSONStore) Save(cfg *domain.Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func emptyConfig() *domain.Config {
	return &domain.Config{
		Version:        1,
		Routes:         map[string]domain.Route{},
		Consentimentos: map[string]domain.Consentimento{},
	}
}
