package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

var profileHandler *ProfileHandler

func init() {
	auroraArn := os.Getenv("AURORA_ARN")
	secretArn := os.Getenv("SECRET_ARN")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := rdsdataservice.New(sess)

	profileHandler = &ProfileHandler{client, aws.String(auroraArn), aws.String(secretArn)}
}

func main() {
	http.HandleFunc("/", profileHandler.HandleFunc)
	http.ListenAndServe(":5000", nil)
}
