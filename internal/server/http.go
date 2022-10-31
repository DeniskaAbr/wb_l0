package server

import (
	"embed"
	"io/fs"
	"net/http"

	order "wb_l0/internal/order"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func (s *server) runHttpServer() {
	go func() {

		orderRoot := EmbedFolder(order.Content, "static")

		s.gin.Use(static.Serve("/", orderRoot))
		s.gin.GET("/ping", func(c *gin.Context) {
			c.String(200, "test")
		})

		if err := s.gin.Run("0.0.0.0:" + s.cfg.Http_port); err != nil {
			s.log.Printf("Error: %v", err)
		}

	}()
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
