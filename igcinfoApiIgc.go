package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/marni/goigc"
)

//Register new track in the database
func registerTrack(w http.ResponseWriter, r *http.Request) {
	//ReturnID is a struct containing id used to return the ID of received  URLs
	type ReturnID struct {
		ID int `json:"id"`
	}
	//Extract request body
	postBody := r.Body
	defer r.Body.Close()
	var temp IgcObject
	//Extract and decode json payload containing URL, handle possible error
	err := json.NewDecoder(postBody).Decode(&temp)
	if err != nil {
		errStatus(w, http.StatusBadRequest, err, "Decoding of Request body failed")
		return
	}

	//Check URL, first with url.parse, then regex
	_, err = url.Parse(temp.URL)
	if err != nil {
		errStatus(w, http.StatusBadRequest, err, "Not a valid URL")
	}
	works, err := regexp.MatchString("^http.+\\.igc$", temp.URL)
	if err != nil {
		errStatus(w, http.StatusBadRequest, err, "Regex matchString failed for url")
	}
	if !works {
		errStatus(w, http.StatusBadRequest, nil, "Not a valid URL(regex)")
	}
	//Extract URL and try to parse it
	uri := temp.URL
	_, err = igc.ParseLocation(uri)
	if err != nil {
		errStatus(w, http.StatusBadRequest, err, "Failed to parse the URL")
	}
	//Loop through the database to find an unused ID, then insert the url and return the ID.
	for i := 1; true; i++ {
		_, ok := database[i]
		if !ok {
			database[i] = temp
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(ReturnID{i})
			if err != nil {
				errStatus(w, http.StatusInternalServerError, err, "Could not encode return ID")
			}
			return
		}
	}
}

//getAll returns an array with all IDs currently in the database
func getAll(w http.ResponseWriter) {
	//ReturnIDs contains an array with all the IDs in the database, to be returned to the user
	type ReturnIDs struct {
		IDs []int `json:"id"`
	}
	//Initialize a ReturnIDs variable
	returnID := ReturnIDs{make([]int, 0)}
	//The response should have header type json no matter what
	w.Header().Set("Content-Type", "application/json")
	//If the database is not empty, loop through all entries and append to return array.
	if len(database) != 0 {
		for id := range database {
			returnID.IDs = append(returnID.IDs, id)
		}
	}
	//Encode response, empty array if database is empty, array of IDs if not. If fail, handle error
	err := json.NewEncoder(w).Encode(returnID.IDs)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Could not encode empty array")
	}
}
