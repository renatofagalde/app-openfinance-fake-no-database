package main

import (
	"log/slog"
	"os"
	"path/filepath"
	
	"github.com/gin-gonic/gin"
	
	"bootstrap/internal/adapter/input/controller"
	"bootstrap/internal/adapter/input/routes"
	"bootstrap/internal/adapter/output/filestore"
	"bootstrap/internal/application/service"
	"bootstrap/internal/config/logger"
)

const (
	defaultMocksFile = "./mocks.json"
	defaultSeedFile  = "./mocks.json"
	defaultPort      = "8080"
)

func main() {
	log := logger.GetLogger()
	
	configPath := envOrDefault("MOCKS_FILE", defaultMocksFile)
	seedPath := envOrDefault("SEED_FILE", defaultSeedFile)
	port := envOrDefault("PORT", defaultPort)
	
	seedIfMissing(log, configPath, seedPath)
	
	storage := filestore.NewJSONStore(configPath)
	
	mockSvc, err := service.NewMockService(storage)
	if err != nil {
		log.Error("erro carregando config", "err", err.Error(), "mocksFile", configPath)
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
