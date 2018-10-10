package main

import (
	"net/http"
	"os"
	"time"
)

//IgcObject has the url for to the igc file
type IgcObject struct {
	URL string `json:"url"`
}

//Create in-memory storage:
var database map[int]IgcObject

//StartTime is service start time
var startTime = time.Now()

func main() {
	//Initialize map
	database = make(map[int]IgcObject)

	//Send all requests to the router
	http.HandleFunc("/", router)

	//Start web server
	port := ":" + os.Getenv("PORT")
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
