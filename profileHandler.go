package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
)

type ProfileHandler struct {
	client    rdsdataserviceiface.RDSDataServiceAPI
	auroraArn *string
	secretArn *string
}

const (
	NAME  = 1
	EMAIL = 2
	PHONE = 3
)

type Profile struct {
	ID           int64
	FullName     string   `json:"full_name"`
	Email        string   `json:"email"`
	PhoneNumbers []string `json:"phone_numbers"`
}

type ProfileResponse struct {
	ProfileID int64 `json:"profile_id"`
}

func (h *ProfileHandler) insertProfile(profile *Profile) (*int64, error) {
	log.Printf("Insert data to DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: h.auroraArn,
		SecretArn:   h.secretArn,
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
	req, resp := h.client.ExecuteStatementRequest(params)
	err := req.Send()
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}
	log.Printf("%s\n", resp.GoString())
	return resp.GeneratedFields[0].LongValue, nil
}

func (h *ProfileHandler) getProfiles() ([]Profile, error) {
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

func (h *ProfileHandler) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		log.Printf("HTTP GET\n")
		// get data from DB
		profiles, err := h.getProfiles()
		if err != nil {
			log.Printf("Error fetching data: %s\n", err)
			http.Error(w, "Something bad happend", http.StatusInternalServerError)
			return
		}
		// Return json
		log.Printf("Response: %v\n", profiles)
		returnJSON(w, profiles, http.StatusOK)

	case "POST":
		log.Printf("HTTP POST\n")
		var profile Profile
		err := json.NewDecoder(r.Body).Decode(&profile)
		// Validate
		if err != nil {
			log.Printf("Error decoding: %s\n", err)
			http.Error(w, "Bad Request.", http.StatusBadRequest)
			return
		}
		log.Printf("Request Body: %v\n", profile)
		// Need to validate phone

		// Save to DB
		profileID, err := h.insertProfile(&profile)
		if err != nil {

			log.Printf("Error inserting: %s\n", err)
			http.Error(w, "Unable to add that user.", http.StatusBadRequest)
			return
		}
		resp := ProfileResponse{ProfileID: *profileID}
		// Return json
		returnJSON(w, resp, http.StatusCreated)
	default:
		http.Error(w, "Only GET and POST methods are supported.", http.StatusBadRequest)
	}
}

func returnJSON(w http.ResponseWriter, jsonObj interface{}, statusCode int) {
	js, err := json.Marshal(jsonObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
}
