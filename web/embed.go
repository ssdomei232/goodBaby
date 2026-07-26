// Package web 把前端构建产物嵌入到二进制中并挂载到 gin
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var distFS embed.FS

// Register 把前端静态资源挂到根路径。
//
// 未匹配任何 API 路由的 GET 请求回退到 index.html，交给前端路由处理。
func Register(r *gin.Engine) {
	dist, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(dist))

	r.NoRoute(func(c *gin.Context) {
		// API 404 保持 JSON 返回
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "data": "接口不存在"})
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		if _, err := fs.Stat(dist, path); err != nil {
			// SPA 回退：交给 index.html 处理前端路由
			c.Request.URL.Path = "/"
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
