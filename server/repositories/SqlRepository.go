package repositories

import (
	"log"
	"paujim/auroraserverless/server/entities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type SqlRepository interface {
	InsertProfile(fullName, email, phoneNumber string) (*int64, error)
	GetProfiles() ([]entities.Profile, error)
}

type DataAPI interface {
	ExecuteStatement(input *rdsdataservice.ExecuteStatementInput) (*rdsdataservice.ExecuteStatementOutput, error)
}

type sqlRepository struct {
	dataAPI   DataAPI
	auroraArn *string
	secretArn *string
}

func NewSqlRepository(auroraArn, secretArn *string, dataAPI DataAPI) SqlRepository {
	return &sqlRepository{
		auroraArn: auroraArn,
		secretArn: secretArn,
		dataAPI:   dataAPI,
	}
}

func (c *sqlRepository) InsertProfile(fullName, email, phoneNumber string) (*int64, error) {
	log.Printf("Insert data to DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: c.auroraArn,
		SecretArn:   c.secretArn,
		Sql:         aws.String("INSERT INTO TestDB.Profiles (FullName, Email, Phone) VALUES (:name, :email, :phone);"),
		Parameters: []*rdsdataservice.SqlParameter{
			{
				Name: aws.String("name"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(fullName),
				},
			},
			{
				Name: aws.String("email"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(email),
				},
			},
			{
				Name: aws.String("phone"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(phoneNumber),
				},
			},
		},
	}
	resp, err := c.dataAPI.ExecuteStatement(params)
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}
	log.Printf("%s\n", resp.GoString())
	return resp.GeneratedFields[0].LongValue, nil
}

func (h *sqlRepository) GetProfiles() ([]entities.Profile, error) {
	log.Printf("Get data from DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: h.auroraArn,
		SecretArn:   h.secretArn,
		Sql:         aws.String("SELECT * FROM TestDB.Profiles"),
	}
	resp, err := h.dataAPI.ExecuteStatement(params)
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}

	profiles := []entities.Profile{}
	for _, record := range resp.Records {
		profiles = append(profiles, entities.Profile{
			ID:          *record[0].LongValue,
			FullName:    *record[entities.NAME].StringValue,
			Email:       *record[entities.EMAIL].StringValue,
			PhoneNumber: *record[entities.PHONE].StringValue,
		})
	}
	return profiles, nil
}
