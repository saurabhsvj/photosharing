package storage

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
)

// GetBucketList lists the buckets used by the user
func GetBucketList(session *session.Session) {

	svc := s3.New(session)
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println("Failed to list buckets", err)
		return
	}

	log.Println("Buckets:")
	for _, bucket := range result.Buckets {
		log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	}
}

// WriteImageToBucket writes a new image object to S3
func WriteImageToBucket(session *session.Session, bucketName string, fileName string, buffer []byte, fileSize int64) (uuid.UUID, error) {
	u1 := uuid.NewV4()

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := fmt.Sprintf("%s%s", "media/images/", u1.String())
	uploader := s3manager.NewUploader(session)
	params := &s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	}
	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err := uploader.Upload(params)

	if err != nil {
		// Print the error and exit.
		fmt.Printf("Unable to upload %q to %q, %v", fileName, bucketName, err)
		return u1, err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucketName)
	return u1, nil
}
