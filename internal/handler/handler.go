package handler

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/model"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger, cfg *config.Config) http.Handler {
	if cfg.Env == model.EnvLocal {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	api := router.Group("/api/v1")
	api.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"success": true}) })

	return router
}
