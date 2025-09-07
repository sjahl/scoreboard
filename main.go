package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"slices"
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

type Competitor struct {
	Uid      string
	HomeAway string
	Form     string
	Score    int
}

type CompetitionStatus struct {
	Clock        int    `json:"clock"`
	DisplayClock string `json:"displayClock"`
}

type Competition struct {
	Status struct {
		Clock        int
		DisplayClock string
	}
	Competitors []Competitor
}

type Event struct {
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Competitions []Competition `json:"competitions"`
}

type ApiResponse struct {
	Events []Event `json:"events"`
}

func fetchScores(req url.URL) []byte {
	res, err := http.Get(req.String())

	if err != nil {
		log.Fatalf("HTTP request failed to %s", req)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatal(err)
	}

	return body
}

func validateLeague(league string) bool {
	keys := slices.Collect(maps.Keys(leagueMap))
	return slices.Contains(keys, league)
}

func main() {

	fmt.Println("Welcome to the scoreboard!")

	var league string
	flag.StringVar(&league, "l", "championship", "specify the league to fetch stores for")

	var date_string string
	flag.StringVar(&date_string, "d", "", "specify the date to fetch scores for, format: YYYYMMDD")

	flag.Parse()

	if !validateLeague(league) {
		log.Fatalf("Unknown league %s", league)
	}

	params := url.Values{}
	if date_string != "" {
		params.Add("dates", date_string)
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

	api_response := fetchScores(u)

	res := ApiResponse{}
	json.Unmarshal(api_response, &res)

	if len(res.Events) > 0 {
		fmt.Printf("%v", res)
	} else {
		fmt.Println("No events on date:", params.Get("dates"))
	}
}
