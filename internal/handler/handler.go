package handler

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/model"
	"apiservice/internal/handler/collection"
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

	collectionHandler := collection.New(serv.CollectionService)
	researchHandler := research.New(serv.ResearchService, serv.CollectionService)

	router := gin.Default()

	api := router.Group("/api/v1")
	api.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"success": true}) })

	collectionApi := api.Group("/collections")
	collectionApi.POST("/new", collectionHandler.Add)
	collectionApi.GET("", collectionHandler.List)
	collectionApi.GET("/:id", collectionHandler.GetOne)
	collectionApi.DELETE("/:id", collectionHandler.Delete)

	researchApi := api.Group(("researches"))
	researchApi.POST("/upload", researchHandler.Upload)
	researchApi.DELETE("/:id", researchHandler.Delete)

	return router
}
