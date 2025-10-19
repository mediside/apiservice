package inference

import (
	"apiservice/internal/domain/inference"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type researchProvider interface {
	InferenceCh() <-chan inference.InferenceProgress
	RunFolderProcessing(collectionId string) error
}

type Handler struct {
	researchProvider researchProvider
	inferenceCh      <-chan inference.InferenceProgress
	upgrader         websocket.Upgrader
	clients          map[*websocket.Conn]bool
}

func New(researchProvider researchProvider) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	infHandler := &Handler{
		researchProvider: researchProvider,
		inferenceCh:      researchProvider.InferenceCh(),
		upgrader:         upgrader,
		clients:          make(map[*websocket.Conn]bool),
	}

	go infHandler.broadcastMessages()

	return infHandler
}

func (h *Handler) Connect(ctx *gin.Context) {
	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.clients[conn] = true
}

func (h *Handler) broadcastMessages() {
	for {
		message := <-h.inferenceCh
		data, err := json.Marshal(message)
		if err != nil {
			return
		}
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}

func (h *Handler) RunOnFolder(ctx *gin.Context) {
	collectionId := ctx.Query("collection_id")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "collection_id query param required"})
		return
	}

	if err := h.researchProvider.RunFolderProcessing(collectionId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot run folder processing"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
