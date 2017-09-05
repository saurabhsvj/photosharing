package main

import (
	"os"

	"github.com/kataras/iris"
	"github.com/saurabhsvj/photosharing/server"
)

const (
	// LOCALUPLOADFOLDER holds the uploaded files in the local server
	LOCALUPLOADFOLDER = "./uploads"
)

func main() {

	// Ensure the uploads folder is available to create local files first

	if _, err := os.Stat(LOCALUPLOADFOLDER); os.IsNotExist(err) {
		os.Mkdir(LOCALUPLOADFOLDER, 0755)
	}

	app := iris.New()
	app.Post("/upload", iris.LimitRequestBodySize(10<<20), server.UploadHandler)
	app.Run(iris.Addr(":8080"))

	// file, _ := os.Open("test.jpeg")
	// u1, _ := storage.WriteImageToBucket(awsSession, "saurabhphotosharingapp", file)
	// print(u1.String())

}
