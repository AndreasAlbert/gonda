/*
gonda-server starts the gonda server.
*/
package main

import (
	"fmt"

	"github.com/AndreasAlbert/gonda/auth"
	"github.com/AndreasAlbert/gonda/storage"
	"github.com/AndreasAlbert/gonda/webserver"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
)

func getAuthHandlers() []auth.OAuthHandler {
	var handlers []auth.OAuthHandler
	handler := auth.OAuthHandler{
		Name: "github",
		Config: oauth2.Config{
			ClientID:     viper.GetString("server.auth.github.client_id"),
			ClientSecret: viper.GetString("server.auth.github.client_secret"),
			Scopes:       []string{"read:user"},
			Endpoint:     github.Endpoint,
			RedirectURL:  viper.GetString("server.auth.github.redirect_url")}}
	handlers = append(handlers, handler)
	return handlers
}

func main() {

	// Log configuration
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Configuration with Viper
	viper.SetConfigName("gonda")
	viper.SetConfigType("yaml")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Dependency: File storage
	store, fileStoreErr := storage.NewLocalFileStore("/tmp/gonda/")
	if fileStoreErr != nil {
		panic(fileStoreErr)
	}

	router := gin.Default()
	s := webserver.NewServer(
		store, router, getAuthHandlers())

	s.Router.Run()
}
