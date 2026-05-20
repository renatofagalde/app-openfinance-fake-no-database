package domain

import "encoding/json"

// Route representa uma rota mockada.
// Status zero e tratado como 200.
// Body e armazenado como JSON cru para preservar a forma exata fornecida.
type Route struct {
	Permission  string          `json:"permission,omitempty"`
	SkipConsent bool            `json:"skipConsent,omitempty"`
	Status      int             `json:"status"`
	Body        json.RawMessage `json:"body,omitempty"`
}

// Config e o estado completo do mock server, persistido em disco.
type Config struct {
	Version        int                      `json:"version"`
	Routes         map[string]Route         `json:"routes"`
	Consentimentos map[string]Consentimento `json:"consentimentos"`
}
