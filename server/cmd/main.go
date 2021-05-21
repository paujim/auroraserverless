package main

import (
	"net/http"
	"os"
	"paujim/auroraserverless/server/controllers"
	"paujim/auroraserverless/server/repositories"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

var repo repositories.SqlRepository

func init() {
	auroraArn := os.Getenv("AURORA_ARN")
	secretArn := os.Getenv("SECRET_ARN")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	repo = repositories.NewSqlRepository(aws.String(auroraArn), aws.String(secretArn), rdsdataservice.New(sess))
}

func main() {
	http.HandleFunc("/", controllers.ProfileHandler(repo))
	http.ListenAndServe(":5000", nil)
}
