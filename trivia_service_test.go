package hud

import (
	"context"
	"testing"

	"github.com/oysterriveroverdrive/hud/hudtest"
	"github.com/oysterriveroverdrive/hud/model"
	"github.com/stretchr/testify/assert"
)

func TestIt(t *testing.T) {
	blueAlliance, _ := hudtest.NewBlueAlliance()
	defer blueAlliance.Close()

	blueAlliance.Districts = []*model.District{
		{
			Abbreviation: "ne",
			DisplayName:  "New England",
			Key:          "2020ne",
			Year:         2020,
		},
	}

	ts := NewTriviaService(blueAlliance.Server.Client(), "")
	ts.URL = blueAlliance.Server.URL
	d, err := ts.Districts(context.Background())
	assert.NoError(t, err)
	t.Log(d)
}
