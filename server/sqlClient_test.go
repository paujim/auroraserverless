package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDataService struct {
	mock.Mock
}

func (m *MockDataService) ExecuteStatement(input *rdsdataservice.ExecuteStatementInput) (*rdsdataservice.ExecuteStatementOutput, error) {
	args := m.Called(input)
	var resp *rdsdataservice.ExecuteStatementOutput
	if args.Get(0) != nil {
		resp = args.Get(0).(*rdsdataservice.ExecuteStatementOutput)
	}
	return resp, args.Error(1)
}

func TestDataService(t *testing.T) {

	t.Run("GetProfiles successfull", func(t *testing.T) {

		mockDS := &MockDataService{}
		output := &rdsdataservice.ExecuteStatementOutput{Records: [][]*rdsdataservice.Field{
			{
				{LongValue: aws.Int64(1)}, {StringValue: aws.String("NAME")}, {StringValue: aws.String("EMAIL")}, {StringValue: aws.String("PHONES")},
			},
		}}
		mockDS.On("ExecuteStatement", mock.Anything).Return(output, nil)
		client := SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
		profiles, err := client.GetProfiles()
		assert.NoError(t, err)
		assert.Len(t, profiles, 1)
		mockDS.AssertExpectations(t)
	})
	t.Run("GetProfiles Fail", func(t *testing.T) {

		mockDS := &MockDataService{}
		mockDS.On("ExecuteStatement", mock.Anything).Return(nil, errors.New("Some Error"))
		client := SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
		_, err := client.GetProfiles()
		assert.Error(t, err, "Some Error")
		mockDS.AssertExpectations(t)
	})

	t.Run("InsertProfile Success", func(t *testing.T) {
		profileId := aws.Int64(100)
		mockDS := &MockDataService{}
		output := &rdsdataservice.ExecuteStatementOutput{GeneratedFields: []*rdsdataservice.Field{
			{LongValue: profileId},
		}}
		mockDS.On("ExecuteStatement", mock.Anything).Return(output, nil)
		client := SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
		id, err := client.InsertProfile("NAME", "EMAIL", "PHONE")
		assert.NoError(t, err)
		assert.Equal(t, *profileId, *id)
		mockDS.AssertExpectations(t)
	})
	t.Run("InsertProfile Fail", func(t *testing.T) {

		mockDS := &MockDataService{}
		mockDS.On("ExecuteStatement", mock.Anything).Return(nil, errors.New("Some Error"))
		client := SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
		_, err := client.InsertProfile("NAME", "EMAIL", "PHONE")
		assert.Error(t, err, "Some Error")
		mockDS.AssertExpectations(t)
	})

}
