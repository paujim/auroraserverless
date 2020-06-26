package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRoute(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDS := &MockDataService{}
	output := &rdsdataservice.ExecuteStatementOutput{Records: [][]*rdsdataservice.Field{
		{
			{LongValue: aws.Int64(1)}, {StringValue: aws.String("NAME")}, {StringValue: aws.String("EMAIL")}, {StringValue: aws.String("PHONE1")},
		},
	}}
	mockDS.On("ExecuteStatement", mock.Anything).Return(output, nil)
	mockClient := &SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler(mockClient))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	expected := `[{"full_name":"NAME","email":"EMAIL","phone_numbers":["PHONE1"]}]`
	assert.Equal(t, expected, rec.Body.String())
}

func TestPostRoute(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"full_name": "NAME","email": "EMAIL","phone_numbers": ["phone1","phone2"]}`))
	if err != nil {
		t.Fatal(err)
	}
	mockDS := &MockDataService{}
	output := &rdsdataservice.ExecuteStatementOutput{GeneratedFields: []*rdsdataservice.Field{
		{LongValue: aws.Int64(1000)},
	}}
	mockDS.On("ExecuteStatement", mock.Anything).Return(output, nil)
	mockClient := &SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler(mockClient))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	expected := `{"profile_id":1000}`
	assert.Equal(t, expected, rec.Body.String())
}
