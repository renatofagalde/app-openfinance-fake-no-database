package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/renatofagalde/app-openfinance-fake-no-database/internal/adapter/input/controller"
	"github.com/renatofagalde/app-openfinance-fake-no-database/internal/adapter/input/routes"
	"github.com/renatofagalde/app-openfinance-fake-no-database/internal/adapter/output/filestore"
	"github.com/renatofagalde/app-openfinance-fake-no-database/internal/application/service"
	"github.com/renatofagalde/app-openfinance-fake-no-database/internal/config/logger"
)

func main() {
	log := logger.GetLogger()

	configPath := envOrDefault("MOCKS_FILE", "/data/mocks.json")
	seedPath := envOrDefault("SEED_FILE", "/seed/mocks.json")
	port := envOrDefault("PORT", "8080")

	seedIfMissing(log, configPath, seedPath)

	storage := filestore.NewJSONStore(configPath)

	mockSvc, err := service.NewMockService(storage)
	if err != nil {
		log.Error("erro carregando config", "err", err.Error())
		os.Exit(1)
	}

	mockCtrl := controller.NewMockController(mockSvc)
	adminCtrl := controller.NewAdminController(mockSvc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	routes.Setup(r, mockCtrl, adminCtrl)

	log.Info("listening", "addr", ":"+port, "mocksFile", configPath)
	if err := r.Run(":" + port); err != nil {
		log.Error("server crashed", "err", err.Error())
		os.Exit(1)
	}
}

// seedIfMissing copia o arquivo de seed para o caminho operacional na primeira
// inicializacao, quando o volume montado em /data ainda esta vazio.
func seedIfMissing(log *slog.Logger, dst, src string) {
	if _, err := os.Stat(dst); err == nil {
		return
	}
	if _, err := os.Stat(src); err != nil {
		return
	}
	b, err := os.ReadFile(src)
	if err != nil {
		log.Warn("falha lendo seed", "err", err.Error())
		return
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		log.Warn("falha criando pasta", "err", err.Error())
		return
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		log.Warn("falha escrevendo seed", "err", err.Error())
		return
	}
	log.Info("seed copiado", "from", src, "to", dst)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
