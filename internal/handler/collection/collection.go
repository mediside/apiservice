package collection

import (
	"apiservice/internal/domain/collection"
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type collectionProvider interface {
	Create() (collection.Collection, error)
	Delete(id string) error
	List() ([]collection.Collection, error)
	GetOne(id string) (collection.CollectionWithResearches, error)
	CreateReport(id string) (*bytes.Buffer, error)
}

type Handler struct {
	collectionProvider collectionProvider
}

func New(collectionProvider collectionProvider) *Handler {
	return &Handler{
		collectionProvider: collectionProvider,
	}
}

func (h *Handler) Add(ctx *gin.Context) {
	res, err := h.collectionProvider.Create()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *Handler) List(ctx *gin.Context) {
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

func (h *Handler) GetOne(ctx *gin.Context) {
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

func (h *Handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	if err := h.collectionProvider.Delete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) Report(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "err"})
		return
	}

	buf, err := h.collectionProvider.CreateReport(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "err"})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=users.xlsx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	ctx.Status(http.StatusOK)
}
