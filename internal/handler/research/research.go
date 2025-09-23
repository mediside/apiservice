package research

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type ResearchProvider interface {
	SaveFile(filename string, src io.Reader) error
}

type ResearchHandler struct {
	researchProvider ResearchProvider
}

func New(researchProvider ResearchProvider) *ResearchHandler {
	return &ResearchHandler{
		researchProvider: researchProvider,
	}
}

func (r *ResearchHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(c.GetHeader("Content-Type"))
		c.JSON(400, gin.H{"error": "error in multipart form"})
		fmt.Println(err)
		return
	}

	files := form.File["files"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "cannot open file"})
			return
		}
		defer src.Close()

		if err := r.researchProvider.SaveFile(file.Filename, src); err != nil {
			c.JSON(500, gin.H{"error": "cannot save file"})
			return
		}

		fmt.Println("success upload file", file.Filename)
	}

	c.JSON(200, gin.H{"message": "success", "count": len(files)})
}
