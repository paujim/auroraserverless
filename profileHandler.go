package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	NAME  = 1
	EMAIL = 2
	PHONE = 3
)

type Profile struct {
	ID           int64    `json:"-"`
	FullName     string   `json:"full_name"`
	Email        string   `json:"email"`
	PhoneNumbers []string `json:"phone_numbers"`
}

type ProfileResponse struct {
	ProfileID int64 `json:"profile_id"`
}

func profileHandler(h *SqlClient) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			log.Printf("HTTP GET\n")
			// get data from DB
			profiles, err := h.GetProfiles()
			if err != nil {
				log.Printf("Error fetching data: %s\n", err)
				http.Error(w, "Something bad happend", http.StatusInternalServerError)
				return
			}
			// Return json
			log.Printf("Response: %v\n", profiles)
			returnJSON(w, profiles, http.StatusOK)

		case http.MethodPost:
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
			profileID, err := h.InsertProfile(&profile)
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
