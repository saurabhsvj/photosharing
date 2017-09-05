package main

import (
	"github.com/kataras/iris"
	"github.com/saurabhsvj/photosharing/server"
)

func main() {

	app := iris.New()
	app.Post("/upload", iris.LimitRequestBodySize(10<<20), server.UploadHandler)
	app.Run(iris.Addr(":8080"))

	// file, _ := os.Open("test.jpeg")
	// u1, _ := storage.WriteImageToBucket(awsSession, "saurabhphotosharingapp", file)
	// print(u1.String())

}
