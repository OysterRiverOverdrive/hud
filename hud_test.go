package hud

import (
	"testing"
	"time"

	"github.com/oysterriveroverdrive/hud/model"
	"github.com/stretchr/testify/assert"
)

type testClock struct {
	now time.Time
}

func (c *testClock) Now() time.Time {
	newYork, _ := time.LoadLocation("America/New_York")
	return c.now.In(newYork)
}

func genMatches() []*model.MatchSimple {
	matches := []*model.MatchSimple{}
	match := &model.MatchSimple{}
	match.Alliances.Red.TeamKeys = []string{"frc8410"}
	match.PredictedTime = 4444
	matches = append(matches, match)

	return matches
}

func TestNextMatch(t *testing.T) {
	matches := genMatches()

	next, err := NextMatch(matches, 8410)
	assert.NoError(t, err, "should have been able to find a match")
	assert.Equal(t, next.PredictedTime, int64(4444))
}

func TestPrintNextMatchSummary(t *testing.T) {
	now := time.Unix(1667779200, 0)
	clock := &testClock{now}
	later := now.Add(10 * time.Minute)
	next := &model.MatchSimple{
		Alliances: struct {
			Red  model.MatchAlliance "json:\"red\""
			Blue model.MatchAlliance "json:\"blue\""
		}{
			Red: model.MatchAlliance{
				TeamKeys: []string{"frc8410", "frc222", "frc333"},
			},
			Blue: model.MatchAlliance{
				TeamKeys: []string{"frc444", "frc555", "frc666"},
			},
		},
		PredictedTime: later.Unix(),
	}
	got := PrintNextMatchSummary(clock, next, 8410)
	want := []string{
		"Match Starts At: 2022-11-07 00:10:00 +0000 UTC",
		"Starts In: 10m0s",
		"Alliance: Red",
		"Allies: frc8410 frc222 frc333",
		"Opponents: frc444 frc555 frc666",
	}
	assert.Equal(t, want, got)
}
