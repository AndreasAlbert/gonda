package webserver

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/AndreasAlbert/gonda/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Store storage.FileStore
	*gin.Engine
}

func NewServer(fs storage.FileStore, engine *gin.Engine) Server {
	s := Server{
		fs,
		engine,
	}

	routes(s)

	return s
}
func routes(s Server) {
	s.GET("/ping", s.HandlePing)

	// s.GET("/packages", s.HandleGetPackages)
	// s.POST("/packages", s.HandlePostPackages)
	// s.GET("/packages/:name", s.HandleGetPackage)

	// s.GET("/packages/:name/version/:version", s.HandleGetPackageVersion)
	s.POST("/uploads", s.HandlePostUploads)
	// s.GET("/uploads/:id", s.HandleUploadGet)

	// s.GET("/channels", s.HandleGetChannels)
	// s.POST("/channels", s.HandlePostChannels)
	// s.GET("/channels/:name", s.HandleGetChannel)

}
func (s Server) HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s Server) HandlePostUploads(c *gin.Context) {
	file, _ := c.FormFile("file")

	// Validate file name
	match, regexError := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9\-_\.]*[a-zA-Z0-9]$`, file.Filename)
	if regexError != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Code is broken."))
		return
	} else if !match {
		c.String(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid file name '%s'", file.Filename))
		return
	}

	// Push file to storage
	f, fileOpenErr := file.Open()
	if fileOpenErr != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to process uploaded file."))
		return
	}

	tmpname := filepath.Join("_upload/", file.Filename)
	putError := s.Store.Put(tmpname, f, false)
	if putError != nil {
		c.String(http.StatusConflict, fmt.Sprintf("File exists."))
		return
	}

	// TODO: Create pending Upload record in DB

	// Send response
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
