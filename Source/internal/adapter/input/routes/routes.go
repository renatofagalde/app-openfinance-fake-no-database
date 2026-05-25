package routes

import (
	"bootstrap/internal/adapter/input/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
