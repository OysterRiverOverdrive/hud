package bamm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func genMatches() []*MatchSimple {
	matches := []*MatchSimple{}
	match := &MatchSimple{}
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
	next := &MatchSimple{}
	next.Alliances.Red.TeamKeys = []string{"frc8410", "frc222", "frc333"}
	next.Alliances.Blue.TeamKeys = []string{"frc444", "frc555", "frc666"}
	next.PredictedTime = 1648315545
	got := PrintNextMatchSummary(next, 8410)
	want := []string{}
	assert.Equal(t, want, got)
}
