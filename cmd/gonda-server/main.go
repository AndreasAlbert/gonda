/*
gonda-server starts the gonda server.
*/
package main

import (
	"fmt"
	"strings"

	"github.com/AndreasAlbert/gonda/auth"
	"github.com/AndreasAlbert/gonda/storage"
	"github.com/AndreasAlbert/gonda/webserver"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
)

// getAuthHandlers creates OAuthHandlers based on the given auth config
// Currently only creates a GitHub handler, but can be extended
func getAuthHandlers(viper *viper.Viper) []auth.OAuthHandler {
	var handlers []auth.OAuthHandler
	handler := auth.OAuthHandler{
		Name: "github",
		Config: oauth2.Config{
			ClientID:     viper.GetString("github.client_id"),
			ClientSecret: viper.GetString("github.client_secret"),
			Scopes:       []string{"read:user"},
			Endpoint:     github.Endpoint,
			RedirectURL:  viper.GetString("github.redirect_url")},
		UserURL: viper.GetString("github.user_url")}

	handlers = append(handlers, handler)
	return handlers
}

func main() {

	// Log configuration
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Configuration with Viper
	// Environment
	v := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "__")))
	v.SetEnvPrefix("GONDA")
	v.AutomaticEnv()

	// Config file
	v.SetConfigName("gonda")
	v.SetConfigType("yaml")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
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
		store, router, getAuthHandlers(v.Sub("server.auth")))

	s.Router.Run()
}
