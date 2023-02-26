package bluealliance

import (
	"context"
	"testing"

	"github.com/oysterriveroverdrive/hud/bluealliance/batest"
	"github.com/oysterriveroverdrive/hud/bluealliance/model"
	"github.com/stretchr/testify/assert"
)

func TestIt(t *testing.T) {
	s, _ := batest.NewService()
	defer s.Close()

	s.Districts = []*model.District{
		{
			Abbreviation: "ne",
			DisplayName:  "New England",
			Key:          "2020ne",
			Year:         2020,
		},
	}

	ts := NewService(s.Server.Client(), "")
	ts.URL = s.Server.URL
	d, err := ts.Districts(context.Background())
	assert.NoError(t, err)
	t.Log(d)
}
