package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/renatofagalde/app-openfinance-fake-no-database/internal/adapter/input/controller"
)

// Setup registra apenas o que e fixo no Gin:
//   - healthz
//   - endpoints administrativos sob /_admin/...
//
// Todas as demais requisicoes caem no MockController via gin.NoRoute,
// permitindo que rotas mockadas sejam adicionadas em runtime sem restart.
func Setup(r *gin.Engine, mockCtrl *controller.MockController, adminCtrl *controller.AdminController) {
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	admin := r.Group("/_admin")
	{
		admin.GET("/config", adminCtrl.GetConfig)
		admin.PUT("/routes", adminCtrl.UpsertRoute)
		admin.DELETE("/routes", adminCtrl.DeleteRoute)
		admin.PUT("/consentimentos", adminCtrl.UpsertConsentimento)
		admin.DELETE("/consentimentos", adminCtrl.DeleteConsentimento)
	}

	r.NoRoute(mockCtrl.Handle)
}
