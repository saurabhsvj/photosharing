package main

import (
	"log"
	"os"

	"github.com/kataras/iris"
	"github.com/saurabhsvj/photosharing/server"
)

const (
	// LOCALUPLOADFOLDER holds the uploaded files on the local server
	LOCALUPLOADFOLDER = "./uploads"

	// DEFAULTPORT indicates the default port used by the API Server
	DEFAULTPORT = ":8080"
)

func main() {

	// Ensure the uploads folder is available to create local files first
	log.Printf("Creating required directories for the app to work...")
	createRequiredDirectories()

	log.Printf("Start the web API...")
	startAPIServer()
}

func createRequiredDirectories() {

	if _, err := os.Stat(LOCALUPLOADFOLDER); os.IsNotExist(err) {
		os.Mkdir(LOCALUPLOADFOLDER, 0755)
	}
}

func startAPIServer() {

	app := iris.New()
	log.Printf("Created a new API instance.")

	log.Printf("Adding routes to the API server...")
	addRoutes(app)

	log.Printf("API server is booting up...")
	app.Run(iris.Addr(DEFAULTPORT))
}

func addRoutes(app *iris.Application) {
	app.Post("/upload", iris.LimitRequestBodySize(10<<20), server.UploadHandler)
}
