package internal

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// loadTemplates loads templates from the embedded file system
func LoadTemplates(filesystem fs.FS) *template.Template {
	templ := template.Must(template.New("").ParseFS(filesystem, "*.html"))
	return templ
}

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
