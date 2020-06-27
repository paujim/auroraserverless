package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/ttacon/libphonenumber"
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

type InsertProfileResponse struct {
	ProfileID *int64 `json:"profile_id"`
}
type GetProfileResponse struct {
	Profiles []Profile `json:"profiles"`
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
				http.Error(w, "Error fetching data", http.StatusInternalServerError)
				return
			}
			// Return json
			log.Printf("Response: %v\n", profiles)
			returnJSON(w, GetProfileResponse{Profiles: profiles}, http.StatusOK)

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
			if !IsValidEmail(profile.Email) {
				log.Printf("Not valid email: %s\n", profile.Email)
				http.Error(w, fmt.Sprintf("Invalid email [%s].", profile.Email), http.StatusBadRequest)
				return
			}

			formatedPhones := []string{}
			for _, phone := range profile.PhoneNumbers {
				num, err := libphonenumber.Parse(phone, "AU")
				if err != nil || !libphonenumber.IsValidNumber(num) {
					log.Printf("Not a valid phone: %s\n", phone)
					http.Error(w, fmt.Sprintf("Invalid phone [%s].", phone), http.StatusBadRequest)
					return
				}
				formatedPhones = append(formatedPhones, libphonenumber.Format(num, libphonenumber.E164))
			}

			// Save to DB
			profileID, err := h.InsertProfile(profile.FullName, profile.Email, formatedPhones)
			if err != nil {
				log.Printf("Error inserting: %s\n", err)
				http.Error(w, "Unable to add the profile.", http.StatusBadRequest)
				return
			}
			// Return json
			returnJSON(w, InsertProfileResponse{ProfileID: profileID}, http.StatusCreated)
		default:
			http.Error(w, "Only GET and POST methods are supported.", http.StatusBadRequest)
		}
	}
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
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
