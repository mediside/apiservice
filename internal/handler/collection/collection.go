package collection

import (
	"apiservice/internal/domain/collection"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CollectionProvider interface {
	Create() (collection.Collection, error)
	Delete(id string) error
	List() ([]collection.Collection, error)
	GetOne(id string) (collection.CollectionWithResearches, error)
}

type CollectionHandler struct {
	collectionProvider CollectionProvider
}

func New(collectionProvider CollectionProvider) *CollectionHandler {
	return &CollectionHandler{
		collectionProvider: collectionProvider,
	}
}

func (h *CollectionHandler) Add(ctx *gin.Context) {
	res, err := h.collectionProvider.Create()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CollectionHandler) List(ctx *gin.Context) {
	list, err := h.collectionProvider.List()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	if len(list) == 0 { // ожидаем на фронте "[]" вместо null
		list = []collection.Collection{}
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *CollectionHandler) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	res, err := h.collectionProvider.GetOne(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CollectionHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	if err := h.collectionProvider.Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.Status(http.StatusOK)
}
