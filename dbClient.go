package main

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
)

type SqlClient struct {
	client    rdsdataserviceiface.RDSDataServiceAPI
	auroraArn *string
	secretArn *string
}

func (c *SqlClient) insertProfile(profile *Profile) (*int64, error) {
	log.Printf("Insert data to DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: c.auroraArn,
		SecretArn:   c.secretArn,
		Sql:         aws.String("INSERT INTO TestDB.Profiles (FullName, Email, Phone) VALUES (:name, :email, :phone);"),
		Parameters: []*rdsdataservice.SqlParameter{
			{
				Name: aws.String("name"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(profile.FullName),
				},
			},
			{
				Name: aws.String("email"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(profile.Email),
				},
			},
			{
				Name: aws.String("phone"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(strings.Join(profile.PhoneNumbers, ";")),
				},
			},
		},
	}
	req, resp := c.client.ExecuteStatementRequest(params)
	err := req.Send()
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}
	log.Printf("%s\n", resp.GoString())
	return resp.GeneratedFields[0].LongValue, nil
}

func (h *SqlClient) getProfiles() ([]Profile, error) {
	log.Printf("Get data from DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: h.auroraArn,
		SecretArn:   h.secretArn,
		Sql:         aws.String("SELECT * FROM TestDB.Profiles"),
	}
	req, resp := h.client.ExecuteStatementRequest(params)
	err := req.Send()
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}

	var profiles []Profile
	for _, record := range resp.Records {
		profiles = append(profiles, Profile{
			ID:           *record[0].LongValue,
			FullName:     *record[NAME].StringValue,
			Email:        *record[EMAIL].StringValue,
			PhoneNumbers: []string{*record[PHONE].StringValue},
		})
	}
	return profiles, nil
}