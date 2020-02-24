package main

import (
	"log"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/tkanos/gonfig"
)

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trends := getTrends()
	trends_boxes := trendsToBoxes(trends)
	potpack(trends_boxes)

	response, _ := json.Marshal(trends_boxes)
    w.Write(response)
}

func main() {
	gonfig.GetConf("./config.json", &configuration)
	r := mux.NewRouter()
	r.HandleFunc("/", get).Methods(http.MethodGet)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{`*`})
	methodsOk := handlers.AllowedMethods(
		[]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	)
    log.Fatal(http.ListenAndServe(
		":8080",
		handlers.CORS(originsOk, headersOk, methodsOk)(r)),
	)
}