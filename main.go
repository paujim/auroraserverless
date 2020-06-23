package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", profileHandler)
	http.ListenAndServe(":5000", nil)
}
