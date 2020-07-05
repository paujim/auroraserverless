package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

var sqlClient *SqlClient

func init() {
	auroraArn := os.Getenv("AURORA_ARN")
	secretArn := os.Getenv("SECRET_ARN")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqlClient = &SqlClient{rdsdataservice.New(sess), aws.String(auroraArn), aws.String(secretArn)}
}

func main() {
	http.HandleFunc("/", profileHandler(sqlClient))
	http.ListenAndServe(":5000", nil)
}
