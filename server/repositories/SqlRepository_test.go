package repositories

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDataAPI struct {
	mock.Mock
}

func (m *MockDataAPI) ExecuteStatement(input *rdsdataservice.ExecuteStatementInput) (*rdsdataservice.ExecuteStatementOutput, error) {
	args := m.Called(input)
	var resp *rdsdataservice.ExecuteStatementOutput
	if args.Get(0) != nil {
		resp = args.Get(0).(*rdsdataservice.ExecuteStatementOutput)
	}
	return resp, args.Error(1)
}

func TestSqlRepository(t *testing.T) {

	t.Run("GetProfiles successfull", func(t *testing.T) {

		mockData := &MockDataAPI{}
		output := &rdsdataservice.ExecuteStatementOutput{Records: [][]*rdsdataservice.Field{
			{
				{LongValue: aws.Int64(1)}, {StringValue: aws.String("NAME")}, {StringValue: aws.String("EMAIL")}, {StringValue: aws.String("PHONES")},
			},
		}}
		mockData.On("ExecuteStatement", mock.Anything).Return(output, nil)
		repo := NewSqlRepository(aws.String("arn"), aws.String("secret"), mockData)
		profiles, err := repo.GetProfiles()
		assert.NoError(t, err)
		assert.Len(t, profiles, 1)
		mockData.AssertExpectations(t)
	})
	t.Run("GetProfiles Fail", func(t *testing.T) {

		mockData := &MockDataAPI{}
		mockData.On("ExecuteStatement", mock.Anything).Return(nil, errors.New("Some Error"))
		repo := NewSqlRepository(aws.String("arn"), aws.String("secret"), mockData)
		_, err := repo.GetProfiles()
		assert.Error(t, err, "Some Error")
		mockData.AssertExpectations(t)
	})

	t.Run("InsertProfile Success", func(t *testing.T) {
		profileId := aws.Int64(100)
		mockData := &MockDataAPI{}
		output := &rdsdataservice.ExecuteStatementOutput{GeneratedFields: []*rdsdataservice.Field{
			{LongValue: profileId},
		}}
		mockData.On("ExecuteStatement", mock.Anything).Return(output, nil)
		repo := NewSqlRepository(aws.String("arn"), aws.String("secret"), mockData)
		id, err := repo.InsertProfile("NAME", "EMAIL", "PHONE")
		assert.NoError(t, err)
		assert.Equal(t, *profileId, *id)
		mockData.AssertExpectations(t)
	})
	t.Run("InsertProfile Fail", func(t *testing.T) {

		mockData := &MockDataAPI{}
		mockData.On("ExecuteStatement", mock.Anything).Return(nil, errors.New("Some Error"))
		client := NewSqlRepository(aws.String("arn"), aws.String("secret"), mockData)
		_, err := client.InsertProfile("NAME", "EMAIL", "PHONE")
		assert.Error(t, err, "Some Error")
		mockData.AssertExpectations(t)
	})

}
