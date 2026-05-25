package controller

import (
	"bootstrap/internal/domain"
	"bootstrap/internal/port/input"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MockController e o handler generico. Toda requisicao que nao bate em
// uma rota registrada explicitamente cai aqui (via gin.NoRoute).
//
// O fluxo replica a logica do mock original:
//  1. Match da rota pela chave "METHOD PATH"
//  2. Se a rota nao for skipConsent, valida header x-consent-id
//  3. Valida status do consentimento (AUTHORISED)
//  4. Valida permissao granular contra a lista Negar
//  5. Devolve o body configurado com o status definido
type MockController struct {
	svc input.Matcher
}

func NewMockController(svc input.Matcher) *MockController {
	return &MockController{svc: svc}
}

func (m *MockController) Handle(c *gin.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path

	result, ok := m.svc.Match(method, path)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"errors": []gin.H{{
				"code":  "ROTA_NAO_CONFIGURADA",
				"title": "Nenhuma rota mockada para " + method + " " + path,
			}},
		})
		return
	}

	if !result.Route.SkipConsent {
		if errStatus, errBody, deny := m.validarConsentimento(c, result.Route); deny {
			c.JSON(errStatus, errBody)
			return
		}
	}

	status := result.Route.Status
	if status == 0 {
		status = http.StatusOK
	}

	if len(result.Route.Body) == 0 {
		c.Status(status)
		return
	}

	var body interface{}
	if err := json.Unmarshal(result.Route.Body, &body); err != nil {
		c.Data(status, "application/json", result.Route.Body)
		return
	}
	c.JSON(status, body)
}

func (m *MockController) validarConsentimento(c *gin.Context, r domain.Route) (int, gin.H, bool) {
	consentId := c.GetHeader("x-consent-id")
	if consentId == "" {
		return http.StatusBadRequest, gin.H{
			"errors": []gin.H{{
				"code":  "BAD_REQUEST",
				"title": "x-consent-id header obrigatorio",
			}},
		}, true
	}

	cons, found := m.svc.GetConsentimento(consentId)
	if !found {
		return http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{
				"code":  "CONSENTIMENTO_NAO_ENCONTRADO",
				"title": "Consentimento " + consentId + " nao encontrado",
			}},
		}, true
	}

	if cons.Status != domain.StatusAuthorised {
		return http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{
				"code":  "CONSENTIMENTO_NAO_AUTORIZADO",
				"title": "Consentimento com status " + cons.Status,
			}},
		}, true
	}

	if r.Permission != "" {
		for _, p := range cons.Negar {
			if p == r.Permission {
				return http.StatusForbidden, gin.H{
					"errors": []gin.H{{
						"code":  "FORBIDDEN",
						"title": "Permissao negada para " + p,
					}},
				}, true
			}
		}
	}

	return 0, nil, false
}
