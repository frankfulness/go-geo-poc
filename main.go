package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8080", r)
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetLocationByIP(ip string) (Location, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}

	var loc Location
	err = json.Unmarshal(body, &loc)
	if err != nil {
		return Location{}, err
	}

	return loc, nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	loc, err := GetLocationByIP(ip)
	if err != nil {
		http.Error(w, "Error retrieving location", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Your location: %v, %v", loc.Lat, loc.Lon)
}
