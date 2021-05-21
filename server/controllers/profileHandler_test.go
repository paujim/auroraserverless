package controllers

import (
	"net/http"
	"net/http/httptest"
	"paujim/auroraserverless/server/entities"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) InsertProfile(fullName, email, phoneNumber string) (*int64, error) {
	args := m.Called(fullName, email, phoneNumber)
	var resp *int64
	if args.Get(0) != nil {
		resp = args.Get(0).(*int64)
	}
	return resp, args.Error(1)
}
func (m *MockRepository) GetProfiles() ([]entities.Profile, error) {
	args := m.Called()
	var resp []entities.Profile
	if args.Get(0) != nil {
		resp = args.Get(0).([]entities.Profile)
	}
	return resp, args.Error(1)
}

func TestGet(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockRepo := &MockRepository{}
	mockRepo.On("GetProfiles").Return([]entities.Profile{
		{
			ID:          1,
			FullName:    "alex",
			Email:       "alex@email.com",
			PhoneNumber: "+61491570156",
		},
	}, nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(ProfileHandler(mockRepo))
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
	mockRepo := &MockRepository{}
	mockRepo.On("InsertProfile", mock.Anything, mock.Anything, mock.Anything).Return(aws.Int64(1000), nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(ProfileHandler(mockRepo))
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
	mockRepo := &MockRepository{}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(ProfileHandler(mockRepo))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, `{"error":"Invalid email [bad email]."}`, rec.Body.String())
}

func TestPostWithInvalidPhone(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"full_name": "hugo","email": "some@email.com","phone_number": "9999 99491570 156"}`))
	if err != nil {
		t.Fatal(err)
	}
	mockRepo := &MockRepository{}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(ProfileHandler(mockRepo))
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, `{"error":"Invalid phone [9999 99491570 156]."}`, rec.Body.String())
}
