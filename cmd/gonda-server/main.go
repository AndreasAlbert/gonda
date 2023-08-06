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
)

func main() {

	viper.SetConfigName("gonda")
	viper.SetConfigType("yaml")
	viper.SetConfigType("yml")

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// Dependency: File storage
	store, fileStoreErr := storage.NewLocalFileStore("/tmp/gonda/")
	if fileStoreErr != nil {
		panic(fileStoreErr)
	}

	githubHandler := auth.OAuthHandler{
		Name: "github",
		Config: oauth2.Config{
			ClientID:     viper.GetString("server.auth.github.client_id"),
			ClientSecret: viper.GetString("server.auth.github.client_secret"),
			Scopes:       []string{"read:user"},
			Endpoint:     github.Endpoint,
			RedirectURL:  "http://localhost:8080/oauth/github/callback"}}

	router := gin.Default()
	s := webserver.NewServer(
		store, router, []auth.OAuthHandler{githubHandler})

	s.Router.Run()
}
