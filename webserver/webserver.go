package webserver

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/AndreasAlbert/gonda/auth"
	fstore "github.com/AndreasAlbert/gonda/storage/files"
	ustore "github.com/AndreasAlbert/gonda/storage/users"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	FileStore     fstore.FileStore
	UserStore     ustore.UserStore
	Router        *gin.Engine
	OAuthHandlers []auth.OAuthHandler
}

func NewServer(fs fstore.FileStore, us ustore.UserStore, engine *gin.Engine, oauthhandlers []auth.OAuthHandler) Server {
	s := Server{
		fs, us, engine, oauthhandlers}
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
func not_logged_in_error(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, gin.H{"reason": "Not logged in."})
}
func middleware_authenticated(s Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Case 1: User already present in session
		session := sessions.Default(ctx)
		user := session.Get("user")
		if user != nil {
			return
		}

		// Case 2: User has logged in through oauth,
		// but user object not pulled from store yet
		user_provider := session.Get("user_provider")
		user_name := session.Get("user_name")
		if (user_provider == nil) || (user_name == nil) {
			not_logged_in_error(ctx)
			return
		}
		user, err := s.UserStore.RetrieveUser(user_provider.(string), user_name.(string))
		if err != nil {
			not_logged_in_error(ctx)
			return
		}

		// Case 3: User has logged in through oauth,
		// but user object not pulled from store yet

		user, err = s.UserStore.CreateUser(user_provider.(string), user_name.(string))
		session.Set("user", user)
		session.Save()

		ctx.Next()

	}
}
func addRoutes(s Server) {

	// Unauthenticated server basics
	s.Router.GET("/ping", s.HandlePing)

	// OAuth routes
	group_oauth := s.Router.Group("/oauth")
	{
		for _, handler := range s.OAuthHandlers {
			group := group_oauth.Group(fmt.Sprintf("/%s", handler.Name))
			group.GET("/login", handler.HandleLogin)
			group.GET("/callback", handler.HandleCallback)
		}
	}

	group_authenticated := s.Router.Group("/", middleware_authenticated(s))
	{
		group_authenticated.POST("/uploads", s.HandlePostUploads)
		group_authenticated.GET("/authping", s.HandlePing)

	}
	// s.GET("/packages", s.HandleGetPackages)
	// s.POST("/packages", s.HandlePostPackages)
	// s.GET("/packages/:name", s.HandleGetPackage)

	// s.GET("/packages/:name/version/:version", s.HandleGetPackageVersion)
	// s.Router
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
	putError := s.FileStore.Put(tmpname, f, false)
	if putError != nil {
		c.String(http.StatusConflict, fmt.Sprintf("File exists."))
		return
	}

	// TODO: Create pending Upload record in DB

	// Send response
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
