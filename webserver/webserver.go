package webserver

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/AndreasAlbert/gonda/auth"
	"github.com/AndreasAlbert/gonda/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Store         storage.FileStore
	Router        *gin.Engine
	OAuthHandlers []auth.OAuthHandler
}

func NewServer(fs storage.FileStore, engine *gin.Engine, oauthhandlers []auth.OAuthHandler) Server {
	s := Server{
		fs, engine, oauthhandlers}
	store := cookie.NewStore([]byte("kdjalskdjalskj"))
	s.Router.Use(sessions.Sessions("gonda", store))

	addRoutes(s)

	return s
}

func (s Server) HandleWhoAmI(ctx *gin.Context) {
	session := sessions.Default(ctx)
	fmt.Printf("%v", session.Get("user_name"))
	ctx.JSON(http.StatusOK, gin.H{"user_name": session.Get("user_name"), "user_provider": session.Get("user_provider"), "test": session.Get("test")})
}
func addRoutes(s Server) {
	for _, handler := range s.OAuthHandlers {
		s.Router.GET(fmt.Sprintf("/oauth/%s/login", handler.Name), handler.HandleLogin)

		s.Router.GET(fmt.Sprintf("/oauth/%s/callback", handler.Name), handler.HandleCallback)

	}
	s.Router.GET("/ping", s.HandlePing)

	// s.GET("/packages", s.HandleGetPackages)
	// s.POST("/packages", s.HandlePostPackages)
	// s.GET("/packages/:name", s.HandleGetPackage)

	// s.GET("/packages/:name/version/:version", s.HandleGetPackageVersion)
	s.Router.POST("/uploads", s.HandlePostUploads)
	// s.GET("/uploads/:id", s.HandleUploadGet)

	// s.GET("/channels", s.HandleGetChannels)
	// s.POST("/channels", s.HandlePostChannels)
	// s.GET("/channels/:name", s.HandleGetChannel)

	s.Router.GET("/me", s.HandleWhoAmI)
}

func (s Server) HandlePing(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"user":    user})
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
