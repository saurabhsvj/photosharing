package server

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kataras/iris"
	"github.com/saurabhsvj/photosharing/aws"
	"github.com/saurabhsvj/photosharing/storage"
)

// UploadHandler handles upload of user content to S3
func UploadHandler(ctx iris.Context) {

	// TODO: Should move out from here
	awsSession, err := aws.GenerateAWSSession()
	// Create a single aws session to be used for all AWS calls

	if err != nil {
		log.Fatalf("App Failed to startup. Cannot create AWS session.")
		panic(err)
	}
	// -----------------------------------------------------------------

	// Get the file from the request
	file, info, err := ctx.FormFile("uploadfile")

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		log.Fatal("Failed to get the uploaded image from request")
		return
	}

	defer file.Close()

	// Create a file with the same name
	// assuming that you have a folder named 'uploads'
	localPath := "./uploads/" + info.Filename
	out, err := os.OpenFile(localPath,
		os.O_WRONLY|os.O_CREATE, 0666)

	// Copy the contents to a local file
	io.Copy(out, file)

	out.Close()

	local, _ := os.Open(localPath)
	defer local.Close()

	fileInfo, _ := local.Stat()
	fileSize := fileInfo.Size()
	print(fileInfo)
	uuid, err := storage.WriteImageToBucket(awsSession, "saurabhphotosharingapp", local, fileSize)

	if err != nil {
		fmt.Println(err)
	}

	log.Print(uuid.String())

	return
}
