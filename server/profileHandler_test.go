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

func TestGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDS := &MockDataService{}
	output := &rdsdataservice.ExecuteStatementOutput{Records: [][]*rdsdataservice.Field{
		{
			{LongValue: aws.Int64(1)}, {StringValue: aws.String("alex")}, {StringValue: aws.String("alex@email.com")}, {StringValue: aws.String("+61491570156")},
		},
	}}
	mockDS.On("ExecuteStatement", mock.Anything).Return(output, nil)
	mockClient := &SqlClient{mockDS, aws.String("arn"), aws.String("secret")}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler(mockClient))
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	expected := `{"profiles":[{"full_name":"alex","email":"alex@email.com","phone_number":"+61491570156"}]}`
	assert.Equal(t, expected, rec.Body.String())
}

func TestPost(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"full_name": "hugo","email": "hugo@email.com","phone_number": "61 491 570 157"}`))
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

func TestPostWithInvalidEmail(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"full_name": "hugo","email": "bad email","phone_number": "0491 570 156"}`))
	if err != nil {
		t.Fatal(err)
	}
	mockClient := &SqlClient{&MockDataService{}, aws.String("arn"), aws.String("secret")}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler(mockClient))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, `{"error":"Invalid email [bad email]."}`, rec.Body.String())
}

func TestPostWithInvalidPhone(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"full_name": "hugo","email": "some@email.com","phone_number": "9999 99491570 156"}`))
	if err != nil {
		t.Fatal(err)
	}
	mockClient := &SqlClient{&MockDataService{}, aws.String("arn"), aws.String("secret")}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(profileHandler(mockClient))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, `{"error":"Invalid phone [9999 99491570 156]."}`, rec.Body.String())
}
