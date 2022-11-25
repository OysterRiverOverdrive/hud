package hud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/oysterriveroverdrive/hud/model"
)

type TriviaService struct {
	URL           string
	TriviaService *http.Client
	AuthKey       string
}

func NewTriviaService(client *http.Client, AuthKey string) *TriviaService {
	return &TriviaService{
		URL:           DEFAULT_SERVER,
		TriviaService: client,
		AuthKey:       AuthKey,
	}
}

func (ts *TriviaService) setHeaders(req *http.Request) {
	req.Header.Set("X-TBA-Auth-Key", ts.AuthKey)
	req.Header.Set("accept", "application/json")
}

func (ts *TriviaService) Get(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}
	ts.setHeaders(req)
	return ts.TriviaService.Do(req)
}

func (ts *TriviaService) Dump(ctx context.Context, year, teamNum int, eventKey, districtKey, matchKey string) {
	teamKey := fmt.Sprintf("frc%d", teamNum)

	endpoints := []string{
		"/district/{district_key}/rankings",
		"/district/{district_key}/teams",
		"/district/{district_key}/teams/keys",
		"/district/{district_key}/teams/simple",
		"/event/{event_key}",
		"/event/{event_key}/alliances",
		"/event/{event_key}/awards",
		"/event/{event_key}/district_points",
		"/event/{event_key}/insights",
		"/event/{event_key}/matches",
		"/event/{event_key}/matches/keys",
		"/event/{event_key}/matches/simple",
		"/event/{event_key}/matches/timeseries", // not implemented
		"/event/{event_key}/oprs",
		"/event/{event_key}/predictions",
		"/event/{event_key}/rankings",
		"/event/{event_key}/simple",
		"/event/{event_key}/teams",
		"/event/{event_key}/teams/keys",
		"/event/{event_key}/teams/simple",
		"/event/{event_key}/teams/statuses",
		"/events/{year}",
		"/events/{year}/keys",
		"/events/{year}/simple",
		"/match/{match_key}",
		"/match/{match_key}/simple",
		"/match/{match_key}/timeseries",
		"/match/{match_key}/zebra_motionworks",
		"/team/{team_key}",
		"/team/{team_key}/awards",
		"/team/{team_key}/awards/{year}",
		"/team/{team_key}/districts",
		"/team/{team_key}/event/{event_key}/awards",
		"/team/{team_key}/event/{event_key}/matches",
		"/team/{team_key}/event/{event_key}/matches/keys",
		"/team/{team_key}/event/{event_key}/matches/simple",
		"/team/{team_key}/event/{event_key}/status",
		"/team/{team_key}/events",
		"/team/{team_key}/events/{year}",
		"/team/{team_key}/events/{year}/keys",
		"/team/{team_key}/events/{year}/simple",
		"/team/{team_key}/events/{year}/statuses",
		"/team/{team_key}/events/keys",
		"/team/{team_key}/events/simple",
		"/team/{team_key}/matches/{year}",
		"/team/{team_key}/matches/{year}/keys",
		"/team/{team_key}/matches/{year}/simple",
		"/team/{team_key}/media/{year}",
		// "/team/{team_key}/media/tag/{media_tag}",
		// "/team/{team_key}/media/tag/{media_tag}/{year}",
		"/team/{team_key}/robots",
		"/team/{team_key}/simple",
		"/team/{team_key}/social_media",
		"/team/{team_key}/years_participated",
	}
	for _, endpoint := range endpoints {
		endpointURL := strings.Replace(endpoint, "{team_key}", teamKey, -1)
		endpointURL = strings.Replace(endpointURL, "{event_key}", eventKey, -1)
		endpointURL = strings.Replace(endpointURL, "{district_key}", districtKey, -1)
		endpointURL = strings.Replace(endpointURL, "{match_key}", matchKey, -1)
		endpointURL = strings.Replace(endpointURL, "{year}", fmt.Sprintf("%d", year), -1)
		resp, err := ts.Get(ts.URL+endpointURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fname := strings.Replace(endpoint, "{team_key}", teamKey, -1)
		fname = strings.Replace(fname, "{event_key}", eventKey, -1)
		fname = strings.Replace(fname, "{district_key}", districtKey, -1)
		fname = strings.Replace(fname, "{year}", fmt.Sprintf("%d", year), -1)
		fname = strings.Replace(fname, "{match_key}", matchKey, -1)
		fname = strings.Replace(fname, "/", "-", -1)
		err = os.WriteFile(fname[1:]+".json", data, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (ts *TriviaService) Districts(ctx context.Context) ([]*model.District, error) {
	resp, err := ts.Get(ts.URL+"/districts", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r []*model.District
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("unable to parse districts response: %w", err)
	}
	return r, nil
}

func (ts *TriviaService) TeamSimple(ctx context.Context, teamKey string) (*model.TeamSimple, error) {
	resp, err := ts.Get(ts.URL+"/team/"+teamKey+"/simple", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &model.TeamSimple{}
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("unable to parse team response: %w", err)
	}
	return r, nil
}

func (ts *TriviaService) TeamSocialMedia(ctx context.Context, teamKey string) ([]*model.Media, error) {
	resp, err := ts.Get(ts.URL+"/team/"+teamKey+"/social_media", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := []*model.Media{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("unable to parse team media response: %w", err)
	}
	return r, nil
}

func (ts *TriviaService) EventTeams(ctx context.Context, eventKey string) ([]*model.Team, error) {
	resp, err := ts.Get(ts.URL+"/event/"+eventKey+"/teams", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := []*model.Team{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("unable to parse event teams response: %w", err)
	}
	return r, nil
}

func (ts *TriviaService) EventMatchesSimple(ctx context.Context, eventKey string) ([]*model.MatchSimple, error) {
	resp, err := ts.Get(ts.URL+"/event/"+eventKey+"/matches/simple", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := []*model.MatchSimple{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("unable to parse event matches simple response: %w", err)
	}
	return r, nil
}
