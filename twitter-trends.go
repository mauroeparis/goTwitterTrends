package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	APIKey string
	APISecret string
	AccessToken string
	AccessSecret string
	WOEID int64
}

type custom_trend struct {
	Name string `json:"name"`
	URL string `json:"url"`
	TweetVolume int64 `json:"tweet_volume"`
}

var (
	cach = cache.New(5*time.Minute, 10*time.Minute)
	configuration = Configuration{}
)

func getTwitterClient() *twitter.Client {
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
	var cus_trends []custom_trend
 	cache_trends, found := cach.Get("cus_trends")
	
	 if found {
		cus_trends = cache_trends.([]custom_trend)
	} else {
		client := getTwitterClient()
		var arg_woeid int64 = configuration.WOEID
		cus_trends = make([]custom_trend, 0)
		trends, _, _ := client.Trends.Place(
			arg_woeid,
			&twitter.TrendsPlaceParams{},
		)

		for _, v := range trends[0].Trends {
			if v.TweetVolume > 0 {
				pop_trend := custom_trend{ v.Name, v.URL, v.TweetVolume }
				cus_trends = append(cus_trends, pop_trend)
			}
		}
		sort.Slice(cus_trends, func(i, j int) bool {
			return cus_trends[i].TweetVolume > cus_trends[j].TweetVolume
		})
		cach.Set("cus_trends", cus_trends, cache.DefaultExpiration)
	}
	return cus_trends
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	trends := getTrends()
	response, _ := json.Marshal(trends)
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
