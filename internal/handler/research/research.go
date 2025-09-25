package research

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResearchProvider interface {
	RunFileProcessing(filename, collectionId string, src io.Reader) error
}

type CollectionProvider interface {
	CheckExists(id string) (bool, error)
}

type ResearchHandler struct {
	researchProvider   ResearchProvider
	collectionProvider CollectionProvider
}

func New(researchProvider ResearchProvider, collectionProvider CollectionProvider) *ResearchHandler {
	return &ResearchHandler{
		researchProvider:   researchProvider,
		collectionProvider: collectionProvider,
	}
}

func (r *ResearchHandler) Upload(c *gin.Context) {
	collectionId := c.Query("collection_id")
	if collectionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection_id query param required"})
		return
	}

	exists, err := r.collectionProvider.CheckExists(collectionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error find collection"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(c.GetHeader("Content-Type"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in multipart form"})
		fmt.Println(err)
		return
	}

	files := form.File["files"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
			return
		}
		defer src.Close()

		if err := r.researchProvider.RunFileProcessing(file.Filename, collectionId, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
			return
		}

		fmt.Println("success upload file", file.Filename)
	}

	c.JSON(200, gin.H{"message": "success", "count": len(files)})
}
