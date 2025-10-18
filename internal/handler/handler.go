package handler

import (
	"apiservice/internal/config"
	"apiservice/internal/domain/model"
	"apiservice/internal/handler/collection"
	"apiservice/internal/handler/inference"
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
	inferenceHandler := inference.New(serv.ResearchService)

	router := gin.Default()

	api := router.Group("/api/v1")

	collectionApi := api.Group("/collections")
	collectionApi.POST("/new", collectionHandler.Add)
	collectionApi.GET("", collectionHandler.List)
	collectionApi.GET("/:id", collectionHandler.GetOne)
	collectionApi.PATCH("/:id", collectionHandler.Update)
	collectionApi.GET("/:id/report", collectionHandler.Report)
	collectionApi.DELETE("/:id", collectionHandler.Delete)

	researchApi := api.Group(("researches"))
	researchApi.POST("/upload", researchHandler.Upload)
	researchApi.GET("/check", researchHandler.Check)
	researchApi.DELETE("/:id", researchHandler.Delete)
	researchApi.GET("/update/ws/", researchHandler.Connect)

	inferenceApi := api.Group("/inference")
	inferenceApi.GET("/progress/ws/", inferenceHandler.Connect)

	return router
}
