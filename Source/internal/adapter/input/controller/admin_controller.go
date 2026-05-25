package controller

import (
	"bootstrap/internal/domain"
	"bootstrap/internal/port/input"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	svc input.Admin
}

func NewAdminController(svc input.Admin) *AdminController {
	return &AdminController{svc: svc}
}

type upsertRouteRequest struct {
	Key   string       `json:"key" binding:"required"`
	Route domain.Route `json:"route"`
}

type upsertConsentRequest struct {
	ConsentId     string               `json:"consentId" binding:"required"`
	Consentimento domain.Consentimento `json:"consentimento"`
}

// GET /_admin/config
func (a *AdminController) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, a.svc.Config())
}

func (a *AdminController) UpsertRoute(c *gin.Context) {
	var req upsertRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.svc.UpsertRoute(req.Key, req.Route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "key": req.Key})
}

func (a *AdminController) DeleteRoute(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'key' obrigatorio"})
		return
	}
	if err := a.svc.DeleteRoute(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "key": key})
}

func (a *AdminController) UpsertConsentimento(c *gin.Context) {
	var req upsertConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.svc.UpsertConsentimento(req.ConsentId, req.Consentimento); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "consentId": req.ConsentId})
}

func (a *AdminController) DeleteConsentimento(c *gin.Context) {
	id := c.Query("consentId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'consentId' obrigatorio"})
		return
	}
	if err := a.svc.DeleteConsentimento(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "consentId": id})
}
