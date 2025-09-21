package handler

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/model"
	"apiservice/internal/handler/research"
	"apiservice/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger, cfg *config.Config, serv service.Service) http.Handler {
	if cfg.Env == model.EnvLocal {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	researchService := research.New(serv.ResearchService)

	router := gin.Default()

	api := router.Group("/api/v1")
	api.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"success": true}) })

	researchApi := api.Group("/researches")
	researchApi.POST("/new", researchService.Add)
	researchApi.GET("", researchService.List)
	researchApi.GET("/:id", researchService.FullInfo)
	researchApi.DELETE("/:id", researchService.Delete)

	return router
}
