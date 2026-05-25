package service

import (
	"bootstrap/internal/adapter/input/matcher"
	"bootstrap/internal/domain"
	"bootstrap/internal/port/input"
	output "bootstrap/internal/port/storage"
	"errors"
	"strings"
	"sync"
)

type MockService struct {
	storage output.ConfigStorage
	cfg     *domain.Config
	mu      sync.RWMutex
}

// NewMockService carrega a config do storage e devolve o servico pronto.
func NewMockService(storage output.ConfigStorage) (*MockService, error) {
	cfg, err := storage.Load()
	if err != nil {
		return nil, err
	}
	return &MockService{storage: storage, cfg: cfg}, nil
}

// --- input.Matcher ---

// Match procura uma rota cuja chave seja "METHOD PATH" e cujo padrao
// case com o path da requisicao.
func (s *MockService) Match(method, path string) (*input.MatchResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key, route := range s.cfg.Routes {
		parts := strings.SplitN(key, " ", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] != method {
			continue
		}
		if params, ok := matcher.MatchPath(parts[1], path); ok {
			return &input.MatchResult{Route: route, Params: params, Key: key}, true
		}
	}
	return nil, false
}

func (s *MockService) GetConsentimento(consentId string) (*domain.Consentimento, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.cfg.Consentimentos[consentId]
	if !ok {
		return nil, false
	}
	return &c, true
}

// --- input.Admin ---

// Config retorna um snapshot da configuracao atual.
// Atencao: o ponteiro devolvido nao deve ser mutado fora deste servico.
func (s *MockService) Config() *domain.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *MockService) UpsertRoute(key string, route domain.Route) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("key vazia")
	}
	s.mu.Lock()
	s.cfg.Routes[key] = route
	err := s.storage.Save(s.cfg)
	s.mu.Unlock()
	return err
}

func (s *MockService) DeleteRoute(key string) error {
	s.mu.Lock()
	delete(s.cfg.Routes, key)
	err := s.storage.Save(s.cfg)
	s.mu.Unlock()
	return err
}

func (s *MockService) UpsertConsentimento(id string, c domain.Consentimento) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("consentId vazio")
	}
	s.mu.Lock()
	s.cfg.Consentimentos[id] = c
	err := s.storage.Save(s.cfg)
	s.mu.Unlock()
	return err
}

func (s *MockService) DeleteConsentimento(id string) error {
	s.mu.Lock()
	delete(s.cfg.Consentimentos, id)
	err := s.storage.Save(s.cfg)
	s.mu.Unlock()
	return err
}
