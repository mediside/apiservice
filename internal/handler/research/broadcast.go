package research

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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
		message := <-h.updateCh
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
