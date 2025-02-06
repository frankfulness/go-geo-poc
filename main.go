package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/uber/h3-go/v4"
)

type Location struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type POI struct {
    ID   int
    Name string
    Lat  float64
    Lon  float64
}

// Global variable declaration
var indexedPOIs map[h3.Cell][]POI

func main() {
	// Initialize sample POIs - all within ~1km of NYC test coordinates
	samplePOIs := []POI{
			{ID: 1, Name: "World Trade Center", Lat: 40.7127, Lon: -74.0134},
			{ID: 2, Name: "City Hall", Lat: 40.7128, Lon: -74.0060},
			{ID: 3, Name: "Brooklyn Bridge", Lat: 40.7061, Lon: -73.9969},
	}

	// Initialize the global indexedPOIs
	indexedPOIs = indexPOIs(samplePOIs)

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", r)
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

	// Print the raw response to see what fields we're getting
	fmt.Printf("Raw API response: %s\n", string(body))

	var response struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
			// Add other fields from the API response
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
			return Location{}, err
	}

	return Location{
			Lat: response.Lat,
			Lon: response.Lon,
	}, nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get IP address, trying X-Forwarded-For first
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
			ip = r.RemoteAddr
			// Strip the port number if present
			if colon := strings.LastIndex(ip, ":"); colon != -1 {
					ip = ip[:colon]
			}
	}

	// For localhost testing, use a dummy IP
	if ip == "[::1]" || ip == "localhost" || ip == "127.0.0.1" {
			// Using New York City coordinates for testing
			loc := Location{
					Lat: 40.7128,
					Lon: -74.0060,
			}
			
			nearbyPOIs := getNearbyPOIs(indexedPOIs, loc.Lat, loc.Lon)
			fmt.Fprintf(w, "Your location (localhost testing - using NYC): %v, %v\n", loc.Lat, loc.Lon)
			fmt.Fprint(w, "Recommended POIs:\n")
			for _, poi := range nearbyPOIs {
					fmt.Fprintf(w, "- %s\n", poi.Name)
			}
			return
	}

	loc, err := GetLocationByIP(ip)
	if err != nil {
			fmt.Printf("Error getting location: %v\n", err)
			http.Error(w, "Error retrieving location", http.StatusInternalServerError)
			return
	}

	nearbyPOIs := getNearbyPOIs(indexedPOIs, loc.Lat, loc.Lon)
	fmt.Fprintf(w, "Your location: %v, %v\n", loc.Lat, loc.Lon)
	fmt.Fprint(w, "Recommended POIs:\n")
	for _, poi := range nearbyPOIs {
			fmt.Fprintf(w, "- %s\n", poi.Name)
	}
}

func indexPOIs(pois []POI) map[h3.Cell][]POI {
	index := make(map[h3.Cell][]POI)

	for _, poi := range pois {
			latLng := h3.LatLng{
					Lat: poi.Lat,
					Lng: poi.Lon,
			}
			cell, err := h3.LatLngToCell(latLng, 9)
			if err != nil {
					continue
			}
			fmt.Printf("POI %s is in cell: %v\n", poi.Name, cell) // Debug print
			index[cell] = append(index[cell], poi)
	}

	return index
}

func getNearbyPOIs(indexedPOIs map[h3.Cell][]POI, lat, lon float64) []POI {
	latLng := h3.LatLng{
			Lat: lat,
			Lng: lon,
	}
	cell, err := h3.LatLngToCell(latLng, 9)
	if err != nil {
			return nil
	}
	fmt.Printf("Looking for POIs in cell: %v\n", cell) // Debug print
	return indexedPOIs[cell]
}
