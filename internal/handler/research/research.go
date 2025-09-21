package research

import (
	"apiservice/internal/domain/research"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResearchProvider interface {
	Create() (research.Research, error)
}

type ResearchHandler struct {
	researchProvider ResearchProvider
}

func New(researchProvider ResearchProvider) *ResearchHandler {
	return &ResearchHandler{
		researchProvider: researchProvider,
	}
}

func (h *ResearchHandler) Add(ctx *gin.Context) {
	res, err := h.researchProvider.Create()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *ResearchHandler) List(ctx *gin.Context) {

}

func (h *ResearchHandler) FullInfo(ctx *gin.Context) {

}

func (h *ResearchHandler) Delete(ctx *gin.Context) {

}
