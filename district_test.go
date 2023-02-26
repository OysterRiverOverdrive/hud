package hud

import (
	"encoding/json"
	"testing"

	"github.com/oysterriveroverdrive/hud/hudtest"
	"github.com/oysterriveroverdrive/hud/model"
	"github.com/stretchr/testify/assert"
)

func TestDistricts(t *testing.T) {
	blueAlliance, _ := hudtest.NewBlueAlliance()
	defer blueAlliance.Close()

	want := []*model.District{
		{
			Abbreviation: "ne",
			DisplayName:  "New England",
			Key:          "2020ne",
			Year:         2020,
		},
	}
	blueAlliance.Districts = want

	t.Log(blueAlliance.Server.URL)
	client := blueAlliance.Server.Client()
	resp, err := client.Get(blueAlliance.Server.URL + "/districts")
	if err != nil {
		assert.NoError(t, err)
	}
	defer resp.Body.Close()
	var got []*model.District
	json.NewDecoder(resp.Body).Decode(&got)
	assert.Equal(t, got, want)
}
