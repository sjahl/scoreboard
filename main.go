package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var leagueMap = map[string]string{
	"premier":      "eng.1",
	"championship": "eng.2",
	"league-1":     "eng.3",
	"league-2":     "eng.4",
	"mls":          "usa.1",
	"bundesliga":   "ger.1",
	"eredivisie":   "ned.1",
	"laliga":       "esp.1",
	"ligue-1":      "fra.1",
}

func main() {

	fmt.Println("Welcome to the scoreboard!")

	var league string
	flag.StringVar(&league, "l", "championship", "specify the league to fetch stores for")

	var date_string string
	flag.StringVar(&date_string, "d", "", "specify the date to fetch scores for, format: YYYYMMDD")

	flag.Parse()

	params := url.Values{}
	if date_string != "" {
		params.Add("date", date_string)
	}

	// TODO: check that we can actually access leagueMap[league]
	req_url, err := url.JoinPath("/apis/site/v2/sports/soccer", leagueMap[league], "scoreboard")
	if err != nil {
		log.Fatalf("")
	}

	u := url.URL{
		Scheme:   "http",
		Host:     "site.api.espn.com",
		Path:     req_url,
		RawQuery: params.Encode(),
	}

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("HTTP request failed to %s", req_url)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}
