package research

import (
	"apiservice/internal/domain/research"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type researchProvider interface {
	RunFileProcessing(filename, collectionId string, src io.Reader) error
	Delete(id string) error
	UpdateCh() <-chan research.ResearchUpdate
}

type collectionProvider interface {
	CheckExists(id string) (bool, error)
}

type Handler struct {
	researchProvider   researchProvider
	collectionProvider collectionProvider
	updateCh           <-chan research.ResearchUpdate
	upgrader           websocket.Upgrader
	clients            map[*websocket.Conn]bool
}

func New(researchProvider researchProvider, collectionProvider collectionProvider) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	resHandler := &Handler{
		researchProvider:   researchProvider,
		collectionProvider: collectionProvider,
		updateCh:           researchProvider.UpdateCh(),
		upgrader:           upgrader,
		clients:            make(map[*websocket.Conn]bool),
	}

	go resHandler.broadcastMessages()

	return resHandler
}

func (r *Handler) Upload(ctx *gin.Context) {
	collectionId := ctx.Query("collection_id")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "collection_id query param required"})
		return
	}

	exists, err := r.collectionProvider.CheckExists(collectionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error find collection"})
		return
	}
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		fmt.Println(ctx.GetHeader("Content-Type"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error in multipart form"})
		fmt.Println(err)
		return
	}

	files := form.File["files"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer src.Close()

		if err := r.researchProvider.RunFileProcessing(file.Filename, collectionId, src); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
			return
		}
	}

	ctx.JSON(200, gin.H{"message": "success", "count": len(files)})
}

func (r *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	if err := r.researchProvider.Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
