package handler

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/model"
	"apiservice/internal/handler/collection"
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

	collectionService := collection.New(serv.CollectionService)

	router := gin.Default()

	api := router.Group("/api/v1")
	api.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"success": true}) })

	collectionApi := api.Group("/collections")
	collectionApi.POST("/new", collectionService.Add)
	collectionApi.GET("", collectionService.List)
	collectionApi.GET("/:id", collectionService.GetOne)
	collectionApi.DELETE("/:id", collectionService.Delete)

	return router
}
