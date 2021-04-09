package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const baseUrl = "http://api.aviationstack.com/v1"

type health struct {
	Status string `json:"status"`
}

type Airline struct {
	Name string `json:"name"`
}

type FlightInfo struct {
	IATACode string `json:"iata"`
}

type Destination struct {
	Airport   string `json:"airport"`
	Scheduled string `json:"scheduled"`
	Estimated string `json:"estimated"`
	Terminal  string `json:"terminal"`
	Gate      string `json:"gate"`
}

type LiveData struct {
	IsGround bool `json:"is_ground"`
}

type Flight struct {
	Airline    Airline     `json:"airline"`
	FlightInfo FlightInfo  `json:"flight"`
	Departure  Destination `json:"departure"`
	Arrival    Destination `json:"arrival"`
	Live       LiveData    `json:"live"`
}

type Response struct {
	Flights []Flight `json:"data"`
}

func main() {
	fmt.Println("Wayaround Flights API 0.1 Powered by Mux")
	handleRequests()
}

func handleRequests() {

	//create a new instance of our router
	wayaroundFlightsRouter := mux.NewRouter().StrictSlash(true)

	//the following handles the differen places for routing on the api
	wayaroundFlightsRouter.HandleFunc("/health", healthCheck)
	wayaroundFlightsRouter.HandleFunc("/flightstatus", flightsCheckHandler)
	log.Fatal(http.ListenAndServe(":10000", wayaroundFlightsRouter))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	healthStatus := health{Status: "Service OK!"}
	json.NewEncoder(w).Encode(healthStatus)
}

func flightsCheckHandler(w http.ResponseWriter, r *http.Request) {
	iata := r.URL.Query().Get("iata")
	flightId := r.URL.Query().Get("flightId")
	response := checkFlightFromSource(iata, flightId)
	json.NewEncoder(w).Encode(response)
}

func checkFlightFromSource(iata string, flightId string) Response {
	apiKey := os.Getenv("AVIATION_STACK_AUTH")
	requestUrl := baseUrl + "/flights?airline_iata=" + iata + "&flight_number=" + flightId + "&access_key=" + apiKey

	//lets doublecheck the url:
	fmt.Println(requestUrl)
	response, err := http.Get(requestUrl)
	var apiResponse Response

	if err != nil {
		fmt.Println(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal((err))
	}

	json.Unmarshal([]byte(responseData), &apiResponse)
	defer response.Body.Close()
	return apiResponse
}
