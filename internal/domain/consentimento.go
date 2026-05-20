package domain

// Status possiveis para consentimentos do Open Finance Brasil.
const (
	StatusAuthorised = "AUTHORISED"
	StatusRejected   = "REJECTED"
	StatusAwaiting   = "AWAITING_AUTHORISATION"
)

// Consentimento representa o estado de autorizacao e as permissoes negadas.
// A lista Negar enumera permissoes que devem retornar 403 mesmo o consentimento
// estando AUTHORISED, simulando o cenario de permissao granular ausente.
type Consentimento struct {
	Status string   `json:"status"`
	Negar  []string `json:"negar"`
}
