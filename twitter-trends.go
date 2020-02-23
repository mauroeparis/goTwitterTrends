package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/tkanos/gonfig"
	"github.com/gorilla/mux"
)

type Configuration struct {
	APIKey string
	APISecret string
	AccessToken string
	AccessSecret string
	WOEID int64
}

type custom_trend struct {
	name string `json:"name"`
	URL string `json:"url"`
	tweetVolume int64 `json:"tweet_volume"`
}

func getTwitterClient() *twitter.Client {
	configuration := Configuration{}
	gonfig.GetConf("./config.json", &configuration)
	config := oauth1.NewConfig(
		configuration.APIKey,
		configuration.APISecret,
	)
	token := oauth1.NewToken(
		configuration.AccessToken,
		configuration.AccessSecret,
	)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client
}

func getTrends() []custom_trend {
	configuration := Configuration{}
	gonfig.GetConf("./config.json", &configuration)
	client := getTwitterClient()
	var arg_woeid int64 = configuration.WOEID

	trends, _, _ := client.Trends.Place(
		arg_woeid,
		&twitter.TrendsPlaceParams{},
	)

	var cus_trends = make([]custom_trend, 0)

	for _, v := range trends[0].Trends { 
		if v.TweetVolume > 0 {
			pop_trend := custom_trend{ v.Name, v.URL, v.TweetVolume }
			cus_trends = append(cus_trends, pop_trend)
		}
	}

	return cus_trends
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trends := getTrends()
	var data []string

	for i, t := range trends { 
		res := fmt.Sprintf(`{
			"name": "%v",
			"URL": "%v",
			"tweetVolume": "%v"
		},`, t.name, t.URL, t.tweetVolume)
		if i==len(trends)-1 {
			res = res[:len(res)-1]
		}
		data = append(data, res)
	}
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"message": %v}`, data)
    w.Write([]byte(response))
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", get).Methods(http.MethodGet)
    log.Fatal(http.ListenAndServe(":8080", r))
}