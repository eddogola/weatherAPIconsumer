package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const file = "key.txt"

func getApiKey(file string) string {
	data, err := ioutil.ReadFile(file)
	check(err)
	key := string(data)

	return key
}

type Temperature struct {
	Value float64
	Unit string
}

type Weather struct{
	IconPhrase string
	Temperature Temperature
}

const locationKey = "224749"

func generateURL() string {
	url := fmt.Sprintf("http://dataservice.accuweather.com/forecasts/v1/hourly/12hour/%s", locationKey)
	req, err := http.NewRequest("GET", url, nil)

	check(err)

	apiKey := getApiKey(file)

	q := req.URL.Query()
	q.Add("apikey", apiKey)
	q.Add("metric", "true")
	q.Add("details", "true")

	req.URL.RawQuery = q.Encode()
	return req.URL.String()
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := generateURL()

	resp, err := http.Get(url)

	check(err)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		
		check(err)

		var hours []Weather

		json.Unmarshal(b, &hours)
	
		w.Header().Set("Content-Type", "text/html")
		for _, h := range hours {
			fmt.Fprintf(w, "<li>%s - %.2f %s</li>", h.IconPhrase, h.Temperature.Value, h.Temperature.Unit)
		}
	} else {
		fmt.Println(fmt.Errorf("failed to get url, %v", url))
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("listening on localhost port 8000")
	http.ListenAndServe(":8000", nil)
}