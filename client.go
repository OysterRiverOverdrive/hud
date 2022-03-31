package bamm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type BAClient struct {
	URL     string
	Client  *http.Client
	AuthKey string
}

func NewBAClient(client *http.Client, AuthKey string) *BAClient {
	return &BAClient{
		URL:     DEFAULT_SERVER,
		Client:  client,
		AuthKey: AuthKey,
	}
}

func (ba *BAClient) setHeaders(req *http.Request) {
	req.Header.Set("X-TBA-Auth-Key", ba.AuthKey)
	req.Header.Set("accept", "application/json")
}

func (ba *BAClient) Get(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}
	ba.setHeaders(req)
	return ba.Client.Do(req)
}

func (c *BAClient) Dump(ctx context.Context, teamNum int, eventKey string) {
	teamKey := fmt.Sprintf("frc%d", teamNum)

	endpoints := []string{
		"/team/{team_key}/event/{event_key}/matches",
		"/team/{team_key}/event/{event_key}/matches/simple",
		"/team/{team_key}/event/{event_key}/matches/keys",
		"/team/{team_key}/event/{event_key}/awards",
		"/team/{team_key}/event/{event_key}/status",
		"/event/{event_key}",
		"/event/{event_key}/simple",
		"/event/{event_key}/alliances",
		"/event/{event_key}/insights",
		"/event/{event_key}/oprs",
		"/event/{event_key}/predictions",
		"/event/{event_key}/rankings",
		"/event/{event_key}/district_points",
		"/event/{event_key}/teams",
		"/event/{event_key}/teams/simple",
		"/event/{event_key}/teams/keys",
		"/event/{event_key}/teams/statuses",
		"/event/{event_key}/teams",
		"/event/{event_key}/teams/simple",
		"/event/{event_key}/teams/keys",
		"/event/{event_key}/teams/statuses",
		"/event/{event_key}/matches",
		"/event/{event_key}/matches/simple",
		"/event/{event_key}/matches/keys",
		"/event/{event_key}/matches/timeseries", // not implemented
		"/event/{event_key}/awards",
	}
	for _, endpoint := range endpoints {
		endpointURL := strings.Replace(endpoint, "{team_key}", teamKey, -1)
		endpointURL = strings.Replace(endpointURL, "{event_key}", eventKey, -1)
		resp, err := c.Get(c.URL+endpointURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fname := strings.Replace(endpoint, "{team_key}", teamKey, -1)
		fname = strings.Replace(fname, "{event_key}", eventKey, -1)
		fname = strings.Replace(fname, "/", "-", -1)
		err = os.WriteFile(fname[1:]+".json", data, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (c *BAClient) TeamSimple(ctx context.Context, teamKey string) (*TeamSimple, error) {
	resp, err := c.Get(c.URL+"/team/"+teamKey+"/simple", nil)
	if err != nil {
		return nil, err
	}

	r := &TeamSimple{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, fmt.Errorf("unable to parse team response: %w", err)
	}
	return r, nil
}

func (c *BAClient) TeamSocialMedia(ctx context.Context, teamKey string) ([]*Media, error) {
	resp, err := c.Get(c.URL+"/team/"+teamKey+"/social_media", nil)
	if err != nil {
		return nil, err
	}

	r := []*Media{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unable to parse team response: %w", err)
	}
	return r, nil
}

func (c *BAClient) EventTeams(ctx context.Context, eventKey string) ([]*Team, error) {
	resp, err := c.Get(c.URL+"/event/"+eventKey+"/teams", nil)
	if err != nil {
		return nil, err
	}

	r := []*Team{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unable to parse team response: %w", err)
	}
	return r, nil
}

func (c *BAClient) EventMatchesSimple(ctx context.Context, eventKey string) ([]*MatchSimple, error) {
	resp, err := c.Get(c.URL+"/event/"+eventKey+"/matches/simple", nil)
	if err != nil {
		return nil, err
	}

	r := []*MatchSimple{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unable to parse team response: %w", err)
	}
	return r, nil
}
