package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/marni/goigc"
)

//getNumber gets a specific ID from the database if it exists
func getNumber(w http.ResponseWriter, r *http.Request) {
	//Track is a struck containing all the data that is relevant
	type Track struct {
		HDate       time.Time `json:"H_date"`
		Pilot       string    `json:"pilot"`
		Glider      string    `json:"glider"`
		GliderID    string    `json:"glider_id"`
		TrackLength float64   `json:"track_length"`
	}
	//Isolate the ID from the url, and convert to int, if that fails, handle error
	id := path.Base(r.URL.Path)
	num, err := strconv.Atoi(id)
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Could not convert ID from string to int")
		return
	}
	//Response should be json, if ID does exist in database, get it. Else write 404 error
	w.Header().Set("Content-Type", "application/json")
	url, ok := database[num]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		//Parse the url and store all data, if parse fails, handle error
		track, err := igc.ParseLocation(url.URL)
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Could not parse url, has the URL changed?")
			return
		}
		//Create a Track object and populate with data from track, including calculated length
		temp := Track{track.Header.Date, track.Pilot, track.GliderType,
			track.GliderID, lengthCalc(track)}
		//Try to encode based on the Track struct format, if fail, handle error
		err = json.NewEncoder(w).Encode(temp)
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Failed to encode return json payload")
			return
		}
	}
}
