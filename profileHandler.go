package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
)

type ProfileHandler struct {
	client    rdsdataserviceiface.RDSDataServiceAPI
	auroraArn *string
	secretArn *string
}

type Profile struct {
	FullName     string   `json:"full_name"`
	Email        string   `json:"email"`
	PhoneNumbers []string `json:"phone_numbers"`
}

type ProfileResponse struct {
	ProfileID string `json:"profile_id"`
}

func (h *ProfileHandler) getProfiles() error {
	log.Printf("Get data from DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: h.auroraArn,
		SecretArn:   h.secretArn,
		Sql:         aws.String("SELECT * FROM database.profile"),
	}
	req, resp := h.client.ExecuteStatementRequest(params)
	err := req.Send()
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return err
	}
	log.Printf(resp.GoString())
	return nil
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
		err := h.getProfiles()
		if err != nil {
			log.Printf("Error fetching data: %s\n", err)
			http.Error(w, "Something bad happend", http.StatusInternalServerError)
			return
		}

		profiles := []Profile{{
			FullName:     "full name",
			Email:        "full@name.com",
			PhoneNumbers: []string{"phone1", "phone2"},
		}}
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
		log.Printf("Put data in DB\n")

		resp := ProfileResponse{"profile_id"}
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
