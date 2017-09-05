package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
)

// GenerateAWSSession is used to create a single session to be used by all calls
func GenerateAWSSession() (*session.Session, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		log.Fatal("Failed to create an AWS session. App will not start up..")
		return nil, err
	}

	return awsSession, nil
}
