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

type Team struct {
	Uid          string `json:"uid"`
	Abbreviation string `json:"abbreviation"`
}

type Competitor struct {
	Uid      string `json:"uid"`
	HomeAway string `json:"homeAway"`
	Form     string `json:"form"`
	Score    string `json:"score"`
	Team     Team   `json:"team"`
}

type CompetitionStatus struct {
	Clock        int    `json:"clock"`
	DisplayClock string `json:"displayClock"`
}

type Competition struct {
	Status      CompetitionStatus `json:"status"`
	Competitors []Competitor      `json:"competitors"`
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
		log.Fatalf("HTTP request failed to %v", req)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	return body
}

func validateLeague(league string) (bool, error) {
	keys := slices.Collect(maps.Keys(leagueMap))
	if !slices.Contains(keys, league) {
		err := fmt.Errorf("please enter one of the supported leagues: %v", keys)
		return false, err
	}

	return true, nil
}

func (e *Event) simpleScore() string {
	homeTeam := e.Competitions[0].Competitors[0]
	awayTeam := e.Competitions[0].Competitors[1]
	return fmt.Sprintf("%s %s - %s %s", homeTeam.Team.Abbreviation, homeTeam.Score, awayTeam.Score, awayTeam.Team.Abbreviation)
}

func main() {

	fmt.Println("Welcome to the scoreboard!")

	var league string
	flag.StringVar(&league, "l", "championship", "specify the league to fetch stores for")

	var date_string string
	flag.StringVar(&date_string, "d", "", "specify the date to fetch scores for, format: YYYYMMDD")

	flag.Parse()

	if _, err := validateLeague(league); err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	if date_string != "" {
		params.Add("dates", date_string)
	}

	req_url, err := url.JoinPath("/apis/site/v2/sports/soccer", leagueMap[league], "scoreboard")
	if err != nil {
		log.Fatalf("")
	}

	u := url.URL{
		Scheme:   "https",
		Host:     "site.api.espn.com",
		Path:     req_url,
		RawQuery: params.Encode(),
	}

	scores := ApiResponse{}
	json.Unmarshal(fetchScores(u), &scores)

	if len(scores.Events) > 0 {
		for _, s := range scores.Events {
			fmt.Println(s.simpleScore())
		}
	} else {
		fmt.Println("No events on date:", params.Get("dates"))
	}
}
