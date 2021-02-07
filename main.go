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

func getDataFromEndpoint() (hours []Weather, err error) {
	url := generateURL()
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		b, err2 := ioutil.ReadAll(resp.Body)
		check(err2)

		json.Unmarshal(b, &hours)

		return
	} else {
		err = fmt.Errorf("failed to get url, %v", url)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	hours, err := getDataFromEndpoint()
	check(err)

	w.Header().Set("Content-Type", "text/html")
	for _, h := range hours {
		fmt.Fprintf(w, "<li>%s - %.2f %s</li>", h.IconPhrase, h.Temperature.Value, h.Temperature.Unit)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("listening on localhost port 8000")
	http.ListenAndServe(":8000", nil)
}