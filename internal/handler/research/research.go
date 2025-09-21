package research

import (
	"apiservice/internal/domain/research"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResearchProvider interface {
	Create() (research.Research, error)
	Delete(id string) error
	List() ([]research.Research, error)
	GetOne(id string) (research.Research, error)
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
	list, err := h.researchProvider.List()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	if len(list) == 0 { // ожидаем на фронте "[]" вместо null
		list = []research.Research{}
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *ResearchHandler) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	res, err := h.researchProvider.GetOne(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *ResearchHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	if err := h.researchProvider.Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.Status(http.StatusOK)
}
