package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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
	flag.StringVar(&date_string, "d", "", "specify the date to fetch scores for")

	flag.Parse()

	fmt.Println("league var has value", league)
	fmt.Println("date_string has value", date_string)

	// TODO: construct query params better later...
	if date_string != "" {
		date_string = fmt.Sprintf("?date=%s", date_string)
	}

	req_url := fmt.Sprintf("http://site.api.espn.com/apis/site/v2/sports/soccer/%s/scoreboard%s", leagueMap[league], date_string)

	res, err := http.Get(req_url)
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
