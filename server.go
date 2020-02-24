package main

import (
	"log"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tkanos/gonfig"
)

type screen_width struct {
	Max_width int `json:max_width`
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	decoder := json.NewDecoder(r.Body)
	var scr_w screen_width
	
	err := decoder.Decode(&scr_w)
	if err != nil {
        panic(err)
	}
	
	trends := getTrends()
	trends_boxes := trendsToBoxes(trends, scr_w.Max_width)
	potpack(trends_boxes)

	response, _ := json.Marshal(trends_boxes)
    w.Write(response)
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	trends := getTrends()
	trends_boxes := trendsToBoxes(trends, 1440)
	potpack(trends_boxes)

	response, _ := json.Marshal(trends_boxes)
    w.Write(response)
}

func main() {
	gonfig.GetConf("./config.json", &configuration)
	r := mux.NewRouter()
	r.HandleFunc("/", post).Methods(http.MethodPost)
	r.HandleFunc("/", get).Methods(http.MethodGet)
	handler := cors.Default().Handler(r)

    log.Fatal(http.ListenAndServe(
		"0.0.0.0:8080",
		handler,
	))
}