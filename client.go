package bamm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
