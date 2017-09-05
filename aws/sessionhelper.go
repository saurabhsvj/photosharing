package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
)

var sharedSession *session.Session

// GetAWSSession is used to create a single session to be used by all calls
func GetAWSSession() *session.Session {

	if sharedSession != nil {
		return sharedSession
	}

	sharedSession, err := session.NewSession()
	if err != nil {
		log.Fatal("Failed to create an AWS session. App will not start up..")
		panic(err)
	}

	return sharedSession
}
