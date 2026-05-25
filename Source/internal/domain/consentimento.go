package domain

const (
	StatusAuthorised = "AUTHORISED"
	StatusRejected   = "REJECTED"
	StatusAwaiting   = "AWAITING_AUTHORISATION"
)

type Consentimento struct {
	Status string   `json:"status"`
	Negar  []string `json:"negar"`
}
