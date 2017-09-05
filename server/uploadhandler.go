package server

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kataras/iris"
	"github.com/saurabhsvj/photosharing/aws"
	"github.com/saurabhsvj/photosharing/storage"
)

// UploadHandler handles upload of user content to S3
func UploadHandler(ctx iris.Context) {

	awsSession := aws.GetAWSSession()

	// Get the file from the request
	file, info, err := ctx.FormFile("uploadfile")

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		log.Fatal("Failed to get the uploaded image from request")
		return
	}

	defer file.Close()

	// Store the file to local and get a handler...
	local, err := getLocalImageFile(file, info.Filename)
	defer local.Close()

	if err != nil {
		log.Fatal("Failed to get a local file handler for the uploaded file")
		log.Fatal(err)

		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	// Check if the upload is a valid image
	fileInfo, _ := local.Stat()
	fileSize := fileInfo.Size()

	// Read the local file into a buffer for upload
	buffer := make([]byte, fileSize)
	local.Read(buffer)

	// Verify if the uploaded file format is allowed
	isValidUpload := checkUploadValidity(buffer)
	if !isValidUpload {
		log.Printf("Client uploaded a non supported file format.")
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	id, err := storage.WriteImageToBucket(awsSession, os.Getenv("BUCKET_NAME"), local.Name(), buffer, fileSize)

	if err != nil {
		log.Println(err)

		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	log.Printf("Successfully uploaded the image to upsteam storage.")
	ctx.StatusCode(iris.StatusOK)
	ctx.Write([]byte(id.String()))
	return
}

func getLocalImageFile(uploadedFile multipart.File, fileName string) (*os.File, error) {
	localPath := "./uploads/" + fileName
	out, err := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	log.Printf("Copying contents to a local file on server...")
	io.Copy(out, uploadedFile)
	out.Close()

	log.Printf("Getting handler for the local file buffer...")
	local, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}

	return local, nil
}

func checkUploadValidity(fileBuffer []byte) bool {

	log.Printf("Verifying validity of the file upload type...")
	allowedUploadTypes := [3]string{"image/jpeg", "image/png"}

	fileType := http.DetectContentType(fileBuffer)
	log.Printf("Uploaded file type is :%s", fileType)

	for _, val := range allowedUploadTypes {
		if fileType == val {

			log.Printf("The uploaded file matches a supported formats.")
			return true
		}
	}

	log.Printf("The uploaded file doesnot match supported formats.")
	return false
}
