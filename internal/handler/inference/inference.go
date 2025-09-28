package inference

import (
	"apiservice/internal/domain/inference"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ResearchProvider interface {
	InferenceCh() <-chan inference.InferenceProgress
}

type InferenceHandler struct {
	researchProvider ResearchProvider
	inferenceCh      <-chan inference.InferenceProgress
	upgrader         websocket.Upgrader
	clients          map[*websocket.Conn]bool
}

func New(researchProvider ResearchProvider) *InferenceHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	infHandler := &InferenceHandler{
		researchProvider: researchProvider,
		inferenceCh:      researchProvider.InferenceCh(),
		upgrader:         upgrader,
		clients:          make(map[*websocket.Conn]bool),
	}

	go infHandler.broadcastMessages()

	return infHandler
}

func (h *InferenceHandler) Connect(ctx *gin.Context) {
	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.clients[conn] = true
}

func (h *InferenceHandler) broadcastMessages() {
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
