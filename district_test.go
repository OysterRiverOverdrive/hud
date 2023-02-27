package hud

import (
	"encoding/json"
	"testing"

	"github.com/oysterriveroverdrive/hud/bluealliance/batest"
	"github.com/oysterriveroverdrive/hud/bluealliance/model"
	"github.com/stretchr/testify/assert"
)

func TestDistricts(t *testing.T) {
	s, _ := batest.NewService()
	defer s.Close()

	want := []*model.District{
		{
			Abbreviation: "ne",
			DisplayName:  "New England",
			Key:          "2020ne",
			Year:         2020,
		},
	}
	s.Districts = want

	t.Log(s.Server.URL)
	client := s.Server.Client()
	resp, err := client.Get(s.Server.URL + "/districts")
	if err != nil {
		assert.NoError(t, err)
	}
	defer resp.Body.Close()
	var got []*model.District
	json.NewDecoder(resp.Body).Decode(&got)
	assert.Equal(t, got, want)
}
