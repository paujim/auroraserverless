package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

var sqlClient *SqlClient

func init() {
	auroraArn := "arn:aws:rds:us-west-2:414215635918:cluster:serverless-cluster"                                    //os.Getenv("AURORA_ARN")
	secretArn := "arn:aws:secretsmanager:us-west-2:414215635918:secret:templatedsecret0EBB07A0-RXStFb83Bspx-xtaK8w" // os.Getenv("SECRET_ARN")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqlClient = &SqlClient{rdsdataservice.New(sess), aws.String(auroraArn), aws.String(secretArn)}
}

func main() {
	http.HandleFunc("/", profileHandler(sqlClient))
	http.ListenAndServe(":5000", nil)
}
