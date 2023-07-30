package main

import (
	"github.com/AndreasAlbert/gonda/storage"
	"github.com/AndreasAlbert/gonda/webserver"

	"github.com/gin-gonic/gin"
)

func main() {

	// Dependency: File storage
	store, fileStoreErr := storage.NewLocalFileStore("/tmp/gonda/")
	if fileStoreErr != nil {
		panic(fileStoreErr)
	}

	s := webserver.NewServer(
		store,
		gin.Default(),
	)

	s.Run()
}
