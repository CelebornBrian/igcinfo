package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

//nolint: gocyclo
func router(w http.ResponseWriter, r *http.Request) {
	//Build regex expressions for the url and handle possible errors
	apiHandler, err := regexp.Compile("^/igcinfo/api/?$")
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Failed to compile api regex")
		return
	}
	apiIgcHandler, err := regexp.Compile("^/igcinfo/api/igc/?$")
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Failed to compile api/igc regex")
		return
	}
	apiIgcNumberHandler, err := regexp.Compile("^/igcinfo/api/igc/[0-9]+/?$")
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Failed to compile api/igc/<id> regex")
		return
	}
	apiIgcNumberFieldHandler, err := regexp.Compile("/igcinfo/api/igc/[0-9]+/(pilot|glider|glider_id|track_length|H_date)$")
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Failed to compile api/igc/<id>/<field> regex")
		return
	}

	//Check if request is GET or POST
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		//Switch on the request url path and select handler
		switch {
		case apiHandler.MatchString(r.URL.Path):
			handleAPI(w)
		case apiIgcHandler.MatchString(r.URL.Path):
			if r.Method == http.MethodPost {
				registerTrack(w, r)
			} else if r.Method == http.MethodGet {
				getAll(w)
			}
		case apiIgcNumberHandler.MatchString(r.URL.Path):
			getNumber(w, r)
		case apiIgcNumberFieldHandler.MatchString(r.URL.Path):
			getField(w, r)
		default:
			errStatus(w, http.StatusNotFound, nil, "")
		}
	} else {
		errStatus(w, http.StatusNotImplemented, nil, "")
	}

}

//Write status header and body with status code, error if exist, and possible extra info
func errStatus(w http.ResponseWriter, status int, err error, extraInfo string) {
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n%s\n%s", http.StatusText(status), err, extraInfo), status)
	} else {
		http.Error(w, fmt.Sprintf("%s\n%s", http.StatusText(status), extraInfo), status)
	}
}

//Handle requests to /api/, respond with json encoded struct
func handleAPI(w http.ResponseWriter) {

	// APIInfo is a struct for /api/ call
	type APIInfo struct {
		Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
	}
	currentTime := time.Now()
	//Create return struct with calculated uptime
	api := APIInfo{uptimeFunc(startTime, currentTime), "Service for IGC tracks", "v1"}
	w.Header().Set("Content-Type", "application/json")
	//Encode the struct and send to response writer, else handle error
	err := json.NewEncoder(w).Encode(api)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Could not encode APIInfo to json")
		return
	}
}
