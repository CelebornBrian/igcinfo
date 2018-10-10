package main

import (
	"net/http"
	"path"
	"strconv"

	"github.com/marni/goigc"
)

//getField handles requests for a specific field given an ID
//nolint: gocyclo
func getField(w http.ResponseWriter, r *http.Request) {
	//Extract ID and convert to int
	temp := path.Dir(r.URL.Path)
	id, err := strconv.Atoi(path.Base(temp))
	if err != nil {
		errStatus(w, http.StatusInternalServerError, err, "Could not convert ID from string to int in <field>")
		return
	}
	//Check if ID exists in database, if not write 404 error
	url, ok := database[id]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	//Parse the URL from the database, handle possible error
	track, err := igc.ParseLocation(url.URL)
	if err != nil {
		errStatus(w, http.StatusNotFound, err, "Could not parse url in database<field>, did the url change?")
	}

	//Switch based on the final part of the url, write the applicable field, handle possible errors
	switch path.Base(r.URL.Path) {
	case "pilot":
		_, err = w.Write([]byte(track.Pilot))
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Error writing Pilot in <field>")
			return
		}
	case "glider":
		_, err = w.Write([]byte(track.GliderType))
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Error writing Pilot in <field>")
			return
		}
	case "glider_id":
		_, err = w.Write([]byte(track.GliderID))
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Error writing Pilot in <field>")
			return
		}
	case "track_length":
		_, err = w.Write([]byte(strconv.Itoa(int(lengthCalc(track)))))
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Error writing Pilot in <field>")
			return
		}
	case "H_date":
		time, err := track.Header.Date.MarshalText()
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Error converting track date to string in <field>")
		}
		_, err = w.Write(time)
		if err != nil {
			errStatus(w, http.StatusInternalServerError, err, "Error writing Pilot in <field>")
			return
		}
	default:
		//Should never happen
		panic("Should never happen regex/field handling has failed, the end is nigh")
	}
}
